package subdomain

import (
	"api/internal/app/ports"
	"context"
	"errors"

	"gorm.io/gorm"
)

type DeleteSubdomain struct {
	subdomainRepo   ports.SubdomainRepo
	applicationRepo ports.ApplicationRepo
	databaseRepo    ports.DatabaseRepo
	serverManager   ports.ServerManagerClient
}

// Execute removes the subdomain's containers on the hosting server before
// deleting the row (§4: nothing may orphan on the hosting server — the row
// cascade would otherwise leave the app container and the database container
// plus its data volume running unowned).
func (d *DeleteSubdomain) Execute(ctx context.Context, subdomainID string) error {
	subdomain, err := d.subdomainRepo.FindByID(subdomainID, ctx)
	if err != nil {
		return err
	}

	application, err := d.applicationRepo.FindBySubdomainID(subdomainID, ctx)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// No application on this subdomain — nothing to remove.
	case err != nil:
		return err
	default:
		if err := d.serverManager.Remove(ctx, application.ID.String()); err != nil && !errors.Is(err, ports.ErrAppNotDeployed) {
			return err
		}
	}

	database, err := d.databaseRepo.FindBySubdomainID(subdomainID, ctx)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// No database on this subdomain — nothing to remove.
	case err != nil:
		return err
	default:
		// The manager's sweep is idempotent; ErrProvisionInFlight propagates
		// so a mid-provision database isn't half-deleted.
		if err := d.serverManager.RemoveDatabase(ctx, database.ID.String()); err != nil {
			return err
		}
	}

	return d.subdomainRepo.Delete(subdomain, ctx)
}
