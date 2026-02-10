import type { ApiResponse } from '../dtos/api_dtos';

function getApiBaseUrl(): string {
	if (
		window.location.hostname !== 'localhost' &&
		window.location.hostname !== '127.0.0.1'
	) {
		return 'https://studententuin.com/api';
	}

	return 'http://localhost:8080';
}

class ApiService {
	private static async request<T>(
		endpoint: string,
		options?: RequestInit
	): Promise<ApiResponse<T>> {
		const url = `${getApiBaseUrl()}${endpoint}`;

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

		const data = await response.json();

		return data;
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
