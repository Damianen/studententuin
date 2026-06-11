import SubdomainController from '../controllers/subdomain_controller';
import type { SubdomainListItemDto } from '../dtos/subdomain_dtos';

export const SUBDOMAIN_SUFFIX = 'studententuin.com';

export function toFullDomain(subdomainName: string): string {
	return `${subdomainName}.${SUBDOMAIN_SUFFIX}`;
}

function findByName(
	subdomains: SubdomainListItemDto[],
	cleanName: string,
): SubdomainListItemDto | undefined {
	return subdomains.find(
		(subdomain) =>
			subdomain.name.toLowerCase() === cleanName ||
			subdomain.fullDomain.toLowerCase() === toFullDomain(cleanName),
	);
}

export async function ensureSubdomain(
	subdomainName: string,
): Promise<SubdomainListItemDto> {
	const cleanName = subdomainName.trim().toLowerCase();
	if (!cleanName) {
		throw new Error('Subdomain is required');
	}

	const current = await SubdomainController.getAll();
	const found = findByName(current, cleanName);

	if (found) {
		return found;
	}

	await SubdomainController.create({
		name: cleanName,
		fullDomain: toFullDomain(cleanName),
	});

	const updated = await SubdomainController.getAll();
	const created = findByName(updated, cleanName);

	if (!created) {
		throw new Error('Failed to resolve created subdomain');
	}

	return created;
}
