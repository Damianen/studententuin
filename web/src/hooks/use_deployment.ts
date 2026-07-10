import { useCallback, useEffect, useRef, useState } from 'react';
import ApplicationController from '@/controllers/application_controller';

const POLL_INTERVAL_MS = 2500;
// Comfortably above the backend's clone+build budgets.
const POLL_BUDGET_MS = 15 * 60 * 1000;

const TERMINAL = new Set(['running', 'failed']);

export interface DeploymentState {
	/** true from the deploy click until the job reaches running/failed */
	deploying: boolean;
	/** queued | cloning | building | starting while in flight */
	stage: string | null;
	/** terminal failure reason */
	error: string | null;
	/** build-log tail attached to a failed deployment */
	buildLog: string | null;
}

const IDLE: DeploymentState = {
	deploying: false,
	stage: null,
	error: null,
	buildLog: null,
};

/**
 * Drives one deployment at a time: POST the deploy, then poll the job until
 * it settles. onSettled fires once per deploy with the outcome so the page
 * can refresh the application record (status badge, last-deployed).
 */
export function useDeployment(
	subdomainId: string,
	appId: string,
	onSettled?: (ok: boolean) => void,
) {
	const [state, setState] = useState<DeploymentState>(IDLE);
	const disposed = useRef(false);
	const timer = useRef<ReturnType<typeof setInterval> | null>(null);
	// The callback rides in a ref so a re-rendering parent doesn't restart
	// or double-fire an in-flight poll loop.
	const settled = useRef(onSettled);
	useEffect(() => {
		settled.current = onSettled;
	}, [onSettled]);

	const stopPolling = useCallback(() => {
		if (timer.current) {
			clearInterval(timer.current);
			timer.current = null;
		}
	}, []);

	useEffect(() => {
		disposed.current = false;
		return () => {
			disposed.current = true;
			stopPolling();
		};
	}, [stopPolling]);

	const deploy = useCallback(async () => {
		if (timer.current) return; // one at a time
		setState({ deploying: true, stage: 'queued', error: null, buildLog: null });

		let deploymentId: string;
		try {
			deploymentId = await ApplicationController.deploy(subdomainId, appId);
		} catch (err) {
			if (disposed.current) return;
			setState({
				deploying: false,
				stage: null,
				error: err instanceof Error ? err.message : 'deploy failed',
				buildLog: null,
			});
			settled.current?.(false);
			return;
		}

		const deadline = Date.now() + POLL_BUDGET_MS;
		timer.current = setInterval(async () => {
			try {
				const status = await ApplicationController.getDeployment(
					subdomainId,
					appId,
					deploymentId,
				);
				if (disposed.current) return;

				if (TERMINAL.has(status.status)) {
					stopPolling();
					const ok = status.status === 'running';
					setState({
						deploying: false,
						stage: null,
						error: ok ? null : (status.error ?? 'deployment failed'),
						buildLog: ok ? null : (status.build_log ?? null),
					});
					settled.current?.(ok);
					return;
				}
				setState((prev) => ({ ...prev, stage: status.status }));
			} catch {
				// Transient poll errors (api restart, brief network loss) are
				// retried until the budget runs out.
			}
			if (Date.now() > deadline && timer.current) {
				stopPolling();
				if (disposed.current) return;
				setState({
					deploying: false,
					stage: null,
					error: 'deployment timed out — check the logs tab',
					buildLog: null,
				});
				settled.current?.(false);
			}
		}, POLL_INTERVAL_MS);
	}, [subdomainId, appId, stopPolling]);

	return { ...state, deploy };
}
