import DatabaseController from './database_controller';
import DatabaseService from '@/services/database_service';

vi.mock('@/services/database_service', () => ({
	default: {
		get: vi.fn(),
		getMetrics: vi.fn(),
		post: vi.fn(),
		patch: vi.fn(),
		delete: vi.fn(),
	},
}));

const db = {
	id: 'db-1',
	name: 'mydb',
	type: 'postgres',
	version: '17',
	status: 'running',
	connection_string: 'postgres://user:pass@host/db',
};

const createDto = {
	name: 'mydb',
	type: 'postgres',
	version: '17',
	db_name: 'app',
	db_password: 'secret123',
};

describe('DatabaseController', () => {
	it('get returns the database on 200', async () => {
		vi.mocked(DatabaseService.get).mockResolvedValue({
			code: 200,
			message: 'success',
			data: db,
		});

		await expect(DatabaseController.get('sub-1', 'db-1')).resolves.toEqual(db);
		expect(DatabaseService.get).toHaveBeenCalledWith('sub-1', 'db-1');
	});

	it('get throws on 404', async () => {
		vi.mocked(DatabaseService.get).mockResolvedValue({
			code: 404,
			message: 'database not found',
		});

		await expect(DatabaseController.get('sub-1', 'db-1')).rejects.toThrow(
			'database not found'
		);
	});

	it('create resolves on 201 and throws otherwise', async () => {
		vi.mocked(DatabaseService.post).mockResolvedValue({
			code: 201,
			message: 'success',
		});
		await expect(
			DatabaseController.create('sub-1', createDto)
		).resolves.toBeUndefined();

		vi.mocked(DatabaseService.post).mockResolvedValue({
			code: 400,
			message: 'Invalid JSON or missing value',
		});
		await expect(DatabaseController.create('sub-1', createDto)).rejects.toThrow(
			'Invalid JSON or missing value'
		);
	});

	it('getMetrics maps the envelope and guarantees all four series', async () => {
		vi.mocked(DatabaseService.getMetrics).mockResolvedValue({
			code: 200,
			message: 'success',
			data: {
				range: '24h',
				series: {
					conn: [{ time: '2026-07-10T12:00:00Z', value: 3 }],
					disk: [{ time: '2026-07-10T12:00:00Z', value: 213.4 }],
				},
			},
		});

		const series = await DatabaseController.getMetrics('sub-1', 'db-1');
		expect(DatabaseService.getMetrics).toHaveBeenCalledWith('sub-1', 'db-1');
		expect(series.conn).toHaveLength(1);
		expect(series.disk[0].value).toBe(213.4);
		// qps/cpu had no samples yet; they still come back as arrays.
		expect(series.qps).toEqual([]);
		expect(series.cpu).toEqual([]);
	});

	it('getMetrics throws on non-200', async () => {
		vi.mocked(DatabaseService.getMetrics).mockResolvedValue({
			code: 502,
			message: 'servermanager unreachable',
		});

		await expect(
			DatabaseController.getMetrics('sub-1', 'db-1')
		).rejects.toThrow('servermanager unreachable');
	});

	it('patch and delete throw on non-200', async () => {
		vi.mocked(DatabaseService.patch).mockResolvedValue({
			code: 404,
			message: 'database not found',
		});
		await expect(
			DatabaseController.patch('sub-1', 'db-1', { name: 'renamed' })
		).rejects.toThrow('database not found');

		vi.mocked(DatabaseService.delete).mockResolvedValue({
			code: 403,
			message: 'unauthorized',
		});
		await expect(DatabaseController.delete('sub-1', 'db-1')).rejects.toThrow(
			'unauthorized'
		);
	});
});
