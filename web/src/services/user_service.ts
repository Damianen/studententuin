import type { UserDto } from '../dtos/user_dtos';
import ApiService from './api_service';

class UserService {
	public static async get() {
		return await ApiService.get<UserDto>('/user');
	}
}

export default UserService;
