import ApplicationController from './application_controller';
import ApplicationService from '@/services/application_service';

vi.mock('@/services/application_service', () => ({
	default: { get: vi.fn(), post: vi.fn(), patch: vi.fn(), delete: vi.fn() },
}));

const app = {
	id: 'app-1',
	name: 'my-app',
	type: 'Nodejs',
	status: 'running',
	repo_url: 'https://github.com/me/my-app',
	branch: 'main',
};

const createDto = {
	name: 'my-app',
	type: 'Nodejs',
	repo_url: 'https://github.com/me/my-app',
	branch: 'main',
	build_command: 'npm run build',
	start_command: 'npm start',
};

describe('ApplicationController', () => {
	it('get returns the application on 200', async () => {
		vi.mocked(ApplicationService.get).mockResolvedValue({
			code: 200,
			message: 'success',
			data: app,
		});

		await expect(ApplicationController.get('sub-1', 'app-1')).resolves.toEqual(
			app
		);
		expect(ApplicationService.get).toHaveBeenCalledWith('sub-1', 'app-1');
	});

	it('get throws on 404', async () => {
		vi.mocked(ApplicationService.get).mockResolvedValue({
			code: 404,
			message: 'application not found',
		});

		await expect(ApplicationController.get('sub-1', 'app-1')).rejects.toThrow(
			'application not found'
		);
	});

	it('create resolves on 201 and throws otherwise', async () => {
		vi.mocked(ApplicationService.post).mockResolvedValue({
			code: 201,
			message: 'success',
		});
		await expect(
			ApplicationController.create('sub-1', createDto)
		).resolves.toBeUndefined();

		vi.mocked(ApplicationService.post).mockResolvedValue({
			code: 400,
			message: 'Unsupported application type',
		});
		await expect(
			ApplicationController.create('sub-1', createDto)
		).rejects.toThrow('Unsupported application type');
	});

	it('patch and delete throw on non-200', async () => {
		vi.mocked(ApplicationService.patch).mockResolvedValue({
			code: 404,
			message: 'application not found',
		});
		await expect(
			ApplicationController.patch('sub-1', 'app-1', { name: 'renamed' })
		).rejects.toThrow('application not found');

		vi.mocked(ApplicationService.delete).mockResolvedValue({
			code: 403,
			message: 'unauthorized',
		});
		await expect(
			ApplicationController.delete('sub-1', 'app-1')
		).rejects.toThrow('unauthorized');
	});
});
