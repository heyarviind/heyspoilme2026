import { writable, get } from 'svelte/store';
import { notifications } from './notifications';

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';

interface WSState {
	connected: boolean;
	reconnecting: boolean;
}

function createWebSocketStore() {
	const { subscribe, set, update } = writable<WSState>({
		connected: false,
		reconnecting: false,
	});

	let ws: WebSocket | null = null;
	let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
	let heartbeatTimer: ReturnType<typeof setInterval> | null = null;

	const messageHandlers = new Map<string, (payload: any) => void>();

	function connect(token: string) {
		if (ws?.readyState === WebSocket.OPEN) return;

		ws = new WebSocket(`${WS_URL}?token=${token}`);

		ws.onopen = () => {
			set({ connected: true, reconnecting: false });
			startHeartbeat();
		};

		ws.onmessage = (event) => {
			try {
				const message = JSON.parse(event.data);
				handleMessage(message);
			} catch {
				console.error('Failed to parse WebSocket message');
			}
		};

		ws.onclose = () => {
			set({ connected: false, reconnecting: true });
			stopHeartbeat();
			scheduleReconnect(token);
		};

		ws.onerror = () => {
			ws?.close();
		};
	}

	function handleMessage(message: { type: string; payload: any }) {
		const handler = messageHandlers.get(message.type);
		if (handler) {
			handler(message.payload);
		}

		// Handle built-in message types
		switch (message.type) {
			case 'notification':
				notifications.addNotification(message.payload);
				break;
			case 'presence':
				// Could update a presence store here
				break;
		}
	}

	function startHeartbeat() {
		heartbeatTimer = setInterval(() => {
			if (ws?.readyState === WebSocket.OPEN) {
				ws.send(JSON.stringify({ type: 'ping' }));
			}
		}, 30000);
	}

	function stopHeartbeat() {
		if (heartbeatTimer) {
			clearInterval(heartbeatTimer);
			heartbeatTimer = null;
		}
	}

	function scheduleReconnect(token: string) {
		if (reconnectTimer) return;
		reconnectTimer = setTimeout(() => {
			reconnectTimer = null;
			connect(token);
		}, 3000);
	}

	function disconnect() {
		if (reconnectTimer) {
			clearTimeout(reconnectTimer);
			reconnectTimer = null;
		}
		stopHeartbeat();
		ws?.close();
		ws = null;
		set({ connected: false, reconnecting: false });
	}

	function send(type: string, payload: any) {
		if (ws?.readyState === WebSocket.OPEN) {
			ws.send(JSON.stringify({ type, payload }));
		}
	}

	function onMessage(type: string, handler: (payload: any) => void) {
		messageHandlers.set(type, handler);
		return () => messageHandlers.delete(type);
	}

	return {
		subscribe,
		connect,
		disconnect,
		send,
		onMessage,
	};
}

export const websocket = createWebSocketStore();

