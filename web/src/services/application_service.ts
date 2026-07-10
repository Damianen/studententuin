import type {
	ApplicationDto,
	CreateApplicationDto,
	DeploymentStatusDto,
	DeployResponseDto,
	LogEntryDto,
	UpdateApplicationDto,
} from '../dtos/application_dtos';
import ApiService from './api_service';

class ApplicationService {
	public static async get(subdomainId: string, appId: string) {
		return await ApiService.get<ApplicationDto>(
			`/subdomain/${subdomainId}/application/${appId}`,
		);
	}

	public static async getLogs(subdomainId: string, appId: string) {
		return await ApiService.get<LogEntryDto[]>(
			`/subdomain/${subdomainId}/application/${appId}/logs`,
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

	public static async deploy(subdomainId: string, appId: string) {
		return await ApiService.post<DeployResponseDto>(
			`/subdomain/${subdomainId}/application/${appId}/deploy`,
			{},
		);
	}

	public static async getDeployment(
		subdomainId: string,
		appId: string,
		deploymentId: string,
	) {
		return await ApiService.get<DeploymentStatusDto>(
			`/subdomain/${subdomainId}/application/${appId}/deployment/${deploymentId}`,
		);
	}

	public static async start(subdomainId: string, appId: string) {
		return await ApiService.post<void>(
			`/subdomain/${subdomainId}/application/${appId}/start`,
			{},
		);
	}

	public static async stop(subdomainId: string, appId: string) {
		return await ApiService.post<void>(
			`/subdomain/${subdomainId}/application/${appId}/stop`,
			{},
		);
	}
}

export default ApplicationService;
