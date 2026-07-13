package db

import (
	"api/internal/app/ports"
	"context"
)

type DeleteDatabase struct {
	databaseRepo  ports.DatabaseRepo
	serverManager ports.ServerManagerClient
}

// Execute removes the database on the hosting server first (container,
// private network, data volume — §4: rows must not outlive their containers,
// and containers must not orphan), then deletes the row. A provision still in
// flight surfaces as ports.ErrProvisionInFlight.
func (d *DeleteDatabase) Execute(ctx context.Context, subdomainID, dbId string) error {
	database, err := d.databaseRepo.FindByID(dbId, ctx)
	if err != nil {
		return err
	}
	if database.SubdomainID.String() != subdomainID {
		return ErrNotInSubdomain
	}

	if err := d.serverManager.RemoveDatabase(ctx, dbId); err != nil {
		return err
	}

	return d.databaseRepo.Delete(database, ctx)
}
