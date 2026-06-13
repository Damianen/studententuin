package apps

import (
	"context"
	"errors"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

type RemoveApp struct {
	runtime ports.ContainerRuntime
}

// Execute removes the app's container. Idempotent: a container that is
// already gone is success — the manager reconciles Docker state to what the
// api tells it (§2).
func (u *RemoveApp) Execute(ctx context.Context, appID string) error {
	err := u.runtime.Remove(ctx, domain.AppContainerName(appID))
	if errors.Is(err, domain.ErrNotFound) {
		return nil
	}
	return err
}
