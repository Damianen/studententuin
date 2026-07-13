package databases

import (
	"context"
	"errors"
	"fmt"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

type DeleteDatabase struct {
	runtime ports.ContainerRuntime
	jobs    *JobStore
}

// Execute sweeps everything the database owns: container, private network
// (force-disconnecting a linked app that is still attached), and the named
// data volume. Every step tolerates "already gone", so the sweep is
// idempotent and also cleans up half-provisioned wreckage.
func (u *DeleteDatabase) Execute(ctx context.Context, dbID string) error {
	if job, ok := u.jobs.Latest(dbID); ok && !job.Terminal() {
		return fmt.Errorf("provision %s still in flight for this database: %w", job.ID, domain.ErrConflict)
	}

	name := domain.DBContainerName(dbID)
	if err := u.runtime.Stop(ctx, name); err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}
	if err := u.runtime.RemoveContainer(ctx, name); err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}
	if err := u.runtime.RemoveNetwork(ctx, domain.DBNetworkName(dbID)); err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}
	// The volume holds the data — last, and only after the container is gone.
	// ErrConflict propagates: an in-use volume means something unexpected
	// still mounts it, and silently keeping the data around would be worse.
	if err := u.runtime.RemoveVolume(ctx, domain.DBVolumeName(dbID)); err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	u.jobs.Drop(dbID)
	return nil
}
