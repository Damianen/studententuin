import UserController from './user_controller';
import UserService from '@/services/user_service';

vi.mock('@/services/user_service', () => ({
	default: { get: vi.fn(), patch: vi.fn(), post: vi.fn(), delete: vi.fn() },
}));

const user = { email: 'a@b.com', name: 'Alice', status: 'active' };

describe('UserController', () => {
	it('get returns the user on 200', async () => {
		vi.mocked(UserService.get).mockResolvedValue({
			code: 200,
			message: 'success',
			data: user,
		});

		await expect(UserController.get()).resolves.toEqual(user);
	});

	it('get throws when data is missing even on 200', async () => {
		vi.mocked(UserService.get).mockResolvedValue({
			code: 200,
			message: 'success',
		});

		await expect(UserController.get()).rejects.toThrow('success');
	});

	it('get throws the API message on error status', async () => {
		vi.mocked(UserService.get).mockResolvedValue({
			code: 401,
			message: 'authentication required',
		});

		await expect(UserController.get()).rejects.toThrow(
			'authentication required'
		);
	});

	it('register resolves on 201 and throws on conflict', async () => {
		vi.mocked(UserService.post).mockResolvedValue({
			code: 201,
			message: 'success',
		});
		await expect(
			UserController.register({ email: 'a@b.com', password: 'x', name: 'A' })
		).resolves.toBeUndefined();

		vi.mocked(UserService.post).mockResolvedValue({
			code: 409,
			message: 'email already in use',
		});
		await expect(
			UserController.register({ email: 'a@b.com', password: 'x', name: 'A' })
		).rejects.toThrow('email already in use');
	});

	it('patch resolves on 200 and throws otherwise', async () => {
		vi.mocked(UserService.patch).mockResolvedValue({
			code: 200,
			message: 'success',
		});
		await expect(
			UserController.patch({ name: 'New' })
		).resolves.toBeUndefined();

		vi.mocked(UserService.patch).mockResolvedValue({
			code: 404,
			message: 'user not found',
		});
		await expect(UserController.patch({ name: 'New' })).rejects.toThrow(
			'user not found'
		);
	});

	it('delete resolves on 200 and throws otherwise', async () => {
		vi.mocked(UserService.delete).mockResolvedValue({
			code: 200,
			message: 'success',
		});
		await expect(UserController.delete()).resolves.toBeUndefined();

		vi.mocked(UserService.delete).mockResolvedValue({
			code: 500,
			message: 'failed to delete user',
		});
		await expect(UserController.delete()).rejects.toThrow(
			'failed to delete user'
		);
	});
});
