<script lang="ts">
	import { page } from '$app/stores';
	import { notifications } from '$lib/stores/notifications';
	import { onMount } from 'svelte';

	let showNotifications = $state(false);
	let notifState = $state<{ notifications: any[]; unreadCount: number; loading: boolean }>({ notifications: [], unreadCount: 0, loading: false });

	notifications.subscribe(s => notifState = s);

	function toggleNotifications() {
		showNotifications = !showNotifications;
		if (showNotifications && notifState.notifications.length === 0) {
			notifications.fetchNotifications();
		}
	}

	function formatNotification(notif: any): string {
		switch (notif.type) {
			case 'new_like':
				return 'Someone liked your profile ❤️';
			case 'new_message':
				return 'You have a new message 💬';
			case 'profile_view':
				return 'Someone viewed your profile 👀';
			default:
				return 'New notification';
		}
	}

	function formatTime(dateStr: string): string {
		const date = new Date(dateStr);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);
		
		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins}m ago`;
		const diffHours = Math.floor(diffMins / 60);
		if (diffHours < 24) return `${diffHours}h ago`;
		return date.toLocaleDateString();
	}

	function handleMarkAllRead() {
		notifications.markAllAsRead();
	}

	let currentPath = $derived($page.url.pathname);

	onMount(() => {
		notifications.fetchUnreadCount();
		// Refresh count periodically
		const interval = setInterval(() => {
			notifications.fetchUnreadCount();
		}, 30000);
		return () => clearInterval(interval);
	});
</script>

<header class="header">
	<a href="/browse" class="logo-link">
		<img src="/img/logo.svg" alt="HeySpoilMe" class="logo" />
	</a>
	
	<nav class="nav">
		<a href="/browse" class="nav-link" class:active={currentPath === '/browse'}>
			Browse
		</a>
		<a href="/messages" class="nav-link" class:active={currentPath.startsWith('/messages')}>
			Messages
		</a>
		<a href="/likes" class="nav-link" class:active={currentPath === '/likes'}>
			Likes
		</a>
		
		<button class="notif-btn" onclick={toggleNotifications}>
			🔔
			{#if notifState?.unreadCount > 0}
				<span class="notif-badge">{notifState.unreadCount}</span>
			{/if}
		</button>

		<a href="/profile" class="nav-link" class:active={currentPath === '/profile'}>
			Profile
		</a>
	</nav>
</header>

{#if showNotifications}
	<div class="notif-overlay" onclick={() => showNotifications = false}></div>
	<div class="notif-dropdown">
		<div class="notif-header">
			<span>Notifications</span>
			{#if notifState?.unreadCount > 0}
				<button class="mark-read-btn" onclick={handleMarkAllRead}>
					Mark all read
				</button>
			{/if}
		</div>
		<div class="notif-list">
			{#if notifState?.loading}
				<div class="notif-loading">Loading...</div>
			{:else if !notifState?.notifications?.length}
				<div class="notif-empty">No notifications yet</div>
			{:else}
				{#each notifState.notifications as notif}
					<div 
						class="notif-item" 
						class:unread={!notif.is_read}
						onclick={() => notifications.markAsRead(notif.id)}
					>
						<p class="notif-text">{formatNotification(notif)}</p>
						<span class="notif-time">{formatTime(notif.created_at)}</span>
					</div>
				{/each}
			{/if}
		</div>
	</div>
{/if}

<style>
	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1rem 2rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
		position: sticky;
		top: 0;
		background: #0a0a0a;
		z-index: 100;
	}

	.logo-link {
		text-decoration: none;
	}

	.logo {
		height: 2.5rem;
	}

	.nav {
		display: flex;
		align-items: center;
		gap: 1.5rem;
	}

	.nav-link {
		color: rgba(255, 255, 255, 0.6);
		text-decoration: none;
		font-size: 0.9rem;
		font-weight: 500;
		transition: color 0.2s ease;
	}

	.nav-link:hover, .nav-link.active {
		color: #fff;
	}

	.notif-btn {
		position: relative;
		background: none;
		border: none;
		font-size: 1.2rem;
		cursor: pointer;
		padding: 0.5rem;
	}

	.notif-badge {
		position: absolute;
		top: 0;
		right: 0;
		background: #ef4444;
		color: #fff;
		font-size: 0.65rem;
		font-weight: 600;
		min-width: 16px;
		height: 16px;
		border-radius: 0;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.notif-overlay {
		position: fixed;
		inset: 0;
		z-index: 90;
	}

	.notif-dropdown {
		position: fixed;
		top: 70px;
		right: 2rem;
		width: 320px;
		max-height: 400px;
		background: #1a1a1a;
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		overflow: hidden;
		z-index: 100;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.5);
	}

	.notif-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
		font-weight: 600;
	}

	.mark-read-btn {
		background: none;
		border: none;
		color: rgba(255, 255, 255, 0.6);
		font-size: 0.75rem;
		cursor: pointer;
	}

	.mark-read-btn:hover {
		color: #fff;
	}

	.notif-list {
		max-height: 340px;
		overflow-y: auto;
	}

	.notif-loading, .notif-empty {
		padding: 2rem;
		text-align: center;
		color: rgba(255, 255, 255, 0.5);
		font-size: 0.9rem;
	}

	.notif-item {
		padding: 1rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
		cursor: pointer;
		transition: background 0.2s ease;
	}

	.notif-item:hover {
		background: rgba(255, 255, 255, 0.03);
	}

	.notif-item.unread {
		background: rgba(255, 255, 255, 0.05);
	}

	.notif-text {
		margin: 0 0 0.25rem 0;
		font-size: 0.9rem;
		color: #fff;
	}

	.notif-time {
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.4);
	}

	@media (max-width: 768px) {
		.header {
			flex-direction: column;
			gap: 1rem;
			padding: 1rem;
		}

		.nav {
			gap: 1rem;
		}

		.notif-dropdown {
			right: 1rem;
			left: 1rem;
			width: auto;
		}
	}
</style>

