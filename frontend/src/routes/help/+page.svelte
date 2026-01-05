<script lang="ts">
	import { Footer } from '$lib/components';
</script>

<svelte:head>
	<title>Help Center — HeySpoilMe</title>
	<meta name="description" content="Get help with HeySpoilMe. Find answers to frequently asked questions, account support, and contact our team." />
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@300;400;500;600&family=Outfit:wght@300;400;500&display=swap" rel="stylesheet">
</svelte:head>

<div class="help-page">
	<nav class="nav">
		<a href="/">
			<img src="/img/logo.svg" alt="HeySpoilMe" class="logo" />
		</a>
		<a href="/auth/login" class="sign-in-btn">Sign In</a>
	</nav>

	<main class="content">
		<header class="page-header">
			<span class="section-label">Support</span>
			<h1>Help Center</h1>
			<p class="subtitle">Find answers to your questions and get the support you need</p>
		</header>

		<div class="search-section">
			<div class="search-box">
				<svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<circle cx="11" cy="11" r="8"/>
					<path d="m21 21-4.35-4.35"/>
				</svg>
				<input type="text" placeholder="Search for help..." class="search-input" />
			</div>
		</div>

		<section class="quick-links">
			<h2>Quick Help</h2>
			<div class="links-grid">
				<a href="#getting-started" class="link-card">
					<div class="icon-wrap">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M12 2L2 7l10 5 10-5-10-5z"/>
							<path d="M2 17l10 5 10-5"/>
							<path d="M2 12l10 5 10-5"/>
						</svg>
					</div>
					<h3>Getting Started</h3>
					<p>New to HeySpoilMe? Learn the basics</p>
				</a>
				<a href="#account" class="link-card">
					<div class="icon-wrap">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<circle cx="12" cy="8" r="4"/>
							<path d="M6 21v-2a4 4 0 0 1 4-4h4a4 4 0 0 1 4 4v2"/>
						</svg>
					</div>
					<h3>Account & Profile</h3>
					<p>Manage your account settings</p>
				</a>
				<a href="#verification" class="link-card">
					<div class="icon-wrap">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
							<path d="m9 12 2 2 4-4"/>
						</svg>
					</div>
					<h3>Verification</h3>
					<p>Get verified and build trust</p>
				</a>
				<a href="#messaging" class="link-card">
					<div class="icon-wrap">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
						</svg>
					</div>
					<h3>Messaging</h3>
					<p>Connect with other members</p>
				</a>
				<a href="#safety" class="link-card">
					<div class="icon-wrap">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
						</svg>
					</div>
					<h3>Safety & Privacy</h3>
					<p>Stay safe on our platform</p>
				</a>
				<a href="#billing" class="link-card">
					<div class="icon-wrap">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<rect x="1" y="4" width="22" height="16" rx="2"/>
							<line x1="1" y1="10" x2="23" y2="10"/>
						</svg>
					</div>
					<h3>Billing & Payments</h3>
					<p>Manage subscriptions and payments</p>
				</a>
			</div>
		</section>

		<article class="faq-content">
			<section id="getting-started">
				<h2>Getting Started</h2>
				
				<div class="faq-item">
					<h3>What is HeySpoilMe?</h3>
					<p>HeySpoilMe is India's premier dating platform designed for discerning adults seeking meaningful connections. We connect successful individuals with those who appreciate the finer things in life, fostering relationships based on mutual respect and genuine interest.</p>
				</div>

				<div class="faq-item">
					<h3>How do I create an account?</h3>
					<p>Creating an account is simple:</p>
					<ol>
						<li>Click "Sign Up" or "Get Started" on our homepage</li>
						<li>Sign in with your Google account or email</li>
						<li>Complete your profile with photos and information about yourself</li>
						<li>Verify your identity to unlock all features</li>
					</ol>
					<p>The entire process takes just a few minutes.</p>
				</div>

				<div class="faq-item">
					<h3>What are the requirements to join?</h3>
					<p>To join HeySpoilMe, you must:</p>
					<ul>
						<li>Be at least 21 years old</li>
						<li>Have a valid email address</li>
						<li>Agree to our Terms of Service and Community Guidelines</li>
						<li>Complete our verification process</li>
					</ul>
				</div>

				<div class="faq-item">
					<h3>How does HeySpoilMe work?</h3>
					<p>Once you've created and verified your profile, you can browse other members, express interest by "liking" profiles, and start conversations when there's mutual interest. On HeySpoilMe, women have the privilege of initiating conversations first, creating a respectful environment for all members.</p>
				</div>
			</section>

			<section id="account">
				<h2>Account & Profile</h2>
				
				<div class="faq-item">
					<h3>How do I edit my profile?</h3>
					<p>To edit your profile:</p>
					<ol>
						<li>Log in to your account</li>
						<li>Click on your profile picture or navigate to "Profile"</li>
						<li>Click "Edit Profile" to update your photos, bio, and preferences</li>
						<li>Save your changes when done</li>
					</ol>
				</div>

				<div class="faq-item">
					<h3>How do I change my password?</h3>
					<p>If you signed up with email, you can reset your password from the login page by clicking "Forgot Password." If you signed up with Google, your password is managed through your Google account.</p>
				</div>

				<div class="faq-item">
					<h3>How do I delete my account?</h3>
					<p>We're sorry to see you go! To delete your account:</p>
					<ol>
						<li>Go to your Profile settings</li>
						<li>Scroll to the bottom and find "Delete Account"</li>
						<li>Confirm your decision</li>
					</ol>
					<p>Please note that account deletion is permanent and cannot be undone. Your profile, messages, and data will be removed within 30 days.</p>
				</div>

				<div class="faq-item">
					<h3>What photos should I upload?</h3>
					<p>Great photos make a difference! We recommend:</p>
					<ul>
						<li><strong>Primary photo:</strong> A clear, recent photo of your face</li>
						<li><strong>Additional photos:</strong> Show your personality, hobbies, and lifestyle</li>
						<li><strong>Quality:</strong> Use well-lit, high-quality images</li>
						<li><strong>Authenticity:</strong> Be yourself – no heavy filters or misleading photos</li>
					</ul>
					<p>Photos must comply with our content guidelines. No nudity is allowed in public photos.</p>
				</div>
			</section>

			<section id="verification">
				<h2>Verification</h2>
				
				<div class="faq-item">
					<h3>Why should I verify my profile?</h3>
					<p>Verification builds trust and shows other members that you're genuine. Verified profiles receive:</p>
					<ul>
						<li>A verified badge on your profile</li>
						<li>Higher visibility in search results</li>
						<li>More trust from potential matches</li>
						<li>Access to all platform features</li>
					</ul>
				</div>

				<div class="faq-item">
					<h3>How do I get verified?</h3>
					<p>Our verification process includes:</p>
					<ol>
						<li><strong>Photo Verification:</strong> Take a selfie matching a specific pose to prove it's really you</li>
						<li><strong>Age Verification:</strong> Submit a government-issued ID (your personal details are kept confidential)</li>
					</ol>
					<p>Verification typically takes 24-48 hours. You'll receive a notification once approved.</p>
				</div>

				<div class="faq-item">
					<h3>Is my verification information secure?</h3>
					<p>Absolutely. Your verification documents are encrypted and handled with the highest security standards. ID photos are used solely for verification purposes and are deleted after the process is complete. We never share your ID information with other members.</p>
				</div>

				<div class="faq-item">
					<h3>What if my verification is rejected?</h3>
					<p>If your verification is rejected, you'll receive an email explaining why. Common reasons include:</p>
					<ul>
						<li>Blurry or unclear photos</li>
						<li>Face not clearly visible</li>
						<li>Pose doesn't match the required pose</li>
						<li>ID document is expired or unreadable</li>
					</ul>
					<p>You can resubmit your verification with corrected materials.</p>
				</div>
			</section>

			<section id="messaging">
				<h2>Messaging</h2>
				
				<div class="faq-item">
					<h3>How do I start a conversation?</h3>
					<p>On HeySpoilMe, women initiate conversations first. If you're a woman, you can message any member you're interested in. If you're a man, you can "like" profiles and wait for women to reach out. Once a woman sends the first message, both parties can chat freely.</p>
				</div>

				<div class="faq-item">
					<h3>Why can't I send the first message?</h3>
					<p>Our "women message first" policy is designed to create a respectful environment where women feel comfortable and in control. This reduces unwanted messages and ensures more meaningful connections.</p>
				</div>

				<div class="faq-item">
					<h3>How do I know if someone is online?</h3>
					<p>Members who are currently online will have a green indicator on their profile. You can also see when a member was last active on their profile page.</p>
				</div>

				<div class="faq-item">
					<h3>Can I video call other members?</h3>
					<p>Video calling features are available for verified members. You can request a video call within a conversation once both parties are comfortable.</p>
				</div>
			</section>

			<section id="safety">
				<h2>Safety & Privacy</h2>
				
				<div class="faq-item">
					<h3>How do I report someone?</h3>
					<p>If you encounter inappropriate behavior, you can report a member by:</p>
					<ol>
						<li>Going to their profile</li>
						<li>Clicking the three-dot menu (⋮)</li>
						<li>Selecting "Report"</li>
						<li>Choosing the reason and providing details</li>
					</ol>
					<p>Our team reviews all reports within 24 hours and takes appropriate action.</p>
				</div>

				<div class="faq-item">
					<h3>How do I block someone?</h3>
					<p>To block a member:</p>
					<ol>
						<li>Go to their profile or your conversation with them</li>
						<li>Click the three-dot menu (⋮)</li>
						<li>Select "Block"</li>
					</ol>
					<p>Blocked members cannot see your profile, send you messages, or interact with you in any way.</p>
				</div>

				<div class="faq-item">
					<h3>How do I hide my profile?</h3>
					<p>You can temporarily hide your profile from browse results in your privacy settings. This allows you to take a break without deleting your account.</p>
				</div>

				<div class="faq-item">
					<h3>Is my personal information safe?</h3>
					<p>Yes. We use industry-standard encryption and security measures to protect your data. Your personal information is never shared with other members. For more details, please read our <a href="/privacy">Privacy Policy</a>.</p>
				</div>
			</section>

			<section id="billing">
				<h2>Billing & Payments</h2>
				
				<div class="faq-item">
					<h3>How do I upgrade to premium?</h3>
					<p>To upgrade your account:</p>
					<ol>
						<li>Go to your Profile or Settings</li>
						<li>Select "Upgrade to Premium"</li>
						<li>Choose your preferred plan</li>
						<li>Complete the payment process</li>
					</ol>
				</div>

				<div class="faq-item">
					<h3>What payment methods do you accept?</h3>
					<p>We accept all major payment methods including:</p>
					<ul>
						<li>Credit/Debit Cards (Visa, Mastercard, RuPay)</li>
						<li>UPI</li>
						<li>Net Banking</li>
						<li>Digital Wallets</li>
					</ul>
				</div>

				<div class="faq-item">
					<h3>How do I cancel my subscription?</h3>
					<p>To cancel your subscription:</p>
					<ol>
						<li>Go to your Profile settings</li>
						<li>Navigate to "Subscription"</li>
						<li>Click "Cancel Subscription"</li>
					</ol>
					<p>Your premium features will remain active until the end of your current billing period.</p>
				</div>

				<div class="faq-item">
					<h3>Can I get a refund?</h3>
					<p>Our general policy is that payments are non-refundable. However, we review refund requests on a case-by-case basis. If you believe you're entitled to a refund, please contact our support team with your order details.</p>
				</div>
			</section>

			<section class="contact-section">
				<h2>Still Need Help?</h2>
				<p>Our support team is here to assist you. We typically respond within 24 hours.</p>
				
				<div class="contact-options">
					<div class="contact-card">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
							<polyline points="22,6 12,13 2,6"/>
						</svg>
						<h3>Email Support</h3>
						<p>support@heyspoilme.com</p>
					</div>
					<div class="contact-card">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
						</svg>
						<h3>In-App Support</h3>
						<p>Use the help button in the app</p>
					</div>
				</div>
			</section>
		</article>
	</main>

	<Footer />
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
	
	.help-page {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
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
		background: rgba(8, 8, 10, 0.95);
		backdrop-filter: blur(10px);
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
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

	/* Content */
	.content {
		flex: 1;
		padding: 8rem 4rem 4rem;
		max-width: 1000px;
		margin: 0 auto;
		width: 100%;
	}

	.page-header {
		margin-bottom: 3rem;
		text-align: center;
	}

	.section-label {
		display: block;
		font-size: 0.7rem;
		font-weight: 500;
		letter-spacing: 0.2em;
		text-transform: uppercase;
		color: rgba(255, 255, 255, 0.4);
		margin-bottom: 1rem;
	}

	.page-header h1 {
		font-family: 'Cormorant Garamond', serif;
		font-size: clamp(2.5rem, 5vw, 3.5rem);
		font-weight: 400;
		margin-bottom: 1rem;
	}

	.subtitle {
		font-size: 1rem;
		color: rgba(255, 255, 255, 0.5);
	}

	/* Search */
	.search-section {
		margin-bottom: 4rem;
	}

	.search-box {
		display: flex;
		align-items: center;
		gap: 1rem;
		max-width: 500px;
		margin: 0 auto;
		padding: 1rem 1.5rem;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 50px;
		transition: border-color 0.3s ease;
	}

	.search-box:focus-within {
		border-color: rgba(255, 255, 255, 0.3);
	}

	.search-icon {
		width: 1.25rem;
		height: 1.25rem;
		color: rgba(255, 255, 255, 0.4);
		flex-shrink: 0;
	}

	.search-input {
		flex: 1;
		background: transparent;
		border: none;
		outline: none;
		color: #fff;
		font-family: 'Outfit', sans-serif;
		font-size: 0.95rem;
	}

	.search-input::placeholder {
		color: rgba(255, 255, 255, 0.4);
	}

	/* Quick Links */
	.quick-links {
		margin-bottom: 4rem;
	}

	.quick-links h2 {
		font-family: 'Cormorant Garamond', serif;
		font-size: 1.5rem;
		font-weight: 500;
		margin-bottom: 1.5rem;
		text-align: center;
	}

	.links-grid {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 1.5rem;
	}

	.link-card {
		padding: 1.5rem;
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid rgba(255, 255, 255, 0.08);
		text-decoration: none;
		color: inherit;
		transition: all 0.3s ease;
	}

	.link-card:hover {
		background: rgba(255, 255, 255, 0.05);
		border-color: rgba(255, 255, 255, 0.15);
		transform: translateY(-2px);
	}

	.icon-wrap {
		width: 2.5rem;
		height: 2.5rem;
		margin-bottom: 1rem;
	}

	.icon-wrap svg {
		width: 100%;
		height: 100%;
		color: rgba(255, 255, 255, 0.6);
	}

	.link-card h3 {
		font-family: 'Outfit', sans-serif;
		font-size: 1rem;
		font-weight: 500;
		margin-bottom: 0.5rem;
	}

	.link-card p {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.5);
		line-height: 1.5;
	}

	/* FAQ Content */
	.faq-content section {
		margin-bottom: 4rem;
		padding-top: 2rem;
		border-top: 1px solid rgba(255, 255, 255, 0.05);
	}

	.faq-content section:first-child {
		border-top: none;
	}

	.faq-content h2 {
		font-family: 'Cormorant Garamond', serif;
		font-size: 1.75rem;
		font-weight: 500;
		margin-bottom: 2rem;
	}

	.faq-item {
		margin-bottom: 2rem;
		padding-bottom: 2rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
	}

	.faq-item:last-child {
		border-bottom: none;
		margin-bottom: 0;
		padding-bottom: 0;
	}

	.faq-item h3 {
		font-family: 'Outfit', sans-serif;
		font-size: 1.1rem;
		font-weight: 500;
		margin-bottom: 1rem;
		color: #fff;
	}

	.faq-item p {
		font-size: 0.95rem;
		line-height: 1.8;
		color: rgba(255, 255, 255, 0.7);
		margin-bottom: 1rem;
	}

	.faq-item p:last-child {
		margin-bottom: 0;
	}

	.faq-item ul,
	.faq-item ol {
		margin: 1rem 0 1rem 1.5rem;
	}

	.faq-item li {
		font-size: 0.95rem;
		line-height: 1.8;
		color: rgba(255, 255, 255, 0.7);
		margin-bottom: 0.5rem;
	}

	.faq-item a {
		color: #fff;
		text-decoration: underline;
		text-underline-offset: 2px;
	}

	.faq-item a:hover {
		color: rgba(255, 255, 255, 0.7);
	}

	.faq-item strong {
		color: #fff;
		font-weight: 500;
	}

	/* Contact Section */
	.contact-section {
		text-align: center;
		padding: 3rem;
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid rgba(255, 255, 255, 0.08);
	}

	.contact-section h2 {
		font-family: 'Cormorant Garamond', serif;
		font-size: 1.75rem;
		font-weight: 500;
		margin-bottom: 1rem;
	}

	.contact-section > p {
		font-size: 0.95rem;
		color: rgba(255, 255, 255, 0.6);
		margin-bottom: 2rem;
	}

	.contact-options {
		display: flex;
		justify-content: center;
		gap: 2rem;
	}

	.contact-card {
		padding: 1.5rem 2rem;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.contact-card svg {
		width: 2rem;
		height: 2rem;
		color: rgba(255, 255, 255, 0.6);
		margin-bottom: 1rem;
	}

	.contact-card h3 {
		font-size: 0.95rem;
		font-weight: 500;
		margin-bottom: 0.5rem;
	}

	.contact-card p {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.5);
	}

	/* Mobile */
	@media (max-width: 768px) {
		.nav {
			padding: 1.25rem 1.5rem;
		}

		.content {
			padding: 7rem 1.5rem 3rem;
		}

		.links-grid {
			grid-template-columns: 1fr;
			gap: 1rem;
		}

		.faq-content section {
			margin-bottom: 3rem;
		}

		.faq-content h2 {
			font-size: 1.5rem;
		}

		.contact-section {
			padding: 2rem 1.5rem;
		}

		.contact-options {
			flex-direction: column;
			gap: 1rem;
		}
	}
</style>


