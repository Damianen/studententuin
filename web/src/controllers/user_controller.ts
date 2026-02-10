import UserService from '../services/user_service';

class UserController {
	public static async get() {
		const response = await UserService.get();
		if (response.code != 200) {
			throw new Error(response.message);
		}
		return response.data;
	}
}

export default UserController;
