package app

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
)

type GetApplication struct {
	applicationRepo ports.ApplicationRepo
}

func (g *GetApplication) Execute(ctx context.Context, appId string) (*domain.Application, error) {
	return g.applicationRepo.FindByID(appId, ctx)
}
