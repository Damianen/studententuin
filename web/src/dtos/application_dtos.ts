export interface ApplicationDto {
	id: string;
	name: string;
	type: string;
	status: string;
	repo_url: string;
	branch: string;
}

export interface CreateApplicationDto {
	name: string;
	type: string;
	repo_url: string;
	branch: string;
	build_command: string;
	start_command: string;
}

export interface UpdateApplicationDto {
	name?: string;
	type?: string;
	repo_url?: string;
	branch?: string;
	build_command?: string;
	start_command?: string;
}

export interface LogEntryDto {
	id: string;
	timestamp: string; // RFC3339
	level: string;
	message: string;
}

export interface DeployResponseDto {
	deployment_id: string;
}

// Status moves queued → cloning → building → starting → running | failed;
// error and a build-log tail are attached on failure.
export interface DeploymentStatusDto {
	id: string;
	status: string;
	error?: string;
	commit_sha?: string;
	build_log?: string;
	created_at: string; // RFC3339
	updated_at: string; // RFC3339
}
