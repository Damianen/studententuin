package ports

import (
	"context"

	"servermanager/internal/domain"
)

// ContainerRuntime abstracts the Docker SDK so services stay unit-testable.
// nameOrID is a derived container name (domain.AppContainerName) or Docker ID.
type ContainerRuntime interface {
	Create(ctx context.Context, spec domain.ContainerSpec) (containerID string, err error)
	Start(ctx context.Context, nameOrID string) error
	Stop(ctx context.Context, nameOrID string) error
	Remove(ctx context.Context, nameOrID string) error
	Inspect(ctx context.Context, nameOrID string) (*domain.ContainerState, error)
	Logs(ctx context.Context, nameOrID string, opts domain.LogOptions) ([]domain.LogLine, error)
}
