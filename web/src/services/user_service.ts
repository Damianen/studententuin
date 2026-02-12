import type {
	RegisterUserDto,
	UpdateUserDto,
	UserDto,
} from '../dtos/user_dtos';
import ApiService from './api_service';

class UserService {
	public static async get() {
		return await ApiService.get<UserDto>('/user');
	}

	public static async patch(values: UpdateUserDto) {
		return await ApiService.patch<void>('/user', values);
	}

	public static async delete() {
		return await ApiService.delete<void>('/user');
	}

	public static async post(user: RegisterUserDto) {
		return await ApiService.post('/user/register', user);
	}
}

export default UserService;
