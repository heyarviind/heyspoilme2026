<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { notifications } from '$lib/stores/notifications';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';

	interface Profile {
		id: string;
		user_id: string;
		display_name: string;
		age: number;
		city: string;
		state: string;
		is_online: boolean;
		images: Array<{ url: string; is_primary: boolean }>;
	}

	interface LikeWithProfile {
		like: {
			id: string;
			liker_id: string;
			liked_id: string;
			created_at: string;
		};
		profile: Profile | null;
	}

	let activeTab = $state<'received' | 'given'>('received');
	let receivedLikes = $state<LikeWithProfile[]>([]);
	let givenLikes = $state<LikeWithProfile[]>([]);
	let loading = $state(true);

	async function loadLikes() {
		loading = true;
		try {
			const [received, given] = await Promise.all([
				api.getReceivedLikes() as Promise<{ likes: LikeWithProfile[]; total: number }>,
				api.getGivenLikes() as Promise<{ likes: LikeWithProfile[]; total: number }>,
			]);
			receivedLikes = received.likes || [];
			givenLikes = given.likes || [];
		} catch (e) {
			console.error('Failed to load likes:', e);
		} finally {
			loading = false;
		}
	}

	function getProfileImage(item: LikeWithProfile): string {
		if (!item.profile?.images?.length) return 'https://via.placeholder.com/200?text=?';
		const primary = item.profile.images.find(img => img.is_primary);
		return primary?.url || item.profile.images[0].url;
	}

	function formatDate(dateStr: string): string {
		if (!dateStr) return '';
		const date = new Date(dateStr);
		if (isNaN(date.getTime())) return '';
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffHours = Math.floor(diffMs / 3600000);
		
		if (diffHours < 24) return 'Today';
		if (diffHours < 48) return 'Yesterday';
		return date.toLocaleDateString();
	}

	onMount(() => {
		loadLikes();
		// Mark all like notifications as read
		notifications.markAllAsRead();
	});

	let currentLikes = $derived(activeTab === 'received' ? receivedLikes : givenLikes);
</script>

<svelte:head>
	<title>Likes | HeySpoilMe</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;500;600;700&family=Montserrat:wght@300;400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="likes-page">
	<Header />

	<main class="main">
		<h1>Likes</h1>

		<div class="tabs">
			<button 
				class="tab" 
				class:active={activeTab === 'received'}
				onclick={() => activeTab = 'received'}
			>
				Received ({receivedLikes.length})
			</button>
			<button 
				class="tab" 
				class:active={activeTab === 'given'}
				onclick={() => activeTab = 'given'}
			>
				Given ({givenLikes.length})
			</button>
		</div>

		{#if loading}
			<div class="loading">
				<div class="spinner"></div>
			</div>
		{:else if currentLikes.length === 0}
			<div class="empty">
				{#if activeTab === 'received'}
					<p>No likes received yet</p>
					<p class="hint">Complete your profile to attract more likes!</p>
				{:else}
					<p>You haven't liked anyone yet</p>
					<p class="hint">Browse profiles and show your interest</p>
				{/if}
				<a href="/browse" class="browse-link">Browse Profiles</a>
			</div>
		{:else}
			<div class="likes-grid">
				{#each currentLikes as item}
					{@const userId = item.profile?.user_id || (activeTab === 'received' ? item.like.liker_id : item.like.liked_id)}
					<a href="/profile/{userId}" class="like-card">
						<div class="image-container">
							<img 
								src={getProfileImage(item)} 
								alt="Profile" 
								class="profile-image" 
							/>
							{#if item.profile?.is_online}
								<span class="online-badge"></span>
							{/if}
						</div>
						<div class="card-content">
							<span class="name">{item.profile?.display_name || 'Unknown'}, {item.profile?.age || '?'}</span>
							<span class="location">{item.profile?.city || 'Unknown'}</span>
							<span class="date">{formatDate(item.like.created_at)}</span>
						</div>
					</a>
				{/each}
			</div>
		{/if}
	</main>

	<Footer />
</div>

<style>
	:global(body) {
		font-family: 'Montserrat', sans-serif;
		background: #0a0a0a;
		color: #fff;
		margin: 0;
	}

	.likes-page {
		min-height: 100vh;
	}

	.main {
		max-width: 800px;
		margin: 0 auto;
		padding: 2rem;
	}

	h1 {
		font-family: 'Playfair Display', serif;
		font-size: 2rem;
		margin: 0 0 1.5rem 0;
	}

	.tabs {
		display: flex;
		gap: 0;
		margin-bottom: 2rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}

	.tab {
		flex: 1;
		padding: 1rem;
		background: none;
		border: none;
		color: rgba(255, 255, 255, 0.5);
		font-family: 'Montserrat', sans-serif;
		font-size: 0.9rem;
		font-weight: 500;
		cursor: pointer;
		position: relative;
		transition: color 0.2s ease;
	}

	.tab:hover {
		color: rgba(255, 255, 255, 0.8);
	}

	.tab.active {
		color: #fff;
	}

	.tab.active::after {
		content: '';
		position: absolute;
		bottom: -1px;
		left: 0;
		right: 0;
		height: 2px;
		background: #fff;
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
		margin: 0 auto;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.hint {
		font-size: 0.9rem;
		color: rgba(255, 255, 255, 0.4);
	}

	.browse-link {
		display: inline-block;
		margin-top: 1rem;
		padding: 0.75rem 1.5rem;
		background: #fff;
		color: #000;
		text-decoration: none;
		border-radius: 0;
		font-weight: 500;
	}

	.likes-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
		gap: 1rem;
	}

	.like-card {
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		overflow: hidden;
		text-decoration: none;
		color: inherit;
		transition: all 0.2s ease;
	}

	.like-card:hover {
		transform: translateY(-4px);
		border-color: rgba(255, 255, 255, 0.2);
	}

	.image-container {
		position: relative;
		aspect-ratio: 1;
	}

	.profile-image {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.online-badge {
		position: absolute;
		top: 0.5rem;
		right: 0.5rem;
		width: 10px;
		height: 10px;
		background: #22c55e;
		border-radius: 0;
		border: 2px solid #0a0a0a;
	}

	.card-content {
		padding: 0.75rem;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.name {
		font-weight: 600;
		font-size: 0.9rem;
	}

	.location {
		font-size: 0.8rem;
		color: rgba(255, 255, 255, 0.6);
	}

	.date {
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.4);
	}

	@media (max-width: 768px) {
		.likes-grid {
			grid-template-columns: repeat(2, 1fr);
		}
	}
</style>

