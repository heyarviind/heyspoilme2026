<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { auth } from '$lib/stores/auth';
	import { validateImage, compressImage, getWebPFilename } from '$lib/utils/image';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import CityAutocomplete from '$lib/components/CityAutocomplete.svelte';

	const MAX_IMAGES = 6;

	interface Profile {
		id: string;
		user_id: string;
		display_name: string;
		gender: 'male' | 'female';
		age: number;
		bio: string;
		salary_range?: string;
		city: string;
		state: string;
		is_complete: boolean;
		is_verified: boolean;
	}

	interface ProfileImage {
		id: string;
		url: string;
		is_primary: boolean;
	}

	let profile = $state<Profile | null>(null);
	let images = $state<ProfileImage[]>([]);
	let loading = $state(true);
	let editing = $state(false);
	let saving = $state(false);
	let uploading = $state(false);
	let showPricingPopup = $state(false);
	let showDeleteConfirm = $state(false);
	let deleteConfirmText = $state('');
	let deleting = $state(false);
	
	let authState = $state<any>(null);
	auth.subscribe(s => authState = s);
	let isEmailVerified = $derived(authState?.user?.email_verified ?? false);
	let isAuthReady = $derived(authState?.initialized && !authState?.loading);

	// Edit form state
	let editDisplayName = $state('');
	let editAge = $state(25);
	let editBio = $state('');
	let editCity = $state('');
	let editState = $state('');
	let editSalaryRange = $state('');
	let editLatitude = $state(0);
	let editLongitude = $state(0);

	const salaryOptions = ['5-10 LPA', '10-20 LPA', '20-50 LPA', '50+ LPA'];

	// Handle city selection from autocomplete
	function handleCitySelect(selectedCity: { city: string; state: string; latitude: number; longitude: number }) {
		editCity = selectedCity.city;
		editState = selectedCity.state;
		editLatitude = selectedCity.latitude;
		editLongitude = selectedCity.longitude;
	}

	async function loadProfile() {
		loading = true;
		try {
			const data = await api.getMyProfile() as { profile: Profile; images: ProfileImage[] };
			profile = data.profile;
			images = data.images || [];

			if (profile) {
				editDisplayName = profile.display_name;
				editAge = profile.age;
				editBio = profile.bio;
				editCity = profile.city;
				editState = profile.state;
				editSalaryRange = profile.salary_range || '';
			}
		} catch (e: any) {
			if (e.message?.includes('not found')) {
				goto('/profile/setup');
			}
		} finally {
			loading = false;
		}
	}

	async function saveProfile() {
		if (!profile) return;
		saving = true;
		try {
			const updates: any = {
				display_name: editDisplayName,
				age: editAge,
				bio: editBio,
				city: editCity,
				state: editState,
				latitude: editLatitude,
				longitude: editLongitude,
			};
			if (profile.gender === 'male') {
				updates.salary_range = editSalaryRange;
			}
			await api.updateProfile(updates);
			await loadProfile();
			await auth.refreshProfile();
			editing = false;
		} catch (e) {
			console.error('Failed to save profile:', e);
		} finally {
			saving = false;
		}
	}

	async function handleImageUpload(e: Event) {
		const input = e.target as HTMLInputElement;
		if (!input.files?.length) return;

		// Check email verification
		if (!isEmailVerified) {
			alert('Please verify your email to upload photos');
			input.value = '';
			return;
		}

		const file = input.files[0];
		
		// Check max images limit
		if (images.length >= MAX_IMAGES) {
			alert(`Maximum ${MAX_IMAGES} photos allowed. Delete one to add more.`);
			input.value = '';
			return;
		}

		// Validate image
		const validation = validateImage(file);
		if (!validation.valid) {
			alert(validation.error);
			input.value = '';
			return;
		}

		uploading = true;
		try {
			// Compress and convert to WebP
			const { blob } = await compressImage(file, {
				maxWidth: 1200,
				maxHeight: 1600,
				quality: 0.85,
			});

			const webpFilename = getWebPFilename(file.name);

			// Get presigned URL for WebP
			const urlData = await api.getPresignedUrl(webpFilename, 'image/webp') as {
				upload_url: string;
				s3_key: string;
				public_url: string;
			};

			// Upload compressed WebP to S3
			const uploadResponse = await fetch(urlData.upload_url, {
				method: 'PUT',
				body: blob,
				headers: { 'Content-Type': 'image/webp' },
			});

			if (!uploadResponse.ok) {
				throw new Error(`Upload failed: ${uploadResponse.status}`);
			}

			// Add to profile
			await api.addProfileImage(urlData.s3_key, urlData.public_url, images.length === 0);
			await loadProfile();
		} catch (e: any) {
			console.error('Failed to upload image:', e);
			if (e.message?.includes('too large')) {
				alert('Image is too large. Please use a smaller image.');
			} else if (e.message?.includes('network') || e.message?.includes('fetch')) {
				alert('Network error. Please check your connection and try again.');
			} else {
				alert('Failed to upload image. Please try again.');
			}
		} finally {
			uploading = false;
			input.value = '';
		}
	}

	async function deleteImage(imageId: string) {
		if (!confirm('Delete this photo?')) return;
		try {
			await api.deleteProfileImage(imageId);
			await loadProfile();
		} catch (e) {
			console.error('Failed to delete image:', e);
		}
	}

	async function logout() {
		auth.logout();
		goto('/auth/login');
	}

	async function deleteAccount() {
		if (deleteConfirmText !== 'DELETE') return;
		
		deleting = true;
		try {
			await api.deleteAccount();
			auth.logout();
			goto('/');
		} catch (e) {
			console.error('Failed to delete account:', e);
			alert('Failed to delete account. Please try again.');
		} finally {
			deleting = false;
		}
	}

	onMount(() => {
		loadProfile();
	});
