import AuthController from './auth_controller';
import AuthService from '@/services/auth_service';

vi.mock('@/services/auth_service', () => ({
	default: { login: vi.fn(), logout: vi.fn() },
}));

const credentials = { email: 'a@b.com', password: 'secret123' };

describe('AuthController', () => {
	it('login resolves on 200', async () => {
		vi.mocked(AuthService.login).mockResolvedValue({
			code: 200,
			message: 'login was successful',
		});

		await expect(AuthController.login(credentials)).resolves.toBeUndefined();
		expect(AuthService.login).toHaveBeenCalledWith(credentials);
	});

	it('login throws the API message on failure', async () => {
		vi.mocked(AuthService.login).mockResolvedValue({
			code: 401,
			message: 'email or password not correct!',
		});

		await expect(AuthController.login(credentials)).rejects.toThrow(
			'email or password not correct!'
		);
	});

	it('logout resolves on 200 and throws otherwise', async () => {
		vi.mocked(AuthService.logout).mockResolvedValue({
			code: 200,
			message: 'logout was successful',
		});
		await expect(AuthController.logout()).resolves.toBeUndefined();

		vi.mocked(AuthService.logout).mockResolvedValue({
			code: 500,
			message: 'boom',
		});
		await expect(AuthController.logout()).rejects.toThrow('boom');
	});
});
