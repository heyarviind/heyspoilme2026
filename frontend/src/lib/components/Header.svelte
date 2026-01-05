<script lang="ts">
	import { page } from '$app/stores';
	import { notifications } from '$lib/stores/notifications';
	import { api } from '$lib/api';
	import { onMount } from 'svelte';

	let notifState = $state<{ unreadCount: number }>({ unreadCount: 0 });
	let unreadMessages = $state(0);

	notifications.subscribe(s => notifState = s);

	async function fetchUnreadMessages() {
		try {
			const data = await api.getUnreadMessageCount() as { count: number };
			unreadMessages = data.count;
		} catch {
			// Ignore errors
		}
	}

	let currentPath = $derived($page.url.pathname);

	onMount(() => {
		notifications.fetchUnreadCount();
		fetchUnreadMessages();
		// Refresh counts periodically
		const interval = setInterval(() => {
			notifications.fetchUnreadCount();
			fetchUnreadMessages();
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
			Messages{#if unreadMessages > 0}<span class="count-badge">{unreadMessages > 99 ? '99+' : unreadMessages}</span>{/if}
		</a>
		<a href="/likes" class="nav-link" class:active={currentPath === '/likes'}>
			Likes{#if notifState?.unreadCount > 0}<span class="count-badge">{notifState.unreadCount > 99 ? '99+' : notifState.unreadCount}</span>{/if}
		</a>
		<a href="/profile" class="nav-link" class:active={currentPath === '/profile'}>
			Profile
		</a>
	</nav>
</header>

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

	.count-badge {
		background: #ef4444;
		color: #fff;
		font-size: 0.7rem;
		font-weight: 600;
		padding: 0.15rem 0.4rem;
		margin-left: 0.35rem;
		min-width: 1.25rem;
		text-align: center;
		display: inline-block;
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
	}
</style>

