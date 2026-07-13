package db

import (
	"api/internal/app/ports"
	"context"
	"errors"
)

type GetDatabaseLogs struct {
	databaseRepo  ports.DatabaseRepo
	serverManager ports.ServerManagerClient
}

func (g *GetDatabaseLogs) Execute(ctx context.Context, subdomainID, dbID string, opts ports.LogOptions) ([]ports.LogEntry, error) {
	database, err := g.databaseRepo.FindByID(dbID, ctx)
	if err != nil {
		return nil, err
	}
	if database.SubdomainID.String() != subdomainID {
		return nil, ErrNotInSubdomain
	}

	entries, err := g.serverManager.DatabaseLogs(ctx, dbID, opts)
	if errors.Is(err, ports.ErrDatabaseNotProvisioned) {
		// Not provisioned (yet) is not an error to the dashboard — just no
		// logs to show.
		return []ports.LogEntry{}, nil
	}
	return entries, err
}

type StreamDatabaseLogs struct {
	databaseRepo  ports.DatabaseRepo
	serverManager ports.ServerManagerClient
}

// Execute opens a live log stream after the same ownership validation as the
// one-shot read. ErrDatabaseNotProvisioned propagates here — the controller
// turns it into the 404 that makes the client fall back to polling.
func (s *StreamDatabaseLogs) Execute(ctx context.Context, subdomainID, dbID string, opts ports.LogOptions) (<-chan ports.LogEntry, error) {
	database, err := s.databaseRepo.FindByID(dbID, ctx)
	if err != nil {
		return nil, err
	}
	if database.SubdomainID.String() != subdomainID {
		return nil, ErrNotInSubdomain
	}
	return s.serverManager.StreamDatabaseLogs(ctx, dbID, opts)
}
