import type { Context } from 'hono';
import type { ContentfulStatusCode } from 'hono/utils/http-status';

export interface ApiResponse<T = any> {
	success: boolean;
	data?: T;
	error?: {
		code: string;
		message: string;
		details?: string[];
	};
}

export function ResSuccess<T>(c: Context, data: T, statusCode: number = 200): Response {
	const response: ApiResponse<T> = {
		success: true,
		data,
	};

	return c.json(response, statusCode as ContentfulStatusCode);
}

export function ResError(
	c: Context,
	code: string,
	message: string,
	statusCode: number = 400,
	details?: string[]
): Response {
	const response: ApiResponse = {
		success: false,
		error: {
			code,
			message,
			...(details && { details }),
		},
	};

	return c.json(response, statusCode as ContentfulStatusCode);
}
