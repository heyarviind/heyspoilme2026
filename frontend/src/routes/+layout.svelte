<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	import { auth } from '$lib/stores/auth';
	import { websocket } from '$lib/stores/websocket';
	import { EmailVerificationBanner, VerificationBanner } from '$lib/components';

	let authState = $state<{ user: any; profile: any; loading: boolean; initialized: boolean }>({ user: null, profile: null, loading: true, initialized: false });
	let unsubscribe: (() => void) | null = null;
	let currentPath = $state('/');
	let bannerReady = $state(false); // Delay banner display to prevent flash

	// Pages where we should NOT show the email verification banner
	const noEmailBannerPaths = [
		'/',
		'/auth/login',
		'/auth/callback',
		'/auth/error',
		'/auth/verify-email',
		'/admin',
	];

	// Pages where we should NOT show the identity verification banner
	const noIdentityBannerPaths = [
		'/',
		'/auth/login',
		'/auth/callback',
		'/auth/error',
		'/auth/verify-email',
		'/profile/verify',
		'/profile/setup',
		'/admin',
	];


	onMount(async () => {
		// Subscribe FIRST to catch all state changes during init
		unsubscribe = auth.subscribe(state => {
			authState = state;
			if (state.user && state.initialized) {
				const token = localStorage.getItem('token');
				if (token) {
					websocket.connect(token);
				}
			}
		});

		// Track current path
		page.subscribe(p => {
			currentPath = p.url.pathname;
		});

		// Initialize auth (subscription will receive loading state updates)
		await auth.init();

		// Delay banner display by 1 second to prevent flash on page load
		setTimeout(() => {
			bannerReady = true;
		}, 1000);
	});

	onDestroy(() => {
		if (unsubscribe) unsubscribe();
		websocket.disconnect();
	});
</script>

{#if bannerReady && authState.initialized && !authState.loading && authState.user && !authState.user.email_verified && !noEmailBannerPaths.some(p => currentPath === p || currentPath.startsWith(p + '?') || currentPath.startsWith(p + '/'))}
	<EmailVerificationBanner />
{:else if bannerReady && authState.initialized && !authState.loading && authState.user && authState.user.email_verified && authState.profile && !authState.profile.is_verified && !noIdentityBannerPaths.some(p => currentPath === p || currentPath.startsWith(p + '?') || currentPath.startsWith(p + '/'))}
	<VerificationBanner />
{/if}

<slot />

<style>
	:global(*) {
		margin: 0;
		padding: 0;
		box-sizing: border-box;
	}

	:global(html) {
		scroll-behavior: smooth;
	}

	:global(body) {
		font-family: 'Montserrat', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
		background: #0a0a0a;
		color: #fff;
		min-height: 100vh;
		-webkit-font-smoothing: antialiased;
		-moz-osx-font-smoothing: grayscale;
	}

	:global(::selection) {
		background: rgba(255, 255, 255, 0.2);
	}

	:global(::-webkit-scrollbar) {
		width: 8px;
	}

	:global(::-webkit-scrollbar-track) {
		background: #0a0a0a;
	}

	:global(::-webkit-scrollbar-thumb) {
		background: rgba(255, 255, 255, 0.2);
		border-radius: 0;
	}

	:global(::-webkit-scrollbar-thumb:hover) {
		background: rgba(255, 255, 255, 0.3);
	}
</style>
