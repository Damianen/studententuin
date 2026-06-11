import type { ApiResponse } from '../dtos/api_dtos';

// Same-origin in dev (Vite proxy) and prod (reverse proxy at /api).
const API_BASE_URL = '/api';

type UnauthorizedHandler = () => void;

let onUnauthorized: UnauthorizedHandler | null = null;

export function setOnUnauthorized(handler: UnauthorizedHandler | null) {
	onUnauthorized = handler;
}

class ApiService {
	private static async request<T>(
		endpoint: string,
		options?: RequestInit
	): Promise<ApiResponse<T>> {
		const url = `${API_BASE_URL}${endpoint}`;

		const headers: Record<string, string> = {
			'Content-Type': 'application/json',
		};

		if (options?.headers) {
			Object.assign(headers, options.headers);
		}

		const response = await fetch(url, {
			...options,
			headers,
			credentials: 'include',
		});

		if (response.status === 401 && endpoint !== '/auth/login') {
			onUnauthorized?.();
		}

		const text = await response.text();
		try {
			return JSON.parse(text);
		} catch {
			return {
				code: response.status,
				message: response.statusText || 'unexpected server response',
				data: undefined as T,
			};
		}
	}

	public static async get<T>(endpoint: string): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, { method: 'GET' });
	}

	public static async post<T>(
		endpoint: string,
		data?: unknown
	): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, {
			method: 'POST',
			body: JSON.stringify(data),
		});
	}

	public static async patch<T>(
		endpoint: string,
		data?: unknown
	): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, {
			method: 'PATCH',
			body: JSON.stringify(data),
		});
	}

	public static async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, { method: 'DELETE' });
	}
}

export default ApiService;