</script>

<svelte:head>
	<title>My Profile | HeySpoilMe</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;500;600;700&family=Montserrat:wght@300;400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="profile-page">
	<Header />

	<main class="main">
		{#if loading}
			<div class="loading">
				<div class="spinner"></div>
			</div>
		{:else if profile}
			<div class="profile-content">
				{#if profile.gender === 'male' && !profile.is_verified}
					<div class="upgrade-banner">
						<img src="/img/crown.png" alt="Crown" class="banner-crown" />
						<div class="banner-content">
							<h2>Upgrade to Trusted Member</h2>
							<p>Stand out as a verified, serious member</p>
						</div>
						<button class="upgrade-btn" onclick={() => showPricingPopup = true}>
							Unlock Trusted Status
						</button>
					</div>
				{/if}

				<div class="photos-section">
					<h2>Photos</h2>
					<div class="photos-grid">
						{#each images as image}
							<div class="photo-item">
								<img src={image.url} alt="Profile" />
								{#if image.is_primary}
									<span class="primary-badge">Primary</span>
								{/if}
								<button class="delete-btn" onclick={() => deleteImage(image.id)}>×</button>
							</div>
						{/each}
						{#if images.length < 6}
							<label class="add-photo" class:disabled={isAuthReady && !isEmailVerified}>
								<input 
									type="file" 
									accept="image/*" 
									onchange={handleImageUpload}
									disabled={uploading || (isAuthReady && !isEmailVerified)}
								/>
								{#if isAuthReady && !isEmailVerified}
									<span class="verify-hint">🔒 Verify email</span>
								{:else if uploading}
									<span class="uploading">Uploading...</span>
								{:else}
									<span>+ Add Photo</span>
								{/if}
							</label>
						{/if}
					</div>
				</div>

				<div class="info-section">
					<div class="section-header">
						<h2>Profile Info</h2>
						{#if !editing}
							<button class="edit-btn" onclick={() => editing = true}>Edit</button>
						{/if}
					</div>

					{#if editing}
						<div class="edit-form">
					<div class="form-group">
						<label for="edit-name">Display Name</label>
						<input id="edit-name" type="text" bind:value={editDisplayName} maxlength="50" placeholder="Your display name" />
					</div>
					<div class="form-group">
						<label for="edit-age">Age</label>
						<input id="edit-age" type="number" bind:value={editAge} min="21" max="100" />
					</div>
					<div class="form-group">
						<label for="edit-bio">Bio</label>
						<textarea id="edit-bio" bind:value={editBio} rows="4" maxlength="500"></textarea>
					</div>
					<div class="form-group">
						<label for="edit-city">City</label>
						<CityAutocomplete 
							bind:value={editCity}
							onSelect={handleCitySelect}
							placeholder="Start typing your city..."
						/>
					</div>
					{#if editState}
						<div class="form-group">
							<label>State</label>
							<div class="state-display">{editState}</div>
						</div>
					{/if}
					{#if profile.gender === 'male'}
						<div class="form-group">
							<label for="edit-salary">Salary Range</label>
							<select id="edit-salary" bind:value={editSalaryRange}>
								{#each salaryOptions as option}
									<option value={option}>{option}</option>
								{/each}
							</select>
						</div>
					{/if}
							<div class="form-actions">
								<button class="btn-secondary" onclick={() => editing = false}>Cancel</button>
								<button class="btn-primary" onclick={saveProfile} disabled={saving}>
									{saving ? 'Saving...' : 'Save Changes'}
								</button>
							</div>
						</div>
					{:else}
						<div class="profile-info">
							<div class="info-row">
								<span class="label">Display Name</span>
								<span class="value">{profile.display_name}</span>
							</div>
							<div class="info-row">
								<span class="label">Gender</span>
								<span class="value">{profile.gender === 'male' ? 'Man' : 'Woman'}</span>
							</div>
							<div class="info-row">
								<span class="label">Age</span>
								<span class="value">{profile.age}</span>
							</div>
							<div class="info-row">
								<span class="label">Location</span>
								<span class="value">{profile.city}, {profile.state}</span>
							</div>
							{#if profile.gender === 'male' && profile.salary_range}
								<div class="info-row">
									<span class="label">Salary Range</span>
									<span class="value">{profile.salary_range}</span>
								</div>
							{/if}
							<div class="info-row bio">
								<span class="label">About</span>
								<span class="value">{profile.bio}</span>
							</div>
						</div>
					{/if}
				</div>

				{#if profile.gender === 'male' && profile.is_verified}
					<div class="verification-section verified">
						<div class="verified-status">
							<span class="verified-badge">✓ TRUSTED MEMBER</span>
							<p>Your profile is highlighted with premium styling</p>
						</div>
					</div>
				{/if}

				<div class="account-section">
					<h2>Account</h2>
					<button class="logout-btn" onclick={logout}>Sign Out</button>
					<button class="delete-account-btn" onclick={() => showDeleteConfirm = true}>Delete Account</button>
				</div>
			</div>
		{/if}
	</main>

	<Footer />

	{#if showDeleteConfirm}
		<div class="popup-overlay" onclick={() => showDeleteConfirm = false} role="button" tabindex="0" onkeypress={(e) => e.key === 'Enter' && (showDeleteConfirm = false)}>
			<div class="popup-content delete-popup" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true">
				<button class="popup-close" onclick={() => showDeleteConfirm = false}>×</button>
				
				<div class="popup-header">
					<span class="warning-icon">⚠️</span>
					<h2>Delete Account</h2>
				</div>

				<div class="delete-warning">
					<p><strong>This action is permanent and cannot be undone.</strong></p>
					<p>All your data will be deleted including:</p>
					<ul>
						<li>Your profile and photos</li>
						<li>All messages and conversations</li>
						<li>Likes given and received</li>
						<li>Notifications</li>
					</ul>
				</div>

				<div class="confirm-input">
					<label for="delete-confirm">Type <strong>DELETE</strong> to confirm:</label>
					<input 
						id="delete-confirm"
						type="text" 
						bind:value={deleteConfirmText}
						placeholder="DELETE"
						autocomplete="off"
					/>
				</div>

				<button 
					class="confirm-delete-btn" 
					onclick={deleteAccount}
					disabled={deleteConfirmText !== 'DELETE' || deleting}
				>
					{deleting ? 'Deleting...' : 'Permanently Delete Account'}
				</button>
			</div>
		</div>
	{/if}

	{#if showPricingPopup}
		<div class="popup-overlay" onclick={() => showPricingPopup = false} role="button" tabindex="0" onkeypress={(e) => e.key === 'Enter' && (showPricingPopup = false)}>
			<div class="popup-content" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true">
				<button class="popup-close" onclick={() => showPricingPopup = false}>×</button>
				
				<div class="popup-header">
					<img src="/img/crown.png" alt="Crown" class="crown-img" />
					<h2>Upgrade to Trusted Member</h2>
				</div>

				<div class="pricing-card">
					<div class="price">
						<span class="currency">₹</span>
						<span class="amount">5,999</span>
						<span class="period">/month</span>
					</div>
				</div>

				<p class="pricing-tagline">Stand out as a verified, serious member in a private community built on trust and discretion.</p>

				<ul class="benefits">
					<li>
						<span class="check">✓</span>
						<span>Trusted badge that signals credibility</span>
					</li>
					<li>
						<span class="check">✓</span>
						<span>Priority placement in member discovery</span>
					</li>
					<li>
						<span class="check">✓</span>
						<span>Increased visibility to verified profiles</span>
					</li>
					<li>
						<span class="check">✓</span>
						<span>Unlock message requests & replies</span>
					</li>
					<li>
						<span class="check">✓</span>
						<span>Designed for members who value quality over volume</span>
					</li>
				</ul>

				<button class="subscribe-btn" onclick={() => goto('/profile/verify')}>
					Unlock Trusted Status
				</button>

				<p class="terms">
					By subscribing, you agree to our terms. Cancel anytime.
				</p>
			</div>
		</div>
	{/if}
</div>

<style>
	:global(body) {
		font-family: 'Montserrat', sans-serif;
		background: #0a0a0a;
		color: #fff;
		margin: 0;
	}

	.profile-page {
		min-height: 100vh;
	}

	.main {
		max-width: 700px;
		margin: 0 auto;
		padding: 2rem;
	}

	.loading {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 50vh;
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 3px solid rgba(255, 255, 255, 0.1);
		border-top-color: #fff;
		border-radius: 0;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.profile-content {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	/* Upgrade Banner - Top of profile for unverified males */
	.upgrade-banner {
		display: flex;
		align-items: center;
		gap: 1.25rem;
		padding: 1.25rem 1.5rem;
		background: linear-gradient(135deg, rgba(251, 191, 36, 0.15) 0%, rgba(251, 191, 36, 0.05) 100%);
		border: 1px solid rgba(251, 191, 36, 0.4);
	}

	.banner-crown {
		width: 48px;
		height: 48px;
		object-fit: contain;
		flex-shrink: 0;
	}

	.banner-content {
		flex: 1;
	}

	.banner-content h2 {
		font-size: 1.1rem;
		margin: 0 0 0.25rem;
		color: #fbbf24;
	}

	.banner-content p {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.7);
		margin: 0;
	}

	.upgrade-btn {
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
		color: #000;
		border: none;
		padding: 0.75rem 1.5rem;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.85rem;
		font-weight: 600;
		cursor: pointer;
		white-space: nowrap;
		transition: transform 0.2s ease, box-shadow 0.2s ease;
	}

	.upgrade-btn:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(251, 191, 36, 0.3);
	}

	@media (max-width: 600px) {
		.upgrade-banner {
			flex-direction: column;
			text-align: center;
			gap: 1rem;
		}

		.upgrade-btn {
			width: 100%;
		}
	}

	.photos-section, .info-section, .account-section {
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		padding: 1.5rem;
	}

	h2 {
		font-family: 'Playfair Display', serif;
		font-size: 1.25rem;
		margin: 0 0 1rem 0;
	}

	.photos-grid {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 1rem;
	}

	.photo-item {
		position: relative;
		aspect-ratio: 1;
		border-radius: 0;
		overflow: hidden;
	}

	.photo-item img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.primary-badge {
		position: absolute;
		bottom: 0.5rem;
		left: 0.5rem;
		background: #fff;
		color: #000;
		font-size: 0.65rem;
		font-weight: 600;
		padding: 0.2rem 0.4rem;
		border-radius: 0;
	}

	.delete-btn {
		position: absolute;
		top: 0.5rem;
		right: 0.5rem;
		width: 24px;
		height: 24px;
		background: rgba(0, 0, 0, 0.7);
		border: none;
		border-radius: 0;
		color: #fff;
		font-size: 1rem;
		cursor: pointer;
		opacity: 0;
		transition: opacity 0.2s ease;
	}

	.photo-item:hover .delete-btn {
		opacity: 1;
	}

	.add-photo {
		aspect-ratio: 1;
		border: 2px dashed rgba(255, 255, 255, 0.2);
		border-radius: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		color: rgba(255, 255, 255, 0.5);
		font-size: 0.85rem;
		transition: border-color 0.2s ease;
	}

	.add-photo:hover:not(.disabled) {
		border-color: rgba(255, 255, 255, 0.4);
	}

	.add-photo.disabled {
		cursor: not-allowed;
		border-color: rgba(255, 56, 92, 0.3);
		background: rgba(255, 56, 92, 0.05);
	}

	.add-photo input {
		display: none;
	}

	.verify-hint {
		font-size: 0.75rem;
		color: #FF385C;
	}

	.section-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
	}

	.section-header h2 {
		margin: 0;
	}

	.edit-btn {
		background: none;
		border: 1px solid rgba(255, 255, 255, 0.2);
		color: #fff;
		padding: 0.5rem 1rem;
		border-radius: 0;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.85rem;
		cursor: pointer;
	}

	.profile-info {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.info-row {
		display: flex;
		justify-content: space-between;
	}

	.info-row.bio {
		flex-direction: column;
		gap: 0.5rem;
	}

	.label {
		color: rgba(255, 255, 255, 0.5);
		font-size: 0.85rem;
	}

	.value {
		color: #fff;
	}

	.edit-form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-group label {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.6);
	}

	.form-group input, .form-group textarea, .form-group select {
		padding: 0.75rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
	}

	.form-group select {
		appearance: none;
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='white' d='M6 8L1 3h10z'/%3E%3C/svg%3E");
		background-repeat: no-repeat;
		background-position: right 1rem center;
	}

	.state-display {
		padding: 0.75rem;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.08);
		color: rgba(255, 255, 255, 0.7);
		font-size: 1rem;
	}

	.form-actions {
		display: flex;
		gap: 1rem;
		margin-top: 0.5rem;
	}

	.btn-primary, .btn-secondary {
		flex: 1;
		padding: 0.75rem;
		border-radius: 0;
		font-family: 'Montserrat', sans-serif;
		font-weight: 500;
		cursor: pointer;
		border: none;
	}

	.btn-primary {
		background: #fff;
		color: #000;
	}

	.btn-secondary {
		background: transparent;
		color: #fff;
		border: 1px solid rgba(255, 255, 255, 0.2);
	}

	.logout-btn {
		width: 100%;
		padding: 0.75rem;
		background: transparent;
		border: 1px solid rgba(255, 100, 100, 0.3);
		border-radius: 0;
		color: #ff6666;
		font-family: 'Montserrat', sans-serif;
		cursor: pointer;
	}

	.logout-btn:hover {
		background: rgba(255, 100, 100, 0.1);
	}

	.delete-account-btn {
		width: 100%;
		padding: 0.75rem;
		margin-top: 0.75rem;
		background: transparent;
		border: 1px solid rgba(255, 50, 50, 0.3);
		border-radius: 0;
		color: #ff4444;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.9rem;
		cursor: pointer;
	}

	.delete-account-btn:hover {
		background: rgba(255, 50, 50, 0.1);
	}

	.delete-popup {
		max-width: 400px;
	}

	.warning-icon {
		font-size: 2.5rem;
		display: block;
		margin-bottom: 0.75rem;
	}

	.delete-popup h2 {
		color: #ff4444;
		margin: 0;
	}

	.delete-warning {
		background: rgba(255, 50, 50, 0.1);
		border: 1px solid rgba(255, 50, 50, 0.2);
		padding: 1rem;
		margin-bottom: 1.5rem;
	}

	.delete-warning p {
		color: rgba(255, 255, 255, 0.8);
		margin: 0 0 0.5rem 0;
		font-size: 0.9rem;
	}

	.delete-warning ul {
		margin: 0.5rem 0 0 0;
		padding-left: 1.25rem;
		color: rgba(255, 255, 255, 0.6);
		font-size: 0.85rem;
	}

	.delete-warning li {
		margin: 0.25rem 0;
	}

	.confirm-input {
		margin-bottom: 1.5rem;
	}

	.confirm-input label {
		display: block;
		color: rgba(255, 255, 255, 0.7);
		font-size: 0.85rem;
		margin-bottom: 0.5rem;
	}

	.confirm-input input {
		width: 100%;
		padding: 0.75rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		font-size: 1rem;
		text-transform: uppercase;
		letter-spacing: 1px;
	}

	.confirm-input input:focus {
		outline: none;
		border-color: #ff4444;
	}

	.confirm-delete-btn {
		width: 100%;
		padding: 1rem;
		background: #ff4444;
		border: none;
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		font-size: 0.9rem;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.confirm-delete-btn:disabled {
		background: rgba(255, 68, 68, 0.3);
		cursor: not-allowed;
	}

	.confirm-delete-btn:not(:disabled):hover {
		background: #ff2222;
	}

	/* Verification Section */
	.verification-section {
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		padding: 1.5rem;
	}

	.verification-section.verified {
		background: linear-gradient(135deg, rgba(251, 191, 36, 0.15) 0%, rgba(251, 191, 36, 0.05) 100%);
		border: 2px solid #fbbf24;
	}

	.verified-status {
		text-align: center;
	}

	.verified-badge {
		display: inline-block;
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
		color: #000;
		font-size: 0.85rem;
		font-weight: 700;
		padding: 0.5rem 1rem;
		letter-spacing: 0.5px;
		margin-bottom: 0.5rem;
	}

	.verified-status p {
		color: rgba(255, 255, 255, 0.6);
		font-size: 0.9rem;
		margin: 0.5rem 0 0 0;
	}

	.unverified-status {
		text-align: center;
	}

	.unverified-status h2 {
		font-family: 'Playfair Display', serif;
		margin: 0 0 0.5rem 0;
	}

	.unverified-status p {
		color: rgba(255, 255, 255, 0.6);
		font-size: 0.9rem;
		margin: 0 0 1rem 0;
	}

	.verify-btn {
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
		color: #000;
		border: none;
		padding: 0.875rem 2rem;
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		font-size: 0.9rem;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.verify-btn:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 20px rgba(251, 191, 36, 0.4);
	}

	/* Pricing Popup */
	.popup-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.85);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 1rem;
	}

	.popup-content {
		background: linear-gradient(180deg, #1a1a1a 0%, #0d0d0d 100%);
		border: 1px solid rgba(251, 191, 36, 0.3);
		max-width: 420px;
		width: 100%;
		padding: 2.5rem;
		position: relative;
		box-shadow: 0 0 60px rgba(251, 191, 36, 0.15);
	}

	.popup-close {
		position: absolute;
		top: 1rem;
		right: 1rem;
		width: 36px;
		height: 36px;
		background: rgba(255, 255, 255, 0.1);
		border: none;
		color: #fff;
		font-size: 1.5rem;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.popup-header {
		text-align: center;
		margin-bottom: 2rem;
	}

	.crown-img {
		width: 64px;
		height: 64px;
		object-fit: contain;
		display: block;
		margin: 0 auto 1rem;
	}

	.popup-header h2 {
		font-family: 'Playfair Display', serif;
		font-size: 1.5rem;
		margin: 0;
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
	}

	.pricing-card {
		background: linear-gradient(135deg, rgba(251, 191, 36, 0.15) 0%, rgba(251, 191, 36, 0.05) 100%);
		border: 1px solid rgba(251, 191, 36, 0.3);
		padding: 1.5rem;
		text-align: center;
		margin-bottom: 1rem;
	}

	.pricing-tagline {
		text-align: center;
		color: rgba(255, 255, 255, 0.7);
		font-size: 0.9rem;
		line-height: 1.6;
		margin: 0 0 1.5rem;
		padding: 0 0.5rem;
	}

	.price {
		display: flex;
		align-items: baseline;
		justify-content: center;
		gap: 0.25rem;
	}

	.currency {
		font-size: 1.5rem;
		color: #fbbf24;
	}

	.amount {
		font-size: 3rem;
		font-weight: 700;
		color: #fbbf24;
		font-family: 'Playfair Display', serif;
	}

	.period {
		font-size: 1rem;
		color: rgba(255, 255, 255, 0.5);
	}

	.benefits {
		list-style: none;
		padding: 0;
		margin: 0 0 2rem 0;
	}

	.benefits li {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 0;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
	}

	.benefits li:last-child {
		border-bottom: none;
	}

	.benefits .check {
		color: #fbbf24;
		font-weight: bold;
	}

	.benefits span:last-child {
		color: rgba(255, 255, 255, 0.8);
		font-size: 0.9rem;
	}

	.subscribe-btn {
		width: 100%;
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
		color: #000;
		border: none;
		padding: 1rem;
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
		font-size: 1rem;
		cursor: pointer;
		transition: all 0.2s ease;
		letter-spacing: 0.5px;
	}

	.subscribe-btn:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 20px rgba(251, 191, 36, 0.5);
	}

	.terms {
		text-align: center;
		color: rgba(255, 255, 255, 0.4);
		font-size: 0.75rem;
		margin: 1rem 0 0 0;
	}

	@media (max-width: 768px) {
		.photos-grid {
			grid-template-columns: repeat(2, 1fr);
		}

		.popup-overlay {
			padding: 0;
			background: #0a0a0a;
		}

		.popup-content {
			max-width: none;
			height: 100%;
			display: flex;
			flex-direction: column;
			padding: 0;
			padding-bottom: 5rem;
			background: #0a0a0a;
			border: none;
			box-shadow: none;
			overflow-y: auto;
			-webkit-overflow-scrolling: touch;
		}

		.popup-close {
			position: fixed;
			top: 1rem;
			right: 1rem;
			z-index: 10;
			background: transparent;
		}

		.popup-header {
			padding: 4rem 1.5rem 1.5rem;
			margin-bottom: 0;
			flex-shrink: 0;
		}

		.crown-img {
			width: 80px;
			height: 80px;
		}

		.popup-header h2 {
			font-size: 1.75rem;
		}

		.pricing-card {
			margin: 0 1.5rem 1rem;
			flex-shrink: 0;
		}

		.pricing-tagline {
			padding: 0 1.5rem;
			margin-bottom: 1rem;
			flex-shrink: 0;
		}

		.benefits {
			padding: 0 1.5rem;
			margin-bottom: 2rem;
			flex-shrink: 0;
		}

		.popup-content .terms {
			display: none;
		}

		.subscribe-btn {
			position: fixed;
			bottom: 0;
			left: 0;
			right: 0;
			padding: 1.25rem;
			font-size: 1rem;
		}

		.amount {
			font-size: 2.5rem;
		}
	}
</style>

