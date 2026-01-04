<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	import { api } from '$lib/api';
	import { websocket } from '$lib/stores/websocket';
	import { auth } from '$lib/stores/auth';

	interface Message {
		id: string;
		conversation_id: string;
		sender_id: string;
		content: string;
		read_at?: string;
		created_at: string;
	}

	interface OtherUser {
		user_id: string;
		age: number;
		bio: string;
		city: string;
		state: string;
		is_online: boolean;
		images: Array<{ url: string; is_primary: boolean }>;
	}

	let messages = $state<Message[]>([]);
	let loading = $state(true);
	let newMessage = $state('');
	let sending = $state(false);
	let messagesContainer: HTMLDivElement;
	let otherUser = $state<OtherUser | null>(null);

	let conversationId = $derived($page.params.id);

	let authState = $state<any>(null);
	auth.subscribe(s => authState = s);
	let currentUserId = $derived(authState?.user?.id);

	async function loadConversation() {
		try {
			const data = await api.getConversations() as { conversations: any[] };
			const conv = data.conversations?.find((c: any) => c.id === conversationId);
			if (conv?.other_user) {
				otherUser = conv.other_user;
			}
		} catch (e) {
			console.error('Failed to load conversation:', e);
		}
	}

	async function loadMessages() {
		loading = true;
		try {
			const data = await api.getMessages(conversationId) as { messages: Message[] };
			messages = (data.messages || []).reverse();
			scrollToBottom();
		} catch (e) {
			console.error('Failed to load messages:', e);
		} finally {
			loading = false;
		}
	}

	function getProfileImage(): string {
		const primary = otherUser?.images?.find(img => img.is_primary);
		if (primary) return primary.url;
		if (otherUser?.images?.length) return otherUser.images[0].url;
		return 'https://via.placeholder.com/40?text=?';
	}

	async function sendMessage() {
		if (!newMessage.trim() || sending) return;
		sending = true;
		try {
			const message = await api.sendMessage(conversationId, newMessage.trim()) as Message;
			messages = [...messages, message];
			newMessage = '';
			scrollToBottom();
		} catch (e) {
			console.error('Failed to send message:', e);
		} finally {
			sending = false;
		}
	}

	function scrollToBottom() {
		setTimeout(() => {
			if (messagesContainer) {
				messagesContainer.scrollTop = messagesContainer.scrollHeight;
			}
		}, 100);
	}

	function formatDate(dateStr: string): string {
		const date = new Date(dateStr);
		const today = new Date();
		const yesterday = new Date(today);
		yesterday.setDate(yesterday.getDate() - 1);

		if (date.toDateString() === today.toDateString()) return 'Today';
		if (date.toDateString() === yesterday.toDateString()) return 'Yesterday';
		return date.toLocaleDateString();
	}

	function shouldShowDate(index: number): boolean {
		if (index === 0) return true;
		const current = new Date(messages[index].created_at).toDateString();
		const previous = new Date(messages[index - 1].created_at).toDateString();
		return current !== previous;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			sendMessage();
		}
	}

	// WebSocket message handler
	let unsubscribe: (() => void) | null = null;

	onMount(() => {
		loadConversation();
		loadMessages();
		
		// Listen for new messages via WebSocket
		unsubscribe = websocket.onMessage('message', (message: Message) => {
			if (message.conversation_id === conversationId) {
				messages = [...messages, message];
				scrollToBottom();
			}
		});
	});

	onDestroy(() => {
		if (unsubscribe) unsubscribe();
	});
</script>

<svelte:head>
	<title>Chat | HeySpoilMe</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Montserrat:wght@300;400;500;600&display=swap" rel="stylesheet">
</svelte:head>

