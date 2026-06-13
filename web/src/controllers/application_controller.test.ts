import ApplicationController from './application_controller';
import ApplicationService from '@/services/application_service';

vi.mock('@/services/application_service', () => ({
	default: {
		get: vi.fn(),
		getLogs: vi.fn(),
		post: vi.fn(),
		patch: vi.fn(),
		delete: vi.fn(),
	},
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

	it('getLogs maps entries to display shape', async () => {
		vi.mocked(ApplicationService.getLogs).mockResolvedValue({
			code: 200,
			message: 'success',
			data: [
				{
					id: '1-0',
					timestamp: '2026-06-13T10:00:01Z',
					level: 'error',
					message: 'tock',
				},
				{ id: '2-1', timestamp: '', level: 'weird', message: 'tick' },
			],
		});

		const logs = await ApplicationController.getLogs('sub-1', 'app-1');
		expect(ApplicationService.getLogs).toHaveBeenCalledWith('sub-1', 'app-1');
		expect(logs).toHaveLength(2);
		expect(logs[0].level).toBe('error');
		expect(logs[0].message).toBe('tock');
		// RFC3339 timestamps are shortened to a local HH:MM:SS for the terminal.
		expect(logs[0].timestamp).toMatch(/^\d{2}:\d{2}:\d{2}$/);
		// Unknown levels narrow to info; missing timestamps get a placeholder.
		expect(logs[1].level).toBe('info');
		expect(logs[1].timestamp).toBe('--:--:--');
	});

	it('getLogs returns empty when the api omits data', async () => {
		vi.mocked(ApplicationService.getLogs).mockResolvedValue({
			code: 200,
			message: 'success',
		});

		await expect(
			ApplicationController.getLogs('sub-1', 'app-1')
		).resolves.toEqual([]);
	});

	it('getLogs throws on non-200', async () => {
		vi.mocked(ApplicationService.getLogs).mockResolvedValue({
			code: 502,
			message: 'log service unavailable',
		});

		await expect(
			ApplicationController.getLogs('sub-1', 'app-1')
		).rejects.toThrow('log service unavailable');
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
