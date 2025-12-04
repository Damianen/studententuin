import type { CreateDatabaseInput, Database, UpdateDatabaseInput } from '../entities';

export interface IDatabaseRepo {
	create(database: CreateDatabaseInput): Promise<Database | null>;
	update(database: UpdateDatabaseInput): Promise<Database | null>;
	delete(id: string): Promise<void>;
	findById(id: string): Promise<Database | null>;
	findBySubdomainId(subdomainId: string): Promise<Database | null>;
}
