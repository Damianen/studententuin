package apps

import (
	"context"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

type GetAppLogs struct {
	runtime ports.ContainerRuntime
}

// Execute returns the container's log lines. Unlike status, a missing
// container is an error here (domain.ErrNotFound → 404): the api decides
// whether "never deployed" means an empty log view.
func (u *GetAppLogs) Execute(ctx context.Context, appID string, opts domain.LogOptions) ([]domain.LogLine, error) {
	return u.runtime.Logs(ctx, domain.AppContainerName(appID), opts)
}

type FollowAppLogs struct {
	runtime ports.ContainerRuntime
}

// Execute opens a live log stream for the app's container. The channel
// closes when the container's log stream ends or ctx is canceled.
func (u *FollowAppLogs) Execute(ctx context.Context, appID string, opts domain.LogOptions) (<-chan domain.LogLine, error) {
	return u.runtime.FollowLogs(ctx, domain.AppContainerName(appID), opts)
}
