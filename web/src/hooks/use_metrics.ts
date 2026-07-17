import { useCallback, useEffect, useRef } from 'react';
import type { MetricPoint, ResourceKind } from '@/lib/mock_telemetry';
import ApplicationController from '@/controllers/application_controller';
import DatabaseController from '@/controllers/database_controller';
import { useAsync } from './use_async';

const REFRESH_INTERVAL_MS = 60_000;

interface Metrics {
	/** every real key present as an array; {} until the first fetch lands */
	series: Record<string, MetricPoint[]>;
	loading: boolean;
	error: string | null;
}

/**
 * Fetches the real metric series (cpu/mem for applications, conn/qps/cpu/disk
 * for databases) and refreshes in the background — the manager samples every
 * ~30s, so a minute keeps charts current without hammering the api. The reload
 * rides in a ref so parent re-renders never restart the interval, the same
 * trick use_provision_watch.ts uses.
 */
export function useMetrics(
	kind: ResourceKind,
	subdomainId: string,
	resourceId: string,
): Metrics {
	const { data, loading, error, reload } = useAsync(
		useCallback(
			() =>
				kind === 'application'
					? ApplicationController.getMetrics(subdomainId, resourceId)
					: DatabaseController.getMetrics(subdomainId, resourceId),
			[kind, subdomainId, resourceId],
		),
	);

	const reloadRef = useRef(reload);
	useEffect(() => {
		reloadRef.current = reload;
	}, [reload]);

	useEffect(() => {
		const timer = setInterval(() => {
			void reloadRef.current();
		}, REFRESH_INTERVAL_MS);
		return () => clearInterval(timer);
	}, []);

	return { series: data ?? {}, loading, error };
}
