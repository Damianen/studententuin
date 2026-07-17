package ports

import "context"

// ContainerExec runs fixed-argv commands inside a managed container and reads
// its configuration — the db metrics sampler's window into postgres. argv is
// always a server-side constant plus re-validated identifiers, never request
// text.
type ContainerExec interface {
	// ExecCapture runs argv in the container and returns its stdout (bounded).
	// A non-zero exit is an error carrying a short stderr tail; a stopped
	// container maps to domain.ErrConflict, a missing one to domain.ErrNotFound.
	ExecCapture(ctx context.Context, nameOrID string, argv []string) (stdout string, err error)
	// ContainerEnv returns only the requested keys from the container's
	// configured environment. Values may be secrets — never log them.
	ContainerEnv(ctx context.Context, nameOrID string, keys ...string) (map[string]string, error)
}
