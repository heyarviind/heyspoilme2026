<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { auth } from '$lib/stores/auth';
	import Footer from '$lib/components/Footer.svelte';
	import HeartIcon from '$lib/components/HeartIcon.svelte';

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
		is_liked: boolean;
		has_liked_me: boolean;
		wealth_status?: string;
		images: Array<{ id: string; url: string; is_primary: boolean }>;
	}

	let profile = $state<Profile | null>(null);
	let loading = $state(true);
	let error = $state('');
	let currentImageIndex = $state(0);
	let showMessageModal = $state(false);
	let showVerificationModal = $state(false);
	let showSubscriptionModal = $state(false);
	let messageText = $state('');
	let sendingMessage = $state(false);
	
	function getWealthLabel(status?: string): string {
		switch (status) {
			case 'low': return 'Trusted';
			case 'medium': return 'Premium';
			case 'high': return 'Elite';
			default: return '';
		}
	}

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
				profile.is_liked = false;
			} else {
				await api.likeProfile(profile.user_id);
				profile.is_liked = true;
			}
		} catch (e: any) {
			const errMsg = e.message || '';
			if (errMsg.includes('verification') || errMsg.includes('person_verification')) {
				showVerificationModal = true;
			} else {
				console.error('Failed to toggle like:', e);
			}
		}
	}

	async function startConversation() {
		if (!profile || !messageText.trim()) return;
		sendingMessage = true;
		try {
			await api.createConversation(profile.user_id, messageText.trim());
			goto('/messages');
		} catch (e: any) {
			const errMsg = e.message || '';
			if (errMsg.includes('verification') || errMsg.includes('person_verification')) {
				showMessageModal = false;
				showVerificationModal = true;
			} else if (errMsg.includes('subscription') || errMsg.includes('wealth_status')) {
				showMessageModal = false;
				showSubscriptionModal = true;
			} else {
				alert(errMsg || 'Failed to send message');
			}
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

			<div class="info" class:verified-profile={profile.is_verified}>
				<div class="info-header">
					<div class="name-status">
						<div class="name-row">
							<h1>
								{profile.display_name}, {profile.age}
								{#if profile.is_verified}
									<svg class="verified-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
										<path fill-rule="evenodd" clip-rule="evenodd" d="M12 1.25C11.4388 1.25 10.9816 1.48611 10.5656 1.80358C10.1759 2.10089 9.74606 2.53075 9.24319 3.03367L9.20782 3.06904C8.69316 3.5837 8.24449 3.78626 7.55208 3.78626C7.4653 3.78626 7.35579 3.78318 7.23212 3.7797C6.91509 3.77078 6.50497 3.75924 6.14615 3.79027C5.62128 3.83566 4.96532 3.97929 4.46745 4.48134C3.9734 4.97955 3.83327 5.63282 3.78923 6.15439C3.75922 6.50995 3.77075 6.91701 3.77966 7.23178L3.77966 7.23181C3.78317 7.35581 3.78628 7.46549 3.78628 7.55206C3.78628 8.24448 3.58371 8.69315 3.06902 9.20784L3.03367 9.24319C2.53075 9.74606 2.10089 10.1759 1.80358 10.5655C1.48612 10.9816 1.25001 11.4388 1.25 12C1.25001 12.5611 1.48613 13.0183 1.80358 13.4344C2.10095 13.8242 2.53091 14.2541 3.03395 14.7571L3.06906 14.7922C3.40272 15.1258 3.56011 15.3422 3.64932 15.5464C3.73619 15.7453 3.78628 15.9971 3.78628 16.4479C3.78628 16.5347 3.7832 16.6442 3.77972 16.7679C3.7708 17.0849 3.75926 17.495 3.79029 17.8539C3.83569 18.3787 3.97933 19.0347 4.48139 19.5326C4.97961 20.0266 5.63287 20.1667 6.15443 20.2107C6.50997 20.2408 6.91703 20.2292 7.23179 20.2203C7.35581 20.2168 7.4655 20.2137 7.55206 20.2137C7.99328 20.2137 8.24126 20.2581 8.43645 20.3386C8.63147 20.4191 8.84006 20.5632 9.15424 20.8774C9.22129 20.9444 9.30963 21.0391 9.41153 21.1483L9.41176 21.1486L9.41179 21.1486L9.4118 21.1486C9.64176 21.3951 9.94071 21.7155 10.22 21.9596C10.6437 22.33 11.2516 22.75 12 22.75C12.7485 22.75 13.3563 22.33 13.7801 21.9596C14.0593 21.7155 14.3583 21.3951 14.5882 21.1486C14.6902 21.0392 14.7787 20.9445 14.8458 20.8773C15.1599 20.5632 15.3685 20.4191 15.5635 20.3386C15.7587 20.2581 16.0067 20.2137 16.4479 20.2137C16.5345 20.2137 16.6442 20.2168 16.7682 20.2203C17.083 20.2292 17.49 20.2408 17.8456 20.2107C18.3671 20.1667 19.0204 20.0266 19.5186 19.5326C20.0207 19.0347 20.1643 18.3787 20.2097 17.8539C20.2407 17.495 20.2292 17.0849 20.2203 16.7679L20.2203 16.7676C20.2168 16.644 20.2137 16.5346 20.2137 16.4479C20.2137 15.9971 20.2638 15.7453 20.3507 15.5464C20.4399 15.3422 20.5973 15.1258 20.9309 14.7922L20.9661 14.7571C21.4691 14.2541 21.8991 13.8242 22.1964 13.4344C22.5139 13.0183 22.75 12.5611 22.75 12C22.75 11.4388 22.5139 10.9816 22.1964 10.5655C21.8991 10.1759 21.4693 9.74607 20.9664 9.24322L20.931 9.20784C20.5973 8.87416 20.4399 8.65779 20.3507 8.45354C20.2638 8.25468 20.2137 8.00288 20.2137 7.55206C20.2137 7.46534 20.2168 7.35593 20.2203 7.23236L20.2203 7.2321C20.2292 6.91507 20.2407 6.50496 20.2097 6.14615C20.1643 5.62129 20.0207 4.96533 19.5187 4.46747C19.0205 3.97339 18.3672 3.83325 17.8456 3.78921C17.49 3.75919 17.083 3.77072 16.7682 3.77964C16.6442 3.78315 16.5345 3.78626 16.4479 3.78626C15.7553 3.78626 15.3067 3.58361 14.7922 3.06904L14.7568 3.03368C14.2539 2.53075 13.8241 2.10089 13.4344 1.80358C13.0184 1.48611 12.5612 1.25 12 1.25ZM15.7657 10.1432C16.1209 9.72033 16.0661 9.08954 15.6432 8.73432C15.2203 8.37909 14.5895 8.43394 14.2343 8.85683L10.6972 13.0676L9.66603 12.1469C9.25406 11.7791 8.6219 11.8149 8.25407 12.2269C7.88624 12.6388 7.92202 13.271 8.33399 13.6388L10.134 15.246C10.3357 15.4261 10.6018 15.5168 10.8716 15.4975C11.1413 15.4781 11.3918 15.3503 11.5657 15.1432L15.7657 10.1432Z"></path>
									</svg>
								{/if}
							</h1>
							{#if !profile.is_verified}
								<span class="not-verified-tag">NOT VERIFIED</span>
							{/if}
						</div>
						{#if profile.is_online}
							<span class="online-status">● Online</span>
						{:else if profile.last_seen}
							<span class="offline-status">Last seen {formatLastSeen(profile.last_seen)}</span>
						{/if}
					</div>
					<button class="like-btn" class:liked={profile.is_liked} onclick={toggleLike}>
						<HeartIcon liked={profile.is_liked} size={28} />
					</button>
				</div>

				<p class="location">{profile.city}, {profile.state}</p>

				{#if profile.has_liked_me}
					<div class="liked-you-badge">
						<span>💕 Liked your profile</span>
					</div>
				{/if}

				{#if profile.gender === 'male'}
					<div class="male-badges">
						{#if profile.wealth_status && profile.wealth_status !== 'none'}
							<span class="wealth-badge {profile.wealth_status}">{getWealthLabel(profile.wealth_status)} Member</span>
						{/if}
						{#if profile.salary_range}
							<span class="salary-badge">💰 {profile.salary_range}</span>
						{/if}
					</div>
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

	<Footer />
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

{#if showVerificationModal}
	<div class="modal-overlay" onclick={() => showVerificationModal = false}>
		<div class="modal verification-modal" onclick={(e) => e.stopPropagation()}>
			<div class="modal-icon">🔐</div>
			<h2>Identity Verification Required</h2>
			<p class="modal-description">
				Complete a quick identity verification to like profiles and send messages. 
				This keeps our community safe and ensures only serious members can connect.
			</p>
			<div class="modal-actions">
				<button class="btn-secondary" onclick={() => showVerificationModal = false}>
					Later
				</button>
				<a href="/profile/verify" class="btn-primary">
					Verify Now
				</a>
			</div>
		</div>
	</div>
{/if}

{#if showSubscriptionModal}
	<div class="modal-overlay" onclick={() => showSubscriptionModal = false}>
		<div class="modal subscription-modal" onclick={(e) => e.stopPropagation()}>
			<div class="modal-icon">✨</div>
			<h2>Become a Trusted Member</h2>
			<p class="modal-description">
				Upgrade to read and reply to messages from women who are interested in you.
			</p>
			<ul class="subscription-benefits">
				<li>✓ Read & reply to message requests</li>
				<li>✓ Initiate conversations</li>
				<li>✓ Priority placement in discovery</li>
				<li>✓ Trusted Member badge on your profile</li>
			</ul>
			<p class="subscription-note">
				We keep the community high-quality by allowing only verified, serious members to initiate conversations.
			</p>
			<div class="modal-actions">
				<button class="btn-secondary" onclick={() => showSubscriptionModal = false}>
					Maybe Later
				</button>
				<a href="/profile" class="btn-gold">
					Upgrade Now
				</a>
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
		/* Verified styling handled by icon */
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
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
	}

	.verified-icon {
		color: #ec4899;
		flex-shrink: 0;
	}

	.not-verified-tag {
		background: rgba(255, 255, 255, 0.1);
		color: rgba(255, 255, 255, 0.6);
		font-size: 0.7rem;
		font-weight: 600;
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
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.2s ease;
		color: #fff;
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

	/* Has Liked You Badge */
	.liked-you-badge {
		display: inline-block;
		padding: 0.5rem 1rem;
		background: linear-gradient(135deg, rgba(236, 72, 153, 0.2) 0%, rgba(236, 72, 153, 0.1) 100%);
		border: 1px solid rgba(236, 72, 153, 0.4);
		color: #ec4899;
		font-size: 0.85rem;
		margin: 0.75rem 0;
	}

	/* Male Badges */
	.male-badges {
		display: flex;
		gap: 0.75rem;
		flex-wrap: wrap;
		margin: 0.75rem 0;
	}

	.wealth-badge {
		padding: 0.4rem 0.8rem;
		font-size: 0.75rem;
		font-weight: 600;
		letter-spacing: 0.5px;
	}

	.wealth-badge.low {
		background: linear-gradient(135deg, rgba(59, 130, 246, 0.2) 0%, rgba(59, 130, 246, 0.1) 100%);
		border: 1px solid rgba(59, 130, 246, 0.4);
		color: #3b82f6;
	}

	.wealth-badge.medium {
		background: linear-gradient(135deg, rgba(168, 85, 247, 0.2) 0%, rgba(168, 85, 247, 0.1) 100%);
		border: 1px solid rgba(168, 85, 247, 0.4);
		color: #a855f7;
	}

	.wealth-badge.high {
		background: linear-gradient(135deg, rgba(212, 175, 55, 0.2) 0%, rgba(212, 175, 55, 0.1) 100%);
		border: 1px solid rgba(212, 175, 55, 0.4);
		color: #d4af37;
	}

	.salary-badge {
		padding: 0.4rem 0.8rem;
		font-size: 0.75rem;
		background: rgba(34, 197, 94, 0.1);
		border: 1px solid rgba(34, 197, 94, 0.3);
		color: #22c55e;
	}

	/* Verification & Subscription Modals */
	.verification-modal,
	.subscription-modal {
		text-align: center;
	}

	.modal-icon {
		font-size: 3rem;
		margin-bottom: 1rem;
	}

	.modal-description {
		color: rgba(255, 255, 255, 0.7);
		line-height: 1.6;
		margin: 1rem 0 1.5rem;
	}

	.subscription-benefits {
		list-style: none;
		padding: 0;
		margin: 1.5rem 0;
		text-align: left;
	}

	.subscription-benefits li {
		padding: 0.5rem 0;
		color: rgba(255, 255, 255, 0.8);
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
	}

	.subscription-note {
		font-size: 0.8rem;
		color: rgba(255, 255, 255, 0.4);
		margin-bottom: 1.5rem;
	}

	.btn-gold {
		flex: 1;
		padding: 0.875rem;
		background: linear-gradient(135deg, #d4af37 0%, #c5a028 100%);
		color: #000;
		text-decoration: none;
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		text-align: center;
		border: none;
		cursor: pointer;
	}

	.btn-gold:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 20px rgba(212, 175, 55, 0.4);
	}
</style>

