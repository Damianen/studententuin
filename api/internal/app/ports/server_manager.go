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

// ServerManagerClient is the api-side view of the servermanager's internal
// HTTP API. It grows per phase; phase 2 only needs logs.
type ServerManagerClient interface {
	Logs(ctx context.Context, appID string, opts LogOptions) ([]LogEntry, error)
	// StreamLogs follows the container's logs live. The channel is closed
	// when the stream ends or ctx is canceled; setup failures (e.g. the app
	// was never deployed) are returned synchronously.
	StreamLogs(ctx context.Context, appID string, opts LogOptions) (<-chan LogEntry, error)
}
