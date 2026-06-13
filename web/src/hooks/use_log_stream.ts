import { useEffect, useState } from 'react';
import { makeLogs } from '@/lib/mock_telemetry';
import type { LogEntry, ResourceKind } from '@/lib/mock_telemetry';
import type { LogEntryDto } from '@/dtos/application_dtos';
import ApplicationController, {
	mapLogEntry,
} from '@/controllers/application_controller';

const POLL_INTERVAL_MS = 5000;
const MAX_BUFFERED_LINES = 2000;

interface LogStream {
	logs: LogEntry[];
	loading: boolean;
	error: string | null;
	/** true while the websocket tail is connected; false in poll fallback */
	live: boolean;
}

/**
 * Application logs arrive over a live websocket tail (history first, then
 * lines as the container writes them). When the stream is unavailable —
 * app not deployed, handshake rejected, connection dropped — it degrades
 * to polling the one-shot endpoint. Database logs stay on seeded sample
 * data until the backend grows database provisioning (phase 5).
 */
export function useLogStream(
	kind: ResourceKind,
	subdomainId: string,
	resourceId: string,
): LogStream {
	// Non-application kinds render their sample data immediately and never
	// touch the network, so their state is final from the first render.
	const isApplication = kind === 'application';
	const [logs, setLogs] = useState<LogEntry[]>(() =>
		isApplication ? [] : makeLogs(resourceId, kind),
	);
	const [loading, setLoading] = useState(isApplication);
	const [error, setError] = useState<string | null>(null);
	const [live, setLive] = useState(false);

	useEffect(() => {
		if (kind !== 'application') {
			return;
		}

		let disposed = false;
		let socket: WebSocket | null = null;
		let pollTimer: ReturnType<typeof setInterval> | null = null;

		const fetchOnce = async () => {
			try {
				const entries = await ApplicationController.getLogs(
					subdomainId,
					resourceId,
				);
				if (disposed) return;
				setLogs(entries);
				setError(null);
			} catch (err) {
				if (disposed) return;
				setError(
					err instanceof Error ? err.message : 'Something went wrong',
				);
			} finally {
				if (!disposed) setLoading(false);
			}
		};

		const startPolling = () => {
			if (disposed || pollTimer) return;
			setLive(false);
			fetchOnce();
			pollTimer = setInterval(fetchOnce, POLL_INTERVAL_MS);
		};

		try {
			const scheme = window.location.protocol === 'https:' ? 'wss' : 'ws';
			socket = new WebSocket(
				`${scheme}://${window.location.host}/api/subdomain/${subdomainId}/application/${resourceId}/logs/stream`,
			);
			socket.onopen = () => {
				if (disposed) return;
				setLive(true);
				setLoading(false);
				setError(null);
			};
			socket.onmessage = (event) => {
				if (disposed) return;
				try {
					const entry = mapLogEntry(
						JSON.parse(event.data as string) as LogEntryDto,
					);
					setLogs((current) =>
						[...current, entry].slice(-MAX_BUFFERED_LINES),
					);
				} catch {
					// Skip malformed frames rather than killing the tail.
				}
			};
			socket.onclose = () => {
				// Rejected handshake (e.g. app not deployed), network drop, or
				// the container's log stream ended — degrade to polling.
				if (!disposed) startPolling();
			};
		} catch {
			startPolling();
		}

		return () => {
			disposed = true;
			socket?.close();
			if (pollTimer) clearInterval(pollTimer);
		};
	}, [kind, subdomainId, resourceId]);

	return { logs, loading, error, live };
}
