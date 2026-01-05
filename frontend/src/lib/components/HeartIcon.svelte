<script lang="ts">
	interface Props {
		liked?: boolean;
		size?: number;
	}

	let { liked = false, size = 24 }: Props = $props();
	
	// Track animation trigger
	let animationKey = $state(0);
	let wasLiked = $state(liked);
	
	$effect(() => {
		if (liked && !wasLiked) {
			animationKey++;
		}
		wasLiked = liked;
	});
</script>

<div class="heart-container" style="--size: {size}px">
	{#key animationKey}
		{#if liked}
			<!-- Burst particles -->
			<div class="particles">
				{#each Array(6) as _, i}
					<span class="particle" style="--i: {i}"></span>
				{/each}
			</div>
			<!-- Ring explosion -->
			<div class="ring"></div>
		{/if}
	{/key}
	
	<svg 
		xmlns="http://www.w3.org/2000/svg" 
		width={size} 
		height={size} 
		viewBox="0 0 24 24" 
		fill={liked ? '#ec4899' : 'currentColor'}
		class="heart-icon"
		class:liked
		class:animate={liked && animationKey > 0}
		role="img"
		aria-label={liked ? 'Liked' : 'Not liked'}
	>
		<path d="M19.4189 15.6602C21.1899 13.624 22.75 11.153 22.75 8.69434C22.7499 5.45202 20.3484 2.75 17 2.75C15.4082 2.75 13.8662 3.26268 12 4.96484C10.1338 3.26268 8.59184 2.75 7 2.75C3.65156 2.75 1.25005 5.45202 1.25 8.69434C1.25 11.153 2.8101 13.624 4.58105 15.6602C6.37954 17.7279 8.5291 19.4969 9.96191 20.5684L10.1943 20.7285C11.3812 21.4741 12.8985 21.4204 14.0381 20.5684L14.6074 20.1348C16.0032 19.05 17.8453 17.4694 19.4189 15.6602Z"></path>
	</svg>
</div>

<style>
	.heart-container {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		width: var(--size);
		height: var(--size);
	}

	.heart-icon {
		display: block;
		transition: fill 0.15s ease;
		position: relative;
		z-index: 1;
	}

	.heart-icon:not(.liked) {
		opacity: 0.6;
	}

	.heart-icon.animate {
		animation: bounce 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
	}

	@keyframes bounce {
		0% { 
			transform: scale(0.2);
		}
		40% { 
			transform: scale(1.3);
		}
		60% {
			transform: scale(0.9);
		}
		80% {
			transform: scale(1.1);
		}
		100% { 
			transform: scale(1);
		}
	}

	/* Particle burst effect */
	.particles {
		position: absolute;
		width: 100%;
		height: 100%;
		pointer-events: none;
	}

	.particle {
		position: absolute;
		left: 50%;
		top: 50%;
		width: 6px;
		height: 6px;
		background: #ec4899;
		border-radius: 50%;
		animation: particle-burst 0.6s cubic-bezier(0.22, 0.61, 0.36, 1) forwards;
		--angle: calc(var(--i) * 60deg);
		transform-origin: center;
	}

	.particle:nth-child(2) { background: #f472b6; width: 4px; height: 4px; }
	.particle:nth-child(3) { background: #fb7185; }
	.particle:nth-child(4) { background: #f472b6; width: 5px; height: 5px; }
	.particle:nth-child(5) { background: #ec4899; width: 4px; height: 4px; }
	.particle:nth-child(6) { background: #fb7185; width: 5px; height: 5px; }

	@keyframes particle-burst {
		0% {
			opacity: 1;
			transform: translate(-50%, -50%) rotate(var(--angle)) translateY(0) scale(0);
		}
		50% {
			opacity: 1;
			transform: translate(-50%, -50%) rotate(var(--angle)) translateY(calc(var(--size) * -1)) scale(1);
		}
		100% {
			opacity: 0;
			transform: translate(-50%, -50%) rotate(var(--angle)) translateY(calc(var(--size) * -1.5)) scale(0.5);
		}
	}

	/* Ring explosion effect */
	.ring {
		position: absolute;
		left: 50%;
		top: 50%;
		width: calc(var(--size) * 0.5);
		height: calc(var(--size) * 0.5);
		border: 3px solid #ec4899;
		border-radius: 50%;
		transform: translate(-50%, -50%) scale(0);
		animation: ring-burst 0.5s cubic-bezier(0.22, 0.61, 0.36, 1) forwards;
		pointer-events: none;
	}

	@keyframes ring-burst {
		0% {
			transform: translate(-50%, -50%) scale(0);
			opacity: 1;
			border-width: 4px;
		}
		100% {
			transform: translate(-50%, -50%) scale(3);
			opacity: 0;
			border-width: 1px;
		}
	}
</style>
