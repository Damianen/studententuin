package app

import (
	"api/internal/app/ports"
	"context"
)

type DeleteApplication struct {
	applicationRepo ports.ApplicationRepo
}

func (d *DeleteApplication)  Execute(ctx context.Context, appId string) error {
	app, err := d.applicationRepo.FindByID(appId, ctx)
	if err != nil {
		return err
	}

	return d.applicationRepo.Delete(app, ctx)
}
