import type { LoginUserDto } from '../dtos/auth_dtos';
import ApiService from './api_service';

class AuthService {
	public static async login(data: LoginUserDto) {
		return await ApiService.post<null>('/auth/login', data);
	}

	public static async logout() {
		return await ApiService.post<null>('/auth/logout');
	}
}

export default AuthService;
