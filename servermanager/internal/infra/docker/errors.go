package docker

import (
	"fmt"

	"servermanager/internal/domain"

	cerrdefs "github.com/containerd/errdefs"
)

// mapErr translates Docker client errors into domain sentinels at the adapter
// boundary, so nothing above this package imports Docker error types.
func mapErr(op string, err error) error {
	switch {
	case err == nil:
		return nil
	case cerrdefs.IsNotFound(err):
		return fmt.Errorf("%s: %w", op, domain.ErrNotFound)
	case cerrdefs.IsConflict(err):
		return fmt.Errorf("%s: %w", op, domain.ErrConflict)
	default:
		return fmt.Errorf("%s: %w", op, err)
	}
}
