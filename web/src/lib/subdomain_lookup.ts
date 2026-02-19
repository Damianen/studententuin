import SubdomainController from '../controllers/subdomain_controller';

const SUBDOMAIN_SUFFIX = 'studententuin.com';

function toFullDomain(subdomainName: string): string {
	return `${subdomainName}.${SUBDOMAIN_SUFFIX}`;
}

export async function ensureSubdomain(subdomainName: string): Promise<string> {
	const cleanName = subdomainName.trim().toLowerCase();
	if (!cleanName) {
		throw new Error('Subdomain is required');
	}

	const current = await SubdomainController.getAll();
	const found = current.find(
		(subdomain) =>
			subdomain.name.toLowerCase() === cleanName ||
			subdomain.fullDomain.toLowerCase() === toFullDomain(cleanName),
	);

	if (found) {
		return found.id;
	}

	await SubdomainController.create({
		name: cleanName,
		fullDomain: toFullDomain(cleanName),
	});

	const updated = await SubdomainController.getAll();
	const created = updated.find(
		(subdomain) =>
			subdomain.name.toLowerCase() === cleanName ||
			subdomain.fullDomain.toLowerCase() === toFullDomain(cleanName),
	);

	if (!created) {
		throw new Error('Failed to resolve created subdomain');
	}

	return created.id;
}
