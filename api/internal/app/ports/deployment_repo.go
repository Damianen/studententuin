package ports

import (
	"api/internal/domain"
	"context"
)

type DeploymentRepo interface {
	Create(deployment *domain.Deployment, context context.Context) error
	// UpdateByManagerID finalizes the row keyed by the manager's deployment
	// ID; a missing row (pre-history deploy) is a no-op, not an error.
	UpdateByManagerID(managerID string, updates map[string]any, context context.Context) error
	// ListByApplication returns the application's deployments newest first,
	// capped at limit.
	ListByApplication(appID string, limit int, context context.Context) ([]domain.Deployment, error)
	// FailInFlight sweeps every in_flight row to failed with the given
	// reason — run at startup, when the in-memory polling state is known
	// lost.
	FailInFlight(reason string, context context.Context) error
}
