export interface DatabaseDto {
	id: string;
	name: string;
	type: string;
	version: string;
	status: string;
	connection_string: string;
}

// No password field: credentials are generated server-side and surface only
// via the connection string once the database is provisioned (phase 5).
export interface CreateDatabaseDto {
	name: string;
	type: string;
	version: string;
	db_name: string;
}

// Engine (type/version) is immutable after creation — changing it would mean
// re-provisioning the container; delete and recreate instead.
export interface UpdateDatabaseDto {
	name?: string;
}
