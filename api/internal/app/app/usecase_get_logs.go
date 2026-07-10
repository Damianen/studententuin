package app

import (
	"api/internal/app/ports"
	"context"
	"errors"
)

// ErrNotInSubdomain means the application exists but belongs to a different
// subdomain than the request claimed — treated as not-found so the route
// doesn't leak which application IDs exist.
var ErrNotInSubdomain = errors.New("application does not belong to this subdomain")

type GetApplicationLogs struct {
	applicationRepo ports.ApplicationRepo
	serverManager   ports.ServerManagerClient
}

func (g *GetApplicationLogs) Execute(ctx context.Context, subdomainID, appID string, opts ports.LogOptions) ([]ports.LogEntry, error) {
	application, err := g.applicationRepo.FindByID(appID, ctx)
	if err != nil {
		return nil, err
	}
	if application.SubdomainID.String() != subdomainID {
		return nil, ErrNotInSubdomain
	}

	entries, err := g.serverManager.Logs(ctx, appID, opts)
	if errors.Is(err, ports.ErrAppNotDeployed) {
		// Never deployed is not an error to the dashboard — just no logs yet.
		return []ports.LogEntry{}, nil
	}
	return entries, err
}

type StreamApplicationLogs struct {
	applicationRepo ports.ApplicationRepo
	serverManager   ports.ServerManagerClient
}

// Execute opens a live log stream after the same ownership validation as the
// one-shot read. ErrAppNotDeployed propagates here — the controller decides
// how a live view of a non-existent container degrades.
func (s *StreamApplicationLogs) Execute(ctx context.Context, subdomainID, appID string, opts ports.LogOptions) (<-chan ports.LogEntry, error) {
	application, err := s.applicationRepo.FindByID(appID, ctx)
	if err != nil {
		return nil, err
	}
	if application.SubdomainID.String() != subdomainID {
		return nil, ErrNotInSubdomain
	}
	return s.serverManager.StreamLogs(ctx, appID, opts)
}
