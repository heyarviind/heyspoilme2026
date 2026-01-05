<script lang="ts">
	import { goto } from '$app/navigation';

	function goToVerify() {
		goto('/profile/verify');
	}
</script>

<div class="verification-banner">
	<div class="banner-content">
		<div class="icon-badge">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
				<path d="M9 12l2 2 4-4"/>
			</svg>
		</div>
		<div class="text">
			<strong>Get Verified</strong>
			<span class="subtitle">Verified profiles get 3x more matches and appear at the top of search results</span>
		</div>
		<button class="verify-btn" onclick={goToVerify}>
			Verify Now
		</button>
	</div>
</div>

<style>
	.verification-banner {
		background: #fff;
		border-bottom: 1px solid rgba(255, 255, 255, 0.2);
		padding: 0.875rem 1.5rem;
	}

	.banner-content {
		max-width: 1200px;
		margin: 0 auto;
		display: flex;
		align-items: center;
		gap: 1rem;
		flex-wrap: wrap;
	}

	.icon-badge {
		width: 36px;
		height: 36px;
		background: rgba(0, 0, 0, 0.1);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: #0a0a0a;
		flex-shrink: 0;
	}

	.text {
		flex: 1;
		min-width: 200px;
	}

	.text strong {
		color: #0a0a0a;
		display: block;
		font-size: 0.9rem;
	}

	.subtitle {
		color: rgba(0, 0, 0, 0.6);
		font-size: 0.8rem;
	}

	.verify-btn {
		background: #0a0a0a;
		border: none;
		color: #fff;
		padding: 0.625rem 1.25rem;
		font-size: 0.85rem;
		font-weight: 600;
		font-family: 'Montserrat', sans-serif;
		cursor: pointer;
		transition: all 0.2s ease;
		white-space: nowrap;
	}

	.verify-btn:hover {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
	}

	@media (max-width: 768px) {
		.verification-banner {
			padding: 0.75rem 1rem;
		}

		.banner-content {
			gap: 0.75rem;
		}

		.icon-badge {
			width: 32px;
			height: 32px;
		}

		.text strong {
			font-size: 0.85rem;
		}

		.subtitle {
			font-size: 0.75rem;
		}

		.verify-btn {
			width: 100%;
			padding: 0.75rem 1rem;
		}
	}
</style>

