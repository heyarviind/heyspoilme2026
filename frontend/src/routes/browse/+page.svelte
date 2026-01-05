<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page as pageStore } from '$app/stores';
	import { api } from '$lib/api';
	import { auth } from '$lib/stores/auth';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import HeartIcon from '$lib/components/HeartIcon.svelte';
	import VerificationModal from '$lib/components/VerificationModal.svelte';

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
		is_online: boolean;
		is_verified: boolean;
		last_seen?: string;
		distance_km?: number;
		is_liked: boolean;
		images: Array<{ id: string; url: string; is_primary: boolean }>;
	}

	let profiles = $state<Profile[]>([]);
	let loading = $state(true);
	let page = $state(1);
	let total = $state(0);
	let userGender = $state<'male' | 'female' | null>(null);
	let filters = $state({
		minAge: 21,
		maxAge: 60,
		maxDistance: 0,
		onlineOnly: false,
	});

	// Get user's gender to show opposite
	let authState = $state<any>(null);
	auth.subscribe(s => {
		authState = s;
		if (s?.profile?.gender) {
			userGender = s.profile.gender;
		}
	});

	let isEmailVerified = $derived(authState?.user?.email_verified ?? false);
	let isAuthReady = $derived(authState?.initialized && !authState?.loading);
	let resending = $state(false);
	let resent = $state(false);
	let showVerificationModal = $state(false);

	async function resendVerification() {
		if (resending || resent) return;
		resending = true;
		try {
			await api.resendVerificationEmail();
			resent = true;
			setTimeout(() => { resent = false; }, 60000);
		} catch (e) {
			console.error('Failed to resend verification:', e);
		} finally {
			resending = false;
		}
	}

	// Read filters and page from URL on page load
	function initFiltersFromUrl() {
		const url = new URL(window.location.href);
		const minAge = url.searchParams.get('minAge');
		const maxAge = url.searchParams.get('maxAge');
		const maxDistance = url.searchParams.get('maxDistance');
		const onlineOnly = url.searchParams.get('onlineOnly');
		const pageParam = url.searchParams.get('page');

		if (minAge) filters.minAge = parseInt(minAge) || 21;
		if (maxAge) filters.maxAge = parseInt(maxAge) || 60;
		if (maxDistance) filters.maxDistance = parseInt(maxDistance) || 0;
		if (onlineOnly === 'true') filters.onlineOnly = true;
		if (pageParam) page = parseInt(pageParam) || 1;
	}

	// Update URL when filters or page change
	function updateUrlParams(options: { replaceState?: boolean } = {}) {
		const url = new URL(window.location.href);
		
		if (filters.minAge > 21) {
			url.searchParams.set('minAge', filters.minAge.toString());
		} else {
			url.searchParams.delete('minAge');
		}
		
		if (filters.maxAge < 60) {
			url.searchParams.set('maxAge', filters.maxAge.toString());
		} else {
			url.searchParams.delete('maxAge');
		}
		
		if (filters.maxDistance > 0) {
			url.searchParams.set('maxDistance', filters.maxDistance.toString());
		} else {
			url.searchParams.delete('maxDistance');
		}
		
		if (filters.onlineOnly) {
			url.searchParams.set('onlineOnly', 'true');
		} else {
			url.searchParams.delete('onlineOnly');
		}

		if (page > 1) {
			url.searchParams.set('page', page.toString());
		} else {
			url.searchParams.delete('page');
		}

		goto(url.pathname + url.search, { replaceState: options.replaceState ?? false, noScroll: true });
	}

	function onFilterChange() {
		page = 1; // Reset to page 1 when filters change
		updateUrlParams({ replaceState: true });
		loadProfiles();
	}

	function loadMoreProfiles() {
		page++;
		updateUrlParams({ replaceState: false }); // Push new history entry for back button
		window.scrollTo({ top: 0, behavior: 'smooth' });
		loadProfiles();
	}

	async function loadProfiles() {
		if (!isEmailVerified) {
			// Show message about verification requirement
			loading = false;
			profiles = [];
			return;
		}
		
		loading = true;
		try {
			const params: Record<string, any> = { page, limit: 12 };
			
			// Show opposite gender only
			if (userGender === 'male') {
				params.gender = 'female';
			} else if (userGender === 'female') {
				params.gender = 'male';
			}
			
			if (filters.minAge > 21) params.min_age = filters.minAge;
			if (filters.maxAge < 60) params.max_age = filters.maxAge;
			if (filters.maxDistance) params.max_distance = filters.maxDistance;
			if (filters.onlineOnly) params.online_only = true;

			const data = await api.listProfiles(params) as { profiles: Profile[]; total: number };
			profiles = data.profiles || [];
			total = data.total;
		} catch (e: any) {
			console.error('Failed to load profiles:', e);
			// Handle email verification error from backend
			if (e.message?.includes('email_not_verified') || e.message?.includes('verify your email')) {
				profiles = [];
			}
		} finally {
			loading = false;
		}
	}

	async function toggleLike(profile: Profile) {
		if (!isEmailVerified) {
			showVerificationModal = true;
			return;
		}
		
		try {
			if (profile.is_liked) {
				await api.unlikeProfile(profile.user_id);
			} else {
				await api.likeProfile(profile.user_id);
			}
			profile.is_liked = !profile.is_liked;
		} catch (e: any) {
			const errMsg = e.message || '';
			if (errMsg.includes('verification') || errMsg.includes('person_verification') || 
			    errMsg.includes('email_not_verified') || errMsg.includes('verify your email')) {
				showVerificationModal = true;
			} else {
				console.error('Failed to toggle like:', e);
			}
		}
	}

	function getPrimaryImage(profile: Profile): string | null {
		const primary = profile.images?.find(img => img.is_primary);
		if (primary) return primary.url;
		if (profile.images?.length > 0) return profile.images[0].url;
		return null;
	}

	function formatDistance(km?: number): string {
		if (!km) return '';
		if (km < 1) return '<1 km away';
		if (km < 100) return `${Math.round(km)} km away`;
		return `${Math.round(km)} km away`;
	}

	function formatLastSeen(lastSeen?: string): string {
		if (!lastSeen) return '';
		const date = new Date(lastSeen);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);
		
		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins}m ago`;
		const diffHours = Math.floor(diffMins / 60);
		if (diffHours < 24) return `${diffHours}h ago`;
		const diffDays = Math.floor(diffHours / 24);
		return `${diffDays}d ago`;
	}

	onMount(() => {
		// Wait for auth state to be ready
		const unsub = auth.subscribe(s => {
			if (s.initialized) {
				initFiltersFromUrl();
				loadProfiles();
				unsub();
			}
		});

		// Handle back/forward browser navigation
		const handlePopState = () => {
			initFiltersFromUrl();
			loadProfiles();
		};
		window.addEventListener('popstate', handlePopState);
		return () => window.removeEventListener('popstate', handlePopState);
	});
</script>

<svelte:head>
	<title>Browse | HeySpoilMe</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;500;600;700&family=Montserrat:wght@300;400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="browse-page">
	<Header />

	<main class="main">
		<div class="filters">
			<div class="filter-group age-filter">
				<label>Age</label>
				<div class="age-inputs">
					<input 
						type="number" 
						min="21" 
						max="60"
						bind:value={filters.minAge}
						onchange={onFilterChange}
						class="age-input"
					/>
					<span>-</span>
					<input 
						type="number" 
						min="21" 
						max="60"
						bind:value={filters.maxAge}
						onchange={onFilterChange}
						class="age-input"
					/>
				</div>
			</div>
			<select bind:value={filters.maxDistance} onchange={onFilterChange} class="filter-select">
				<option value={0}>Any Distance</option>
				<option value={25}>Within 25 km</option>
				<option value={50}>Within 50 km</option>
				<option value={100}>Within 100 km</option>
			</select>
			<label class="online-filter">
				<input 
					type="checkbox" 
					bind:checked={filters.onlineOnly}
					onchange={onFilterChange}
				/>
				<span>Online Now</span>
			</label>
		</div>

		{#if !isAuthReady}
			<div class="loading">
				<div class="spinner"></div>
				<p>Loading...</p>
			</div>
		{:else if !isEmailVerified}
			<div class="verification-required">
				<div class="verification-box">
					<span class="icon">🔒</span>
					<h2>Verify Your Email</h2>
					<p>Please verify your email address to browse profiles. Check your inbox for the verification link.</p>
					<button class="resend-link" onclick={resendVerification} disabled={resending || resent}>
						{#if resending}
							Sending...
						{:else if resent}
							Email Sent ✓
						{:else}
							Resend verification email
						{/if}
					</button>
				</div>
			</div>
		{:else if loading}
			<div class="loading">
				<div class="spinner"></div>
				<p>Finding people near you...</p>
			</div>
		{:else if profiles.length === 0}
			<div class="empty">
				<p>No profiles found</p>
				<p class="hint">Try adjusting your filters</p>
			</div>
		{:else}
			<div class="profiles-grid">
				{#each profiles as profile}
					<div class="profile-card" class:verified={profile.is_verified}>
						<a href="/profile/{profile.user_id}" class="card-link">
							<div class="image-container">
								{#if getPrimaryImage(profile)}
									<img src={getPrimaryImage(profile)} alt="{profile.age}" class="profile-image" />
								{:else}
									<div class="no-image-placeholder">
										<span>No image uploaded</span>
									</div>
								{/if}
								{#if !profile.is_verified}
									<span class="not-verified-badge">NOT VERIFIED</span>
								{/if}
								{#if profile.is_online}
									<span class="online-badge"></span>
								{/if}
								{#if profile.distance_km}
									<span class="distance-badge">{formatDistance(profile.distance_km)}</span>
								{/if}
							</div>
							<div class="card-content">
								<div class="card-header">
									<span class="name">
										{profile.display_name}, {profile.age}
										{#if profile.is_verified}
											<svg class="verified-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
												<path fill-rule="evenodd" clip-rule="evenodd" d="M12 1.25C11.4388 1.25 10.9816 1.48611 10.5656 1.80358C10.1759 2.10089 9.74606 2.53075 9.24319 3.03367L9.20782 3.06904C8.69316 3.5837 8.24449 3.78626 7.55208 3.78626C7.4653 3.78626 7.35579 3.78318 7.23212 3.7797C6.91509 3.77078 6.50497 3.75924 6.14615 3.79027C5.62128 3.83566 4.96532 3.97929 4.46745 4.48134C3.9734 4.97955 3.83327 5.63282 3.78923 6.15439C3.75922 6.50995 3.77075 6.91701 3.77966 7.23178L3.77966 7.23181C3.78317 7.35581 3.78628 7.46549 3.78628 7.55206C3.78628 8.24448 3.58371 8.69315 3.06902 9.20784L3.03367 9.24319C2.53075 9.74606 2.10089 10.1759 1.80358 10.5655C1.48612 10.9816 1.25001 11.4388 1.25 12C1.25001 12.5611 1.48613 13.0183 1.80358 13.4344C2.10095 13.8242 2.53091 14.2541 3.03395 14.7571L3.06906 14.7922C3.40272 15.1258 3.56011 15.3422 3.64932 15.5464C3.73619 15.7453 3.78628 15.9971 3.78628 16.4479C3.78628 16.5347 3.7832 16.6442 3.77972 16.7679C3.7708 17.0849 3.75926 17.495 3.79029 17.8539C3.83569 18.3787 3.97933 19.0347 4.48139 19.5326C4.97961 20.0266 5.63287 20.1667 6.15443 20.2107C6.50997 20.2408 6.91703 20.2292 7.23179 20.2203C7.35581 20.2168 7.4655 20.2137 7.55206 20.2137C7.99328 20.2137 8.24126 20.2581 8.43645 20.3386C8.63147 20.4191 8.84006 20.5632 9.15424 20.8774C9.22129 20.9444 9.30963 21.0391 9.41153 21.1483L9.41176 21.1486L9.41179 21.1486L9.4118 21.1486C9.64176 21.3951 9.94071 21.7155 10.22 21.9596C10.6437 22.33 11.2516 22.75 12 22.75C12.7485 22.75 13.3563 22.33 13.7801 21.9596C14.0593 21.7155 14.3583 21.3951 14.5882 21.1486C14.6902 21.0392 14.7787 20.9445 14.8458 20.8773C15.1599 20.5632 15.3685 20.4191 15.5635 20.3386C15.7587 20.2581 16.0067 20.2137 16.4479 20.2137C16.5345 20.2137 16.6442 20.2168 16.7682 20.2203C17.083 20.2292 17.49 20.2408 17.8456 20.2107C18.3671 20.1667 19.0204 20.0266 19.5186 19.5326C20.0207 19.0347 20.1643 18.3787 20.2097 17.8539C20.2407 17.495 20.2292 17.0849 20.2203 16.7679L20.2203 16.7676C20.2168 16.644 20.2137 16.5346 20.2137 16.4479C20.2137 15.9971 20.2638 15.7453 20.3507 15.5464C20.4399 15.3422 20.5973 15.1258 20.9309 14.7922L20.9661 14.7571C21.4691 14.2541 21.8991 13.8242 22.1964 13.4344C22.5139 13.0183 22.75 12.5611 22.75 12C22.75 11.4388 22.5139 10.9816 22.1964 10.5655C21.8991 10.1759 21.4693 9.74607 20.9664 9.24322L20.931 9.20784C20.5973 8.87416 20.4399 8.65779 20.3507 8.45354C20.2638 8.25468 20.2137 8.00288 20.2137 7.55206C20.2137 7.46534 20.2168 7.35593 20.2203 7.23236L20.2203 7.2321C20.2292 6.91507 20.2407 6.50496 20.2097 6.14615C20.1643 5.62129 20.0207 4.96533 19.5187 4.46747C19.0205 3.97339 18.3672 3.83325 17.8456 3.78921C17.49 3.75919 17.083 3.77072 16.7682 3.77964C16.6442 3.78315 16.5345 3.78626 16.4479 3.78626C15.7553 3.78626 15.3067 3.58361 14.7922 3.06904L14.7568 3.03368C14.2539 2.53075 13.8241 2.10089 13.4344 1.80358C13.0184 1.48611 12.5612 1.25 12 1.25ZM15.7657 10.1432C16.1209 9.72033 16.0661 9.08954 15.6432 8.73432C15.2203 8.37909 14.5895 8.43394 14.2343 8.85683L10.6972 13.0676L9.66603 12.1469C9.25406 11.7791 8.6219 11.8149 8.25407 12.2269C7.88624 12.6388 7.92202 13.271 8.33399 13.6388L10.134 15.246C10.3357 15.4261 10.6018 15.5168 10.8716 15.4975C11.1413 15.4781 11.3918 15.3503 11.5657 15.1432L15.7657 10.1432Z"></path>
											</svg>
										{/if}
									</span>
									{#if !profile.is_online && profile.last_seen}
										<span class="last-seen">{formatLastSeen(profile.last_seen)}</span>
									{/if}
								</div>
								<p class="location">{profile.city}, {profile.state}</p>
								{#if profile.gender === 'male' && profile.salary_range}
									<p class="salary">{profile.salary_range}</p>
								{/if}
							</div>
						</a>
						<button 
							class="like-btn" 
							class:liked={profile.is_liked}
							onclick={() => toggleLike(profile)}
						>
							<HeartIcon liked={profile.is_liked} size={20} />
						</button>
					</div>
				{/each}
			</div>

			{#if total > profiles.length}
				<div class="pagination">
					<button 
						onclick={loadMoreProfiles}
						class="load-more"
					>
						Load More
					</button>
				</div>
			{/if}
		{/if}
	</main>

	<Footer />
</div>

<VerificationModal bind:show={showVerificationModal} onClose={() => showVerificationModal = false} />

<style>
	:global(body) {
		font-family: 'Montserrat', sans-serif;
		background: #0a0a0a;
		color: #fff;
		margin: 0;
	}

	.browse-page {
		min-height: 100vh;
	}

	.main {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;
	}

	.filters {
		display: flex;
		gap: 1rem;
		margin-bottom: 2rem;
		flex-wrap: wrap;
	}

	.filters {
		align-items: center;
	}

	.filter-group {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.filter-group label {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.6);
	}

	.age-inputs {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.age-inputs span {
		color: rgba(255, 255, 255, 0.4);
	}

	.age-input {
		width: 60px;
		height: 44px;
		padding: 0 0.5rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.85rem;
		text-align: center;
		box-sizing: border-box;
	}

	.age-input:focus {
		outline: none;
		border-color: rgba(255, 255, 255, 0.3);
	}

	.filter-select {
		height: 44px;
		padding: 0 2.5rem 0 1rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.9rem;
		min-width: 140px;
		appearance: none;
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='white' d='M6 8L1 3h10z'/%3E%3C/svg%3E");
		background-repeat: no-repeat;
		background-position: right 1rem center;
		box-sizing: border-box;
	}

	.online-filter {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		cursor: pointer;
		height: 44px;
		padding: 0 1rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		transition: all 0.2s ease;
		box-sizing: border-box;
	}

	.online-filter:hover {
		border-color: rgba(255, 255, 255, 0.2);
	}

	.online-filter input {
		accent-color: #22c55e;
	}

	.online-filter span {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.7);
	}

	.online-filter:has(input:checked) {
		border-color: #22c55e;
		background: rgba(34, 197, 94, 0.1);
	}

	.online-filter:has(input:checked) span {
		color: #22c55e;
	}

	.loading, .empty {
		text-align: center;
		padding: 4rem 2rem;
		color: rgba(255, 255, 255, 0.6);
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 3px solid rgba(255, 255, 255, 0.1);
		border-top-color: #fff;
		border-radius: 0;
		animation: spin 1s linear infinite;
		margin: 0 auto 1rem;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.hint {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.4);
	}

	.verification-required {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 50vh;
		padding: 2rem;
	}

	.verification-box {
		text-align: center;
		max-width: 400px;
		padding: 3rem 2rem;
		background: rgba(255, 56, 92, 0.08);
		border: 1px solid rgba(255, 56, 92, 0.3);
	}

	.verification-box .icon {
		font-size: 3rem;
		display: block;
		margin-bottom: 1rem;
	}

	.verification-box h2 {
		font-family: 'Playfair Display', serif;
		font-size: 1.5rem;
		margin: 0 0 1rem 0;
		color: #FF385C;
	}

	.verification-box p {
		color: rgba(255, 255, 255, 0.6);
		line-height: 1.6;
		margin: 0 0 1.5rem 0;
	}

	.resend-link {
		background: none;
		border: none;
		color: #FF385C;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.85rem;
		cursor: pointer;
		text-decoration: underline;
		padding: 0;
		transition: opacity 0.2s ease;
	}

	.resend-link:hover:not(:disabled) {
		opacity: 0.8;
	}

	.resend-link:disabled {
		cursor: default;
		text-decoration: none;
		opacity: 0.6;
	}

	.profiles-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
		gap: 1.5rem;
	}

	.profile-card {
		position: relative;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		overflow: hidden;
		transition: all 0.2s ease;
	}

	.profile-card.verified {
		/* Verified styling handled by icon */
	}

	.profile-card:hover {
		transform: translateY(-4px);
		border-color: rgba(255, 255, 255, 0.2);
	}

	.profile-card.verified:hover {
		/* Verified styling handled by icon */
	}

	.card-link {
		text-decoration: none;
		color: inherit;
	}

	.image-container {
		position: relative;
		aspect-ratio: 3/4;
		overflow: hidden;
	}

	.profile-image {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.no-image-placeholder {
		width: 100%;
		height: 100%;
		background: rgba(255, 255, 255, 0.05);
		display: flex;
		align-items: center;
		justify-content: center;
		text-align: center;
		padding: 1rem;
	}

	.no-image-placeholder span {
		color: rgba(255, 255, 255, 0.4);
		font-size: 0.85rem;
	}

	.online-badge {
		position: absolute;
		top: 0.75rem;
		right: 0.75rem;
		width: 12px;
		height: 12px;
		background: #22c55e;
		border-radius: 0;
		border: 2px solid #0a0a0a;
	}

	.verified-badge {
		position: absolute;
		top: 0.75rem;
		left: 0.75rem;
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
		color: #000;
		font-size: 0.65rem;
		font-weight: 700;
		padding: 0.25rem 0.5rem;
		letter-spacing: 0.5px;
	}

	.not-verified-badge {
		position: absolute;
		bottom: 0.75rem;
		right: 0.75rem;
		background: rgba(0, 0, 0, 0.7);
		color: #fff;
		font-size: 0.75rem;
		padding: 0.25rem 0.5rem;
		border-radius: 0;
	}

	.distance-badge {
		position: absolute;
		bottom: 0.75rem;
		left: 0.75rem;
		background: rgba(0, 0, 0, 0.7);
		padding: 0.25rem 0.5rem;
		border-radius: 0;
		font-size: 0.75rem;
		color: #fff;
	}

	.card-content {
		padding: 1rem;
	}

	.card-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 0.25rem;
	}

	.name {
		font-size: 1rem;
		font-weight: 600;
		display: flex;
		align-items: center;
		gap: 0.35rem;
	}

	.verified-icon {
		color: #ec4899;
		flex-shrink: 0;
	}

	.last-seen {
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.4);
	}

	.location {
		margin: 0;
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.6);
	}

	.salary {
		margin: 0.5rem 0 0 0;
		font-size: 0.8rem;
		color: #22c55e;
	}

	.like-btn {
		position: absolute;
		top: 0.75rem;
		left: 0.75rem;
		width: 40px;
		height: 40px;
		background: rgba(0, 0, 0, 0.5);
		border: none;
		border-radius: 0;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.2s ease;
		color: #fff;
	}

	.like-btn:hover {
		background: rgba(0, 0, 0, 0.7);
		transform: scale(1.1);
	}

	.pagination {
		text-align: center;
		margin-top: 2rem;
	}

	.load-more {
		padding: 0.875rem 2rem;
		background: transparent;
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.9rem;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.load-more:hover {
		border-color: rgba(255, 255, 255, 0.4);
	}

	@media (max-width: 768px) {
		.main {
			padding: 1rem;
		}

		.filters {
			display: flex;
			flex-wrap: wrap;
			gap: 0.75rem;
		}

		.filter-group label {
			display: none;
		}

		.age-input {
			width: 50px;
			height: 38px;
			font-size: 0.8rem;
		}

		.filter-select {
			min-width: 0;
			flex: 1;
			height: 38px;
			font-size: 0.8rem;
			padding: 0 1.5rem 0 0.5rem;
			background-position: right 0.5rem center;
		}

		.online-filter {
			height: 38px;
			padding: 0 0.75rem;
		}

		.online-filter span {
			font-size: 0.8rem;
		}

		.profiles-grid {
			grid-template-columns: repeat(2, 1fr);
			gap: 1rem;
		}

		.not-verified-badge {
			bottom: auto;
			top: 0.5rem;
			right: 0.5rem;
			font-size: 0.6rem;
			padding: 0.15rem 0.35rem;
		}

		.distance-badge {
			bottom: 0.5rem;
			left: 0.5rem;
			font-size: 0.65rem;
			padding: 0.15rem 0.35rem;
		}

		.online-badge {
			top: 0.5rem;
			right: auto;
			left: 0.5rem;
			width: 10px;
			height: 10px;
		}
	}
</style>

