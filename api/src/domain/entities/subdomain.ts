export interface Subdomain {
	id: string;
	name: string;
	fullDomain: string;
	userId: string;
	applicationId?: string;
	databaseId?: string;
	isActive: boolean;
	createdAt: Date;
	updatedAt: Date;
}

export type CreateSubdomainInput = Omit<Subdomain, 'id' | 'createdAt' | 'updatedAt' | 'fullDomain'>;

export type UpdateSubdomainInput = Partial<
	Omit<Subdomain, 'id' | 'userId' | 'createdAt' | 'updatedAt'>
>;
