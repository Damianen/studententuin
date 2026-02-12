import type {
	CreateSubdomainDto,
	UpdateSubdomainDto,
} from '../dtos/subdomain_dtos';
import SubdomainService from '../services/subdomain_service';

class SubdomainController {
	public static async getAll() {
		const response = await SubdomainService.getAll();
		if (response.code != 200) {
			throw new Error(response.message);
		}
		return response.data;
	}

	public static async create(values: CreateSubdomainDto) {
		const response = await SubdomainService.post(values);
		if (response.code != 201) {
			throw new Error(response.message);
		}
		return response.data;
	}

	public static async patch(id: string, values: UpdateSubdomainDto) {
		const response = await SubdomainService.patch(id, values);
		if (response.code != 200) {
			throw new Error(response.message);
		}
		return response.data;
	}

	public static async delete(id: string) {
		const response = await SubdomainService.delete(id);
		if (response.code != 200) {
			throw new Error(response.message);
		}
	}
}

export default SubdomainController;
