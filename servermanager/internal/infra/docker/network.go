package docker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"servermanager/internal/domain"

	"github.com/moby/moby/client"
)

// ensureNetwork makes sure an isolated bridge network exists.
// Idempotent: a concurrent create racing us is fine.
func (r *Runtime) ensureNetwork(ctx context.Context, name string, labels map[string]string) error {
	_, err := r.cli.NetworkInspect(ctx, name, client.NetworkInspectOptions{})
	if err == nil {
		return nil
	}
	if !errors.Is(mapErr("", err), domain.ErrNotFound) {
		return mapErr("inspect network "+name, err)
	}

	withManaged := map[string]string{"studententuin.managed": "true"}
	for k, v := range labels {
		withManaged[k] = v
	}
	_, err = r.cli.NetworkCreate(ctx, name, client.NetworkCreateOptions{
		Driver: "bridge",
		Labels: withManaged,
	})
	if err != nil && !errors.Is(mapErr("", err), domain.ErrConflict) {
		return mapErr("create network "+name, err)
	}
	return nil
}

// ConnectNetwork attaches a container to an existing network. It deliberately
// never creates the network: a missing database network means the database
// isn't provisioned, and the caller must surface that, not paper over it.
func (r *Runtime) ConnectNetwork(ctx context.Context, networkName, nameOrID string) error {
	_, err := r.cli.NetworkConnect(ctx, networkName, client.NetworkConnectOptions{Container: nameOrID})
	return mapErr("connect "+nameOrID+" to "+networkName, err)
}

// RemoveNetwork force-disconnects whatever is still attached (e.g. a linked
// app container outliving its database), then removes the network.
func (r *Runtime) RemoveNetwork(ctx context.Context, networkName string) error {
	res, err := r.cli.NetworkInspect(ctx, networkName, client.NetworkInspectOptions{})
	if err != nil {
		return mapErr("inspect network "+networkName, err)
	}
	for id := range res.Network.Containers {
		_, err := r.cli.NetworkDisconnect(ctx, networkName, client.NetworkDisconnectOptions{Container: id, Force: true})
		if err != nil && !errors.Is(mapErr("", err), domain.ErrNotFound) {
			return mapErr("disconnect "+id+" from "+networkName, err)
		}
	}
	return r.removeNetworks(ctx, []string{networkName})
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
