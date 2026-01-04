import { writable } from 'svelte/store';
import api from '$lib/api';

interface User {
	id: string;
	email: string;
	google_id?: string;
	created_at: string;
}

interface Profile {
	id: string;
	user_id: string;
	gender: 'male' | 'female';
	age: number;
	bio: string;
	salary_range?: string;
	city: string;
	state: string;
	is_complete: boolean;
}

interface AuthState {
	user: User | null;
	profile: Profile | null;
	loading: boolean;
	initialized: boolean;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		profile: null,
		loading: true,
		initialized: false,
	});

	return {
		subscribe,
		
		async init() {
			// Check for token
			if (typeof window === 'undefined') return;
			
			const token = localStorage.getItem('token');
			if (!token) {
				set({ user: null, profile: null, loading: false, initialized: true });
				return;
			}

			update(state => ({ ...state, loading: true }));
			try {
				const data = await api.getCurrentUser() as { user: User; profile: Profile | null };
				set({ 
					user: data.user, 
					profile: data.profile, 
					loading: false, 
					initialized: true 
				});
			} catch {
				// Token invalid, clear it
				localStorage.removeItem('token');
				localStorage.removeItem('user');
				set({ user: null, profile: null, loading: false, initialized: true });
			}
		},

		async refreshProfile() {
			try {
				const profileData = await api.getMyProfile() as { profile: Profile };
				update(state => ({ ...state, profile: profileData.profile }));
			} catch {
				// Profile might not exist yet
			}
		},

		logout() {
			api.logout().catch(() => {});
			localStorage.removeItem('token');
			localStorage.removeItem('user');
			set({ user: null, profile: null, loading: false, initialized: true });
		},

		setUser(user: User) {
			update(state => ({ ...state, user }));
		},

		setProfile(profile: Profile) {
			update(state => ({ ...state, profile }));
		},
	};
}

export const auth = createAuthStore();
