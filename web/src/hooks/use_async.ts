import { useCallback, useEffect, useState } from 'react';

interface AsyncState<T> {
	data: T | null;
	loading: boolean;
	error: string | null;
	reload: () => Promise<void>;
}

export function useAsync<T>(fn: () => Promise<T>): AsyncState<T> {
	const [data, setData] = useState<T | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	const reload = useCallback(async () => {
		try {
			const result = await fn();
			setData(result);
			setError(null);
		} catch (err) {
			setError(err instanceof Error ? err.message : 'Something went wrong');
		} finally {
			setLoading(false);
		}
		// fn is expected to be stable (controller static method or wrapped in useCallback)
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	useEffect(() => {
		// Intentional fetch-on-mount; all setState calls happen after an await.
		// eslint-disable-next-line react-hooks/set-state-in-effect
		reload();
	}, [reload]);

	return { data, loading, error, reload };
}
