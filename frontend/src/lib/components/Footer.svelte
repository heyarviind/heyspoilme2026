<script lang="ts">
</script>

<footer class="footer">
	<div class="footer-content">
		<div class="footer-brand">
			<img src="/img/logo.svg" alt="HeySpoilMe" class="logo" />
			<p>A private platform for premium connections</p>
		</div>
		
		<div class="footer-links">
			<div class="footer-col">
				<h5>Company</h5>
				<a href="/about">About Us</a>
				<a href="/contact">Contact</a>
			</div>
			<div class="footer-col">
				<h5>Legal</h5>
				<a href="/terms">Terms of Service</a>
				<a href="/privacy">Privacy Policy</a>
			</div>
			<div class="footer-col">
				<h5>Support</h5>
				<a href="/help">Help Center</a>
				<a href="/safety">Safety Guidelines</a>
			</div>
		</div>
	</div>
	
	<div class="footer-bottom">
		<span>© 2026 HeySpoilMe. All rights reserved.</span>
		<span>21+ Adults Only</span>
	</div>
</footer>

<style>
	.footer {
		padding: 4rem 2rem;
		border-top: 1px solid rgba(255, 255, 255, 0.05);
		margin-top: 4rem;
	}

	.footer-content {
		max-width: 1100px;
		margin: 0 auto;
		display: flex;
		justify-content: space-between;
		padding-bottom: 3rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
	}

	.footer-brand {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 0.75rem;
	}

	.logo {
		height: 1.75rem;
	}

	.footer-brand p {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.4);
		margin: 0;
	}

	.footer-links {
		display: flex;
		gap: 5rem;
	}

	.footer-col {
		display: flex;
		flex-direction: column;
	}

	.footer-col h5 {
		font-size: 0.75rem;
		font-weight: 500;
		letter-spacing: 0.1em;
		text-transform: uppercase;
		color: rgba(255, 255, 255, 0.5);
		margin: 0 0 1rem 0;
	}

	.footer-col a {
		display: block;
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.6);
		text-decoration: none;
		margin-bottom: 0.75rem;
		transition: color 0.2s ease;
	}

	.footer-col a:hover {
		color: #fff;
	}

	.footer-bottom {
		max-width: 1100px;
		margin: 0 auto;
		display: flex;
		justify-content: space-between;
		padding-top: 2rem;
	}

	.footer-bottom span {
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.3);
	}

	/* Mobile */
	@media (max-width: 768px) {
		.footer {
			padding: 3rem 1.5rem;
			margin-top: 2rem;
		}

		.footer-content {
			flex-direction: column;
			gap: 2.5rem;
		}

		.footer-links {
			flex-direction: column;
			gap: 1.5rem;
		}

		.footer-bottom {
			flex-direction: column;
			gap: 0.5rem;
			text-align: center;
		}
	}
</style>

