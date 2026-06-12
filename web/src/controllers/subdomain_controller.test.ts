import SubdomainController from './subdomain_controller';
import SubdomainService from '@/services/subdomain_service';

vi.mock('@/services/subdomain_service', () => ({
	default: { getAll: vi.fn(), post: vi.fn(), patch: vi.fn(), delete: vi.fn() },
}));

const item = {
	id: 'abc',
	name: 'myapp',
	fullDomain: 'myapp.studententuin.com',
	isActive: true,
};

describe('SubdomainController', () => {
	it('getAll returns the list on 200', async () => {
		vi.mocked(SubdomainService.getAll).mockResolvedValue({
			code: 200,
			message: 'success',
			data: [item],
		});

		await expect(SubdomainController.getAll()).resolves.toEqual([item]);
	});

	it('getAll returns [] when the API omits data for an empty list', async () => {
		vi.mocked(SubdomainService.getAll).mockResolvedValue({
			code: 200,
			message: 'success',
		});

		await expect(SubdomainController.getAll()).resolves.toEqual([]);
	});

	it('getAll throws the API message on error', async () => {
		vi.mocked(SubdomainService.getAll).mockResolvedValue({
			code: 500,
			message: 'failed to get subdomains',
		});

		await expect(SubdomainController.getAll()).rejects.toThrow(
			'failed to get subdomains'
		);
	});

	it('create resolves on 201 and throws on conflict', async () => {
		vi.mocked(SubdomainService.post).mockResolvedValue({
			code: 201,
			message: 'success',
		});
		await expect(
			SubdomainController.create({
				name: 'myapp',
				fullDomain: 'myapp.studententuin.com',
			})
		).resolves.toBeUndefined();

		vi.mocked(SubdomainService.post).mockResolvedValue({
			code: 409,
			message: 'domain already in use',
		});
		await expect(
			SubdomainController.create({
				name: 'myapp',
				fullDomain: 'myapp.studententuin.com',
			})
		).rejects.toThrow('domain already in use');
	});

	it('patch and delete throw on non-200', async () => {
		vi.mocked(SubdomainService.patch).mockResolvedValue({
			code: 403,
			message: 'unauthorized',
		});
		await expect(
			SubdomainController.patch('abc', { name: 'renamed' })
		).rejects.toThrow('unauthorized');

		vi.mocked(SubdomainService.delete).mockResolvedValue({
			code: 404,
			message: 'subdomain not found',
		});
		await expect(SubdomainController.delete('abc')).rejects.toThrow(
			'subdomain not found'
		);
	});
});
