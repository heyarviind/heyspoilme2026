<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	interface Conversation {
		id: string;
		initiated_by: string;
		other_user?: {
			user_id: string;
			age: number;
			city: string;
			is_online: boolean;
			images: Array<{ url: string; is_primary: boolean }>;
		};
		last_message?: {
			content: string;
			created_at: string;
			sender_id: string;
		};
		unread_count: number;
	}

	let conversations = $state<Conversation[]>([]);
	let loading = $state(true);

	async function loadConversations() {
		loading = true;
		try {
			const data = await api.getConversations() as { conversations: Conversation[] };
			conversations = data.conversations || [];
		} catch (e) {
			console.error('Failed to load conversations:', e);
		} finally {
			loading = false;
		}
	}

	function getProfileImage(conversation: Conversation): string {
		const primary = conversation.other_user?.images?.find(img => img.is_primary);
		if (primary) return primary.url;
		if (conversation.other_user?.images?.length) return conversation.other_user.images[0].url;
		return 'https://via.placeholder.com/80?text=?';
	}

	function formatTime(dateStr?: string): string {
		if (!dateStr) return '';
		const date = new Date(dateStr);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);
		
		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins}m`;
		const diffHours = Math.floor(diffMins / 60);
		if (diffHours < 24) return `${diffHours}h`;
		const diffDays = Math.floor(diffHours / 24);
		if (diffDays < 7) return `${diffDays}d`;
		return date.toLocaleDateString();
	}

	function truncateMessage(content?: string): string {
		if (!content) return '';
		if (content.length <= 50) return content;
		return content.substring(0, 50) + '...';
	}

	onMount(() => {
		loadConversations();
	});
</script>

<svelte:head>
	<title>Messages | HeySpoilMe</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;500;600;700&family=Montserrat:wght@300;400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="messages-page">
	<header class="header">
		<a href="/browse" class="logo-link">
			<img src="/img/logo.svg" alt="HeySpoilMe" class="logo" />
		</a>
		<nav class="nav">
			<a href="/browse" class="nav-link">Browse</a>
			<a href="/messages" class="nav-link active">Messages</a>
			<a href="/likes" class="nav-link">Likes</a>
			<a href="/profile" class="nav-link">Profile</a>
		</nav>
	</header>

	<main class="main">
		<h1>Messages</h1>

		{#if loading}
			<div class="loading">
				<div class="spinner"></div>
			</div>
		{:else if conversations.length === 0}
			<div class="empty">
				<p>No messages yet</p>
				<p class="hint">Start a conversation by visiting someone's profile</p>
				<a href="/browse" class="browse-link">Browse Profiles</a>
			</div>
		{:else}
			<div class="conversation-list">
				{#each conversations as conversation}
					<a href="/messages/{conversation.id}" class="conversation-item">
						<div class="avatar-container">
							<img 
								src={getProfileImage(conversation)} 
								alt="Profile" 
								class="avatar" 
							/>
							{#if conversation.other_user?.is_online}
								<span class="online-dot"></span>
							{/if}
						</div>
						<div class="conversation-info">
							<div class="top-row">
								<span class="name">{conversation.other_user?.age || 'Unknown'}</span>
								<span class="time">{formatTime(conversation.last_message?.created_at)}</span>
							</div>
							<div class="bottom-row">
								<p class="preview">{truncateMessage(conversation.last_message?.content)}</p>
								{#if conversation.unread_count > 0}
									<span class="unread-badge">{conversation.unread_count}</span>
								{/if}
							</div>
						</div>
					</a>
				{/each}
			</div>
		{/if}
	</main>
</div>

<style>
	:global(body) {
		font-family: 'Montserrat', sans-serif;
		background: #0a0a0a;
		color: #fff;
		margin: 0;
	}

	.messages-page {
		min-height: 100vh;
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1rem 2rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}

	.logo-link {
		text-decoration: none;
	}

	.logo {
		height: 2.5rem;
	}

	.nav {
		display: flex;
		gap: 2rem;
	}

	.nav-link {
		color: rgba(255, 255, 255, 0.6);
		text-decoration: none;
		font-size: 0.9rem;
		font-weight: 500;
	}

	.nav-link:hover, .nav-link.active {
		color: #fff;
	}

	.main {
		max-width: 600px;
		margin: 0 auto;
		padding: 2rem;
	}

	h1 {
		font-family: 'Playfair Display', serif;
		font-size: 2rem;
		margin: 0 0 2rem 0;
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

	.conversation-list {
		display: flex;
		flex-direction: column;
	}

	.conversation-item {
		display: flex;
		gap: 1rem;
		padding: 1rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
		text-decoration: none;
		color: inherit;
		transition: background 0.2s ease;
	}

	.conversation-item:hover {
		background: rgba(255, 255, 255, 0.03);
	}

	.avatar-container {
		position: relative;
		flex-shrink: 0;
	}

	.avatar {
		width: 56px;
		height: 56px;
		border-radius: 0;
		object-fit: cover;
	}

	.online-dot {
		position: absolute;
		bottom: 2px;
		right: 2px;
		width: 14px;
		height: 14px;
		background: #22c55e;
		border-radius: 0;
		border: 2px solid #0a0a0a;
	}

	.conversation-info {
		flex: 1;
		min-width: 0;
		display: flex;
		flex-direction: column;
		justify-content: center;
	}

	.top-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.25rem;
	}

	.name {
		font-weight: 600;
	}

	.time {
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.4);
	}

	.bottom-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.preview {
		margin: 0;
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.6);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.unread-badge {
		background: #fff;
		color: #000;
		font-size: 0.7rem;
		font-weight: 600;
		padding: 0.2rem 0.5rem;
		border-radius: 0;
		flex-shrink: 0;
	}

	@media (max-width: 768px) {
		.header {
			flex-direction: column;
			gap: 1rem;
		}

		.main {
			padding: 1rem;
		}
	}
</style>

