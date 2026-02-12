import type {
	CreateSubdomainDto,
	SubdomainListItemDto,
	UpdateSubdomainDto,
} from '../dtos/subdomain_dtos';
import ApiService from './api_service';

class SubdomainService {
	public static async getAll() {
		return await ApiService.get<SubdomainListItemDto[]>('/subdomain');
	}

	public static async post(subdomain: CreateSubdomainDto) {
		return await ApiService.post<void>('/subdomain', subdomain);
	}

	public static async patch(id: string, values: UpdateSubdomainDto) {
		return await ApiService.patch<void>(`/subdomain/${id}`, values);
	}

	public static async delete(id: string) {
		return await ApiService.delete<void>(`/subdomain/${id}`);
	}
}

export default SubdomainService;
