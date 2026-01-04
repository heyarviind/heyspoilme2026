<script lang="ts">
</script>

<svelte:head>
	<title>HeySpoilMe — Where Luxury Meets Connection</title>
	<meta name="description" content="A private, consent-first platform for modern sugar dating in India. Verified profiles, complete privacy." />
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@300;400;500;600&family=Outfit:wght@300;400;500&display=swap" rel="stylesheet">
</svelte:head>

<div class="landing">
	<!-- Hero Section -->
	<section class="hero">
		<nav class="nav">
			<img src="/img/logo.svg" alt="HeySpoilMe" class="logo" />
			<a href="/auth/login" class="sign-in-btn">Sign In</a>
		</nav>

		<div class="hero-content">
			<h1>
				<span class="line1">Where Luxury</span>
				<span class="line2">Meets Connection</span>
			</h1>
			
			<p class="tagline">
				India's most private platform for modern sugar dating. 
				Verified profiles. Complete discretion.
			</p>

			<div class="hero-cta">
				<a href="/auth/login" class="btn-primary">Get Started</a>
				<span class="cta-note">Free to join · 21+ only</span>
			</div>
		</div>

		<div class="scroll-hint">
			<span>Scroll</span>
			<div class="scroll-line"></div>
		</div>
	</section>

	<!-- Features Section -->
	<section class="features">
		<div class="section-header">
			<span class="section-label">Why HeySpoilMe</span>
			<h2>Built for Genuine Connections</h2>
		</div>

		<div class="features-grid">
			<div class="feature-card">
				<div class="feature-icon">
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
						<path d="M8 11.9997L12.9999 5.99994L13 10.9999L15.9999 10.9997L11 16.9999V11.9999L8 11.9997Z"/>
						<path d="M3.43067 4.93443C6.56932 5.16025 8.59312 2.00103 12.0266 2.00103C15.4086 1.93395 17.0635 4.97662 20.5307 4.97662C22.3893 14.3836 18.5226 20.2927 12.0549 21.9999C6.0595 20.583 1.5547 14.6133 3.43067 4.93443Z"/>
					</svg>
				</div>
				<h3>Verified Profiles</h3>
				<p>Every member goes through our verification process. Real people, real intentions.</p>
			</div>

			<div class="feature-card">
				<div class="feature-icon">
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
						<path d="M21.9598 10.9707C21.9711 11.146 21.98 11.3226 21.9866 11.5C22.011 12.1627 22.0021 12.836 21.9598 13.4909C21.6856 17.7332 18.3536 21.1125 14.1706 21.3905C12.7435 21.4854 11.2536 21.4852 9.8294 21.3905C9.33896 21.3579 8.8044 21.2409 8.34401 21.0513C7.83177 20.8403 7.5756 20.7348 7.44544 20.7508C7.31527 20.7668 7.1264 20.9061 6.74868 21.1846C6.08268 21.6757 5.24367 22.0285 3.99943 21.9982C3.37026 21.9829 3.05568 21.9752 2.91484 21.7351C2.77401 21.495 2.94941 21.1626 3.30021 20.4978C3.78674 19.5758 4.09501 18.5203 3.62791 17.6746C2.82343 16.4666 2.1401 15.036 2.04024 13.4909C1.98659 12.6607 1.98659 11.8009 2.04024 10.9707C2.31441 6.72836 5.64639 3.34912 9.8294 3.07106C10.3789 3.03453 10.4383 3.0121 11 3.00374"/>
						<path d="M11.9953 12.5H12.0043M15.9908 12.5H15.9998M7.99982 12.5H8.00879"/>
						<path d="M15.015 2.3866C16.0876 1.74691 17.0238 2.0047 17.5863 2.41532C17.8169 2.58369 17.9322 2.66787 18 2.66787C18.0678 2.66787 18.1831 2.58369 18.4137 2.41532C18.9762 2.0047 19.9124 1.74691 20.985 2.3866C22.3928 3.22613 22.7113 5.99576 19.4642 8.3324C18.8457 8.77746 18.5365 8.99998 18 8.99998C17.4635 8.99998 17.1543 8.77746 16.5358 8.3324C13.2887 5.99576 13.6072 3.22613 15.015 2.3866Z"/>
					</svg>
				</div>
				<h3>Women Message First</h3>
				<p>Ladies are in control. No unwanted messages, no harassment. Quality over quantity.</p>
			</div>

			<div class="feature-card">
				<div class="feature-icon">
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linejoin="round">
						<path d="M12 20.5C12 20.5 2 14.5 2 8.69444C2 5.82563 4.10526 3.5 7 3.5C8.5 3.5 10 4 12 6C14 4 15.5 3.5 17 3.5C19.8947 3.5 22 5.82563 22 8.69444C22 14.5 12 20.5 12 20.5Z"/>
					</svg>
				</div>
				<h3>Genuine Matches</h3>
				<p>Find meaningful connections based on compatibility, not just appearances.</p>
			</div>

			<div class="feature-card">
				<div class="feature-icon">
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
						<path d="M6 18C4.75527 18.396 4 19.184 4 19.7537C4 20.9943 7.58172 22 12 22C16.4183 22 20 20.9943 20 19.7537C20 19.184 19.2447 18.396 18 18"/>
						<path d="M14.5 9C14.5 10.3807 13.3807 11.5 12 11.5C10.6193 11.5 9.5 10.3807 9.5 9C9.5 7.61929 10.6193 6.5 12 6.5C13.3807 6.5 14.5 7.61929 14.5 9Z"/>
						<path d="M19 9.2C19 13.9809 13.8979 17.0126 12.4054 17.7971C12.1485 17.9322 11.8515 17.9322 11.5946 17.7971C10.1021 17.0126 5 13.9809 5 9.2C5 5.22355 8.13401 2 12 2C15.866 2 19 5.22355 19 9.2Z"/>
					</svg>
				</div>
				<h3>Meet Nearby</h3>
				<p>Connect with verified members in your city. Distance-based matching for real meetups.</p>
			</div>
		</div>
	</section>

	<!-- Verification Section -->
	<section class="verification">
		<div class="verification-content">
			<div class="verification-badge">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
					<path d="M8 11.9997L12.9999 5.99994L13 10.9999L15.9999 10.9997L11 16.9999V11.9999L8 11.9997Z"/>
					<path d="M3.43067 4.93443C6.56932 5.16025 8.59312 2.00103 12.0266 2.00103C15.4086 1.93395 17.0635 4.97662 20.5307 4.97662C22.3893 14.3836 18.5226 20.2927 12.0549 21.9999C6.0595 20.583 1.5547 14.6133 3.43067 4.93443Z"/>
				</svg>
			</div>
			<span class="section-label">Trust & Safety</span>
			<h2>Every Profile is Verified</h2>
			<p>We take verification seriously. Each member undergoes identity checks to ensure you're connecting with real, genuine people.</p>
			
			<ul class="verification-list">
				<li>
					<span class="check">✓</span>
					Photo verification to match profile pictures
				</li>
				<li>
					<span class="check">✓</span>
					Age verification (21+ only)
				</li>
				<li>
					<span class="check">✓</span>
					Income verification for Sugar Daddies
				</li>
				<li>
					<span class="check">✓</span>
					Manual review of every profile
				</li>
			</ul>
		</div>
	</section>

	<!-- Privacy Section -->
	<section class="privacy">
		<div class="privacy-grid">
			<div class="privacy-content">
				<span class="section-label">Your Privacy Matters</span>
				<h2>Complete Discretion, Always</h2>
				<p>Your privacy is our top priority. We've built every feature with discretion in mind.</p>
			</div>
			
			<div class="privacy-features">
				<div class="privacy-item">
					<div class="privacy-icon">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
						</svg>
					</div>
					<div>
						<h4>Private Photos</h4>
						<p>Control who sees your photos</p>
					</div>
				</div>
				
				<div class="privacy-item">
					<div class="privacy-icon">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
						</svg>
					</div>
					<div>
						<h4>Browse Anonymously</h4>
						<p>View profiles without being seen</p>
					</div>
				</div>
				
				<div class="privacy-item">
					<div class="privacy-icon">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
						</svg>
					</div>
					<div>
						<h4>Discreet Notifications</h4>
						<p>No revealing app names or icons</p>
					</div>
				</div>
				
				<div class="privacy-item">
					<div class="privacy-icon">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
						</svg>
					</div>
					<div>
						<h4>Delete Anytime</h4>
						<p>Full control over your data</p>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- CTA Section -->
	<section class="cta">
		<h2>Ready to Find Your Match?</h2>
		<p>Join thousands of verified members already on HeySpoilMe</p>
		<a href="/auth/login" class="btn-primary">Create Free Account</a>
	</section>

	<!-- Footer -->
	<footer class="footer">
		<div class="footer-content">
			<div class="footer-brand">
				<img src="/img/logo.svg" alt="HeySpoilMe" class="logo" />
				<p>India's premier sugar dating platform</p>
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
					<a href="/safety">Safety Tips</a>
				</div>
			</div>
		</div>
		
		<div class="footer-bottom">
			<span>© 2026 HeySpoilMe. All rights reserved.</span>
			<span>21+ Adults Only</span>
		</div>
		</footer>
