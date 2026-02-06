package db

import (
	"api/internal/app/ports"
	"context"
)

type DeleteDatabase struct {
	databaseRepo ports.DatabaseRepo
}

func (d *DeleteDatabase) Execute(ctx context.Context, dbId string) error {
	db, err := d.databaseRepo.FindByID(dbId, ctx)
	if err != nil {
		return err
	}

	return d.databaseRepo.Delete(db, ctx)
}
