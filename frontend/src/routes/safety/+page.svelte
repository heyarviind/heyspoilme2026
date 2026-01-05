<script lang="ts">
	import { Footer } from '$lib/components';
</script>

<svelte:head>
	<title>Safety Guidelines — HeySpoilMe</title>
	<meta name="description" content="Safety guidelines for HeySpoilMe members. Learn how to stay safe while dating online and meeting people in person." />
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@300;400;500;600&family=Outfit:wght@300;400;500&display=swap" rel="stylesheet">
</svelte:head>

<div class="safety-page">
	<nav class="nav">
		<a href="/">
			<img src="/img/logo.svg" alt="HeySpoilMe" class="logo" />
		</a>
		<a href="/auth/login" class="sign-in-btn">Sign In</a>
	</nav>

	<main class="content">
		<header class="page-header">
			<span class="section-label">Your Safety Matters</span>
			<h1>Safety Guidelines</h1>
			<p class="subtitle">Essential tips to protect yourself while connecting on HeySpoilMe</p>
		</header>

		<div class="safety-intro">
			<div class="intro-icon">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
					<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
					<path d="m9 12 2 2 4-4"/>
				</svg>
			</div>
			<p>At HeySpoilMe, your safety is our priority. While we work hard to create a secure platform, online dating comes with inherent risks. This guide will help you navigate our platform safely and protect yourself when meeting people.</p>
		</div>

		<article class="safety-content">
			<section class="safety-section">
				<div class="section-header">
					<div class="section-icon">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
							<path d="M7 11V7a5 5 0 0 1 10 0v4"/>
						</svg>
					</div>
					<h2>Protect Your Personal Information</h2>
				</div>
				
				<div class="safety-tips">
					<div class="tip-card warning">
						<h3>Never Share Sensitive Information</h3>
						<p>Never share your home address, workplace address, financial details, or government ID numbers with someone you haven't met in person. Scammers may use this information to steal your identity or locate you.</p>
					</div>

					<div class="tip-card">
						<h3>Keep Conversations on the Platform</h3>
						<p>Until you're comfortable, keep your conversations within HeySpoilMe. Our messaging system helps us protect you—if someone asks to move to WhatsApp, Telegram, or other platforms immediately, proceed with caution.</p>
					</div>

					<div class="tip-card">
						<h3>Be Cautious with Photos</h3>
						<p>Avoid sharing intimate photos, especially with someone you haven't met. These can be used for blackmail or shared without your consent. If you do share personal photos, be aware of the risks involved.</p>
					</div>

					<div class="tip-card">
						<h3>Protect Your Social Media</h3>
						<p>Be careful about sharing your social media handles early on. Someone could use this to find more personal information about you, including your full name, workplace, and social connections.</p>
					</div>
				</div>
			</section>

			<section class="safety-section">
				<div class="section-header">
					<div class="section-icon">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<circle cx="12" cy="12" r="10"/>
							<line x1="12" y1="8" x2="12" y2="12"/>
							<line x1="12" y1="16" x2="12.01" y2="16"/>
						</svg>
					</div>
					<h2>Recognize Warning Signs</h2>
				</div>
				
				<div class="warning-grid">
					<div class="warning-item">
						<span class="warning-icon">💸</span>
						<h3>Asks for Money</h3>
						<p>Never send money to someone you haven't met in person, no matter how compelling their story is</p>
					</div>
					<div class="warning-item">
						<span class="warning-icon">⚡</span>
						<h3>Moves Too Fast</h3>
						<p>Professing love or wanting commitment very quickly can be a manipulation tactic</p>
					</div>
					<div class="warning-item">
						<span class="warning-icon">📷</span>
						<h3>Refuses Video Calls</h3>
						<p>Consistently avoiding video calls may indicate they're not who they claim to be</p>
					</div>
					<div class="warning-item">
						<span class="warning-icon">📍</span>
						<h3>Inconsistent Stories</h3>
						<p>Details about their life, job, or location that keep changing are red flags</p>
					</div>
					<div class="warning-item">
						<span class="warning-icon">🔒</span>
						<h3>Requests Private Info</h3>
						<p>Asking for passwords, bank details, or other sensitive information early on</p>
					</div>
					<div class="warning-item">
						<span class="warning-icon">😤</span>
						<h3>Pressures or Threatens</h3>
						<p>Any form of pressure, manipulation, or threats is unacceptable behavior</p>
					</div>
				</div>
			</section>

			<section class="safety-section">
				<div class="section-header">
					<div class="section-icon">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
							<circle cx="9" cy="7" r="4"/>
							<path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
							<path d="M16 3.13a4 4 0 0 1 0 7.75"/>
						</svg>
					</div>
					<h2>Meeting in Person</h2>
				</div>
				
				<div class="meeting-steps">
					<div class="step">
						<div class="step-number">1</div>
						<div class="step-content">
							<h3>Video Call First</h3>
							<p>Before meeting in person, have at least one video call to verify they are who they claim to be. This simple step can prevent many unpleasant surprises.</p>
						</div>
					</div>
					<div class="step">
						<div class="step-number">2</div>
						<div class="step-content">
							<h3>Meet in Public</h3>
							<p>Always choose a public place for your first several meetings—a popular restaurant, café, or bar. Avoid private residences, secluded areas, or being picked up from your home.</p>
						</div>
					</div>
					<div class="step">
						<div class="step-number">3</div>
						<div class="step-content">
							<h3>Tell Someone</h3>
							<p>Inform a friend or family member about your date—where you're going, who you're meeting, and when you expect to be back. Consider sharing your live location with them.</p>
						</div>
					</div>
					<div class="step">
						<div class="step-number">4</div>
						<div class="step-content">
							<h3>Arrange Your Own Transport</h3>
							<p>Drive yourself, book your own cab, or use public transport. Never get into your date's car for the first few meetings. Maintain your independence and ability to leave.</p>
						</div>
					</div>
					<div class="step">
						<div class="step-number">5</div>
						<div class="step-content">
							<h3>Stay Sober and Alert</h3>
							<p>Keep your alcohol consumption minimal on first dates. Never leave your drink unattended, and if you feel unusually intoxicated, seek help immediately.</p>
						</div>
					</div>
					<div class="step">
						<div class="step-number">6</div>
						<div class="step-content">
							<h3>Trust Your Instincts</h3>
							<p>If something feels wrong, leave. You don't owe anyone an explanation. Your safety is more important than politeness. Have an exit plan ready.</p>
						</div>
					</div>
				</div>
			</section>

			<section class="safety-section">
				<div class="section-header">
					<div class="section-icon">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z"/>
							<path d="M12 6v6l4 2"/>
						</svg>
					</div>
					<h2>Financial Safety</h2>
				</div>
				
				<div class="safety-tips">
					<div class="tip-card warning">
						<h3>Never Send Money</h3>
						<p>This cannot be stressed enough: never send money, gift cards, or cryptocurrency to someone you haven't met in person, regardless of the reason. Romance scammers are extremely convincing and often build emotional connections over weeks or months before asking for money.</p>
					</div>

					<div class="tip-card">
						<h3>Don't Accept Unusual Financial Requests</h3>
						<p>If someone asks you to receive money on their behalf, cash checks, or transfer funds—even temporarily—this is likely a scam that could make you legally liable.</p>
					</div>

					<div class="tip-card">
						<h3>Verify Before Investing</h3>
						<p>If someone recommends investments, cryptocurrency opportunities, or business ventures, be extremely cautious. "Pig butchering" scams involve building romantic relationships to lure victims into fraudulent investments.</p>
					</div>

					<div class="tip-card">
						<h3>Keep Financial Arrangements Transparent</h3>
						<p>If you're in a mutually beneficial arrangement, establish clear terms upfront. Never feel pressured to provide financial support before meeting or beyond your comfort level.</p>
					</div>
				</div>
			</section>

			<section class="safety-section">
				<div class="section-header">
					<div class="section-icon">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<path d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z"/>
							<line x1="4" y1="22" x2="4" y2="15"/>
						</svg>
					</div>
					<h2>Reporting & Getting Help</h2>
				</div>
				
				<div class="report-info">
					<div class="report-card">
						<h3>Report Within the App</h3>
						<p>If someone violates our community guidelines, behaves inappropriately, or makes you uncomfortable, report them immediately. Go to their profile, tap the three-dot menu, and select "Report." Our safety team reviews all reports within 24 hours.</p>
						<ul>
							<li>Harassment or threatening behavior</li>
							<li>Spam or scam attempts</li>
							<li>Fake profiles or impersonation</li>
							<li>Inappropriate content</li>
							<li>Any behavior that violates our terms</li>
						</ul>
					</div>

					<div class="report-card">
						<h3>Block Users</h3>
						<p>You can block any user at any time. Blocked users cannot see your profile, send you messages, or interact with you. Don't hesitate to use this feature if someone makes you uncomfortable.</p>
					</div>

					<div class="report-card emergency">
						<h3>Emergency Situations</h3>
						<p>If you are in immediate danger, contact local emergency services immediately:</p>
						<div class="emergency-numbers">
							<div class="emergency-item">
								<span class="number">112</span>
								<span class="label">Emergency (India)</span>
							</div>
							<div class="emergency-item">
								<span class="number">181</span>
								<span class="label">Women Helpline</span>
							</div>
							<div class="emergency-item">
								<span class="number">1091</span>
								<span class="label">Women in Distress</span>
							</div>
						</div>
					</div>
				</div>
			</section>

			<section class="safety-section">
				<div class="section-header">
					<div class="section-icon">
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
							<circle cx="12" cy="12" r="10"/>
							<path d="M12 16v-4"/>
							<path d="M12 8h.01"/>
						</svg>
					</div>
					<h2>Consent & Boundaries</h2>
				</div>
				
				<div class="consent-content">
					<div class="consent-principle">
						<h3>Respect is Non-Negotiable</h3>
						<p>Every interaction on HeySpoilMe should be respectful. Harassment, unsolicited explicit content, and pressure of any kind are strictly prohibited. Everyone has the right to set their own boundaries.</p>
					</div>

					<div class="consent-points">
						<div class="consent-item">
							<span class="check-icon">✓</span>
							<p><strong>Consent is ongoing</strong> — It can be withdrawn at any time. A "yes" yesterday doesn't mean "yes" today.</p>
						</div>
						<div class="consent-item">
							<span class="check-icon">✓</span>
							<p><strong>"No" means "No"</strong> — Respect when someone declines, whether it's a date, a request for photos, or anything else.</p>
						</div>
						<div class="consent-item">
							<span class="check-icon">✓</span>
							<p><strong>Clear communication</strong> — When in doubt, ask. Don't assume someone's boundaries or comfort level.</p>
						</div>
						<div class="consent-item">
							<span class="check-icon">✓</span>
							<p><strong>No means no</strong> — Silence or lack of resistance is not consent. Enthusiastic, clear agreement is the standard.</p>
						</div>
					</div>
				</div>
			</section>

			<section class="commitment-section">
				<h2>Our Commitment to Your Safety</h2>
				<div class="commitment-grid">
					<div class="commitment-item">
						<h3>Verification</h3>
						<p>We verify members to ensure authenticity and reduce fake profiles</p>
					</div>
					<div class="commitment-item">
						<h3>24/7 Moderation</h3>
						<p>Our team reviews reports and monitors for suspicious activity around the clock</p>
					</div>
					<div class="commitment-item">
						<h3>Secure Platform</h3>
						<p>Industry-standard encryption protects your data and conversations</p>
					</div>
					<div class="commitment-item">
						<h3>Women Message First</h3>
						<p>Our policy reduces unwanted messages and empowers women to control their experience</p>
					</div>
				</div>
			</section>

			<section class="final-note">
				<p>Remember: Your safety is ultimately in your hands. While we do everything we can to create a safe platform, online dating requires vigilance. Trust your instincts, take your time, and never let anyone pressure you into anything you're not comfortable with.</p>
				<p>If you have safety concerns or suggestions, please reach out to us at <strong>safety@heyspoilme.com</strong>.</p>
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
	
	.safety-page {
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

	/* Safety Intro */
	.safety-intro {
		display: flex;
		align-items: flex-start;
		gap: 1.5rem;
		padding: 2rem;
		background: linear-gradient(135deg, rgba(76, 175, 80, 0.1) 0%, rgba(76, 175, 80, 0.03) 100%);
		border: 1px solid rgba(76, 175, 80, 0.2);
		margin-bottom: 4rem;
	}

	.intro-icon {
		flex-shrink: 0;
		width: 3rem;
		height: 3rem;
	}

	.intro-icon svg {
		width: 100%;
		height: 100%;
		color: #4CAF50;
	}

	.safety-intro p {
		font-size: 1rem;
		line-height: 1.8;
		color: rgba(255, 255, 255, 0.8);
	}

	/* Safety Sections */
	.safety-section {
		margin-bottom: 4rem;
	}

	.section-header {
		display: flex;
		align-items: center;
		gap: 1rem;
		margin-bottom: 2rem;
		padding-bottom: 1rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}

	.section-icon {
		width: 2.5rem;
		height: 2.5rem;
		flex-shrink: 0;
	}

	.section-icon svg {
		width: 100%;
		height: 100%;
		color: rgba(255, 255, 255, 0.6);
	}

	.section-header h2 {
		font-family: 'Cormorant Garamond', serif;
		font-size: 1.75rem;
		font-weight: 500;
	}

	/* Tip Cards */
	.safety-tips {
		display: grid;
		gap: 1.5rem;
	}

	.tip-card {
		padding: 1.5rem;
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid rgba(255, 255, 255, 0.08);
	}

	.tip-card.warning {
		background: rgba(255, 193, 7, 0.05);
		border-color: rgba(255, 193, 7, 0.2);
	}

	.tip-card h3 {
		font-family: 'Outfit', sans-serif;
		font-size: 1.05rem;
		font-weight: 500;
		margin-bottom: 0.75rem;
	}

	.tip-card p {
		font-size: 0.95rem;
		line-height: 1.8;
		color: rgba(255, 255, 255, 0.7);
	}

	/* Warning Grid */
	.warning-grid {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 1.5rem;
	}

	.warning-item {
		padding: 1.5rem;
		background: rgba(244, 67, 54, 0.05);
		border: 1px solid rgba(244, 67, 54, 0.15);
		text-align: center;
	}

	.warning-icon {
		font-size: 2rem;
		display: block;
		margin-bottom: 1rem;
	}

	.warning-item h3 {
		font-size: 0.95rem;
		font-weight: 500;
		margin-bottom: 0.5rem;
	}

	.warning-item p {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.6);
		line-height: 1.6;
	}

	/* Meeting Steps */
	.meeting-steps {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.step {
		display: flex;
		gap: 1.5rem;
		padding: 1.5rem;
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid rgba(255, 255, 255, 0.08);
	}

	.step-number {
		width: 2.5rem;
		height: 2.5rem;
		background: rgba(255, 255, 255, 0.1);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.1rem;
		font-weight: 500;
		flex-shrink: 0;
	}

	.step-content h3 {
		font-size: 1.05rem;
		font-weight: 500;
		margin-bottom: 0.5rem;
	}

	.step-content p {
		font-size: 0.95rem;
		line-height: 1.7;
		color: rgba(255, 255, 255, 0.7);
	}

	/* Report Info */
	.report-info {
		display: grid;
		gap: 1.5rem;
	}

	.report-card {
		padding: 1.5rem;
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid rgba(255, 255, 255, 0.08);
	}

	.report-card.emergency {
		background: rgba(244, 67, 54, 0.05);
		border-color: rgba(244, 67, 54, 0.2);
	}

	.report-card h3 {
		font-size: 1.1rem;
		font-weight: 500;
		margin-bottom: 1rem;
	}

	.report-card p {
		font-size: 0.95rem;
		line-height: 1.7;
		color: rgba(255, 255, 255, 0.7);
		margin-bottom: 1rem;
	}

	.report-card ul {
		margin: 1rem 0 0 1.5rem;
	}

	.report-card li {
		font-size: 0.9rem;
		color: rgba(255, 255, 255, 0.6);
		margin-bottom: 0.5rem;
	}

	.emergency-numbers {
		display: flex;
		gap: 2rem;
		margin-top: 1rem;
	}

	.emergency-item {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.emergency-item .number {
		font-size: 1.5rem;
		font-weight: 600;
		color: #F44336;
	}

	.emergency-item .label {
		font-size: 0.8rem;
		color: rgba(255, 255, 255, 0.5);
		margin-top: 0.25rem;
	}

	/* Consent */
	.consent-content {
		padding: 2rem;
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid rgba(255, 255, 255, 0.08);
	}

	.consent-principle {
		margin-bottom: 2rem;
		padding-bottom: 2rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
	}

	.consent-principle h3 {
		font-size: 1.1rem;
		font-weight: 500;
		margin-bottom: 0.75rem;
	}

	.consent-principle p {
		font-size: 0.95rem;
		line-height: 1.8;
		color: rgba(255, 255, 255, 0.7);
	}

	.consent-points {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.consent-item {
		display: flex;
		gap: 1rem;
		align-items: flex-start;
	}

	.check-icon {
		width: 1.5rem;
		height: 1.5rem;
		background: rgba(76, 175, 80, 0.2);
		color: #4CAF50;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 0.85rem;
		flex-shrink: 0;
	}

	.consent-item p {
		font-size: 0.95rem;
		line-height: 1.7;
		color: rgba(255, 255, 255, 0.7);
	}

	.consent-item strong {
		color: #fff;
		font-weight: 500;
	}

	/* Commitment Section */
	.commitment-section {
		padding: 3rem;
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid rgba(255, 255, 255, 0.08);
		margin-bottom: 3rem;
	}

	.commitment-section h2 {
		font-family: 'Cormorant Garamond', serif;
		font-size: 1.75rem;
		font-weight: 500;
		text-align: center;
		margin-bottom: 2rem;
	}

	.commitment-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 2rem;
	}

	.commitment-item h3 {
		font-size: 1rem;
		font-weight: 500;
		margin-bottom: 0.5rem;
	}

	.commitment-item p {
		font-size: 0.9rem;
		color: rgba(255, 255, 255, 0.6);
		line-height: 1.6;
	}

	/* Final Note */
	.final-note {
		padding: 2rem;
		background: linear-gradient(135deg, rgba(255, 255, 255, 0.03) 0%, rgba(255, 255, 255, 0.01) 100%);
		border: 1px solid rgba(255, 255, 255, 0.1);
		text-align: center;
	}

	.final-note p {
		font-size: 0.95rem;
		line-height: 1.8;
		color: rgba(255, 255, 255, 0.7);
		margin-bottom: 1rem;
	}

	.final-note p:last-child {
		margin-bottom: 0;
	}

	.final-note strong {
		color: #fff;
		font-weight: 500;
	}

	/* Mobile */
	@media (max-width: 768px) {
		.nav {
			padding: 1.25rem 1.5rem;
		}

		.content {
			padding: 7rem 1.5rem 3rem;
		}

		.safety-intro {
			flex-direction: column;
			text-align: center;
		}

		.intro-icon {
			margin: 0 auto;
		}

		.warning-grid {
			grid-template-columns: 1fr;
		}

		.step {
			flex-direction: column;
			gap: 1rem;
		}

		.emergency-numbers {
			flex-direction: column;
			gap: 1rem;
		}

		.commitment-grid {
			grid-template-columns: 1fr;
		}

		.commitment-section {
			padding: 2rem 1.5rem;
		}

		.section-header {
			flex-direction: column;
			text-align: center;
		}

		.section-header h2 {
			font-size: 1.5rem;
		}
	}
</style>

