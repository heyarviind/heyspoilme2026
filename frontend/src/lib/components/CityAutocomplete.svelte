<script lang="ts">
	import { api } from '$lib/api';

	interface City {
		city: string;
		state: string;
		latitude: number;
		longitude: number;
	}

	interface Props {
		value: string;
		onSelect: (city: City) => void;
		placeholder?: string;
		country?: string;
	}

	let { value = $bindable(), onSelect, placeholder = 'Start typing city name...', country = 'IN' }: Props = $props();

	let suggestions = $state<City[]>([]);
	let showSuggestions = $state(false);
	let loading = $state(false);
	let debounceTimer: ReturnType<typeof setTimeout> | null = null;
	let inputElement: HTMLInputElement;

	async function searchCities(query: string) {
		if (query.length < 2) {
			suggestions = [];
			return;
		}

		loading = true;
		try {
			const result = await api.searchCities(query, country);
			suggestions = result.cities || [];
		} catch (e) {
			console.error('Failed to search cities:', e);
			suggestions = [];
		} finally {
			loading = false;
		}
	}

	function handleInput(e: Event) {
		const target = e.target as HTMLInputElement;
		value = target.value;

		if (debounceTimer) {
			clearTimeout(debounceTimer);
		}

		debounceTimer = setTimeout(() => {
			searchCities(value);
		}, 200);
	}

	function handleSelect(city: City) {
		value = city.city;
		suggestions = [];
		showSuggestions = false;
		onSelect(city);
	}

	function handleFocus() {
		showSuggestions = true;
		if (value.length >= 2) {
			searchCities(value);
		}
	}

	function handleBlur() {
		// Delay hiding so click events on suggestions work
		setTimeout(() => {
			showSuggestions = false;
		}, 200);
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			showSuggestions = false;
			suggestions = [];
		}
	}
</script>

<div class="city-autocomplete">
	<input
		bind:this={inputElement}
		type="text"
		{value}
		{placeholder}
		oninput={handleInput}
		onfocus={handleFocus}
		onblur={handleBlur}
		onkeydown={handleKeydown}
		autocomplete="off"
		class="city-input"
	/>
	
	{#if showSuggestions && (suggestions.length > 0 || loading)}
		<div class="suggestions">
			{#if loading}
				<div class="suggestion-item loading">Searching...</div>
			{:else}
				{#each suggestions as city}
					<button 
						type="button"
						class="suggestion-item"
						onmousedown={() => handleSelect(city)}
					>
						<span class="city-name">{city.city}</span>
						<span class="state-name">{city.state}</span>
					</button>
				{/each}
			{/if}
		</div>
	{/if}
</div>

<style>
	.city-autocomplete {
		position: relative;
		width: 100%;
	}

	.city-input {
		width: 100%;
		padding: 0.75rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		font-size: 1rem;
	}

	.city-input:focus {
		outline: none;
		border-color: rgba(255, 255, 255, 0.3);
	}

	.suggestions {
		position: absolute;
		top: 100%;
		left: 0;
		right: 0;
		background: #1a1a1a;
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-top: none;
		max-height: 250px;
		overflow-y: auto;
		z-index: 100;
	}

	.suggestion-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		width: 100%;
		padding: 0.75rem 1rem;
		background: transparent;
		border: none;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
		color: #fff;
		cursor: pointer;
		text-align: left;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.9rem;
		transition: background 0.15s ease;
	}

	.suggestion-item:hover {
		background: rgba(255, 255, 255, 0.1);
	}

	.suggestion-item:last-child {
		border-bottom: none;
	}

	.suggestion-item.loading {
		color: rgba(255, 255, 255, 0.5);
		cursor: default;
	}

	.city-name {
		font-weight: 500;
	}

	.state-name {
		font-size: 0.8rem;
		color: rgba(255, 255, 255, 0.5);
	}
</style>

