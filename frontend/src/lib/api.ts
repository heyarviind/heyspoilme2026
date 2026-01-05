const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

interface FetchOptions extends RequestInit {
	body?: any;
}

function getAuthHeader(): Record<string, string> {
	if (typeof window === 'undefined') return {};
	const token = localStorage.getItem('token');
	return token ? { Authorization: `Bearer ${token}` } : {};
}

async function fetchAPI<T>(endpoint: string, options: FetchOptions = {}): Promise<T> {
	const { body, ...rest } = options;
	
	const config: RequestInit = {
		...rest,
		headers: {
			'Content-Type': 'application/json',
			...getAuthHeader(),
			...rest.headers,
		},
		credentials: 'include',
	};

	if (body) {
		config.body = JSON.stringify(body);
	}

	const response = await fetch(`${API_BASE}${endpoint}`, config);
	
	if (!response.ok) {
		if (response.status === 401) {
			// Don't redirect for auth endpoints (login/signup failures should show error, not redirect)
			const isAuthEndpoint = endpoint.startsWith('/api/auth/signin') || endpoint.startsWith('/api/auth/signup');
			if (!isAuthEndpoint && typeof window !== 'undefined') {
				localStorage.removeItem('token');
				window.location.href = '/auth/login';
				throw new Error('Unauthorized');
			}
		}
		const error = await response.json().catch(() => ({ error: 'Request failed' }));
		throw new Error(error.error || error.message || 'Request failed');
	}

	return response.json();
}

export const api = {
	// Auth
	signup: (email: string, password: string) => 
		fetchAPI('/api/auth/signup', { method: 'POST', body: { email, password } }),
	signin: (email: string, password: string) => 
		fetchAPI('/api/auth/signin', { method: 'POST', body: { email, password } }),
	getCurrentUser: () => fetchAPI('/api/auth/me'),
	logout: () => fetchAPI('/api/auth/logout', { method: 'POST' }),
	getGoogleAuthUrl: () => `${API_BASE}/api/auth/google`,
	verifyEmail: (token: string) => fetchAPI(`/api/auth/verify-email?token=${token}`),
	resendVerificationEmail: () => fetchAPI('/api/auth/resend-verification', { method: 'POST' }),
	deleteAccount: () => fetchAPI('/api/auth/account', { method: 'DELETE' }),

	// Profile
	createProfile: (data: {
		display_name: string;
		gender: string;
		age: number;
		bio: string;
		salary_range?: string;
		city: string;
		state: string;
		latitude: number;
		longitude: number;
	}) => fetchAPI('/api/profile', { method: 'POST', body: data }),
	getMyProfile: () => fetchAPI('/api/profile'),
	updateProfile: (data: any) => fetchAPI('/api/profile', { method: 'PUT', body: data }),
	getProfile: (id: string) => fetchAPI(`/api/profiles/${id}`),
	listProfiles: (params?: Record<string, any>) => {
		const searchParams = new URLSearchParams();
		if (params) {
			Object.entries(params).forEach(([key, value]) => {
				if (value !== undefined && value !== '') {
					searchParams.append(key, String(value));
				}
			});
		}
		return fetchAPI(`/api/profiles?${searchParams.toString()}`);
	},

	// Profile Images
	getPresignedUrl: (fileExt: string, contentType: string) => 
		fetchAPI('/api/upload/presigned-url', { method: 'POST', body: { file_ext: fileExt, content_type: contentType } }),
	addProfileImage: (s3Key: string, url: string, isPrimary: boolean) =>
		fetchAPI('/api/profile/images', { method: 'POST', body: { s3_key: s3Key, url, is_primary: isPrimary } }),
	deleteProfileImage: (id: string) => fetchAPI(`/api/profile/images/${id}`, { method: 'DELETE' }),

	// Likes
	likeProfile: (id: string) => fetchAPI(`/api/profiles/${id}/like`, { method: 'POST' }),
	unlikeProfile: (id: string) => fetchAPI(`/api/profiles/${id}/like`, { method: 'DELETE' }),
	getReceivedLikes: (limit = 20, offset = 0) => fetchAPI(`/api/likes/received?limit=${limit}&offset=${offset}`),
	getGivenLikes: (limit = 20, offset = 0) => fetchAPI(`/api/likes/given?limit=${limit}&offset=${offset}`),

	// Conversations
	getConversations: () => fetchAPI('/api/conversations'),
	getInbox: () => fetchAPI('/api/inbox'), // Returns locked message support for males
	createConversation: (recipientId: string, message: string) => 
		fetchAPI('/api/conversations', { method: 'POST', body: { recipient_id: recipientId, message } }),
	getMessages: (conversationId: string, limit = 50, offset = 0) => 
		fetchAPI(`/api/conversations/${conversationId}/messages?limit=${limit}&offset=${offset}`),
	sendMessage: (conversationId: string, content: string, imageUrl?: string) =>
		fetchAPI(`/api/conversations/${conversationId}/messages`, { method: 'POST', body: { content, image_url: imageUrl } }),
	getUnreadMessageCount: () => fetchAPI('/api/messages/unread-count'),
	getChatImagePresignedUrl: (conversationId: string, fileExt: string, contentType: string) =>
		fetchAPI('/api/upload/chat-image-url', { method: 'POST', body: { conversation_id: conversationId, file_ext: fileExt, content_type: contentType } }),

	// Notifications
	getNotifications: (limit = 20, offset = 0) => fetchAPI(`/api/notifications?limit=${limit}&offset=${offset}`),
	markNotificationAsRead: (id: string) => fetchAPI(`/api/notifications/${id}/read`, { method: 'PUT' }),
	markAllNotificationsAsRead: () => fetchAPI('/api/notifications/read-all', { method: 'PUT' }),
	getUnreadNotificationCount: () => fetchAPI('/api/notifications/unread-count'),

	// Identity Verification
	getVerificationCode: () => fetchAPI('/api/verification/code'),
	submitVerification: (documentType: string, documentUrl: string, videoUrl: string, code: string) =>
		fetchAPI(`/api/verification/submit?code=${code}`, { 
			method: 'POST', 
			body: { document_type: documentType, document_url: documentUrl, video_url: videoUrl } 
		}),
	getVerificationStatus: () => fetchAPI('/api/verification/status'),
};

export default api;
