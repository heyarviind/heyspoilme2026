import { writable } from 'svelte/store';
import api from '$lib/api';

interface Notification {
	id: string;
	user_id: string;
	type: 'new_like' | 'new_message' | 'profile_view';
	data: any;
	is_read: boolean;
	created_at: string;
}

interface NotificationState {
	notifications: Notification[];
	unreadCount: number;
	loading: boolean;
}

function createNotificationStore() {
	const { subscribe, set, update } = writable<NotificationState>({
		notifications: [],
		unreadCount: 0,
		loading: false,
	});

	return {
		subscribe,

		async fetchNotifications() {
			update(state => ({ ...state, loading: true }));
			try {
				const data = await api.getNotifications() as { notifications: Notification[]; total: number };
				update(state => ({ ...state, notifications: data.notifications || [], loading: false }));
			} catch {
				update(state => ({ ...state, loading: false }));
			}
		},

		async fetchUnreadCount() {
			try {
				const data = await api.getUnreadNotificationCount() as { count: number };
				update(state => ({ ...state, unreadCount: data.count }));
			} catch {
				// Ignore errors
			}
		},

		async markAsRead(id: string) {
			try {
				await api.markNotificationAsRead(id);
				update(state => ({
					...state,
					notifications: state.notifications.map(n => 
						n.id === id ? { ...n, is_read: true } : n
					),
					unreadCount: Math.max(0, state.unreadCount - 1),
				}));
			} catch {
				// Ignore errors
			}
		},

		async markAllAsRead() {
			try {
				await api.markAllNotificationsAsRead();
				update(state => ({
					...state,
					notifications: state.notifications.map(n => ({ ...n, is_read: true })),
					unreadCount: 0,
				}));
			} catch {
				// Ignore errors
			}
		},

		addNotification(notification: Notification) {
			update(state => ({
				...state,
				notifications: [notification, ...state.notifications],
				unreadCount: state.unreadCount + 1,
			}));
		},
	};
}

export const notifications = createNotificationStore();


