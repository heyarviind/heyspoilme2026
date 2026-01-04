<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { websocket } from '$lib/stores/websocket';

	let authState = $state<{ user: any; loading: boolean; initialized: boolean }>({ user: null, loading: true, initialized: false });
	let unsubscribe: (() => void) | null = null;

	onMount(async () => {
		await auth.init();
		
		// Connect WebSocket if authenticated
		unsubscribe = auth.subscribe(state => {
			authState = state;
			if (state.user && state.initialized) {
				const token = localStorage.getItem('token');
				if (token) {
					websocket.connect(token);
				}
			}
		});
	});

	onDestroy(() => {
		if (unsubscribe) unsubscribe();
		websocket.disconnect();
	});
</script>

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
