import type { Application, CreateApplicationInput, UpdateApplicationInput } from '../entities';

export interface IApplicationRepo {
	create(app: CreateApplicationInput): Promise<Application | null>;
	update(app: UpdateApplicationInput): Promise<Application | null>;
	delete(id: string): Promise<void>;
	findById(id: string): Promise<Application | null>;
	findBySubdomainId(subdomainId: string): Promise<Application | null>;
}
