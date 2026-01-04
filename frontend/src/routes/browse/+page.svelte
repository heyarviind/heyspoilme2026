<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page as pageStore } from '$app/stores';
	import { api } from '$lib/api';
	import { auth } from '$lib/stores/auth';

	interface Profile {
		id: string;
		user_id: string;
		gender: 'male' | 'female';
		age: number;
		bio: string;
		salary_range?: string;
		city: string;
		state: string;
		is_online: boolean;
		is_verified: boolean;
		last_seen?: string;
		distance_km?: number;
		is_liked: boolean;
		images: Array<{ id: string; url: string; is_primary: boolean }>;
	}

	let profiles = $state<Profile[]>([]);
	let loading = $state(true);
	let page = $state(1);
	let total = $state(0);
	let userGender = $state<'male' | 'female' | null>(null);
	let filters = $state({
		minAge: 21,
		maxAge: 60,
		maxDistance: 0,
		onlineOnly: false,
	});

	// Get user's gender to show opposite
	let authState = $state<any>(null);
	auth.subscribe(s => {
		authState = s;
		if (s?.profile?.gender) {
			userGender = s.profile.gender;
		}
	});

	// Read filters from URL on page load
	function initFiltersFromUrl() {
		const url = new URL(window.location.href);
		const minAge = url.searchParams.get('minAge');
		const maxAge = url.searchParams.get('maxAge');
		const maxDistance = url.searchParams.get('maxDistance');
		const onlineOnly = url.searchParams.get('onlineOnly');

		if (minAge) filters.minAge = parseInt(minAge) || 21;
		if (maxAge) filters.maxAge = parseInt(maxAge) || 60;
		if (maxDistance) filters.maxDistance = parseInt(maxDistance) || 0;
		if (onlineOnly === 'true') filters.onlineOnly = true;
	}

	// Update URL when filters change
	function updateUrlParams() {
		const url = new URL(window.location.href);
		
		if (filters.minAge > 21) {
			url.searchParams.set('minAge', filters.minAge.toString());
		} else {
			url.searchParams.delete('minAge');
		}
		
		if (filters.maxAge < 60) {
			url.searchParams.set('maxAge', filters.maxAge.toString());
		} else {
			url.searchParams.delete('maxAge');
		}
		
		if (filters.maxDistance > 0) {
			url.searchParams.set('maxDistance', filters.maxDistance.toString());
		} else {
			url.searchParams.delete('maxDistance');
		}
		
		if (filters.onlineOnly) {
			url.searchParams.set('onlineOnly', 'true');
		} else {
			url.searchParams.delete('onlineOnly');
		}

		goto(url.pathname + url.search, { replaceState: true, noScroll: true });
	}

	function onFilterChange() {
		updateUrlParams();
		loadProfiles();
	}

	async function loadProfiles() {
		loading = true;
		try {
			const params: Record<string, any> = { page, limit: 12 };
			
			// Show opposite gender only
			if (userGender === 'male') {
				params.gender = 'female';
			} else if (userGender === 'female') {
				params.gender = 'male';
			}
			
			if (filters.minAge > 21) params.min_age = filters.minAge;
			if (filters.maxAge < 60) params.max_age = filters.maxAge;
			if (filters.maxDistance) params.max_distance = filters.maxDistance;
			if (filters.onlineOnly) params.online_only = true;

			const data = await api.listProfiles(params) as { profiles: Profile[]; total: number };
			profiles = data.profiles || [];
			total = data.total;
		} catch (e) {
			console.error('Failed to load profiles:', e);
		} finally {
			loading = false;
		}
	}

	async function toggleLike(profile: Profile) {
		try {
			if (profile.is_liked) {
				await api.unlikeProfile(profile.user_id);
			} else {
				await api.likeProfile(profile.user_id);
			}
			profile.is_liked = !profile.is_liked;
		} catch (e) {
			console.error('Failed to toggle like:', e);
		}
	}

	function getPrimaryImage(profile: Profile): string {
		const primary = profile.images?.find(img => img.is_primary);
		if (primary) return primary.url;
		if (profile.images?.length > 0) return profile.images[0].url;
		return 'https://via.placeholder.com/300x400?text=No+Photo';
	}

	function formatDistance(km?: number): string {
		if (!km) return '';
		if (km < 1) return '<1 km away';
		if (km < 100) return `${Math.round(km)} km away`;
		return `${Math.round(km)} km away`;
	}

	function formatLastSeen(lastSeen?: string): string {
		if (!lastSeen) return '';
		const date = new Date(lastSeen);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);
		
		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins}m ago`;
		const diffHours = Math.floor(diffMins / 60);
		if (diffHours < 24) return `${diffHours}h ago`;
		const diffDays = Math.floor(diffHours / 24);
		return `${diffDays}d ago`;
	}

	onMount(() => {
		initFiltersFromUrl();
		loadProfiles();
	});
</script>

<svelte:head>
	<title>Browse | HeySpoilMe</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;500;600;700&family=Montserrat:wght@300;400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="browse-page">
	<header class="header">
		<a href="/browse" class="logo-link">
			<img src="/img/logo.svg" alt="HeySpoilMe" class="logo" />
		</a>
		<nav class="nav">
			<a href="/browse" class="nav-link active">Browse</a>
			<a href="/messages" class="nav-link">Messages</a>
			<a href="/likes" class="nav-link">Likes</a>
			<a href="/profile" class="nav-link">Profile</a>
		</nav>
	</header>

	<main class="main">
		<div class="filters">
			<div class="filter-group age-filter">
				<label>Age</label>
				<div class="age-inputs">
					<input 
						type="number" 
						min="21" 
						max="60"
						bind:value={filters.minAge}
						onchange={onFilterChange}
						class="age-input"
					/>
					<span>-</span>
					<input 
						type="number" 
						min="21" 
						max="60"
						bind:value={filters.maxAge}
						onchange={onFilterChange}
						class="age-input"
					/>
				</div>
			</div>
			<select bind:value={filters.maxDistance} onchange={onFilterChange} class="filter-select">
				<option value={0}>Any Distance</option>
				<option value={25}>Within 25 km</option>
				<option value={50}>Within 50 km</option>
				<option value={100}>Within 100 km</option>
			</select>
			<label class="online-filter">
				<input 
					type="checkbox" 
					bind:checked={filters.onlineOnly}
					onchange={onFilterChange}
				/>
				<span>Online Now</span>
			</label>
		</div>

		{#if loading}
			<div class="loading">
				<div class="spinner"></div>
				<p>Finding people near you...</p>
			</div>
		{:else if profiles.length === 0}
			<div class="empty">
				<p>No profiles found</p>
				<p class="hint">Try adjusting your filters</p>
			</div>
		{:else}
			<div class="profiles-grid">
				{#each profiles as profile}
					<div class="profile-card" class:verified={profile.is_verified && profile.gender === 'male'}>
						<a href="/profile/{profile.user_id}" class="card-link">
							<div class="image-container">
								<img src={getPrimaryImage(profile)} alt="{profile.age}" class="profile-image" />
								{#if profile.is_verified && profile.gender === 'male'}
									<span class="verified-badge">✓ VERIFIED</span>
								{/if}
								{#if profile.is_online}
									<span class="online-badge"></span>
								{/if}
								{#if profile.distance_km}
									<span class="distance-badge">{formatDistance(profile.distance_km)}</span>
								{/if}
							</div>
							<div class="card-content">
								<div class="card-header">
									<span class="age">{profile.age}</span>
									{#if !profile.is_online && profile.last_seen}
										<span class="last-seen">{formatLastSeen(profile.last_seen)}</span>
									{/if}
								</div>
								<p class="location">{profile.city}, {profile.state}</p>
								{#if profile.gender === 'male' && profile.salary_range}
									<p class="salary">{profile.salary_range}</p>
								{/if}
							</div>
						</a>
						<button 
							class="like-btn" 
							class:liked={profile.is_liked}
							onclick={() => toggleLike(profile)}
						>
							{profile.is_liked ? '❤️' : '🤍'}
						</button>
					</div>
				{/each}
			</div>

			{#if total > profiles.length}
				<div class="pagination">
					<button 
						onclick={() => { page++; loadProfiles(); }}
						class="load-more"
					>
						Load More
					</button>
				</div>
			{/if}
		{/if}
	</main>
</div>

<style>
	:global(body) {
		font-family: 'Montserrat', sans-serif;
		background: #0a0a0a;
		color: #fff;
		margin: 0;
	}

	.browse-page {
		min-height: 100vh;
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1rem 2rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
		position: sticky;
		top: 0;
		background: #0a0a0a;
		z-index: 100;
	}

	.logo-link {
		text-decoration: none;
	}

	.logo {
		height: 2.5rem;
	}

	.nav {
		display: flex;
		gap: 2rem;
	}

	.nav-link {
		color: rgba(255, 255, 255, 0.6);
		text-decoration: none;
		font-size: 0.9rem;
		font-weight: 500;
		transition: color 0.2s ease;
	}

	.nav-link:hover, .nav-link.active {
		color: #fff;
	}

	.main {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;
	}

	.filters {
		display: flex;
		gap: 1rem;
		margin-bottom: 2rem;
		flex-wrap: wrap;
	}

	.filters {
		align-items: center;
	}

	.filter-group {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.filter-group label {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.6);
	}

	.age-inputs {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.age-inputs span {
		color: rgba(255, 255, 255, 0.4);
	}

	.age-input {
		width: 60px;
		height: 44px;
		padding: 0 0.5rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.85rem;
		text-align: center;
		box-sizing: border-box;
	}

	.age-input:focus {
		outline: none;
		border-color: rgba(255, 255, 255, 0.3);
	}

	.filter-select {
		height: 44px;
		padding: 0 2.5rem 0 1rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.9rem;
		min-width: 140px;
		appearance: none;
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='white' d='M6 8L1 3h10z'/%3E%3C/svg%3E");
		background-repeat: no-repeat;
		background-position: right 1rem center;
		box-sizing: border-box;
	}

	.online-filter {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		cursor: pointer;
		height: 44px;
		padding: 0 1rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		transition: all 0.2s ease;
		box-sizing: border-box;
	}

	.online-filter:hover {
		border-color: rgba(255, 255, 255, 0.2);
	}

	.online-filter input {
		accent-color: #22c55e;
	}

	.online-filter span {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.7);
	}

	.online-filter:has(input:checked) {
		border-color: #22c55e;
		background: rgba(34, 197, 94, 0.1);
	}

	.online-filter:has(input:checked) span {
		color: #22c55e;
	}

	.loading, .empty {
		text-align: center;
		padding: 4rem 2rem;
		color: rgba(255, 255, 255, 0.6);
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 3px solid rgba(255, 255, 255, 0.1);
		border-top-color: #fff;
		border-radius: 0;
		animation: spin 1s linear infinite;
		margin: 0 auto 1rem;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.hint {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.4);
	}

	.profiles-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
		gap: 1.5rem;
	}

	.profile-card {
		position: relative;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		overflow: hidden;
		transition: all 0.2s ease;
	}

	.profile-card.verified {
		border: 2px solid #fbbf24;
		background: linear-gradient(135deg, rgba(251, 191, 36, 0.1) 0%, rgba(251, 191, 36, 0.05) 100%);
		box-shadow: 0 0 20px rgba(251, 191, 36, 0.15);
	}

	.profile-card:hover {
		transform: translateY(-4px);
		border-color: rgba(255, 255, 255, 0.2);
	}

	.profile-card.verified:hover {
		border-color: #fbbf24;
		box-shadow: 0 4px 30px rgba(251, 191, 36, 0.25);
	}

	.card-link {
		text-decoration: none;
		color: inherit;
	}

	.image-container {
		position: relative;
		aspect-ratio: 3/4;
		overflow: hidden;
	}

	.profile-image {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.online-badge {
		position: absolute;
		top: 0.75rem;
		right: 0.75rem;
		width: 12px;
		height: 12px;
		background: #22c55e;
		border-radius: 0;
		border: 2px solid #0a0a0a;
	}

	.verified-badge {
		position: absolute;
		top: 0.75rem;
		left: 0.75rem;
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
		color: #000;
		font-size: 0.65rem;
		font-weight: 700;
		padding: 0.25rem 0.5rem;
		letter-spacing: 0.5px;
	}

	.distance-badge {
		position: absolute;
		bottom: 0.75rem;
		left: 0.75rem;
		background: rgba(0, 0, 0, 0.7);
		padding: 0.25rem 0.5rem;
		border-radius: 0;
		font-size: 0.75rem;
		color: #fff;
	}

	.card-content {
		padding: 1rem;
	}

	.card-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 0.25rem;
	}

	.age {
		font-size: 1.1rem;
		font-weight: 600;
	}

	.last-seen {
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.4);
	}

	.location {
		margin: 0;
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.6);
	}

	.salary {
		margin: 0.5rem 0 0 0;
		font-size: 0.8rem;
		color: #22c55e;
	}

	.like-btn {
		position: absolute;
		top: 0.75rem;
		left: 0.75rem;
		width: 40px;
		height: 40px;
		background: rgba(0, 0, 0, 0.5);
		border: none;
		border-radius: 0;
		cursor: pointer;
		font-size: 1.2rem;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.2s ease;
	}

	.like-btn:hover {
		background: rgba(0, 0, 0, 0.7);
		transform: scale(1.1);
	}

	.pagination {
		text-align: center;
		margin-top: 2rem;
	}

	.load-more {
		padding: 0.875rem 2rem;
		background: transparent;
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.9rem;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.load-more:hover {
		border-color: rgba(255, 255, 255, 0.4);
	}

	@media (max-width: 768px) {
		.header {
			flex-direction: column;
			gap: 1rem;
			padding: 1rem;
		}

		.nav {
			gap: 1rem;
		}

		.main {
			padding: 1rem;
		}

		.filters {
			display: flex;
			flex-wrap: wrap;
			gap: 0.75rem;
		}

		.filter-group label {
			display: none;
		}

		.age-input {
			width: 50px;
			height: 38px;
			font-size: 0.8rem;
		}

		.filter-select {
			min-width: 0;
			flex: 1;
			height: 38px;
			font-size: 0.8rem;
			padding: 0 1.5rem 0 0.5rem;
			background-position: right 0.5rem center;
		}

		.online-filter {
			height: 38px;
			padding: 0 0.75rem;
		}

		.online-filter span {
			font-size: 0.8rem;
		}

		.profiles-grid {
			grid-template-columns: repeat(2, 1fr);
			gap: 1rem;
		}
	}
</style>

