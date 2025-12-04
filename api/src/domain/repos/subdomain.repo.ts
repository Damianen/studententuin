import type { CreateSubdomainInput, Subdomain, UpdateSubdomainInput } from '../entities';

export interface ISubdomainRepo {
	create(subdomain: CreateSubdomainInput): Promise<Subdomain | null>;
	update(subdomain: UpdateSubdomainInput): Promise<Subdomain | null>;
	delete(id: string): Promise<void>;
	findById(id: string): Promise<Subdomain | null>;
	findByName(name: string): Promise<Subdomain | null>;
	findByUserId(userId: string): Promise<Subdomain[] | null>;
}
