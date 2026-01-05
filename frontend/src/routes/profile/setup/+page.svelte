<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { auth } from '$lib/stores/auth';
	import { generateDisplayName } from '$lib/utils';
	import CityAutocomplete from '$lib/components/CityAutocomplete.svelte';

	let step = $state(1);
	let loading = $state(false);
	let error = $state('');

	// Profile data
	let gender = $state<'male' | 'female' | ''>('');
	let displayName = $state('');
	let age = $state(25);
	let bio = $state('');
	let salaryRange = $state('');
	let city = $state('');
	let state_ = $state('');
	let latitude = $state(0);
	let longitude = $state(0);

	// Generate a display name when gender is selected
	function selectGender(selectedGender: 'male' | 'female') {
		gender = selectedGender;
		displayName = generateDisplayName(selectedGender);
	}

	// Regenerate name with same gender
	function regenerateName() {
		if (gender) {
			displayName = generateDisplayName(gender);
		}
	}

	// Handle city selection from autocomplete
	function handleCitySelect(selectedCity: { city: string; state: string; latitude: number; longitude: number }) {
		city = selectedCity.city;
		state_ = selectedCity.state;
		latitude = selectedCity.latitude;
		longitude = selectedCity.longitude;
	}

	const salaryOptions = [
		'5-10 LPA',
		'10-20 LPA',
		'20-50 LPA',
		'50+ LPA',
	];

	function nextStep() {
		error = '';
		
		if (step === 1 && !gender) {
			error = 'Please select your gender';
			return;
		}
		
		if (step === 2) {
			if (age < 21 || age > 100) {
				error = 'Age must be between 21 and 100';
				return;
			}
			if (!bio.trim()) {
				error = 'Please write something about yourself';
				return;
			}
		}

		if (step === 3) {
			if (!city.trim() || !state_ || latitude === 0) {
				error = 'Please select a city from the suggestions';
				return;
			}
		}

		if (step === 4 && gender === 'male' && !salaryRange) {
			error = 'Please select your salary range';
			return;
		}

		step++;
	}

	function prevStep() {
		if (step > 1) step--;
	}

	// No longer need browser geolocation - we get lat/lng from city database

	async function submitProfile() {
		loading = true;
		error = '';

		try {
			const profileData: any = {
				display_name: displayName,
				gender,
				age,
				bio: bio.trim(),
				city: city.trim(),
				state: state_,
				latitude,
				longitude,
			};

			if (gender === 'male') {
				profileData.salary_range = salaryRange;
			}

			await api.createProfile(profileData);
			await auth.refreshProfile();
			goto('/browse');
		} catch (e: any) {
			error = e.message || 'Failed to save profile';
		} finally {
			loading = false;
		}
	}


	// Skip salary step for females
	$effect(() => {
		if (step === 4 && gender === 'female') {
			submitProfile();
		}
	});

	const totalSteps = $derived(gender === 'male' ? 4 : 3);
</script>

