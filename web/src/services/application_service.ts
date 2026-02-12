import type {
	ApplicationDto,
	CreateApplicationDto,
	UpdateApplicationDto,
} from '../dtos/application_dtos';
import ApiService from './api_service';

class ApplicationService {
	public static async get(subdomainId: string, appId: string) {
		return await ApiService.get<ApplicationDto>(
			`/subdomain/${subdomainId}/application/${appId}`,
		);
	}

	public static async post(
		subdomainId: string,
		application: CreateApplicationDto,
	) {
		return await ApiService.post<void>(
			`/subdomain/${subdomainId}/application`,
			application,
		);
	}

	public static async patch(
		subdomainId: string,
		appId: string,
		values: UpdateApplicationDto,
	) {
		return await ApiService.patch<void>(
			`/subdomain/${subdomainId}/application/${appId}`,
			values,
		);
	}

	public static async delete(subdomainId: string, appId: string) {
		return await ApiService.delete<void>(
			`/subdomain/${subdomainId}/application/${appId}`,
		);
	}
}

export default ApplicationService;
