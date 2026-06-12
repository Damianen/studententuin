import { ensureSubdomain, toFullDomain } from './subdomain_lookup';
import SubdomainController from '@/controllers/subdomain_controller';

vi.mock('@/controllers/subdomain_controller', () => ({
	default: { getAll: vi.fn(), create: vi.fn() },
}));

const existing = {
	id: 'sub-1',
	name: 'myapp',
	fullDomain: 'myapp.studententuin.com',
	isActive: true,
};

beforeEach(() => {
	vi.mocked(SubdomainController.getAll).mockReset();
	vi.mocked(SubdomainController.create).mockReset();
});

describe('toFullDomain', () => {
	it('appends the studententuin suffix', () => {
		expect(toFullDomain('myapp')).toBe('myapp.studententuin.com');
	});
});

describe('ensureSubdomain', () => {
	it('throws when the name is empty or whitespace', async () => {
		await expect(ensureSubdomain('   ')).rejects.toThrow(
			'Subdomain is required'
		);
		expect(SubdomainController.getAll).not.toHaveBeenCalled();
	});

	it('returns an existing subdomain without creating (case-insensitive)', async () => {
		vi.mocked(SubdomainController.getAll).mockResolvedValue([existing]);

		const result = await ensureSubdomain('  MyApp ');

		expect(result).toEqual(existing);
		expect(SubdomainController.create).not.toHaveBeenCalled();
	});

	it('matches on full domain too', async () => {
		vi.mocked(SubdomainController.getAll).mockResolvedValue([
			{ ...existing, name: 'different-name' },
		]);

		const result = await ensureSubdomain('myapp');

		expect(result.id).toBe('sub-1');
		expect(SubdomainController.create).not.toHaveBeenCalled();
	});

	it('creates the subdomain when missing and resolves it from a re-fetch', async () => {
		vi.mocked(SubdomainController.getAll)
			.mockResolvedValueOnce([])
			.mockResolvedValueOnce([existing]);
		vi.mocked(SubdomainController.create).mockResolvedValue(undefined);

		const result = await ensureSubdomain('myapp');

		expect(SubdomainController.create).toHaveBeenCalledWith({
			name: 'myapp',
			fullDomain: 'myapp.studententuin.com',
		});
		expect(result).toEqual(existing);
	});

	it('throws when the created subdomain cannot be resolved', async () => {
		vi.mocked(SubdomainController.getAll).mockResolvedValue([]);
		vi.mocked(SubdomainController.create).mockResolvedValue(undefined);

		await expect(ensureSubdomain('myapp')).rejects.toThrow(
			'Failed to resolve created subdomain'
		);
	});

	it('propagates create errors', async () => {
		vi.mocked(SubdomainController.getAll).mockResolvedValue([]);
		vi.mocked(SubdomainController.create).mockRejectedValue(
			new Error('domain already in use')
		);

		await expect(ensureSubdomain('myapp')).rejects.toThrow(
			'domain already in use'
		);
	});
});
