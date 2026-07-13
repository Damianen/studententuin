import { useEffect, useState } from 'react';
import type { LogEntry, ResourceKind } from '@/lib/mock_telemetry';
import type { LogEntryDto } from '@/dtos/application_dtos';
import ApplicationController, {
	mapLogEntry,
} from '@/controllers/application_controller';
import DatabaseController from '@/controllers/database_controller';

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
 * Logs arrive over a live websocket tail (history first, then lines as the
 * container writes them). When the stream is unavailable — container not
 * provisioned/deployed, handshake rejected, connection dropped — it degrades
 * to polling the one-shot endpoint. Applications and databases share the
 * machinery and differ only in endpoint paths (phase 2 / phase 5).
 */
export function useLogStream(
	kind: ResourceKind,
	subdomainId: string,
	resourceId: string,
): LogStream {
	const [logs, setLogs] = useState<LogEntry[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [live, setLive] = useState(false);

	useEffect(() => {
		let disposed = false;
		let socket: WebSocket | null = null;
		let pollTimer: ReturnType<typeof setInterval> | null = null;

		const segment = kind === 'application' ? 'application' : 'database';
		const getLogs =
			kind === 'application'
				? ApplicationController.getLogs
				: DatabaseController.getLogs;

		const fetchOnce = async () => {
			try {
				const entries = await getLogs(subdomainId, resourceId);
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
				`${scheme}://${window.location.host}/api/subdomain/${subdomainId}/${segment}/${resourceId}/logs/stream`,
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
				// Rejected handshake (e.g. nothing provisioned yet), network
				// drop, or the container's log stream ended — degrade to polling.
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
