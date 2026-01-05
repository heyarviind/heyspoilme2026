<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';

	let showPricingPopup = $state(false);

	interface Conversation {
		id: string;
		initiated_by: string;
		other_user?: {
			user_id: string;
			display_name: string;
			age: number;
			city: string;
			is_online: boolean;
			images: Array<{ url: string; is_primary: boolean }>;
		};
		last_message?: {
			content: string;
			image_url?: string;
			created_at: string;
			sender_id: string;
		};
		unread_count: number;
	}

	interface LockedConversation {
		id: string;
		blurred_preview: string;
		sender_image?: string;
		sender_age?: number;
		sender_city?: string;
		created_at: string;
	}

	interface InboxResponse {
		conversations: Conversation[];
		locked_count: number;
		locked_previews?: LockedConversation[];
		can_view_all_messages: boolean;
	}

	let conversations = $state<Conversation[]>([]);
	let lockedCount = $state(0);
	let lockedPreviews = $state<LockedConversation[]>([]);
	let canViewAllMessages = $state(true);
	let loading = $state(true);

	async function loadConversations() {
		loading = true;
		try {
			const data = await api.getInbox() as InboxResponse;
			conversations = data.conversations || [];
			lockedCount = data.locked_count || 0;
			lockedPreviews = data.locked_previews || [];
			canViewAllMessages = data.can_view_all_messages;
		} catch (e) {
			console.error('Failed to load inbox:', e);
			// Fallback to old API if inbox endpoint not available
			try {
				const fallback = await api.getConversations() as { conversations: Conversation[] };
				conversations = fallback.conversations || [];
				canViewAllMessages = true;
			} catch {
				// ignore
			}
		} finally {
			loading = false;
		}
	}

	function getProfileImage(conversation: Conversation): string {
		const primary = conversation.other_user?.images?.find(img => img.is_primary);
		if (primary) return primary.url;
		if (conversation.other_user?.images?.length) return conversation.other_user.images[0].url;
		return 'https://via.placeholder.com/80?text=?';
	}

	function formatTime(dateStr?: string): string {
		if (!dateStr) return '';
		const date = new Date(dateStr);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);
		
		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins}m`;
		const diffHours = Math.floor(diffMins / 60);
		if (diffHours < 24) return `${diffHours}h`;
		const diffDays = Math.floor(diffHours / 24);
		if (diffDays < 7) return `${diffDays}d`;
		return date.toLocaleDateString();
	}

	function getMessagePreview(message?: { content?: string; image_url?: string }): string {
		if (!message) return '';
		if (message.image_url && !message.content) {
			return '📷 Photo';
		}
		if (message.image_url && message.content) {
			return '📷 ' + truncateMessage(message.content);
		}
		return truncateMessage(message.content);
	}

	function truncateMessage(content?: string): string {
		if (!content) return '';
		if (content.length <= 50) return content;
		return content.substring(0, 50) + '...';
	}

	onMount(() => {
		loadConversations();
	});
</script>

<svelte:head>
	<title>Messages | HeySpoilMe</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;500;600;700&family=Montserrat:wght@300;400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="messages-page">
	<Header />

	<main class="main">
		<h1>Messages</h1>

		{#if loading}
			<div class="loading">
				<div class="spinner"></div>
			</div>
		{:else if !canViewAllMessages && lockedCount > 0}
			<!-- Locked inbox for males with wealth_status=none -->
			<div class="locked-inbox">
				<div class="locked-header">
					<div class="locked-count">
						<span class="count-number">{lockedCount}</span>
						<span class="count-label">Message Request{lockedCount === 1 ? '' : 's'}</span>
					</div>
					<p class="locked-subtitle">Women have messaged you! Unlock to read and reply.</p>
				</div>

				<div class="locked-previews">
					{#each lockedPreviews as preview}
						<div class="locked-item">
							<div class="locked-avatar-container">
								{#if preview.sender_image}
									<img src={preview.sender_image} alt="Profile" class="locked-avatar blurred" />
								{:else}
									<div class="locked-avatar placeholder">?</div>
								{/if}
							</div>
							<div class="locked-info">
								<div class="locked-top-row">
									<span class="locked-name">{preview.sender_age ? `${preview.sender_age}` : ''}{preview.sender_city ? `, ${preview.sender_city}` : ''}</span>
								</div>
								<p class="locked-preview blurred-text">{preview.blurred_preview}</p>
							</div>
							<svg class="lock-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
								<rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
								<path d="m7 11v-4a5 5 0 0 1 10 0v4"></path>
							</svg>
						</div>
					{/each}
				</div>

				<div class="unlock-cta">
					<button class="unlock-button" onclick={() => showPricingPopup = true}>
						Unlock Messages
					</button>
					<p class="unlock-benefits">
						✓ Read & reply to messages<br/>
						✓ Priority placement in discovery<br/>
						✓ Trusted Member badge
					</p>
					<p class="unlock-note">We keep the community high-quality by allowing only verified, serious members to initiate conversations.</p>
				</div>
			</div>
		{:else if conversations.length === 0}
			<div class="empty">
				<p>No messages yet</p>
				<p class="hint">Start a conversation by visiting someone's profile</p>
				<a href="/browse" class="browse-link">Browse Profiles</a>
			</div>
		{:else}
			<div class="conversation-list">
				{#each conversations as conversation}
					<a href="/messages/{conversation.id}" class="conversation-item">
						<div class="avatar-container">
							<img 
								src={getProfileImage(conversation)} 
								alt="Profile" 
								class="avatar" 
							/>
							{#if conversation.other_user?.is_online}
								<span class="online-dot"></span>
							{/if}
						</div>
						<div class="conversation-info">
							<div class="top-row">
								<span class="name">{conversation.other_user?.display_name || 'Unknown'}</span>
								<span class="time">{formatTime(conversation.last_message?.created_at)}</span>
							</div>
							<div class="bottom-row">
								<p class="preview">{getMessagePreview(conversation.last_message)}</p>
								{#if conversation.unread_count > 0}
									<span class="unread-badge">{conversation.unread_count}</span>
								{/if}
							</div>
						</div>
					</a>
				{/each}
			</div>
		{/if}
	</main>

	<Footer />
</div>

{#if showPricingPopup}
	<div class="popup-overlay" onclick={() => showPricingPopup = false} role="button" tabindex="0" onkeypress={(e) => e.key === 'Enter' && (showPricingPopup = false)}>
		<div class="popup-content" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true">
			<button class="popup-close" onclick={() => showPricingPopup = false}>×</button>
			
			<div class="popup-header">
				<img src="/img/crown.png" alt="Crown" class="crown-img" />
				<h2>Upgrade to Trusted Member</h2>
			</div>

			<div class="pricing-card">
				<div class="price">
					<span class="currency">₹</span>
					<span class="amount">5,999</span>
					<span class="period">/month</span>
				</div>
			</div>

			<p class="pricing-tagline">Stand out as a verified, serious member in a private community built on trust and discretion.</p>

			<ul class="benefits-list">
				<li>
					<span class="check">✓</span>
					<span>Trusted badge that signals credibility</span>
				</li>
				<li>
					<span class="check">✓</span>
					<span>Priority placement in member discovery</span>
				</li>
				<li>
					<span class="check">✓</span>
					<span>Increased visibility to verified profiles</span>
				</li>
				<li>
					<span class="check">✓</span>
					<span>Unlock message requests & replies</span>
				</li>
				<li>
					<span class="check">✓</span>
					<span>Designed for members who value quality over volume</span>
				</li>
			</ul>

			<button class="subscribe-btn" onclick={() => goto('/profile/verify')}>
				Unlock Trusted Status
			</button>

			<p class="terms">
				By subscribing, you agree to our terms. Cancel anytime.
			</p>
		</div>
	</div>
{/if}

<style>
	:global(body) {
		font-family: 'Montserrat', sans-serif;
		background: #0a0a0a;
		color: #fff;
		margin: 0;
	}

	.messages-page {
		min-height: 100vh;
	}

	.main {
		max-width: 600px;
		margin: 0 auto;
		padding: 2rem;
	}

	h1 {
		font-family: 'Playfair Display', serif;
		font-size: 2rem;
		margin: 0 0 2rem 0;
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
		margin: 0 auto;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.hint {
		font-size: 0.9rem;
		color: rgba(255, 255, 255, 0.4);
	}

	.browse-link {
		display: inline-block;
		margin-top: 1rem;
		padding: 0.75rem 1.5rem;
		background: #fff;
		color: #000;
		text-decoration: none;
		border-radius: 0;
		font-weight: 500;
	}

	.conversation-list {
		display: flex;
		flex-direction: column;
	}

	.conversation-item {
		display: flex;
		gap: 1rem;
		padding: 1rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
		text-decoration: none;
		color: inherit;
		transition: background 0.2s ease;
	}

	.conversation-item:hover {
		background: rgba(255, 255, 255, 0.03);
	}

	.avatar-container {
		position: relative;
		flex-shrink: 0;
	}

	.avatar {
		width: 56px;
		height: 56px;
		border-radius: 0;
		object-fit: cover;
	}

	.online-dot {
		position: absolute;
		bottom: 2px;
		right: 2px;
		width: 14px;
		height: 14px;
		background: #22c55e;
		border-radius: 0;
		border: 2px solid #0a0a0a;
	}

	.conversation-info {
		flex: 1;
		min-width: 0;
		display: flex;
		flex-direction: column;
		justify-content: center;
	}

	.top-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.25rem;
	}

	.name {
		font-weight: 600;
	}

	.time {
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.4);
	}

	.bottom-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.preview {
		margin: 0;
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.6);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.unread-badge {
		background: #fff;
		color: #000;
		font-size: 0.7rem;
		font-weight: 600;
		padding: 0.2rem 0.5rem;
		border-radius: 0;
		flex-shrink: 0;
	}

	@media (max-width: 768px) {
		.main {
			padding: 1rem;
		}
	}

	/* Locked Inbox Styles */
	.locked-inbox {
		padding: 2rem 1rem;
	}

	.locked-header {
		text-align: center;
		margin-bottom: 2rem;
	}

	.locked-count {
		display: flex;
		flex-direction: column;
		align-items: center;
		margin-bottom: 1rem;
	}

	.count-number {
		font-family: 'Playfair Display', serif;
		font-size: 4rem;
		font-weight: 700;
		background: linear-gradient(135deg, #d4af37 0%, #f4e4ba 50%, #d4af37 100%);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
		line-height: 1;
	}

	.count-label {
		font-size: 1.2rem;
		color: rgba(255, 255, 255, 0.8);
		margin-top: 0.5rem;
	}

	.locked-subtitle {
		color: rgba(255, 255, 255, 0.6);
		font-size: 0.95rem;
	}

	.locked-previews {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		margin-bottom: 2rem;
	}

	.locked-item {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1rem;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.08);
	}

	.locked-avatar-container {
		flex-shrink: 0;
	}

	.locked-avatar {
		width: 56px;
		height: 56px;
		object-fit: cover;
	}

	.locked-avatar.blurred {
		filter: blur(8px);
	}

	.locked-avatar.placeholder {
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(255, 255, 255, 0.1);
		color: rgba(255, 255, 255, 0.3);
		font-size: 1.5rem;
	}

	.locked-info {
		flex: 1;
		min-width: 0;
	}

	.locked-top-row {
		margin-bottom: 0.25rem;
	}

	.locked-name {
		font-weight: 500;
		color: rgba(255, 255, 255, 0.9);
	}

	.locked-preview {
		margin: 0;
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.4);
	}

	.blurred-text {
		filter: blur(3px);
		user-select: none;
	}

	.lock-icon {
		width: 20px;
		height: 20px;
		color: rgba(255, 255, 255, 0.3);
		flex-shrink: 0;
	}

	.unlock-cta {
		text-align: center;
		padding: 2rem;
		background: linear-gradient(135deg, rgba(212, 175, 55, 0.1) 0%, rgba(212, 175, 55, 0.05) 100%);
		border: 1px solid rgba(212, 175, 55, 0.3);
	}

	.unlock-button {
		display: inline-block;
		padding: 1rem 3rem;
		background: linear-gradient(135deg, #d4af37 0%, #c5a028 100%);
		color: #000;
		text-decoration: none;
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		font-size: 1.1rem;
		letter-spacing: 0.5px;
		border: none;
		cursor: pointer;
		transition: transform 0.2s ease, box-shadow 0.2s ease;
	}

	.unlock-button:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 20px rgba(212, 175, 55, 0.4);
	}

	.unlock-benefits {
		margin: 1.5rem 0 1rem;
		color: rgba(255, 255, 255, 0.7);
		font-size: 0.9rem;
		line-height: 1.8;
	}

	.unlock-note {
		color: rgba(255, 255, 255, 0.4);
		font-size: 0.8rem;
		max-width: 400px;
		margin: 0 auto;
	}

	/* Pricing Popup Styles */
	.popup-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.85);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 1rem;
	}

	.popup-content {
		background: #0a0a0a;
		border: 1px solid rgba(255, 255, 255, 0.15);
		max-width: 420px;
		width: 100%;
		padding: 2rem;
		position: relative;
	}

	.popup-close {
		position: absolute;
		top: 1rem;
		right: 1rem;
		background: none;
		border: none;
		color: rgba(255, 255, 255, 0.5);
		font-size: 1.5rem;
		cursor: pointer;
		padding: 0.5rem;
		line-height: 1;
	}

	.popup-close:hover {
		color: #fff;
	}

	.popup-header {
		text-align: center;
		margin-bottom: 1.5rem;
	}

	.crown-img {
		width: 64px;
		height: 64px;
		object-fit: contain;
		display: block;
		margin: 0 auto 1rem;
	}

	.popup-header h2 {
		font-family: 'Playfair Display', serif;
		font-size: 1.5rem;
		margin: 0;
	}

	.pricing-card {
		background: linear-gradient(135deg, rgba(251, 191, 36, 0.15) 0%, rgba(251, 191, 36, 0.05) 100%);
		border: 1px solid rgba(251, 191, 36, 0.3);
		padding: 1.5rem;
		text-align: center;
		margin-bottom: 1rem;
	}

	.price {
		display: flex;
		align-items: baseline;
		justify-content: center;
		gap: 0.25rem;
	}

	.currency {
		font-size: 1.5rem;
		color: #fbbf24;
	}

	.amount {
		font-size: 2.5rem;
		font-weight: 700;
		color: #fbbf24;
	}

	.period {
		font-size: 1rem;
		color: rgba(255, 255, 255, 0.5);
	}

	.pricing-tagline {
		text-align: center;
		color: rgba(255, 255, 255, 0.7);
		font-size: 0.9rem;
		line-height: 1.6;
		margin: 0 0 1.5rem;
	}

	.benefits-list {
		list-style: none;
		padding: 0;
		margin: 0 0 2rem 0;
	}

	.benefits-list li {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 0;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
	}

	.benefits-list li:last-child {
		border-bottom: none;
	}

	.benefits-list .check {
		color: #fbbf24;
		font-weight: bold;
	}

	.benefits-list span:last-child {
		color: rgba(255, 255, 255, 0.8);
		font-size: 0.9rem;
	}

	.subscribe-btn {
		width: 100%;
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
		color: #000;
		border: none;
		padding: 1rem;
		font-family: 'Montserrat', sans-serif;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: transform 0.2s ease, box-shadow 0.2s ease;
	}

	.subscribe-btn:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 20px rgba(251, 191, 36, 0.4);
	}

	.terms {
		text-align: center;
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.4);
		margin: 1rem 0 0;
	}
</style>

