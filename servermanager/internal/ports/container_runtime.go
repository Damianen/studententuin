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
	// RemoveContainer force-removes the container (and anonymous volumes) but
	// keeps its per-app networks — deploy cutover replaces the container in
	// place, and the network must survive for the app's other attachments.
	RemoveContainer(ctx context.Context, nameOrID string) error
	Inspect(ctx context.Context, nameOrID string) (*domain.ContainerState, error)
	Logs(ctx context.Context, nameOrID string, opts domain.LogOptions) ([]domain.LogLine, error)
	// FollowLogs streams lines as the container produces them, until ctx is
	// canceled or the container's log stream ends. The channel is closed when
	// the stream ends; setup failures are returned synchronously.
	FollowLogs(ctx context.Context, nameOrID string, opts domain.LogOptions) (<-chan domain.LogLine, error)
	// PullImage fetches ref from its registry unless it is already present on
	// the host. ref comes from the server-side image allowlist, never from a
	// request.
	PullImage(ctx context.Context, ref string) error
	// ConnectNetwork attaches a container to an existing network; it never
	// creates the network (a missing one is the caller's error to surface).
	ConnectNetwork(ctx context.Context, networkName, nameOrID string) error
	// RemoveNetwork force-disconnects any remaining endpoints, then removes
	// the network. Missing networks map to domain.ErrNotFound.
	RemoveNetwork(ctx context.Context, networkName string) error
	// RemoveVolume removes a named volume; in-use volumes map to
	// domain.ErrConflict, missing ones to domain.ErrNotFound.
	RemoveVolume(ctx context.Context, name string) error
	// ListManagedContainers returns the running containers whose names match
	// the manager's naming scheme (domain.ParseManagedName). Stopped
	// containers are excluded — there is nothing to sample.
	ListManagedContainers(ctx context.Context) ([]domain.ManagedContainer, error)
}
