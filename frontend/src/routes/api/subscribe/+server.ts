import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const API_URL = process.env.PUBLIC_API_URL || 'http://localhost:8081';

export const POST: RequestHandler = async ({ request }) => {
	try {
		const body = await request.json();
		
		const response = await fetch(`${API_URL}/api/subscribe`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(body),
		});
		
		const data = await response.json();
		return json(data, { status: response.status });
	} catch (error) {
		console.error('API proxy error:', error);
		return json(
			{ success: false, message: 'Unable to connect to server' },
			{ status: 500 }
		);
	}
};

