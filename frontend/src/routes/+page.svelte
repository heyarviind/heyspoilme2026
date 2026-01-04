<script lang="ts">
	import { Button, Input } from '$lib/components';

	let email = $state('');
	let gender = $state('');
	let isSubmitting = $state(false);
	let message = $state('');
	let isSuccess = $state(false);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		
		if (!email || !gender || isSubmitting) return;
		
		isSubmitting = true;
		message = '';
		
		try {
			const response = await fetch('/api/subscribe', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({ email, gender }),
			});
			
			const data = await response.json();
			
			if (data.success) {
				isSuccess = true;
				message = data.message;
				email = '';
				gender = '';
			} else {
				isSuccess = false;
				message = data.message || 'Something went wrong. Please try again.';
			}
		} catch (error) {
			isSuccess = false;
			message = 'Unable to connect. Please try again later.';
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>HeySpoilMe</title>
	<meta name="description" content="HeySpoilMe - Where luxury meets connection." />
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;500;600;700&family=Montserrat:wght@300;400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="page">
	<div class="image-panel"></div>
	
	<main class="content-panel">
		<div class="content">
			<img src="/img/logo.svg" alt="HeySpoilMe" class="logo" />
			
			<div class="badge">Coming Soon</div>
			
			<h2 class="title">A Private Space for Modern Sugar Dating in <img src="/img/indian-flag.svg" alt="India" class="flag" /> India</h2>
			
			<p class="description">
				A discreet, adult-only platform for mutually agreed sugar relationships 
				built on respect, transparency, and boundaries.
			</p>
			
			
			<ul class="features">
				<li>21+ only</li>
				<li>Consent-first & safety-focused</li>
                <li>Verified members</li>
				<li>No escorts. No transactions. No drama.</li>
			</ul>
			
			<form class="subscribe-form" onsubmit={handleSubmit}>
				<div class="gender-select">
					<button 
						type="button" 
						class="gender-option" 
						class:selected={gender === 'man'}
						onclick={() => gender = 'man'}
					>
						I am a Man
					</button>
					<button 
						type="button" 
						class="gender-option" 
						class:selected={gender === 'woman'}
						onclick={() => gender = 'woman'}
					>
						I am a Woman
					</button>
				</div>
				
				<div class="input-row">
					<Input 
						type="email" 
						placeholder="Enter your email"
						bind:value={email}
						disabled={isSubmitting}
						required
					/>
					<Button type="submit" variant="primary" loading={isSubmitting} disabled={!email || !gender}>
						Notify Me
					</Button>
				</div>
				
				{#if message}
					<p class="form-message" class:success={isSuccess} class:error={!isSuccess}>
						{message}
					</p>
				{/if}
			</form>
		</div>
		
		<footer class="footer">
			<p>2026 HeySpoilMe</p>
		</footer>
	</main>
</div>

<style>
	:global(*) {
		margin: 0;
		padding: 0;
		box-sizing: border-box;
	}
	
	:global(body) {
		font-family: 'Montserrat', sans-serif;
		background: #000;
		color: #fff;
		min-height: 100vh;
	}
	
	.page {
		display: flex;
		min-height: 100vh;
		background: #000;
	}
	
	.image-panel {
		position: fixed;
		left: 0;
		top: 0;
		bottom: 0;
		width: 30%;
		background-image: url('/img/background.png');
		background-size: cover;
		background-position: center;
		background-repeat: no-repeat;
	}
	
	.content-panel {
		margin-left: 30%;
		flex: 1;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 2rem;
		min-height: 100vh;
	}
	
	.content {
		text-align: left;
		width: 100%;
		max-width: 480px;
	}
	
	.logo {
		display: block;
		height: 5rem;
		margin-bottom: 1.5rem;
		
	}
	
	.badge {
		display: inline-block;
		vertical-align: top;
		background: #000;
		border: 1px solid #fff;
		padding: 0.4rem 1rem;
		font-size: 0.75rem;
		font-weight: 500;
		letter-spacing: 0.15em;
		text-transform: uppercase;
		color: #fff;
		margin-bottom: 1.5rem;
	}
	
	.title {
		font-family: 'Playfair Display', serif;
		font-size: 1.75rem;
		font-weight: 500;
		color: #fff;
		margin-bottom: 1.5rem;
		line-height: 1.3;
	}
	
	.flag {
		display: inline-block;
		height: 1em;
		vertical-align: middle;
		margin: 0 0.15em;
	}
	
	.description {
		font-size: 0.95rem;
		line-height: 1.7;
		color: rgba(255, 255, 255, 0.6);
		margin-bottom: 1rem;
	}
	
	.features {
		list-style: none;
		margin-bottom: 2rem;
	}
	
	.features li {
		font-size: 0.9rem;
		color: rgba(255, 255, 255, 0.7);
		margin-bottom: 0.5rem;
	}
	
	.features li::before {
		content: '';
		display: inline-block;
		width: 12px;
		height: 12px;
		margin-right: 0.5rem;
		background: #fff;
		mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='3'%3E%3Cpath d='M20 6L9 17l-5-5'/%3E%3C/svg%3E") center/contain no-repeat;
		-webkit-mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='3'%3E%3Cpath d='M20 6L9 17l-5-5'/%3E%3C/svg%3E") center/contain no-repeat;
	}
	
	.subscribe-form {
		width: 100%;
	}
	
	.gender-select {
		display: flex;
		gap: 0;
		margin-bottom: 1rem;
	}
	
	.gender-option {
		flex: 1;
		padding: 1rem;
		background: #000;
		border: 1px solid rgba(255, 255, 255, 0.3);
		color: rgba(255, 255, 255, 0.5);
		font-family: 'Montserrat', sans-serif;
		font-size: 0.9rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
	}
	
	.gender-option:first-child {
		border-right: none;
	}
	
	.gender-option:hover {
		color: #fff;
		border-color: rgba(255, 255, 255, 0.5);
	}
	
	.gender-option.selected {
		background: #fff;
		color: #000;
		border-color: #fff;
	}
	
	.input-row {
		display: flex;
		gap: 0;
	}
	
	.input-row :global(.input) {
		flex: 1;
		border-right: none;
	}
	
	.input-row :global(.btn) {
		flex-shrink: 0;
	}
	
	.form-message {
		margin-top: 1rem;
		font-size: 0.85rem;
	}
	
	.form-message.success {
		color: #4ade80;
	}
	
	.form-message.error {
		color: #f87171;
	}
	
	.footer {
		position: absolute;
		bottom: 0;
		padding: 2rem;
		text-align: center;
	}
	
	.footer p {
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.25);
		letter-spacing: 0.05em;
	}
	
	@media (max-width: 768px) {
		.image-panel {
			display: none;
		}
		
		.content-panel {
			margin-left: 0;
			padding: 2rem;
		}
		
		.content {
			text-align: center;
		}
		
		.logo {
			height: 3.5rem;
			margin-left: auto;
			margin-right: auto;
		}
		
		.title {
			font-size: 1.5rem;
		}
		
		.description {
			font-size: 0.9rem;
		}
		
		.gender-select {
			flex-direction: column;
		}
		
		.gender-option:first-child {
			border-right: 1px solid rgba(255, 255, 255, 0.3);
			border-bottom: none;
		}
		
		.gender-option.selected:first-child {
			border-right: 1px solid #fff;
		}
		
		.input-row {
			flex-direction: column;
		}
		
		.input-row :global(.input) {
			border-right: 1px solid #fff;
			border-bottom: none;
			text-align: center;
		}
		
		.input-row :global(.btn) {
			width: 100%;
		}
		
		.footer {
			position: relative;
		}
	}
</style>
