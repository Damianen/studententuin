package db

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
	"errors"
)

// ErrNotInSubdomain means the database exists but belongs to a different
// subdomain than the request claimed — treated as not-found so the route
// doesn't leak which database IDs exist. Without this check, Get would leak a
// foreign connection string and Delete would destroy a foreign data volume.
var ErrNotInSubdomain = errors.New("database does not belong to this subdomain")

type GetDatabase struct {
	databaseRepo ports.DatabaseRepo
}

func (g *GetDatabase) Execute(ctx context.Context, subdomainID, dbId string) (*domain.Database, error) {
	database, err := g.databaseRepo.FindByID(dbId, ctx)
	if err != nil {
		return nil, err
	}
	if database.SubdomainID.String() != subdomainID {
		return nil, ErrNotInSubdomain
	}
	return database, nil
}