<svelte:head>
	<title>Complete Your Profile | HeySpoilMe</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;500;600;700&family=Montserrat:wght@300;400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="setup-page">
	<div class="setup-card">
		<img src="/img/logo.svg" alt="HeySpoilMe" class="logo" />
		
		<div class="progress">
			{#each Array(totalSteps) as _, i}
				<div class="progress-dot" class:active={i + 1 <= step}></div>
			{/each}
		</div>

		{#if step === 1}
			<div class="step">
				<h2>What's your gender?</h2>
				<p class="description">This helps us personalize your experience</p>
				
				<div class="gender-options">
					<button 
						class="gender-btn" 
						class:selected={gender === 'male'}
						onclick={() => selectGender('male')}
					>
						<span class="icon">👨</span>
						<span>I am a Man</span>
					</button>
					<button 
						class="gender-btn" 
						class:selected={gender === 'female'}
						onclick={() => selectGender('female')}
					>
						<span class="icon">👩</span>
						<span>I am a Woman</span>
					</button>
				</div>

				{#if displayName}
					<div class="generated-name">
						<p class="name-label">Your display name:</p>
						<div class="name-display">
							<span class="name-text">{displayName}</span>
							<button class="regenerate-btn" onclick={regenerateName} type="button">
								🎲
							</button>
						</div>
						<p class="name-hint">You can change this later in settings</p>
					</div>
				{/if}
			</div>
		{:else if step === 2}
			<div class="step">
				<h2>Tell us about yourself</h2>
				<p class="description">Help others get to know you</p>

				<div class="form-group">
					<label for="age">Your Age</label>
					<input 
						type="number" 
						id="age" 
						bind:value={age} 
						min="21" 
						max="100"
						class="input"
					/>
				</div>

				<div class="form-group">
					<label for="bio">About You</label>
					<textarea 
						id="bio" 
						bind:value={bio} 
						placeholder="Write a few words about yourself..."
						rows="4"
						maxlength="500"
						class="input"
					></textarea>
					<span class="char-count">{bio.length}/500</span>
				</div>
			</div>
		{:else if step === 3}
			<div class="step">
				<h2>Where are you located?</h2>
				<p class="description">Find people near you</p>

				<div class="form-group">
					<label for="city">City</label>
					<CityAutocomplete 
						bind:value={city}
						onSelect={handleCitySelect}
						placeholder="Start typing your city..."
					/>
				</div>

				{#if state_}
					<div class="form-group">
						<label>State</label>
						<div class="state-display">{state_}</div>
					</div>
				{/if}
			</div>
		{:else if step === 4 && gender === 'male'}
			<div class="step">
				<h2>Your Monthly Salary Range</h2>
				<p class="description">This is shown to potential matches as a range</p>

				<div class="salary-options">
					{#each salaryOptions as option}
						<button 
							class="salary-btn"
							class:selected={salaryRange === option}
							onclick={() => salaryRange = option}
						>
							{option}
						</button>
					{/each}
				</div>
			</div>
		{/if}

		{#if error}
			<p class="error">{error}</p>
		{/if}

		<div class="actions">
			{#if step > 1}
				<button class="btn-secondary" onclick={prevStep} disabled={loading}>
					Back
				</button>
			{/if}

			{#if step < totalSteps}
				<button class="btn-primary" onclick={nextStep}>
					Continue
				</button>
			{:else}
				<button class="btn-primary" onclick={submitProfile} disabled={loading}>
					{loading ? 'Saving...' : 'Complete Setup'}
				</button>
			{/if}
		</div>
	</div>
</div>

<style>
	:global(body) {
		font-family: 'Montserrat', sans-serif;
		background: #0a0a0a;
		color: #fff;
		margin: 0;
	}

	.setup-page {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 2rem;
		background: linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 100%);
	}

	.setup-card {
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		padding: 2.5rem;
		max-width: 480px;
		width: 100%;
		text-align: center;
	}

	.logo {
		height: 3rem;
		margin-bottom: 2rem;
	}

	.progress {
		display: flex;
		justify-content: center;
		gap: 0.5rem;
		margin-bottom: 2rem;
	}

	.progress-dot {
		width: 8px;
		height: 8px;
		border-radius: 0;
		background: rgba(255, 255, 255, 0.2);
		transition: all 0.3s ease;
	}

	.progress-dot.active {
		background: #fff;
		width: 24px;
		border-radius: 0;
	}

	.step h2 {
		font-family: 'Playfair Display', serif;
		font-size: 1.5rem;
		font-weight: 500;
		margin: 0 0 0.5rem 0;
	}

	.description {
		color: rgba(255, 255, 255, 0.5);
		margin: 0 0 2rem 0;
		font-size: 0.9rem;
	}

	.gender-options {
		display: flex;
		gap: 1rem;
	}

	.gender-btn {
		flex: 1;
		padding: 1.5rem;
		background: rgba(255, 255, 255, 0.05);
		border: 2px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		color: #fff;
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		font-family: 'Montserrat', sans-serif;
	}

	.gender-btn:hover {
		border-color: rgba(255, 255, 255, 0.3);
	}

	.gender-btn.selected {
		background: rgba(255, 255, 255, 0.1);
		border-color: #fff;
	}

	.gender-btn .icon {
		font-size: 2rem;
	}

	.generated-name {
		margin-top: 2rem;
		padding: 1.5rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.name-label {
		margin: 0 0 0.75rem 0;
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.6);
	}

	.name-display {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.75rem;
	}

	.name-text {
		font-family: 'Playfair Display', serif;
		font-size: 1.5rem;
		font-weight: 600;
		color: #fff;
	}

	.regenerate-btn {
		width: 40px;
		height: 40px;
		background: rgba(255, 255, 255, 0.1);
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 0;
		font-size: 1.2rem;
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.regenerate-btn:hover {
		background: rgba(255, 255, 255, 0.15);
		transform: rotate(180deg);
	}

	.name-hint {
		margin: 0.75rem 0 0 0;
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.4);
	}

	.form-group {
		margin-bottom: 1.5rem;
		text-align: left;
	}

	.form-group label {
		display: block;
		margin-bottom: 0.5rem;
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.7);
	}

	.input {
		width: 100%;
		padding: 0.875rem 1rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		font-size: 1rem;
		transition: border-color 0.2s ease;
		box-sizing: border-box;
	}

	.input:focus {
		outline: none;
		border-color: rgba(255, 255, 255, 0.3);
	}

	textarea.input {
		resize: vertical;
		min-height: 100px;
	}

	select.input {
		appearance: none;
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='white' d='M6 8L1 3h10z'/%3E%3C/svg%3E");
		background-repeat: no-repeat;
		background-position: right 1rem center;
		padding-right: 2.5rem;
	}

	.char-count {
		display: block;
		text-align: right;
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.4);
		margin-top: 0.25rem;
	}

	.salary-options {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 1rem;
	}

	.salary-btn {
		padding: 1rem;
		background: rgba(255, 255, 255, 0.05);
		border: 2px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		color: #fff;
		cursor: pointer;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.9rem;
		transition: all 0.2s ease;
	}

	.salary-btn:hover {
		border-color: rgba(255, 255, 255, 0.3);
	}

	.salary-btn.selected {
		background: rgba(255, 255, 255, 0.1);
		border-color: #fff;
	}

	.error {
		color: #f87171;
		margin: 1rem 0;
		font-size: 0.9rem;
	}

	.actions {
		display: flex;
		gap: 1rem;
		margin-top: 2rem;
	}

	.btn-primary, .btn-secondary {
		flex: 1;
		padding: 0.875rem 1.5rem;
		border-radius: 0;
		font-family: 'Montserrat', sans-serif;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
		border: none;
	}

	.btn-primary {
		background: #fff;
		color: #000;
	}

	.btn-primary:hover:not(:disabled) {
		transform: translateY(-2px);
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-secondary {
		background: transparent;
		color: #fff;
		border: 1px solid rgba(255, 255, 255, 0.2);
	}

	.btn-secondary:hover:not(:disabled) {
		border-color: rgba(255, 255, 255, 0.4);
	}

	.state-display {
		padding: 0.875rem 1rem;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.08);
		color: rgba(255, 255, 255, 0.7);
		font-size: 1rem;
	}
</style>

