<script lang="ts">
	import { page } from '$app/stores';

	let message = $derived($page.url.searchParams.get('message') || 'An error occurred during authentication');
</script>

<svelte:head>
	<title>Authentication Error | HeySpoilMe</title>
</svelte:head>

<div class="error-page">
	<div class="error-card">
		<div class="icon">⚠️</div>
		<h1>Authentication Failed</h1>
		<p class="message">{message.replace(/_/g, ' ')}</p>
		<a href="/auth/login" class="btn">Try Again</a>
	</div>
</div>

<style>
	.error-page {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #0a0a0a;
		padding: 2rem;
	}

	.error-card {
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		padding: 3rem;
		max-width: 400px;
		width: 100%;
		text-align: center;
	}

	.icon {
		font-size: 3rem;
		margin-bottom: 1rem;
	}

	h1 {
		font-family: 'Playfair Display', serif;
		font-size: 1.5rem;
		color: #fff;
		margin: 0 0 1rem 0;
	}

	.message {
		color: rgba(255, 255, 255, 0.6);
		margin: 0 0 2rem 0;
		text-transform: capitalize;
	}

	.btn {
		display: inline-block;
		padding: 0.75rem 2rem;
		background: #fff;
		color: #000;
		text-decoration: none;
		border-radius: 0;
		font-family: 'Montserrat', sans-serif;
		font-weight: 500;
		transition: all 0.2s ease;
	}

	.btn:hover {
		transform: translateY(-2px);
	}
</style>

