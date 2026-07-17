package db

import (
	"api/internal/app/ports"
	"context"
	"errors"
)

type GetDatabaseMetrics struct {
	databaseRepo  ports.DatabaseRepo
	serverManager ports.ServerManagerClient
}

func (g *GetDatabaseMetrics) Execute(ctx context.Context, subdomainID, dbID, rng string) (*ports.MetricsResponse, error) {
	database, err := g.databaseRepo.FindByID(dbID, ctx)
	if err != nil {
		return nil, err
	}
	if database.SubdomainID.String() != subdomainID {
		return nil, ErrNotInSubdomain
	}

	metrics, err := g.serverManager.DatabaseMetrics(ctx, dbID, rng)
	if errors.Is(err, ports.ErrDatabaseNotProvisioned) {
		// Not provisioned (yet) is not an error to the dashboard — just no
		// samples to show.
		return &ports.MetricsResponse{Range: rng, Series: map[string][]ports.MetricPoint{}}, nil
	}
	return metrics, err
}
