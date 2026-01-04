<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';

	onMount(async () => {
		// The backend handles the OAuth callback and sets cookies
		// This page is just for handling the redirect flow
		await auth.init();
		
		// Redirect based on auth state
		auth.subscribe(state => {
			if (state.initialized && !state.loading) {
				if (state.user && !state.profile?.is_complete) {
					goto('/profile/setup');
				} else if (state.user) {
					goto('/browse');
				} else {
					goto('/auth/login');
				}
			}
		});
	});
</script>

<svelte:head>
	<title>Authenticating... | HeySpoilMe</title>
</svelte:head>

<div class="callback-page">
	<div class="loader">
		<div class="spinner"></div>
		<p>Signing you in...</p>
	</div>
</div>

<style>
	.callback-page {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #0a0a0a;
	}

	.loader {
		text-align: center;
		color: #fff;
	}

	.spinner {
		width: 48px;
		height: 48px;
		border: 3px solid rgba(255, 255, 255, 0.1);
		border-top-color: #fff;
		border-radius: 0;
		animation: spin 1s linear infinite;
		margin: 0 auto 1rem;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	p {
		color: rgba(255, 255, 255, 0.6);
		font-family: 'Montserrat', sans-serif;
	}
</style>