</div>

<style>
	:global(*) {
		margin: 0;
		padding: 0;
		box-sizing: border-box;
	}
	
	:global(body) {
		font-family: 'Outfit', sans-serif;
		background: #08080a;
		color: #fff;
		overflow-x: hidden;
	}
	
	.landing {
		display: flex;
		flex-direction: column;
	}
	
	/* Nav */
	.nav {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1.5rem 4rem;
		z-index: 100;
		background: linear-gradient(180deg, rgba(8, 8, 10, 0.95) 0%, transparent 100%);
	}

	.logo {
		height: 2rem;
	}

	.sign-in-btn {
		padding: 0.6rem 1.5rem;
		background: transparent;
		border: 1px solid rgba(255, 255, 255, 0.3);
		color: #fff;
		font-family: 'Outfit', sans-serif;
		font-size: 0.85rem;
		text-decoration: none;
		transition: all 0.3s ease;
	}

	.sign-in-btn:hover {
		background: #fff;
		color: #000;
	}

	/* Hero */
	.hero {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 8rem 2rem 4rem;
		position: relative;
		background: radial-gradient(ellipse at top, rgba(255, 255, 255, 0.03) 0%, transparent 50%);
	}

	.hero-content {
		text-align: center;
		max-width: 700px;
	}

	h1 {
		font-family: 'Cormorant Garamond', serif;
		font-size: clamp(3rem, 8vw, 5rem);
		font-weight: 300;
		line-height: 1.1;
		letter-spacing: -0.02em;
		margin-bottom: 1.5rem;
	}

	h1 .line1, h1 .line2 {
		display: block;
	}

	h1 .line2 {
		color: rgba(255, 255, 255, 0.5);
	}

	.tagline {
		font-size: 1.1rem;
		font-weight: 300;
		color: rgba(255, 255, 255, 0.5);
		line-height: 1.7;
		max-width: 500px;
		margin: 0 auto 2.5rem;
	}

	.hero-cta {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
	}

	.btn-primary {
		padding: 1rem 3rem;
		background: #fff;
		border: none;
		color: #000;
		font-family: 'Outfit', sans-serif;
		font-size: 0.95rem;
		font-weight: 500;
		text-decoration: none;
		transition: all 0.3s ease;
	}

	.btn-primary:hover {
		transform: translateY(-2px);
		box-shadow: 0 10px 40px rgba(255, 255, 255, 0.15);
	}

	.cta-note {
		font-size: 0.8rem;
		color: rgba(255, 255, 255, 0.4);
	}

	.scroll-hint {
		position: absolute;
		bottom: 3rem;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.75rem;
	}

	.scroll-hint span {
		font-size: 0.7rem;
		letter-spacing: 0.15em;
		text-transform: uppercase;
		color: rgba(255, 255, 255, 0.3);
	}

	.scroll-line {
		width: 1px;
		height: 40px;
		background: linear-gradient(180deg, rgba(255, 255, 255, 0.3), transparent);
	}

	/* Sections Common */
	.section-label {
		display: block;
		font-size: 0.7rem;
		font-weight: 500;
		letter-spacing: 0.2em;
		text-transform: uppercase;
		color: rgba(255, 255, 255, 0.4);
		margin-bottom: 1rem;
	}

	section h2 {
		font-family: 'Cormorant Garamond', serif;
		font-size: clamp(2rem, 4vw, 2.75rem);
		font-weight: 400;
		margin-bottom: 1rem;
	}

	/* Features */
	.features {
		padding: 8rem 4rem;
		border-top: 1px solid rgba(255, 255, 255, 0.05);
	}

	.section-header {
		text-align: center;
		margin-bottom: 4rem;
	}

	.features-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 2rem;
		max-width: 1200px;
		margin: 0 auto;
	}

	.feature-card {
		padding: 2rem;
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid rgba(255, 255, 255, 0.05);
		transition: all 0.3s ease;
	}

	.feature-card:hover {
		background: rgba(255, 255, 255, 0.04);
		border-color: rgba(255, 255, 255, 0.1);
	}

	.feature-icon {
		width: 48px;
		height: 48px;
		margin-bottom: 1.5rem;
	}
	
	.feature-icon svg {
		width: 100%;
		height: 100%;
		color: rgba(255, 255, 255, 0.6);
	}

	.feature-card h3 {
		font-family: 'Outfit', sans-serif;
		font-size: 1rem;
		font-weight: 500;
		margin-bottom: 0.75rem;
	}

	.feature-card p {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.5);
		line-height: 1.6;
	}

	/* Verification */
	.verification {
		padding: 8rem 4rem;
		background: rgba(255, 255, 255, 0.02);
		border-top: 1px solid rgba(255, 255, 255, 0.05);
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
	}

	.verification-content {
		max-width: 600px;
		margin: 0 auto;
		text-align: center;
	}

	.verification-badge {
		width: 80px;
		height: 80px;
		margin: 0 auto 2rem;
		padding: 1.5rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.verification-badge svg {
		width: 100%;
		height: 100%;
		color: #fff;
	}

	.verification-content > p {
		font-size: 1rem;
		color: rgba(255, 255, 255, 0.5);
		line-height: 1.7;
		margin-bottom: 2.5rem;
	}
	
	.verification-list {
		list-style: none;
		text-align: left;
		display: inline-block;
	}

	.verification-list li {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 0.75rem 0;
		font-size: 0.95rem;
		color: rgba(255, 255, 255, 0.7);
	}

	.verification-list .check {
		color: #fff;
		font-weight: 600;
	}

	/* Privacy */
	.privacy {
		padding: 8rem 4rem;
	}

	.privacy-grid {
		max-width: 1100px;
		margin: 0 auto;
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 4rem;
		align-items: center;
	}

	.privacy-content p {
		font-size: 1rem;
		color: rgba(255, 255, 255, 0.5);
		line-height: 1.7;
	}

	.privacy-features {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1.5rem;
	}

	.privacy-item {
		display: flex;
		gap: 1rem;
		padding: 1.25rem;
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid rgba(255, 255, 255, 0.05);
	}

	.privacy-icon {
		width: 40px;
		height: 40px;
		flex-shrink: 0;
	}

	.privacy-icon svg {
		width: 100%;
		height: 100%;
		color: rgba(255, 255, 255, 0.5);
	}

	.privacy-item h4 {
		font-size: 0.9rem;
		font-weight: 500;
		margin-bottom: 0.25rem;
	}

	.privacy-item p {
		font-size: 0.8rem;
		color: rgba(255, 255, 255, 0.4);
	}

	/* CTA */
	.cta {
		padding: 8rem 4rem;
		text-align: center;
		background: linear-gradient(180deg, transparent 0%, rgba(255, 255, 255, 0.02) 100%);
		border-top: 1px solid rgba(255, 255, 255, 0.05);
	}

	.cta p {
		font-size: 1rem;
		color: rgba(255, 255, 255, 0.5);
		margin-bottom: 2rem;
	}

	/* Footer */
	.footer {
		padding: 4rem;
		border-top: 1px solid rgba(255, 255, 255, 0.05);
	}

	.footer-content {
		max-width: 1100px;
		margin: 0 auto;
		display: flex;
		justify-content: space-between;
		padding-bottom: 3rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
	}

	.footer-brand p {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.4);
		margin-top: 0.75rem;
	}

	.footer-links {
		display: flex;
		gap: 5rem;
	}

	.footer-col h5 {
		font-size: 0.75rem;
		font-weight: 500;
		letter-spacing: 0.1em;
		text-transform: uppercase;
		color: rgba(255, 255, 255, 0.5);
		margin-bottom: 1rem;
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
	@media (max-width: 1024px) {
		.features-grid {
			grid-template-columns: repeat(2, 1fr);
		}

		.privacy-grid {
			grid-template-columns: 1fr;
			gap: 3rem;
		}

		.privacy-content {
			text-align: center;
		}
	}

	@media (max-width: 768px) {
		.nav {
			padding: 1.25rem 1.5rem;
		}

		.hero {
			padding: 7rem 1.5rem 4rem;
		}

		.features, .verification, .privacy, .cta {
			padding: 5rem 1.5rem;
		}

		.features-grid {
			grid-template-columns: 1fr;
			gap: 1.5rem;
		}

		.privacy-features {
			grid-template-columns: 1fr;
		}

		.footer {
			padding: 3rem 1.5rem;
		}

		.footer-content {
			flex-direction: column;
			gap: 3rem;
		}

		.footer-links {
			flex-direction: column;
			gap: 2rem;
		}

		.footer-bottom {
			flex-direction: column;
			gap: 0.5rem;
			text-align: center;
		}
		
		.scroll-hint {
			display: none;
		}
	}
</style>
