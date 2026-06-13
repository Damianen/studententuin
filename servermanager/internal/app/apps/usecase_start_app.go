package apps

import (
	"context"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

type StartApp struct {
	runtime ports.ContainerRuntime
}

func (u *StartApp) Execute(ctx context.Context, appID string) error {
	return u.runtime.Start(ctx, domain.AppContainerName(appID))
}
