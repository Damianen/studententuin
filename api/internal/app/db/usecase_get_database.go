package db

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
)

type GetDatabase struct {
	databaseRepo ports.DatabaseRepo
}

func (g *GetDatabase) Execute(ctx context.Context, dbId string) (*domain.Database, error) {
	return g.databaseRepo.FindByID(dbId, ctx)
}
