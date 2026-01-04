<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { auth } from '$lib/stores/auth';

	interface Profile {
		id: string;
		user_id: string;
		gender: 'male' | 'female';
		age: number;
		bio: string;
		salary_range?: string;
		city: string;
		state: string;
		is_online: boolean;
		is_verified: boolean;
		last_seen?: string;
		is_liked: boolean;
		images: Array<{ id: string; url: string; is_primary: boolean }>;
	}

	let profile = $state<Profile | null>(null);
	let loading = $state(true);
	let error = $state('');
	let currentImageIndex = $state(0);
	let showMessageModal = $state(false);
	let messageText = $state('');
	let sendingMessage = $state(false);

	let profileId = $derived($page.params.id);

	async function loadProfile() {
		loading = true;
		error = '';
		try {
			const data = await api.getProfile(profileId) as Profile;
			profile = data;
		} catch (e: any) {
			error = e.message || 'Failed to load profile';
		} finally {
			loading = false;
		}
	}

	async function toggleLike() {
		if (!profile) return;
		try {
			if (profile.is_liked) {
				await api.unlikeProfile(profile.user_id);
			} else {
				await api.likeProfile(profile.user_id);
			}
			profile.is_liked = !profile.is_liked;
		} catch (e) {
			console.error('Failed to toggle like:', e);
		}
	}

	async function startConversation() {
		if (!profile || !messageText.trim()) return;
		sendingMessage = true;
		try {
			await api.createConversation(profile.user_id, messageText.trim());
			goto('/messages');
		} catch (e: any) {
			alert(e.message || 'Failed to send message');
		} finally {
			sendingMessage = false;
		}
	}

	function nextImage() {
		if (profile?.images && currentImageIndex < profile.images.length - 1) {
			currentImageIndex++;
		}
	}

	function prevImage() {
		if (currentImageIndex > 0) {
			currentImageIndex--;
		}
	}

	function formatLastSeen(lastSeen?: string): string {
		if (!lastSeen) return '';
		const date = new Date(lastSeen);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);
		
		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins} minutes ago`;
		const diffHours = Math.floor(diffMins / 60);
		if (diffHours < 24) return `${diffHours} hours ago`;
		const diffDays = Math.floor(diffHours / 24);
		return `${diffDays} days ago`;
	}

	onMount(() => {
		loadProfile();
	});

	let authState = $state<any>(null);
	auth.subscribe(s => authState = s);
	let canMessage = $derived(authState?.profile?.gender === 'female' && profile?.gender === 'male');
</script>

<svelte:head>
	<title>Profile | HeySpoilMe</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;500;600;700&family=Montserrat:wght@300;400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="profile-page">
	<header class="header">
		<button class="back-btn" onclick={() => history.back()}>
			← Back
		</button>
	</header>

	{#if loading}
		<div class="loading">
			<div class="spinner"></div>
		</div>
	{:else if error}
		<div class="error-container">
			<p>{error}</p>
			<button onclick={loadProfile}>Try Again</button>
		</div>
	{:else if profile}
		<div class="profile-content">
			<div class="gallery">
				{#if profile.images && profile.images.length > 0}
					<div class="image-viewer">
						<img 
							src={profile.images[currentImageIndex]?.url} 
							alt="Profile" 
							class="main-image" 
						/>
						{#if profile.images.length > 1}
							<button class="nav-btn prev" onclick={prevImage} disabled={currentImageIndex === 0}>
								‹
							</button>
							<button class="nav-btn next" onclick={nextImage} disabled={currentImageIndex === profile.images.length - 1}>
								›
							</button>
							<div class="dots">
								{#each profile.images as _, i}
									<span class="dot" class:active={i === currentImageIndex}></span>
								{/each}
							</div>
						{/if}
					</div>
				{:else}
					<div class="no-image">
						<span>No photos</span>
					</div>
				{/if}
			</div>

			<div class="info" class:verified-profile={profile.is_verified && profile.gender === 'male'}>
				<div class="info-header">
					<div class="name-status">
						<div class="name-row">
							<h1>{profile.age} years old</h1>
							{#if profile.is_verified && profile.gender === 'male'}
								<span class="verified-tag">✓ VERIFIED</span>
							{/if}
						</div>
						{#if profile.is_online}
							<span class="online-status">● Online</span>
						{:else if profile.last_seen}
							<span class="offline-status">Last seen {formatLastSeen(profile.last_seen)}</span>
						{/if}
					</div>
					<button class="like-btn" class:liked={profile.is_liked} onclick={toggleLike}>
						{profile.is_liked ? '❤️' : '🤍'}
					</button>
				</div>

				<p class="location">{profile.city}, {profile.state}</p>

				{#if profile.gender === 'male' && profile.salary_range}
					<p class="salary">💰 {profile.salary_range}</p>
				{/if}

				<div class="bio-section">
					<h3>About</h3>
					<p class="bio">{profile.bio}</p>
				</div>

				{#if canMessage}
					<button class="message-btn" onclick={() => showMessageModal = true}>
						Send Message
					</button>
				{:else if authState?.profile?.gender === 'male'}
					<p class="message-hint">
						Wait for her to message you first 💫
					</p>
				{/if}
			</div>
		</div>
	{/if}
</div>

{#if showMessageModal}
	<div class="modal-overlay" onclick={() => showMessageModal = false}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>Start a Conversation</h2>
			<p class="modal-hint">Make a great first impression!</p>
			<textarea 
				bind:value={messageText}
				placeholder="Write your message..."
				rows="4"
				maxlength="500"
			></textarea>
			<div class="modal-actions">
				<button class="btn-secondary" onclick={() => showMessageModal = false}>
					Cancel
				</button>
				<button 
					class="btn-primary" 
					onclick={startConversation}
					disabled={!messageText.trim() || sendingMessage}
				>
					{sendingMessage ? 'Sending...' : 'Send Message'}
				</button>
			</div>
		</div>
	</div>
{/if}

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

	.header {
		padding: 1rem 2rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}

	.back-btn {
		background: none;
		border: none;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.9rem;
		cursor: pointer;
		padding: 0.5rem 1rem;
		margin: -0.5rem -1rem;
	}

	.loading, .error-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 60vh;
		color: rgba(255, 255, 255, 0.6);
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
		display: grid;
		grid-template-columns: 1fr 1fr;
		max-width: 1000px;
		margin: 0 auto;
		gap: 2rem;
		padding: 2rem;
	}

	.gallery {
		position: sticky;
		top: 100px;
	}

	.image-viewer {
		position: relative;
		aspect-ratio: 3/4;
		border-radius: 0;
		overflow: hidden;
		background: rgba(255, 255, 255, 0.05);
	}

	.main-image {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.nav-btn {
		position: absolute;
		top: 50%;
		transform: translateY(-50%);
		width: 48px;
		height: 48px;
		background: rgba(0, 0, 0, 0.5);
		border: none;
		border-radius: 0;
		color: #fff;
		font-size: 1.5rem;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.nav-btn:disabled {
		opacity: 0.3;
		cursor: not-allowed;
	}

	.nav-btn.prev {
		left: 1rem;
	}

	.nav-btn.next {
		right: 1rem;
	}

	.dots {
		position: absolute;
		bottom: 1rem;
		left: 50%;
		transform: translateX(-50%);
		display: flex;
		gap: 0.5rem;
	}

	.dot {
		width: 8px;
		height: 8px;
		border-radius: 0;
		background: rgba(255, 255, 255, 0.4);
	}

	.dot.active {
		background: #fff;
	}

	.no-image {
		aspect-ratio: 3/4;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(255, 255, 255, 0.05);
		border-radius: 0;
		color: rgba(255, 255, 255, 0.4);
	}

	.info {
		padding: 1rem 0;
	}

	.info.verified-profile {
		background: linear-gradient(135deg, rgba(251, 191, 36, 0.1) 0%, transparent 50%);
		padding: 1.5rem;
		border-left: 3px solid #fbbf24;
	}

	.info-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		margin-bottom: 0.5rem;
	}

	.name-row {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		flex-wrap: wrap;
	}

	.name-status h1 {
		font-family: 'Playfair Display', serif;
		font-size: 2rem;
		font-weight: 500;
		margin: 0;
	}

	.verified-tag {
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
		color: #000;
		font-size: 0.7rem;
		font-weight: 700;
		padding: 0.3rem 0.6rem;
		letter-spacing: 0.5px;
	}

	.online-status {
		color: #22c55e;
		font-size: 0.85rem;
	}

	.offline-status {
		color: rgba(255, 255, 255, 0.5);
		font-size: 0.85rem;
	}

	.like-btn {
		width: 56px;
		height: 56px;
		background: rgba(255, 255, 255, 0.1);
		border: none;
		border-radius: 0;
		font-size: 1.5rem;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.like-btn:hover {
		transform: scale(1.1);
	}

	.location {
		color: rgba(255, 255, 255, 0.6);
		margin: 0.5rem 0;
	}

	.salary {
		color: #22c55e;
		margin: 0.5rem 0;
		font-weight: 500;
	}

	.bio-section {
		margin-top: 2rem;
	}

	.bio-section h3 {
		font-size: 0.85rem;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		color: rgba(255, 255, 255, 0.5);
		margin: 0 0 0.5rem 0;
	}

	.bio {
		line-height: 1.7;
		color: rgba(255, 255, 255, 0.8);
	}

	.message-btn {
		width: 100%;
		padding: 1rem;
		margin-top: 2rem;
		background: #fff;
		color: #000;
		border: none;
		border-radius: 0;
		font-family: 'Montserrat', sans-serif;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.message-btn:hover {
		transform: translateY(-2px);
	}

	.message-hint {
		text-align: center;
		color: rgba(255, 255, 255, 0.5);
		margin-top: 2rem;
		font-style: italic;
	}

	/* Modal */
	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.8);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 1rem;
	}

	.modal {
		background: #1a1a1a;
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		padding: 2rem;
		max-width: 400px;
		width: 100%;
	}

	.modal h2 {
		font-family: 'Playfair Display', serif;
		margin: 0 0 0.5rem 0;
	}

	.modal-hint {
		color: rgba(255, 255, 255, 0.5);
		margin: 0 0 1.5rem 0;
		font-size: 0.9rem;
	}

	.modal textarea {
		width: 100%;
		padding: 1rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		resize: vertical;
		box-sizing: border-box;
	}

	.modal-actions {
		display: flex;
		gap: 1rem;
		margin-top: 1.5rem;
	}

	.btn-primary, .btn-secondary {
		flex: 1;
		padding: 0.875rem;
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

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-secondary {
		background: transparent;
		color: #fff;
		border: 1px solid rgba(255, 255, 255, 0.2);
	}

	@media (max-width: 768px) {
		.profile-content {
			grid-template-columns: 1fr;
			padding: 1rem;
		}

		.gallery {
			position: static;
		}
	}
</style>

