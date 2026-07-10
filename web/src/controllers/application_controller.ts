import type {
	CreateApplicationDto,
	DeploymentStatusDto,
	LogEntryDto,
	UpdateApplicationDto,
} from '../dtos/application_dtos';
import type { LogEntry, LogLevel } from '../lib/mock_telemetry';
import ApplicationService from '../services/application_service';

const LOG_LEVELS: LogLevel[] = ['info', 'warn', 'error', 'debug'];

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
