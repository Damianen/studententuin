export interface DatabaseDto {
	id: string;
	name: string;
	type: string;
	version: string;
	status: string;
	connection_string: string;
}

export interface CreateDatabaseDto {
	name: string;
	type: string;
	version: string;
	db_name: string;
	db_password: string;
}

export interface UpdateDatabaseDto {
	name?: string;
	type?: string;
	version?: string;
}
