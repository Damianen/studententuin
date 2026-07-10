package app

import (
	"api/internal/app/ports"
	"context"
	"errors"
)

type DeleteApplication struct {
	applicationRepo ports.ApplicationRepo
	serverManager   ports.ServerManagerClient
}

// Execute removes the app's container (and network/volumes) on the hosting
// server before deleting the record, so nothing orphans there (§4). A
// never-deployed app has nothing to remove; any other manager failure aborts
// the delete — better a retryable error than an orphaned container.
func (d *DeleteApplication) Execute(ctx context.Context, appId string) error {
	app, err := d.applicationRepo.FindByID(appId, ctx)
	if err != nil {
		return err
	}

	if err := d.serverManager.Remove(ctx, appId); err != nil && !errors.Is(err, ports.ErrAppNotDeployed) {
		return err
	}

	return d.applicationRepo.Delete(app, ctx)
}
