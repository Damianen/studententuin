import type { LoginUserDto } from '../dtos/auth_dtos';
import AuthService from '../services/auth_service';

class AuthController {
	public static async login(credentials: LoginUserDto) {
		const response = await AuthService.login(credentials);
		if (response.code != 200) {
			throw new Error(response.message);
		}
	}

	public static async logout() {
		const response = await AuthService.logout();
		if (response.code != 200) {
			throw new Error(response.message);
		}
	}
}

export default AuthController;
