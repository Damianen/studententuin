import { renderHook, waitFor } from '@testing-library/react';
import { act } from 'react';
import { useAsync } from './use_async';

describe('useAsync', () => {
	it('starts loading, then exposes the resolved data', async () => {
		const fn = vi.fn().mockResolvedValue(['item']);
		const { result } = renderHook(() => useAsync(fn));

		expect(result.current.loading).toBe(true);
		expect(result.current.data).toBeNull();

		await waitFor(() => expect(result.current.loading).toBe(false));
		expect(result.current.data).toEqual(['item']);
		expect(result.current.error).toBeNull();
	});

	it('exposes the error message on rejection', async () => {
		const fn = vi.fn().mockRejectedValue(new Error('boom'));
		const { result } = renderHook(() => useAsync(fn));

		await waitFor(() => expect(result.current.loading).toBe(false));
		expect(result.current.error).toBe('boom');
		expect(result.current.data).toBeNull();
	});

	it('falls back to a generic message for non-Error rejections', async () => {
		const fn = vi.fn().mockRejectedValue('string failure');
		const { result } = renderHook(() => useAsync(fn));

		await waitFor(() => expect(result.current.loading).toBe(false));
		expect(result.current.error).toBe('Something went wrong');
	});

	it('reload refetches and clears a previous error', async () => {
		const fn = vi
			.fn()
			.mockRejectedValueOnce(new Error('boom'))
			.mockResolvedValueOnce(['recovered']);
		const { result } = renderHook(() => useAsync(fn));

		await waitFor(() => expect(result.current.error).toBe('boom'));

		await act(async () => {
			await result.current.reload();
		});

		expect(fn).toHaveBeenCalledTimes(2);
		expect(result.current.data).toEqual(['recovered']);
		expect(result.current.error).toBeNull();
	});
});
