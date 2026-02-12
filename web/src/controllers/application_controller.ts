import type {
	CreateApplicationDto,
	UpdateApplicationDto,
} from '../dtos/application_dtos';
import ApplicationService from '../services/application_service';

class ApplicationController {
	public static async get(subdomainId: string, appId: string) {
		const response = await ApplicationService.get(subdomainId, appId);
		if (response.code != 200) {
			throw new Error(response.message);
		}
		return response.data;
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
}

export default ApplicationController;
