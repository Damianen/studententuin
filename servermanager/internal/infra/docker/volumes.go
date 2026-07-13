package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// RemoveVolume removes a named volume. The daemon auto-creates named volumes
// on container create, so removal is the only volume operation the manager
// needs. In-use volumes surface as ErrConflict, missing ones as ErrNotFound.
func (r *Runtime) RemoveVolume(ctx context.Context, name string) error {
	_, err := r.cli.VolumeRemove(ctx, name, client.VolumeRemoveOptions{})
	return mapErr("remove volume "+name, err)
}
