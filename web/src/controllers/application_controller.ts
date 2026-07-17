import type {
	CreateApplicationDto,
	DeploymentRecordDto,
	DeploymentStatusDto,
	LogEntryDto,
	UpdateApplicationDto,
} from '../dtos/application_dtos';
import type { MetricPointDto, MetricsDto } from '../dtos/metrics_dtos';
import type {
	DeploymentRecord,
	DeploymentStatus,
	LogEntry,
	LogLevel,
	MetricPoint,
} from '../lib/mock_telemetry';
import { formatDuration, formatTimeAgo } from '../lib/format_time';
import ApplicationService from '../services/application_service';

const LOG_LEVELS: LogLevel[] = ['info', 'warn', 'error', 'debug'];

// Series the api serves for applications today; resp/req arrive with the
// edge proxy (phase 7).
const APP_METRIC_KEYS = ['cpu', 'mem'];

const DEPLOYMENT_STATUSES: Record<string, DeploymentStatus> = {
	succeeded: 'success',
	failed: 'failed',
	in_flight: 'in_progress',
};

// Shared by the one-shot fetch and the live websocket tail.
export function mapLogEntry(entry: LogEntryDto): LogEntry {
	return {
		id: entry.id,
		timestamp: entry.timestamp
			? new Date(entry.timestamp).toLocaleTimeString([], { hour12: false })
			: '--:--:--',
		level: (LOG_LEVELS as string[]).includes(entry.level)
			? (entry.level as LogLevel)
			: 'info',
		message: entry.message,
	};
}

// The api hands the manager's RFC3339 UTC times through untouched; chart
// labels are shortened browser-local, same as mapLogEntry.
export function mapMetricPoint(point: MetricPointDto): MetricPoint {
	return {
		time: new Date(point.time).toLocaleTimeString([], {
			hour: '2-digit',
			minute: '2-digit',
			hour12: false,
		}),
		value: point.value,
	};
}

// Every requested key comes back as an array — the api omits series that
// have no samples yet. Shared with the database controller.
export function mapMetrics(
	dto: MetricsDto | undefined,
	keys: string[],
): Record<string, MetricPoint[]> {
	const series: Record<string, MetricPoint[]> = {};
	for (const key of keys) {
		series[key] = (dto?.series?.[key] ?? []).map(mapMetricPoint);
	}
	return series;
}

export function mapDeployment(record: DeploymentRecordDto): DeploymentRecord {
	const status = DEPLOYMENT_STATUSES[record.status] ?? 'in_progress';
	return {
		id: record.id,
		status,
		// A failed clone never learns its commit.
		commit: record.commit_sha ? record.commit_sha.slice(0, 7) : '—',
		branch: record.branch || '—',
		message:
			record.commit_message ||
			record.error ||
			(status === 'in_progress' ? 'Deploy in progress' : 'Deployment'),
		author: record.commit_author || '—',
		timeAgo: formatTimeAgo(record.started_at),
		duration:
			record.duration_seconds != null
				? formatDuration(record.duration_seconds)
				: '—',
		durationSeconds: record.duration_seconds,
	};
}

class ApplicationController {
	public static async get(subdomainId: string, appId: string) {
		const response = await ApplicationService.get(subdomainId, appId);
		if (response.code != 200 || !response.data) {
			throw new Error(response.message);
		}
		return response.data;
	}

	public static async getLogs(
		subdomainId: string,
		appId: string,
	): Promise<LogEntry[]> {
		const response = await ApplicationService.getLogs(subdomainId, appId);
		if (response.code != 200) {
			throw new Error(response.message);
		}
		// The api omits `data` entirely when there are no log lines.
		return (response.data ?? []).map(mapLogEntry);
	}

	public static async getMetrics(
		subdomainId: string,
		appId: string,
	): Promise<Record<string, MetricPoint[]>> {
		const response = await ApplicationService.getMetrics(subdomainId, appId);
		if (response.code != 200) {
			throw new Error(response.message);
		}
		return mapMetrics(response.data, APP_METRIC_KEYS);
	}

	public static async create(
		subdomainId: string,
		values: CreateApplicationDto,
	) {
		const response = await ApplicationService.post(subdomainId, values);
		if (response.code != 201) {
			throw new Error(response.message);
		}
		return response.data;
	}

	public static async patch(
		subdomainId: string,
		appId: string,
		values: UpdateApplicationDto,
	) {
		const response = await ApplicationService.patch(
			subdomainId,
			appId,
			values,
		);
		if (response.code != 200) {
			throw new Error(response.message);
		}
		return response.data;
	}

	public static async delete(subdomainId: string, appId: string) {
		const response = await ApplicationService.delete(subdomainId, appId);
		if (response.code != 200) {
			throw new Error(response.message);
		}
	}

	// Queues a deploy of the stored repository; resolves to the deployment id
	// to poll with getDeployment.
	public static async deploy(
		subdomainId: string,
		appId: string,
	): Promise<string> {
		const response = await ApplicationService.deploy(subdomainId, appId);
		if (response.code != 202 || !response.data) {
			throw new Error(response.message);
		}
		return response.data.deployment_id;
	}

	public static async getDeployment(
		subdomainId: string,
		appId: string,
		deploymentId: string,
	): Promise<DeploymentStatusDto> {
		const response = await ApplicationService.getDeployment(
			subdomainId,
			appId,
			deploymentId,
		);
		if (response.code != 200 || !response.data) {
			throw new Error(response.message);
		}
		return response.data;
	}

	public static async getDeployments(
		subdomainId: string,
		appId: string,
	): Promise<DeploymentRecord[]> {
		const response = await ApplicationService.getDeployments(
			subdomainId,
			appId,
		);
		if (response.code != 200) {
			throw new Error(response.message);
		}
		// A never-deployed application has no history; the api may omit `data`.
		return (response.data ?? []).map(mapDeployment);
	}

	public static async start(subdomainId: string, appId: string) {
		const response = await ApplicationService.start(subdomainId, appId);
		if (response.code != 200) {
			throw new Error(response.message);
		}
	}

	public static async stop(subdomainId: string, appId: string) {
		const response = await ApplicationService.stop(subdomainId, appId);
		if (response.code != 200) {
			throw new Error(response.message);
		}
	}
}

export default ApplicationController;
