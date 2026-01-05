<script lang="ts">
	import { api } from '$lib/api';

	let sending = $state(false);
	let sent = $state(false);
	let error = $state<string | null>(null);

	async function resendEmail() {
		if (sending || sent) return;
		sending = true;
		error = null;

		try {
			await api.resendVerificationEmail();
			sent = true;
			setTimeout(() => { sent = false; }, 60000); // Allow resend after 1 minute
		} catch (e: any) {
			error = e.message || 'Failed to send verification email';
		} finally {
			sending = false;
		}
	}
</script>

<div class="verification-banner">
	<div class="banner-content">
		<span class="icon">✉️</span>
		<div class="text">
			<strong>Verify your email</strong>
			<span class="subtitle">Check your inbox to unlock messaging, profile browsing, and image uploads.</span>
		</div>
		<button 
			class="resend-btn" 
			onclick={resendEmail}
			disabled={sending || sent}
		>
			{#if sending}
				Sending...
			{:else if sent}
				Email Sent ✓
			{:else}
				Resend Email
			{/if}
		</button>
	</div>
	{#if error}
		<p class="error">{error}</p>
	{/if}
</div>

<style>
	.verification-banner {
		background: rgba(255, 56, 92, 0.08);
		border-bottom: 1px solid rgba(255, 56, 92, 0.3);
		padding: 0.875rem 1.5rem;
		position: sticky;
		top: 0;
		z-index: 1000;
	}

	.banner-content {
		max-width: 1200px;
		margin: 0 auto;
		display: flex;
		align-items: center;
		gap: 1rem;
		flex-wrap: wrap;
	}

	.icon {
		font-size: 1.25rem;
	}

	.text {
		flex: 1;
		min-width: 200px;
	}

	.text strong {
		color: #FF385C;
		display: block;
		font-size: 0.9rem;
	}

	.subtitle {
		color: rgba(255, 255, 255, 0.6);
		font-size: 0.8rem;
	}

	.resend-btn {
		background: rgba(255, 56, 92, 0.15);
		border: 1px solid rgba(255, 56, 92, 0.4);
		color: #FF385C;
		padding: 0.5rem 1rem;
		font-size: 0.8rem;
		font-weight: 500;
		font-family: 'Montserrat', sans-serif;
		cursor: pointer;
		transition: all 0.2s ease;
		white-space: nowrap;
	}

	.resend-btn:hover:not(:disabled) {
		background: rgba(255, 56, 92, 0.25);
		border-color: rgba(255, 56, 92, 0.6);
	}

	.resend-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.error {
		color: #ef4444;
		font-size: 0.8rem;
		margin: 0.5rem 0 0 0;
		text-align: center;
	}

	@media (max-width: 768px) {
		.verification-banner {
			padding: 0.75rem 1rem;
		}

		.banner-content {
			gap: 0.75rem;
		}

		.icon {
			display: none;
		}

		.text strong {
			font-size: 0.85rem;
		}

		.subtitle {
			font-size: 0.75rem;
		}

		.resend-btn {
			width: 100%;
			padding: 0.625rem 1rem;
		}
	}
</style>