<div class="chat-wrapper">
	<header class="header">
		<a href="/messages" class="back-btn">←</a>
		{#if otherUser}
			<div class="header-profile">
				<img src={getProfileImage()} alt="Profile" class="header-avatar" />
				<div class="header-info">
					<span class="header-name">{otherUser.age}, {otherUser.city}</span>
					<span class="header-status" class:online={otherUser.is_online}>
						{otherUser.is_online ? 'Online' : 'Offline'}
					</span>
				</div>
			</div>
		{:else}
			<span class="title">Conversation</span>
		{/if}
		<a href="/profile/{otherUser?.user_id || ''}" class="view-profile-btn">View</a>
	</header>

	<div class="chat-page">
		<div class="messages-container" bind:this={messagesContainer}>
		{#if loading}
			<div class="loading">
				<div class="spinner"></div>
			</div>
		{:else if messages.length === 0}
			<div class="empty">
				<p>No messages yet</p>
				<p class="hint">Send a message to start the conversation</p>
			</div>
		{:else}
			{#each messages as message, index}
				{#if shouldShowDate(index)}
					<div class="date-divider">
						<span>{formatDate(message.created_at)}</span>
					</div>
				{/if}
				<div 
					class="message" 
					class:sent={message.sender_id === currentUserId}
					class:received={message.sender_id !== currentUserId}
				>
					<div class="bubble">
						<p>{message.content}</p>
					</div>
				</div>
			{/each}
		{/if}
		</div>

		<div class="input-area">
			<textarea 
				bind:value={newMessage}
				placeholder="Type a message..."
				onkeydown={handleKeydown}
				rows="1"
			></textarea>
			<button 
				class="send-btn" 
				onclick={sendMessage}
				disabled={!newMessage.trim() || sending}
			>
				{sending ? '...' : '→'}
			</button>
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

	.chat-wrapper {
		display: flex;
		flex-direction: column;
		height: 100vh;
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1rem 2rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
		flex-shrink: 0;
		gap: 1rem;
	}

	.chat-page {
		display: flex;
		flex-direction: column;
		flex: 1;
		max-width: 600px;
		margin: 0 auto;
		width: 100%;
		overflow: hidden;
	}

	.back-btn {
		color: #fff;
		text-decoration: none;
		font-size: 1.2rem;
		padding: 0.5rem;
	}

	.title {
		font-weight: 600;
	}

	.header-profile {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		flex: 1;
	}

	.header-avatar {
		width: 40px;
		height: 40px;
		object-fit: cover;
	}

	.header-info {
		display: flex;
		flex-direction: column;
	}

	.header-name {
		font-weight: 600;
		font-size: 0.95rem;
	}

	.header-status {
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.5);
	}

	.header-status.online {
		color: #22c55e;
	}

	.view-profile-btn {
		color: rgba(255, 255, 255, 0.6);
		text-decoration: none;
		font-size: 0.85rem;
		padding: 0.5rem 1rem;
		border: 1px solid rgba(255, 255, 255, 0.2);
	}

	.view-profile-btn:hover {
		color: #fff;
		border-color: rgba(255, 255, 255, 0.4);
	}

	.messages-container {
		flex: 1;
		overflow-y: auto;
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.loading, .empty {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		flex: 1;
		color: rgba(255, 255, 255, 0.5);
	}

	.spinner {
		width: 32px;
		height: 32px;
		border: 2px solid rgba(255, 255, 255, 0.1);
		border-top-color: #fff;
		border-radius: 0;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.hint {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.4);
	}

	.date-divider {
		text-align: center;
		margin: 1rem 0;
	}

	.date-divider span {
		background: rgba(255, 255, 255, 0.1);
		padding: 0.25rem 0.75rem;
		border-radius: 0;
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.5);
	}

	.message {
		display: flex;
		max-width: 80%;
	}

	.message.sent {
		align-self: flex-end;
	}

	.message.received {
		align-self: flex-start;
	}

	.bubble {
		padding: 0.75rem 1rem;
		border-radius: 0;
		max-width: 100%;
	}

	.sent .bubble {
		background: #fff;
		color: #000;
		border-bottom-right-radius: 4px;
	}

	.received .bubble {
		background: rgba(255, 255, 255, 0.1);
		border-bottom-left-radius: 4px;
	}

	.bubble p {
		margin: 0;
		line-height: 1.4;
		word-wrap: break-word;
	}

	.input-area {
		display: flex;
		gap: 0.5rem;
		padding: 1rem;
		border-top: 1px solid rgba(255, 255, 255, 0.1);
		background: #0a0a0a;
		flex-shrink: 0;
		position: sticky;
		bottom: 0;
	}

	@media (max-width: 768px) {
		.header {
			padding: 1rem;
		}

		.input-area {
			position: fixed;
			bottom: 0;
			left: 0;
			right: 0;
			padding: 0.75rem;
			padding-bottom: calc(0.75rem + env(safe-area-inset-bottom));
		}

		.messages-container {
			padding-bottom: 80px;
		}
	}

	.input-area textarea {
		flex: 1;
		padding: 0.75rem 1rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0;
		color: #fff;
		font-family: 'Montserrat', sans-serif;
		font-size: 0.95rem;
		resize: none;
		max-height: 120px;
	}

	.input-area textarea:focus {
		outline: none;
		border-color: rgba(255, 255, 255, 0.2);
	}

	.send-btn {
		width: 48px;
		height: 48px;
		background: #fff;
		border: none;
		border-radius: 0;
		color: #000;
		font-size: 1.2rem;
		cursor: pointer;
		flex-shrink: 0;
		transition: all 0.2s ease;
	}

	.send-btn:disabled {
		opacity: 0.3;
		cursor: not-allowed;
	}

	.send-btn:not(:disabled):hover {
		transform: scale(1.05);
	}
</style>

