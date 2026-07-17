import ApplicationController, {
	mapDeployment,
	mapMetrics,
} from './application_controller';
import ApplicationService from '@/services/application_service';

vi.mock('@/services/application_service', () => ({
	default: {
		get: vi.fn(),
		getLogs: vi.fn(),
		getMetrics: vi.fn(),
		post: vi.fn(),
		patch: vi.fn(),
		delete: vi.fn(),
		deploy: vi.fn(),
		getDeployment: vi.fn(),
		getDeployments: vi.fn(),
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

	it('getMetrics maps the envelope and guarantees both app series', async () => {
		vi.mocked(ApplicationService.getMetrics).mockResolvedValue({
			code: 200,
			message: 'success',
			data: {
				range: '24h',
				series: { cpu: [{ time: '2026-07-10T12:00:00Z', value: 40 }] },
			},
		});

		const series = await ApplicationController.getMetrics('sub-1', 'app-1');
		expect(ApplicationService.getMetrics).toHaveBeenCalledWith(
			'sub-1',
			'app-1'
		);
		expect(series.cpu).toHaveLength(1);
		expect(series.cpu[0].value).toBe(40);
		// mem had no samples yet; it still comes back as an array.
		expect(series.mem).toEqual([]);
	});

	it('getMetrics throws on non-200', async () => {
		vi.mocked(ApplicationService.getMetrics).mockResolvedValue({
			code: 502,
			message: 'servermanager unreachable',
		});

		await expect(
			ApplicationController.getMetrics('sub-1', 'app-1')
		).rejects.toThrow('servermanager unreachable');
	});

	it('getDeployments maps records and tolerates omitted data', async () => {
		vi.mocked(ApplicationService.getDeployments).mockResolvedValue({
			code: 200,
			message: 'success',
			data: [
				{
					id: 'dep-1',
					status: 'succeeded',
					branch: 'main',
					commit_sha: 'abcdef1234567890',
					commit_message: 'Ship it',
					commit_author: 'damian',
					started_at: '2026-07-10T12:00:00Z',
					finished_at: '2026-07-10T12:02:00Z',
					duration_seconds: 120,
				},
			],
		});

		const deployments = await ApplicationController.getDeployments(
			'sub-1',
			'app-1'
		);
		expect(ApplicationService.getDeployments).toHaveBeenCalledWith(
			'sub-1',
			'app-1'
		);
		expect(deployments).toHaveLength(1);
		expect(deployments[0].status).toBe('success');
		expect(deployments[0].commit).toBe('abcdef1');

		vi.mocked(ApplicationService.getDeployments).mockResolvedValue({
			code: 200,
			message: 'success',
		});
		await expect(
			ApplicationController.getDeployments('sub-1', 'app-1')
		).resolves.toEqual([]);
	});

	it('getDeployments throws on non-200', async () => {
		vi.mocked(ApplicationService.getDeployments).mockResolvedValue({
			code: 404,
			message: 'application not found',
		});

		await expect(
			ApplicationController.getDeployments('sub-1', 'app-1')
		).rejects.toThrow('application not found');
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

describe('mapMetrics', () => {
	it('shortens RFC3339 times to browser-local HH:MM labels', () => {
		const series = mapMetrics(
			{
				range: '24h',
				series: { cpu: [{ time: '2026-07-10T12:34:00Z', value: 12.3 }] },
			},
			['cpu']
		);

		// toLocaleTimeString renders in the runner's zone, exactly like the
		// browser will — derive the expected hour the same (local) way.
		const localHour = String(
			new Date('2026-07-10T12:34:00Z').getHours()
		).padStart(2, '0');
		expect(series.cpu).toEqual([{ time: `${localHour}:34`, value: 12.3 }]);
	});

	it('returns every requested key even when the api omits series', () => {
		expect(mapMetrics({ range: '24h' }, ['cpu', 'mem'])).toEqual({
			cpu: [],
			mem: [],
		});
		expect(mapMetrics({ range: '24h', series: {} }, ['conn'])).toEqual({
			conn: [],
		});
		expect(mapMetrics(undefined, ['disk'])).toEqual({ disk: [] });
	});
});

describe('mapDeployment', () => {
	beforeEach(() => {
		vi.useFakeTimers();
		vi.setSystemTime(new Date('2026-07-13T12:00:00Z'));
	});

	afterEach(() => {
		vi.useRealTimers();
	});

	it('maps a succeeded record to the display shape', () => {
		const record = mapDeployment({
			id: 'dep-1',
			status: 'succeeded',
			branch: 'main',
			commit_sha: 'abcdef1234567890abcdef1234567890abcdef12',
			commit_message: 'Fix login redirect',
			commit_author: 'damian',
			started_at: '2026-07-13T10:00:00Z',
			finished_at: '2026-07-13T10:01:30Z',
			duration_seconds: 90,
		});

		expect(record.status).toBe('success');
		expect(record.commit).toBe('abcdef1');
		expect(record.message).toBe('Fix login redirect');
		expect(record.author).toBe('damian');
		expect(record.timeAgo).toBe('2 hours ago');
		expect(record.duration).toBe('1m 30s');
		expect(record.durationSeconds).toBe(90);
	});

	it('falls back to em-dashes when a failed clone has no commit', () => {
		const record = mapDeployment({
			id: 'dep-2',
			status: 'failed',
			branch: 'main',
			error: 'clone failed: repository not found',
			started_at: '2026-07-12T11:00:00Z',
			finished_at: '2026-07-12T11:00:05Z',
			duration_seconds: 5,
		});

		expect(record.status).toBe('failed');
		expect(record.commit).toBe('—');
		expect(record.author).toBe('—');
		// The failure reason stands in for the missing commit message.
		expect(record.message).toBe('clone failed: repository not found');
		expect(record.timeAgo).toBe('yesterday');
		expect(record.duration).toBe('5s');
	});

	it('maps in_flight to the in-progress spinner state', () => {
		const record = mapDeployment({
			id: 'dep-3',
			status: 'in_flight',
			branch: 'main',
			started_at: '2026-07-13T11:59:30Z',
		});

		expect(record.status).toBe('in_progress');
		expect(record.message).toBe('Deploy in progress');
		expect(record.timeAgo).toBe('just now');
		expect(record.duration).toBe('—');
		expect(record.durationSeconds).toBeUndefined();
	});
});
