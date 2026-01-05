<script lang="ts">
	import { api } from '$lib/api';
	import { goto } from '$app/navigation';

	let mode: 'signin' | 'signup' = $state('signin');
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleSubmit() {
		error = '';
		
		if (!email || !password) {
			error = 'Please fill in all fields';
			return;
		}

		if (mode === 'signup') {
			if (password.length < 8) {
				error = 'Password must be at least 8 characters';
				return;
			}
			if (password !== confirmPassword) {
				error = 'Passwords do not match';
				return;
			}
		}

		loading = true;
		try {
			const response: any = mode === 'signup' 
				? await api.signup(email, password)
				: await api.signin(email, password);
			
			localStorage.setItem('token', response.token);
			localStorage.setItem('user', JSON.stringify(response.user));
			
			if (response.is_new || !response.profile) {
				goto('/profile/setup');
			} else {
				goto('/browse');
			}
		} catch (e: any) {
			error = e.message || 'Authentication failed';
		} finally {
			loading = false;
		}
	}

	function toggleMode() {
		mode = mode === 'signin' ? 'signup' : 'signin';
		error = '';
	}
</script>

<svelte:head>
	<title>{mode === 'signin' ? 'Sign In' : 'Sign Up'} | HeySpoilMe</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;500;600;700&family=Montserrat:wght@300;400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="login-page">
	<div class="login-card">
		<img src="/img/logo.svg" alt="HeySpoilMe" class="logo" />
		
		<h1>{mode === 'signin' ? 'Welcome Back' : 'Create Account'}</h1>
		<p class="subtitle">
			{mode === 'signin' ? 'Sign in to connect with amazing people' : 'Join our exclusive community'}
		</p>

		<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
			{#if error}
				<div class="error-msg">{error}</div>
			{/if}

			<div class="input-group">
				<label for="email">Email</label>
				<input 
					type="email" 
					id="email"
					bind:value={email} 
					placeholder="you@example.com"
					autocomplete="email"
				/>
			</div>

			<div class="input-group">
				<label for="password">Password</label>
				<input 
					type="password" 
					id="password"
					bind:value={password} 
					placeholder="••••••••"
					autocomplete={mode === 'signin' ? 'current-password' : 'new-password'}
				/>
			</div>

			{#if mode === 'signup'}
				<div class="input-group">
					<label for="confirmPassword">Confirm Password</label>
					<input 
						type="password" 
						id="confirmPassword"
						bind:value={confirmPassword} 
						placeholder="••••••••"
						autocomplete="new-password"
					/>
				</div>
			{/if}

			<button type="submit" class="submit-btn" disabled={loading}>
				{#if loading}
					<span class="spinner"></span>
				{:else}
					{mode === 'signin' ? 'Sign In' : 'Create Account'}
				{/if}
			</button>
		</form>

		<div class="divider">
			<span>or</span>
		</div>

		<a href={api.getGoogleAuthUrl()} class="google-btn">
			<svg viewBox="0 0 24 24" width="24" height="24" xmlns="http://www.w3.org/2000/svg">
				<g transform="matrix(1, 0, 0, 1, 27.009001, -39.238998)">
					<path fill="#4285F4" d="M -3.264 51.509 C -3.264 50.719 -3.334 49.969 -3.454 49.239 L -14.754 49.239 L -14.754 53.749 L -8.284 53.749 C -8.574 55.229 -9.424 56.479 -10.684 57.329 L -10.684 60.329 L -6.824 60.329 C -4.564 58.239 -3.264 55.159 -3.264 51.509 Z"/>
					<path fill="#34A853" d="M -14.754 63.239 C -11.514 63.239 -8.804 62.159 -6.824 60.329 L -10.684 57.329 C -11.764 58.049 -13.134 58.489 -14.754 58.489 C -17.884 58.489 -20.534 56.379 -21.484 53.529 L -25.464 53.529 L -25.464 56.619 C -23.494 60.539 -19.444 63.239 -14.754 63.239 Z"/>
					<path fill="#FBBC05" d="M -21.484 53.529 C -21.734 52.809 -21.864 52.039 -21.864 51.239 C -21.864 50.439 -21.724 49.669 -21.484 48.949 L -21.484 45.859 L -25.464 45.859 C -26.284 47.479 -26.754 49.299 -26.754 51.239 C -26.754 53.179 -26.284 54.999 -25.464 56.619 L -21.484 53.529 Z"/>
					<path fill="#EA4335" d="M -14.754 43.989 C -12.984 43.989 -11.404 44.599 -10.154 45.789 L -6.734 42.369 C -8.804 40.429 -11.514 39.239 -14.754 39.239 C -19.444 39.239 -23.494 41.939 -25.464 45.859 L -21.484 48.949 C -20.534 46.099 -17.884 43.989 -14.754 43.989 Z"/>
				</g>
			</svg>
			<span>Continue with Google</span>
		</a>

		<button class="toggle-mode" onclick={toggleMode}>
			{mode === 'signin' ? "Don't have an account? Sign Up" : 'Already have an account? Sign In'}
		</button>

		<p class="terms">
			By signing in, you agree to our Terms of Service and Privacy Policy.
			You must be 21+ to use this platform.
		</p>
	</div>
</div>

<style>
	:global(body) {
		font-family: 'Montserrat', sans-serif;
		background: #0a0a0a;
		color: #fff;
		margin: 0;
	}

	.login-page {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 2rem;
		background: linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 100%);
	}

	.login-card {
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		padding: 3rem;
		max-width: 400px;
		width: 100%;
		text-align: center;
		backdrop-filter: blur(10px);
	}

	.logo {
		height: 2.5rem;
		margin-bottom: 2rem;
	}

	h1 {
		font-family: 'Playfair Display', serif;
		font-size: 2rem;
		font-weight: 500;
		margin: 0 0 0.5rem 0;
	}

	.subtitle {
		color: rgba(255, 255, 255, 0.6);
		margin: 0 0 2rem 0;
		font-size: 0.95rem;
	}

	form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.input-group {
		text-align: left;
	}

	.input-group label {
		display: block;
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.7);
		margin-bottom: 0.5rem;
	}

	.input-group input {
		width: 100%;
		padding: 0.875rem 1rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		font-size: 1rem;
		transition: border-color 0.2s, background 0.2s;
		box-sizing: border-box;
	}

	.input-group input::placeholder {
		color: rgba(255, 255, 255, 0.3);
	}

	.input-group input:focus {
		outline: none;
		border-color: rgba(255, 255, 255, 0.4);
		background: rgba(255, 255, 255, 0.08);
	}

	.error-msg {
		background: rgba(220, 38, 38, 0.15);
		border: 1px solid rgba(220, 38, 38, 0.3);
		color: #fca5a5;
		padding: 0.75rem 1rem;
		border-radius: 0;
		font-size: 0.9rem;
	}

	.submit-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		padding: 1rem 1.5rem;
		background: #fff;
		color: #0a0a0a;
		border: none;
		border-radius: 0;
		font-family: 'Montserrat', sans-serif;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
		margin-top: 0.5rem;
	}

	.submit-btn:hover:not(:disabled) {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(255, 255, 255, 0.2);
	}

	.submit-btn:disabled {
		opacity: 0.7;
		cursor: not-allowed;
	}

	.spinner {
		width: 20px;
		height: 20px;
		border: 2px solid #0a0a0a;
		border-top-color: transparent;
		border-radius: 0;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.divider {
		display: flex;
		align-items: center;
		margin: 1.5rem 0;
		color: rgba(255, 255, 255, 0.4);
	}

	.divider::before,
	.divider::after {
		content: '';
		flex: 1;
		height: 1px;
		background: rgba(255, 255, 255, 0.1);
	}

	.divider span {
		padding: 0 1rem;
		font-size: 0.85rem;
	}

	.google-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.75rem;
		width: 100%;
		padding: 1rem 1.5rem;
		background: #fff;
		color: #333;
		border: none;
		border-radius: 0;
		font-family: 'Montserrat', sans-serif;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		text-decoration: none;
		transition: all 0.2s ease;
	}

	.google-btn:hover {
		background: #f5f5f5;
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
	}

	.toggle-mode {
		background: none;
		border: none;
		color: rgba(255, 255, 255, 0.7);
		font-family: 'Montserrat', sans-serif;
		font-size: 0.9rem;
		cursor: pointer;
		margin-top: 1.5rem;
		padding: 0;
	}

	.toggle-mode:hover {
		color: #fff;
	}


	.terms {
		margin-top: 1.5rem;
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.4);
		line-height: 1.6;
	}
</style>
