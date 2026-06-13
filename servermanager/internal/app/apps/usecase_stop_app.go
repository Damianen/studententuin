package apps

import (
	"context"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

type StopApp struct {
	runtime ports.ContainerRuntime
}

func (u *StopApp) Execute(ctx context.Context, appID string) error {
	return u.runtime.Stop(ctx, domain.AppContainerName(appID))
}
