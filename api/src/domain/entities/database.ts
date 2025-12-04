export type DatabaseType = 'postgresql' | 'mysql' | 'mongodb' | 'redis';

export type DatabaseStatus = 'creating' | 'running' | 'stopped' | 'failed';

export interface Database {
	id: string;
	name: string;
	userId: string;
	subdomainId: string;
	type: DatabaseType;
	version: string;
	status: DatabaseStatus;
	connectionString?: string;
	host?: string;
	port?: number;
	username?: string;
	password?: string;
	dockerImage?: string;
	dockerContainerId?: string;
	dockerContainerName?: string;
	volumes?: string[];
	memoryLimit?: string;
	cpuLimit?: string;
	createdAt: Date;
	updatedAt: Date;
}

export type CreateDatabaseInput = Omit<
	Database,
	'id' | 'createdAt' | 'updatedAt' | 'connectionString' | 'status'
>;

export type UpdateDatabaseInput = Partial<
	Omit<Database, 'id' | 'userId' | 'type' | 'createdAt' | 'updatedAt'>
>;
