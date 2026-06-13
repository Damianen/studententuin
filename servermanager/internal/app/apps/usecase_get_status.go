package apps

import (
	"context"
	"errors"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

type GetAppStatus struct {
	runtime ports.ContainerRuntime
}

// Execute reports the actual Docker-side state (§3.2). A missing container
// is not an error — that's what ContainerState.Exists is for.
func (u *GetAppStatus) Execute(ctx context.Context, appID string) (*domain.ContainerState, error) {
	state, err := u.runtime.Inspect(ctx, domain.AppContainerName(appID))
	if errors.Is(err, domain.ErrNotFound) {
		return &domain.ContainerState{Exists: false}, nil
	}
	if err != nil {
		return nil, err
	}
	return state, nil
}
