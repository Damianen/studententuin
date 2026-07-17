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
	// DatabaseID links the app to the subdomain's database: the manager
	// attaches the new container to that database's private network.
	DatabaseID string `json:"database_id,omitempty"`
}

// DeploymentStatus is the manager's job resource the api polls. Status moves
// queued → cloning → building → starting → running | failed. The commit
// fields exist only after the clone and only on phase-6 managers — absent is
// fine.
type DeploymentStatus struct {
	ID            string    `json:"id"`
	AppID         string    `json:"app_id"`
	Status        string    `json:"status"`
	Error         string    `json:"error,omitempty"`
	CommitSHA     string    `json:"commit_sha,omitempty"`
	CommitMessage string    `json:"commit_message,omitempty"`
	CommitAuthor  string    `json:"commit_author,omitempty"`
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

// ErrProvisionInFlight maps the manager's 409 on database provisioning: one
// provision per database at a time, and an existing container means DELETE
// first.
var ErrProvisionInFlight = errors.New("a provision is already in flight for this database")

// ErrDatabaseNotProvisioned means the manager has no container for this
// database.
var ErrDatabaseNotProvisioned = errors.New("no container for this database")

// DBProvisionSpec is the provision body. Credentials are api-generated; the
// manager passes them into the container env and never logs or stores them.
type DBProvisionSpec struct {
	Type        string `json:"type"`
	Version     string `json:"version"`
	DBName      string `json:"db_name"`
	DBUser      string `json:"db_user"`
	DBPassword  string `json:"db_password"`
	MemoryLimit string `json:"memory_limit,omitempty"`
	CpuLimit    string `json:"cpu_limit,omitempty"`
	Runtime     string `json:"runtime,omitempty"`
}

// DBProvisionState is the provision job folded into the status response.
// Status moves queued → pulling → starting → running | failed.
type DBProvisionState struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	Error         string    `json:"error,omitempty"`
	Image         string    `json:"image,omitempty"`
	ContainerID   string    `json:"container_id,omitempty"`
	ContainerName string    `json:"container_name,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Terminal reports whether the provision reached an end state.
func (p *DBProvisionState) Terminal() bool {
	return p.Status == "running" || p.Status == "failed"
}

// MetricPoint is one bucket-averaged sample as the manager returns it.
type MetricPoint struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

// MetricsResponse is the manager's metrics envelope. Every requested series
// key is present with a (possibly empty) slice, never null.
type MetricsResponse struct {
	Range  string                   `json:"range"`
	Series map[string][]MetricPoint `json:"series"`
}

// DBStatus is the manager's database resource: actual Docker state plus the
// latest provision job. Provision is nil when the manager lost the job (e.g.
// restart) — the poller treats exists=false with a nil Provision as a strike.
type DBStatus struct {
	Exists       bool              `json:"exists"`
	Running      bool              `json:"running"`
	StartedAt    *time.Time        `json:"started_at,omitempty"`
	RestartCount int               `json:"restart_count"`
	Status       string            `json:"status"`
	Health       string            `json:"health,omitempty"`
	Provision    *DBProvisionState `json:"provision,omitempty"`
}

// ServerManagerClient is the api-side view of the servermanager's internal
// HTTP API. It grows per phase; phase 3 added deploy and lifecycle, phase 5
// adds databases.
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
	// ProvisionDatabase queues the async pull→create→start job (manager 202).
	// ErrProvisionInFlight on a concurrent provision or existing container;
	// ErrSpecRejected when the manager refuses the spec (allowlist).
	ProvisionDatabase(ctx context.Context, dbID string, spec DBProvisionSpec) (string, error)
	// DatabaseStatus reports Docker state + provision job; poll it after
	// ProvisionDatabase.
	DatabaseStatus(ctx context.Context, dbID string) (*DBStatus, error)
	// RemoveDatabase sweeps the database's container, network, and data
	// volume. Safe to call when nothing exists.
	RemoveDatabase(ctx context.Context, dbID string) error
	// DatabaseLogs reads the database container's logs;
	// ErrDatabaseNotProvisioned when there is no container.
	DatabaseLogs(ctx context.Context, dbID string, opts LogOptions) ([]LogEntry, error)
	// StreamDatabaseLogs follows the database container's logs live,
	// mirroring StreamLogs.
	StreamDatabaseLogs(ctx context.Context, dbID string, opts LogOptions) (<-chan LogEntry, error)
	// Metrics reads the app's sampled series for the range;
	// ErrAppNotDeployed when the manager has no container.
	Metrics(ctx context.Context, appID string, rng string) (*MetricsResponse, error)
	// DatabaseMetrics mirrors Metrics for the database container;
	// ErrDatabaseNotProvisioned when there is none.
	DatabaseMetrics(ctx context.Context, dbID string, rng string) (*MetricsResponse, error)
}
