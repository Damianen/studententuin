import { useEffect, useRef } from 'react';

const PROVISION_POLL_MS = 2500;

/**
 * Refetches while a database is provisioning, so the status badge flips and
 * the connection string appears without a manual reload (phase 5). The reload
 * callback is kept in a ref so parent re-renders never restart the interval —
 * the same trick use_deployment.ts uses for its callbacks.
 */
export function useProvisionWatch(
	status: string | undefined,
	reload: () => Promise<void>,
) {
	const reloadRef = useRef(reload);
	useEffect(() => {
		reloadRef.current = reload;
	}, [reload]);
	const provisioning = status?.toLowerCase() === 'provisioning';

	useEffect(() => {
		if (!provisioning) return;
		const timer = setInterval(() => {
			void reloadRef.current();
		}, PROVISION_POLL_MS);
		return () => clearInterval(timer);
	}, [provisioning]);
}
