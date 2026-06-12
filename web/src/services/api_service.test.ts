import ApiService, { setOnUnauthorized } from './api_service';

function jsonResponse(status: number, body: unknown, statusText = '') {
	return {
		status,
		statusText,
		text: async () => JSON.stringify(body),
	};
}

const fetchMock = vi.fn();

beforeEach(() => {
	vi.stubGlobal('fetch', fetchMock);
});

afterEach(() => {
	fetchMock.mockReset();
	setOnUnauthorized(null);
	vi.unstubAllGlobals();
});

describe('ApiService', () => {
	it('prefixes endpoints with /api and sends credentials + JSON header', async () => {
		fetchMock.mockResolvedValue(jsonResponse(200, { code: 200, message: 'ok' }));

		await ApiService.get('/user');

		expect(fetchMock).toHaveBeenCalledWith('/api/user', {
			method: 'GET',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
		});
	});

	it('serializes the body for POST and returns the parsed envelope', async () => {
		fetchMock.mockResolvedValue(
			jsonResponse(201, { code: 201, message: 'success' })
		);

		const result = await ApiService.post('/subdomain', { name: 'myapp' });

		expect(fetchMock).toHaveBeenCalledWith(
			'/api/subdomain',
			expect.objectContaining({
				method: 'POST',
				body: JSON.stringify({ name: 'myapp' }),
			})
		);
		expect(result).toEqual({ code: 201, message: 'success' });
	});

	it('sends PATCH and DELETE with the right methods', async () => {
		fetchMock.mockResolvedValue(jsonResponse(200, { code: 200, message: 'ok' }));

		await ApiService.patch('/user', { name: 'x' });
		await ApiService.delete('/user');

		expect(fetchMock).toHaveBeenNthCalledWith(
			1,
			'/api/user',
			expect.objectContaining({ method: 'PATCH' })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			2,
			'/api/user',
			expect.objectContaining({ method: 'DELETE' })
		);
	});

	it('falls back to a synthetic envelope when the body is not JSON', async () => {
		fetchMock.mockResolvedValue({
			status: 502,
			statusText: 'Bad Gateway',
			text: async () => '<html>nginx error</html>',
		});

		const result = await ApiService.get('/user');

		expect(result).toEqual({
			code: 502,
			message: 'Bad Gateway',
			data: undefined,
		});
	});

	it('uses a generic message when statusText is empty', async () => {
		fetchMock.mockResolvedValue({
			status: 500,
			statusText: '',
			text: async () => 'not json',
		});

		const result = await ApiService.get('/user');

		expect(result.message).toBe('unexpected server response');
	});

	it('notifies the unauthorized handler on 401', async () => {
		const onUnauthorized = vi.fn();
		setOnUnauthorized(onUnauthorized);
		fetchMock.mockResolvedValue(
			jsonResponse(401, { code: 401, message: 'authentication required' })
		);

		await ApiService.get('/user');

		expect(onUnauthorized).toHaveBeenCalledTimes(1);
	});

	it('does not notify the unauthorized handler for /auth/login', async () => {
		const onUnauthorized = vi.fn();
		setOnUnauthorized(onUnauthorized);
		fetchMock.mockResolvedValue(
			jsonResponse(401, { code: 401, message: 'email or password not correct!' })
		);

		await ApiService.post('/auth/login', { email: 'a@b.c', password: 'x' });

		expect(onUnauthorized).not.toHaveBeenCalled();
	});
});
