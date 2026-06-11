import type {
	CreateDatabaseDto,
	UpdateDatabaseDto,
} from '../dtos/database_dtos';
import DatabaseService from '../services/database_service';

class DatabaseController {
	public static async get(subdomainId: string, dbId: string) {
		const response = await DatabaseService.get(subdomainId, dbId);
		if (response.code != 200 || !response.data) {
			throw new Error(response.message);
		}
		return response.data;
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
