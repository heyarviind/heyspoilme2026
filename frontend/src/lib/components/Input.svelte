<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';

	interface Props extends HTMLInputAttributes {
		variant?: 'primary' | 'secondary';
		value?: string;
	}

	let { 
		variant = 'primary',
		value = $bindable(''),
		...rest 
	}: Props = $props();
</script>

<input 
	class="input input-{variant}" 
	bind:value
	{...rest}
/>

<style>
	.input {
		padding: 1rem 1.5rem;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.95rem;
		outline: none;
		transition: all 0.2s ease;
		border-radius: 0;
		width: 100%;
	}

	/* Primary: Black bg, white border */
	.input-primary {
		background: #000;
		color: #fff;
		border: 1px solid #fff;
	}

	.input-primary::placeholder {
		color: rgba(255, 255, 255, 0.4);
	}

	.input-primary:focus {
		border-color: #fff;
		box-shadow: 0 0 0 1px #fff;
	}

	/* Secondary: Subtle border */
	.input-secondary {
		background: #000;
		color: #fff;
		border: 1px solid rgba(255, 255, 255, 0.3);
	}

	.input-secondary::placeholder {
		color: rgba(255, 255, 255, 0.4);
	}

	.input-secondary:focus {
		border-color: #fff;
	}

	.input:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}
</style>



