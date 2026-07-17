package app

import (
	"api/internal/app/ports"
	"context"
	"errors"
)

type GetApplicationMetrics struct {
	applicationRepo ports.ApplicationRepo
	serverManager   ports.ServerManagerClient
}

func (g *GetApplicationMetrics) Execute(ctx context.Context, subdomainID, appID, rng string) (*ports.MetricsResponse, error) {
	application, err := g.applicationRepo.FindByID(appID, ctx)
	if err != nil {
		return nil, err
	}
	if application.SubdomainID.String() != subdomainID {
		return nil, ErrNotInSubdomain
	}

	metrics, err := g.serverManager.Metrics(ctx, appID, rng)
	if errors.Is(err, ports.ErrAppNotDeployed) {
		// Never deployed is not an error to the dashboard — just no samples
		// yet.
		return &ports.MetricsResponse{Range: rng, Series: map[string][]ports.MetricPoint{}}, nil
	}
	return metrics, err
}
