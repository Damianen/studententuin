import type {
	CreateDatabaseDto,
	DatabaseDto,
	UpdateDatabaseDto,
} from '../dtos/database_dtos';
import type { LogEntryDto } from '../dtos/application_dtos';
import ApiService from './api_service';

class DatabaseService {
	public static async get(subdomainId: string, dbId: string) {
		return await ApiService.get<DatabaseDto>(
			`/subdomain/${subdomainId}/database/${dbId}`,
		);
	}

	public static async getLogs(subdomainId: string, dbId: string) {
		return await ApiService.get<LogEntryDto[]>(
			`/subdomain/${subdomainId}/database/${dbId}/logs`,
		);
	}

	public static async post(subdomainId: string, database: CreateDatabaseDto) {
		return await ApiService.post<void>(
			`/subdomain/${subdomainId}/database`,
			database,
		);
	}

	public static async patch(
		subdomainId: string,
		dbId: string,
		values: UpdateDatabaseDto,
	) {
		return await ApiService.patch<void>(
			`/subdomain/${subdomainId}/database/${dbId}`,
			values,
		);
	}

	public static async delete(subdomainId: string, dbId: string) {
		return await ApiService.delete<void>(
			`/subdomain/${subdomainId}/database/${dbId}`,
		);
	}
}

export default DatabaseService;
