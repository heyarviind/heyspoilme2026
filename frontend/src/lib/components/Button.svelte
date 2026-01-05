<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	interface Props extends HTMLButtonAttributes {
		variant?: 'primary' | 'secondary';
		loading?: boolean;
		children: Snippet;
	}

	let { 
		variant = 'primary', 
		loading = false, 
		disabled = false,
		children,
		...rest 
	}: Props = $props();
</script>

<button 
	class="btn btn-{variant}" 
	disabled={disabled || loading}
	{...rest}
>
	{#if loading}
		<span class="spinner"></span>
	{:else}
		{@render children()}
	{/if}
</button>

<style>
	.btn {
		padding: 1rem 2rem;
		border: none;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.9rem;
		font-weight: 600;
		letter-spacing: 0.05em;
		cursor: pointer;
		transition: all 0.2s ease;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		min-width: 130px;
		border-radius: 0;
	}

	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	/* Primary: White solid */
	.btn-primary {
		background: #fff;
		color: #000;
		border: 1px solid #fff;
	}

	.btn-primary:hover:not(:disabled) {
		background: rgba(255, 255, 255, 0.9);
	}

	/* Secondary: Black bg, white border */
	.btn-secondary {
		background: #000;
		color: #fff;
		border: 1px solid #fff;
	}

	.btn-secondary:hover:not(:disabled) {
		background: rgba(255, 255, 255, 0.1);
	}

	.spinner {
		width: 20px;
		height: 20px;
		border: 2px solid rgba(128, 128, 128, 0.3);
		border-top-color: currentColor;
		border-radius: 0;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>



