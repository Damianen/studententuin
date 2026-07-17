import type {
	CreateDatabaseDto,
	UpdateDatabaseDto,
} from '../dtos/database_dtos';
import type { LogEntry, MetricPoint } from '../lib/mock_telemetry';
import { mapLogEntry, mapMetrics } from './application_controller';
import DatabaseService from '../services/database_service';

const DB_METRIC_KEYS = ['conn', 'qps', 'cpu', 'disk'];

class DatabaseController {
	public static async get(subdomainId: string, dbId: string) {
		const response = await DatabaseService.get(subdomainId, dbId);
		if (response.code != 200 || !response.data) {
			throw new Error(response.message);
		}
		return response.data;
	}

	public static async getLogs(
		subdomainId: string,
		dbId: string,
	): Promise<LogEntry[]> {
		const response = await DatabaseService.getLogs(subdomainId, dbId);
		if (response.code != 200) {
			throw new Error(response.message);
		}
		// The api omits `data` entirely when there are no log lines.
		return (response.data ?? []).map(mapLogEntry);
	}

	public static async getMetrics(
		subdomainId: string,
		dbId: string,
	): Promise<Record<string, MetricPoint[]>> {
		const response = await DatabaseService.getMetrics(subdomainId, dbId);
		if (response.code != 200) {
			throw new Error(response.message);
		}
		return mapMetrics(response.data, DB_METRIC_KEYS);
	}

	public static async create(
		subdomainId: string,
		values: CreateDatabaseDto,
	) {
		const response = await DatabaseService.post(subdomainId, values);
		if (response.code != 201) {
			throw new Error(response.message);
		}
		return response.data;
	}

	public static async patch(
		subdomainId: string,
		dbId: string,
		values: UpdateDatabaseDto,
	) {
		const response = await DatabaseService.patch(subdomainId, dbId, values);
		if (response.code != 200) {
			throw new Error(response.message);
		}
		return response.data;
	}

	public static async delete(subdomainId: string, dbId: string) {
		const response = await DatabaseService.delete(subdomainId, dbId);
		if (response.code != 200) {
			throw new Error(response.message);
		}
	}
}

export default DatabaseController;
