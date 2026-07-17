package dtos

type CreateApplicationRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
	RepoUrl string `json:"repo_url" binding:"required"`
	Branch string `json:"branch" binding:"required"`
	BuildCommand string `json:"build_command" binding:"required"`
	StartCommand string `json:"start_command" binding:"required"`
}

type UpdateApplicationRequest struct {
	Name *string `json:"name,omitempty" binding:"omitempty"`
	Type *string `json:"type,omitempty" binding:"omitempty"`
	RepoUrl *string `json:"repo_url,omitempty" binding:"omitempty"`
	Branch *string `json:"branch,omitempty" binding:"omitempty"`
	BuildCommand *string `json:"build_command,omitempty" binding:"omitempty"`
	StartCommand *string `json:"start_command,omitempty" binding:"omitempty"`
	// Pointer so "not provided" and "provided empty" (= clear all) differ.
	EnvironmentVariables *map[string]string `json:"environment_variables,omitempty" binding:"omitempty"`
}

type ApplicationListResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Status string `json:"status"`
	RepoUrl string `json:"repo_url"`
	Branch string `json:"branch"`
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty"`
}

type LogEntryResponse struct {
	ID string `json:"id"`
	Timestamp string `json:"timestamp"` // RFC3339Nano, UTC
	Level string `json:"level"`
	Message string `json:"message"`
}

type DeployApplicationResponse struct {
	DeploymentID string `json:"deployment_id"`
}

// DeploymentRecordResponse is one persisted deploy from the history table.
// DurationSeconds exists only once the deployment finished.
type DeploymentRecordResponse struct {
	ID string `json:"id"`
	Status string `json:"status"`
	Branch string `json:"branch,omitempty"`
	CommitSHA string `json:"commit_sha,omitempty"`
	CommitMessage string `json:"commit_message,omitempty"`
	CommitAuthor string `json:"commit_author,omitempty"`
	Error string `json:"error,omitempty"`
	StartedAt string `json:"started_at"` // RFC3339, UTC
	FinishedAt string `json:"finished_at,omitempty"` // RFC3339, UTC
	DurationSeconds *float64 `json:"duration_seconds,omitempty"`
}

// DeploymentStatusResponse is the poll shape for deploy progress: status
// moves queued → cloning → building → starting → running | failed, with the
// error and a build-log tail attached on failure.
type DeploymentStatusResponse struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
	CommitSHA string `json:"commit_sha,omitempty"`
	BuildLog  string `json:"build_log,omitempty"`
	CreatedAt string `json:"created_at"` // RFC3339, UTC
	UpdatedAt string `json:"updated_at"` // RFC3339, UTC
}
