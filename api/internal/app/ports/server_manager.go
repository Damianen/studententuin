package ports

import (
	"context"
	"errors"
	"time"
)

// LogOptions filters a container log read.
type LogOptions struct {
	Tail  int
	Since time.Time
}

// LogEntry is one log line as the servermanager returns it (already in the
// frontend's shape: stdout→info, stderr→error).
type LogEntry struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

// ErrAppNotDeployed means the manager has no container for this application —
// the app exists in the database but was never deployed (or was removed).
var ErrAppNotDeployed = errors.New("no container for this application")

// ErrDeployInFlight maps the manager's 409: one deploy per app at a time.
var ErrDeployInFlight = errors.New("a deployment is already in flight for this application")

// ErrDeploymentNotFound maps the manager's 404 on the deployment resource —
// unknown ID, or the job was lost to a manager restart.
var ErrDeploymentNotFound = errors.New("deployment not found")

// ErrSpecRejected maps the manager's 400: it re-validates the deploy spec
// server-side and refused it. The wrapped message carries the reason.
var ErrSpecRejected = errors.New("servermanager rejected the deploy spec")

// DeploySpec is the §3.2 deploy body, built from the Application record.
type DeploySpec struct {
	RepositoryURL string            `json:"repository_url"`
	Branch        string            `json:"branch,omitempty"`
	Type          string            `json:"type,omitempty"`
	BuildCommand  string            `json:"build_command,omitempty"`
	StartCommand  string            `json:"start_command,omitempty"`
	Env           map[string]string `json:"env,omitempty"`
	Port          int               `json:"port,omitempty"`
	MemoryLimit   string            `json:"memory_limit,omitempty"`
	CpuLimit      string            `json:"cpu_limit,omitempty"`
}

// DeploymentStatus is the manager's job resource the api polls. Status moves
// queued → cloning → building → starting → running | failed.
type DeploymentStatus struct {
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

// Terminal reports whether the deployment reached an end state.
func (d *DeploymentStatus) Terminal() bool {
	return d.Status == "running" || d.Status == "failed"
}

// ServerManagerClient is the api-side view of the servermanager's internal
// HTTP API. It grows per phase; phase 3 adds deploy and lifecycle.
type ServerManagerClient interface {
	Logs(ctx context.Context, appID string, opts LogOptions) ([]LogEntry, error)
	// StreamLogs follows the container's logs live. The channel is closed
	// when the stream ends or ctx is canceled; setup failures (e.g. the app
	// was never deployed) are returned synchronously.
	StreamLogs(ctx context.Context, appID string, opts LogOptions) (<-chan LogEntry, error)
	// Deploy queues the async clone→build→start job (manager 202) and
	// returns its deployment ID. ErrDeployInFlight on a concurrent deploy.
	Deploy(ctx context.Context, appID string, spec DeploySpec) (string, error)
	// DeploymentStatus polls the job. ErrDeploymentNotFound for unknown IDs.
	DeploymentStatus(ctx context.Context, deploymentID string) (*DeploymentStatus, error)
	// Start/Stop drive the existing container; ErrAppNotDeployed when there
	// is none.
	Start(ctx context.Context, appID string) error
	Stop(ctx context.Context, appID string) error
	// Remove deletes the container, its per-app networks, and anonymous
	// volumes. ErrAppNotDeployed when there is nothing to remove.
	Remove(ctx context.Context, appID string) error
}
