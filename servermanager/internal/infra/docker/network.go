package docker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"servermanager/internal/domain"

	"github.com/moby/moby/client"
)

// ensureNetwork makes sure the app's isolated bridge network exists.
// Idempotent: a concurrent create racing us is fine.
func (r *Runtime) ensureNetwork(ctx context.Context, appID string) error {
	name := domain.AppNetworkName(appID)

	_, err := r.cli.NetworkInspect(ctx, name, client.NetworkInspectOptions{})
	if err == nil {
		return nil
	}
	if !errors.Is(mapErr("", err), domain.ErrNotFound) {
		return mapErr("inspect network "+name, err)
	}

	_, err = r.cli.NetworkCreate(ctx, name, client.NetworkCreateOptions{
		Driver: "bridge",
		Labels: map[string]string{
			"studententuin.managed": "true",
			"studententuin.app-id":  appID,
		},
	})
	if err != nil && !errors.Is(mapErr("", err), domain.ErrConflict) {
		return mapErr("create network "+name, err)
	}
	return nil
}

// removeNetworks tears down per-app networks after their container is gone.
// Removal can briefly race the container's force-remove ("active endpoints"),
// so one retry is allowed; a missing network counts as removed.
func (r *Runtime) removeNetworks(ctx context.Context, names []string) error {
	for _, name := range names {
		err := r.removeNetworkOnce(ctx, name)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return ctx.Err()
			}
			err = r.removeNetworkOnce(ctx, name)
		}
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return fmt.Errorf("remove network %s: %w", name, err)
		}
	}
	return nil
}

func (r *Runtime) removeNetworkOnce(ctx context.Context, name string) error {
	_, err := r.cli.NetworkRemove(ctx, name, client.NetworkRemoveOptions{})
	return mapErr("", err)
}
