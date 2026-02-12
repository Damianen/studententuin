import type { RegisterUserDto, UpdateUserDto } from '../dtos/user_dtos';
import UserService from '../services/user_service';

class UserController {
	public static async get() {
		const response = await UserService.get();
		if (response.code != 200) {
			throw new Error(response.message);
		}
		return response.data;
	}
	public static async patch(values: UpdateUserDto) {
		const response = await UserService.patch(values);
		if (response.code != 200) {
			throw new Error(response.message);
		}
		return response.data;
	}
	public static async register(values: RegisterUserDto) {
		const response = await UserService.post(values);
		if (response.code != 201) {
			throw new Error(response.message);
		}
		return response.data;
	}
	public static async delete() {
		const response = await UserService.delete();
		if (response.code != 200) {
			throw new Error(response.message);
		}
	}
}

export default UserController;
