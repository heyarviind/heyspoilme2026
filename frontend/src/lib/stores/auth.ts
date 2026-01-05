import { writable } from 'svelte/store';
import api from '$lib/api';

export type WealthStatus = 'none' | 'low' | 'medium' | 'high';

interface User {
	id: string;
	email: string;
	google_id?: string;
	email_verified: boolean;
	wealth_status: WealthStatus;
	wealth_status_expires_at?: string;
	created_at: string;
}

interface Profile {
	id: string;
	user_id: string;
	display_name: string;
	gender: 'male' | 'female';
	age: number;
	bio: string;
	salary_range?: string;
	city: string;
	state: string;
	is_complete: boolean;
	is_verified: boolean; // person_verified
}

// Utility function to get wealth status label
export function getWealthStatusLabel(status: WealthStatus): string {
	switch (status) {
		case 'low': return 'Trusted';
		case 'medium': return 'Premium';
		case 'high': return 'Elite';
		default: return 'Standard';
	}
}

// Check if user can send messages (male with wealth_status != 'none' or female)
export function canSendMessages(user: User | null, profile: Profile | null): boolean {
	if (!user || !profile) return false;
	if (profile.gender === 'female') return profile.is_verified;
	// Male needs wealth_status != 'none' AND person_verified
	return profile.is_verified && user.wealth_status !== 'none';
}

// Check if user can view all messages
export function canViewMessages(user: User | null, profile: Profile | null): boolean {
	if (!user || !profile) return false;
	if (profile.gender === 'female') return true;
	return user.wealth_status !== 'none';
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

		setEmailVerified(verified: boolean) {
			update(state => {
				if (state.user) {
					return { ...state, user: { ...state.user, email_verified: verified } };
				}
				return state;
			});
		},

		async refreshUser() {
			try {
				const data = await api.getCurrentUser() as { user: User; profile: Profile | null };
				update(state => ({ ...state, user: data.user, profile: data.profile }));
			} catch {
				// Ignore errors
			}
		},
	};
}

export const auth = createAuthStore();
