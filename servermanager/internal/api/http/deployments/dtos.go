package deployments

import "time"

// DeployRequest is the §3.2 deploy body. StartCommand is a shell string here
// — nixpacks bakes it into the image CMD; only /run takes an argv array.
type DeployRequest struct {
	RepositoryURL string            `json:"repository_url"`
	Branch        string            `json:"branch"`
	Type          string            `json:"type"`
	BuildCommand  string            `json:"build_command"`
	StartCommand  string            `json:"start_command"`
	Env           map[string]string `json:"env"`
	Port          int               `json:"port"`
	MemoryLimit   string            `json:"memory_limit"`
	CpuLimit      string            `json:"cpu_limit"`
	Runtime       string            `json:"runtime"`
}

type DeployResponse struct {
	DeploymentID string `json:"deployment_id"`
}

// DeploymentResponse is the poll shape the api consumes. BuildLog is already
// capped to its last 64KiB in the job store — it IS the tail.
type DeploymentResponse struct {
	ID            string    `json:"id"`
	AppID         string    `json:"app_id"`
	Status        string    `json:"status"`
	Error         string    `json:"error,omitempty"`
	CommitSHA     string    `json:"commit_sha,omitempty"`
	Image         string    `json:"image,omitempty"`
	ContainerID   string    `json:"container_id,omitempty"`
	ContainerName string    `json:"container_name,omitempty"`
	BuildLog      string    `json:"build_log,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
