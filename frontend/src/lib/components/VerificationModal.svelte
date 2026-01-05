<script lang="ts">
	interface Props {
		show: boolean;
		onClose: () => void;
	}

	let { show = $bindable(), onClose }: Props = $props();
</script>

{#if show}
	<div class="modal-overlay" onclick={onClose} role="button" tabindex="0" onkeypress={(e) => e.key === 'Enter' && onClose()}>
		<div class="modal verification-modal" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true">
			<button class="mobile-close" onclick={onClose} aria-label="Close">
				<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M18 6L6 18M6 6l12 12"/>
				</svg>
			</button>
			<div class="modal-content">
				<div class="modal-icon">🛡️</div>
				<h2>Become a Trusted Member</h2>
				<p class="modal-description">
					Verify your identity to unlock all features and connect with other members.
				</p>
				<div class="benefits-list">
					<div class="benefit-item">
						<span class="benefit-check">✓</span>
						<span>Like and match with profiles</span>
					</div>
					<div class="benefit-item">
						<span class="benefit-check">✓</span>
						<span>Send and receive messages</span>
					</div>
					<div class="benefit-item">
						<span class="benefit-check">✓</span>
						<span>Get a verified badge on your profile</span>
					</div>
					<div class="benefit-item">
						<span class="benefit-check">✓</span>
						<span>Appear higher in search results</span>
					</div>
				</div>
			</div>
			<div class="modal-actions">
				<a href="/profile/verify" class="btn-primary">
					Get Verified Now
				</a>
				<button class="btn-text" onclick={onClose}>
					Maybe Later
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.8);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 1rem;
	}

	.modal {
		background: #1a1a1a;
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		padding: 2rem;
		max-width: 400px;
		width: 100%;
		position: relative;
	}

	.verification-modal {
		text-align: center;
	}

	.mobile-close {
		display: none;
	}

	.modal-content {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.modal-icon {
		font-size: 3.5rem;
		margin-bottom: 1.5rem;
	}

	.modal h2 {
		font-family: 'Playfair Display', serif;
		font-size: 1.5rem;
		margin: 0 0 0.75rem 0;
	}

	.modal-description {
		color: rgba(255, 255, 255, 0.7);
		line-height: 1.6;
		margin: 0 0 1.5rem;
		font-size: 0.95rem;
	}

	.benefits-list {
		width: 100%;
		text-align: left;
		margin-bottom: 1rem;
	}

	.benefit-item {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.625rem 0;
		color: rgba(255, 255, 255, 0.9);
		font-size: 0.9rem;
	}

	.benefit-check {
		color: #22c55e;
		font-weight: 600;
	}

	.modal-actions {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		width: 100%;
		margin-top: 1rem;
	}

	.btn-primary {
		width: 100%;
		padding: 1rem;
		border-radius: 0;
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		cursor: pointer;
		border: none;
		text-align: center;
		text-decoration: none;
		background: #fff;
		color: #000;
		font-size: 0.95rem;
	}

	.btn-text {
		background: none;
		border: none;
		color: rgba(255, 255, 255, 0.5);
		font-family: 'Montserrat', sans-serif;
		font-size: 0.85rem;
		cursor: pointer;
		padding: 0.5rem;
	}

	.btn-text:hover {
		color: rgba(255, 255, 255, 0.8);
	}

	@media (max-width: 768px) {
		.modal-overlay {
			padding: 0;
			background: #0a0a0a;
		}

		.modal {
			max-width: none;
			height: 100%;
			display: flex;
			flex-direction: column;
			padding: 0;
			padding-bottom: 7rem;
			background: #0a0a0a;
			border: none;
			overflow-y: auto;
			-webkit-overflow-scrolling: touch;
		}

		.mobile-close {
			display: flex;
			position: fixed;
			top: 1rem;
			right: 1rem;
			width: 44px;
			height: 44px;
			align-items: center;
			justify-content: center;
			color: rgba(255, 255, 255, 0.7);
			background: none;
			border: none;
			cursor: pointer;
			z-index: 10;
		}

		.mobile-close:hover {
			color: #fff;
		}

		.modal-content {
			flex: 1;
			justify-content: center;
			padding: 4rem 1.5rem 2rem;
			flex-shrink: 0;
		}

		.modal-icon {
			font-size: 4rem;
			margin-bottom: 2rem;
		}

		.modal h2 {
			font-size: 1.75rem;
		}

		.modal-description {
			font-size: 1rem;
			margin-bottom: 2rem;
		}

		.benefits-list {
			max-width: 280px;
		}

		.benefit-item {
			font-size: 0.95rem;
			padding: 0.75rem 0;
		}

		.modal-actions {
			position: fixed;
			bottom: 0;
			left: 0;
			right: 0;
			padding: 1rem 1.5rem 2rem;
			background: #0a0a0a;
			border-top: 1px solid rgba(255, 255, 255, 0.1);
			margin-top: 0;
		}

		.btn-primary {
			padding: 1rem;
			font-size: 1rem;
		}
	}
</style>

