import ApplicationController from './application_controller';
import ApplicationService from '@/services/application_service';

vi.mock('@/services/application_service', () => ({
	default: {
		get: vi.fn(),
		getLogs: vi.fn(),
		post: vi.fn(),
		patch: vi.fn(),
		delete: vi.fn(),
		deploy: vi.fn(),
		getDeployment: vi.fn(),
		start: vi.fn(),
		stop: vi.fn(),
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

	it('deploy returns the deployment id on 202', async () => {
		vi.mocked(ApplicationService.deploy).mockResolvedValue({
			code: 202,
			message: 'deployment queued',
			data: { deployment_id: 'dep-123' },
		});

		await expect(ApplicationController.deploy('sub-1', 'app-1')).resolves.toBe(
			'dep-123'
		);
		expect(ApplicationService.deploy).toHaveBeenCalledWith('sub-1', 'app-1');
	});

	it('deploy throws with the api message on conflict', async () => {
		vi.mocked(ApplicationService.deploy).mockResolvedValue({
			code: 409,
			message: 'a deployment is already in progress for this application',
		});

		await expect(
			ApplicationController.deploy('sub-1', 'app-1')
		).rejects.toThrow('already in progress');
	});

	it('getDeployment returns the job status', async () => {
		vi.mocked(ApplicationService.getDeployment).mockResolvedValue({
			code: 200,
			message: 'success',
			data: {
				id: 'dep-123',
				status: 'failed',
				error: 'build failed: exit 1',
				build_log: 'npm ERR! boom\n',
				created_at: '2026-07-10T12:00:00Z',
				updated_at: '2026-07-10T12:01:00Z',
			},
		});

		const status = await ApplicationController.getDeployment(
			'sub-1',
			'app-1',
			'dep-123'
		);
		expect(status.status).toBe('failed');
		expect(status.error).toContain('exit 1');
		expect(ApplicationService.getDeployment).toHaveBeenCalledWith(
			'sub-1',
			'app-1',
			'dep-123'
		);
	});

	it('start and stop resolve on 200 and throw otherwise', async () => {
		vi.mocked(ApplicationService.start).mockResolvedValue({
			code: 200,
			message: 'application started',
		});
		await expect(
			ApplicationController.start('sub-1', 'app-1')
		).resolves.toBeUndefined();

		vi.mocked(ApplicationService.stop).mockResolvedValue({
			code: 400,
			message: 'application has not been deployed yet',
		});
		await expect(ApplicationController.stop('sub-1', 'app-1')).rejects.toThrow(
			'not been deployed'
		);
	});
});
