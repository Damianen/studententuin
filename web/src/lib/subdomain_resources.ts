import type { SubdomainListItemDto } from '../dtos/subdomain_dtos';

export interface ResourceSummary {
	resourceId: string;
	subdomainId: string;
	resourceName: string;
	subdomainName: string;
	fullDomain: string;
	typeLabel: string;
	kind: 'application' | 'database';
}

export function mapSubdomainResources(
	subdomains: SubdomainListItemDto[],
): ResourceSummary[] {
	return subdomains.flatMap((subdomain) => {
		const resources: ResourceSummary[] = [];

		if (subdomain.application) {
			resources.push({
				resourceId: subdomain.application.id,
				subdomainId: subdomain.id,
				resourceName: subdomain.application.name,
				subdomainName: subdomain.name,
				fullDomain: subdomain.fullDomain,
				typeLabel: subdomain.application.type,
				kind: 'application',
			});
		}

		if (subdomain.database) {
			resources.push({
				resourceId: subdomain.database.id,
				subdomainId: subdomain.id,
				resourceName: subdomain.database.name,
				subdomainName: subdomain.name,
				fullDomain: subdomain.fullDomain,
				typeLabel: subdomain.database.type,
				kind: 'database',
			});
		}

		return resources;
	});
}

export function findResourceById(
	subdomains: SubdomainListItemDto[],
	resourceId: string,
): ResourceSummary | undefined {
	return mapSubdomainResources(subdomains).find(
		(resource) => resource.resourceId === resourceId,
	);
}
