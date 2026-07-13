package databases

import (
	"context"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

type GetDatabaseLogs struct {
	runtime ports.ContainerRuntime
}

// Execute returns the database container's log lines. A missing container is
// an error (domain.ErrNotFound → 404), mirroring the app logs endpoint.
func (u *GetDatabaseLogs) Execute(ctx context.Context, dbID string, opts domain.LogOptions) ([]domain.LogLine, error) {
	return u.runtime.Logs(ctx, domain.DBContainerName(dbID), opts)
}

type FollowDatabaseLogs struct {
	runtime ports.ContainerRuntime
}

// Execute opens a live log stream for the database's container. The channel
// closes when the container's log stream ends or ctx is canceled.
func (u *FollowDatabaseLogs) Execute(ctx context.Context, dbID string, opts domain.LogOptions) (<-chan domain.LogLine, error) {
	return u.runtime.FollowLogs(ctx, domain.DBContainerName(dbID), opts)
}
