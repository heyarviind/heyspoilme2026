<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { api } from '$lib/api';
	import { auth } from '$lib/stores/auth';

	let status = $state<'loading' | 'success' | 'error'>('loading');
	let errorMessage = $state('');

	onMount(async () => {
		const token = $page.url.searchParams.get('token');

		if (!token) {
			status = 'error';
			errorMessage = 'Invalid verification link';
			return;
		}

		try {
			await api.verifyEmail(token);
			status = 'success';
			
			// Update auth store to reflect verified status
			auth.setEmailVerified(true);
			
			// Redirect to browse after a short delay
			setTimeout(() => {
				goto('/browse');
			}, 3000);
		} catch (e: any) {
			status = 'error';
			errorMessage = e.message || 'Verification failed';
		}
	});
</script>

<svelte:head>
	<title>Verify Email | HeySpoilMe</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;500;600;700&family=Montserrat:wght@300;400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="verify-page">
	<div class="container">
		{#if status === 'loading'}
			<div class="spinner"></div>
			<h1>Verifying your email...</h1>
			<p>Please wait while we verify your email address.</p>
		{:else if status === 'success'}
			<div class="success-icon">✓</div>
			<h1>Email Verified!</h1>
			<p>Your email has been successfully verified. You now have full access to HeySpoilMe.</p>
			<p class="redirect">Redirecting to browse...</p>
			<button class="btn" onclick={() => goto('/browse')}>Go to Browse</button>
		{:else}
			<div class="error-icon">✕</div>
			<h1>Verification Failed</h1>
			<p class="error-message">{errorMessage}</p>
			<p>The verification link may have expired or already been used.</p>
			<button class="btn" onclick={() => goto('/browse')}>Go to App</button>
		{/if}
	</div>
</div>

<style>
	:global(body) {
		font-family: 'Montserrat', sans-serif;
		background: #0a0a0a;
		color: #fff;
		margin: 0;
	}

	.verify-page {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 2rem;
	}

	.container {
		text-align: center;
		max-width: 400px;
	}

	h1 {
		font-family: 'Playfair Display', serif;
		font-size: 1.75rem;
		font-weight: 500;
		margin: 1.5rem 0 0.75rem 0;
	}

	p {
		color: rgba(255, 255, 255, 0.7);
		margin: 0 0 1rem 0;
		line-height: 1.5;
	}

	.error-message {
		color: #ef4444;
		font-weight: 500;
	}

	.redirect {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.5);
	}

	.spinner {
		width: 48px;
		height: 48px;
		border: 3px solid rgba(255, 255, 255, 0.1);
		border-top-color: #fff;
		animation: spin 1s linear infinite;
		margin: 0 auto;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.success-icon {
		width: 64px;
		height: 64px;
		background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
		color: #fff;
		font-size: 2rem;
		font-weight: bold;
		display: flex;
		align-items: center;
		justify-content: center;
		margin: 0 auto;
	}

	.error-icon {
		width: 64px;
		height: 64px;
		background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
		color: #fff;
		font-size: 2rem;
		font-weight: bold;
		display: flex;
		align-items: center;
		justify-content: center;
		margin: 0 auto;
	}

	.btn {
		display: inline-block;
		background: #fff;
		color: #000;
		padding: 0.875rem 2rem;
		text-decoration: none;
		font-weight: 600;
		font-size: 0.9rem;
		font-family: inherit;
		margin-top: 1rem;
		border: none;
		cursor: pointer;
		transition: transform 0.2s ease;
	}

	.btn:hover {
		transform: scale(1.02);
	}
</style>


