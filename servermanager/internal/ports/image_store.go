package ports

import "context"

// ImageStore lists and removes the manager-built images of one app.
// Primitives only — the "keep the newest" GC policy lives in the deploy
// usecase so it stays unit-testable.
type ImageStore interface {
	// ListAppImages returns the repo:tag refs of every image built for the
	// app (repository domain.AppImageRepo(appID)).
	ListAppImages(ctx context.Context, appID string) ([]string, error)
	// RemoveImage untags/removes one ref. An in-use image maps to
	// domain.ErrConflict, a missing one to domain.ErrNotFound.
	RemoveImage(ctx context.Context, ref string) error
}
