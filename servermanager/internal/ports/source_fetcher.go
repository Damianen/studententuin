package ports

import (
	"context"

	"servermanager/internal/domain"
)

// SourceFetcher clones a repository into a throwaway temp directory.
// Implementations re-validate the URL (defense in depth — the deploy usecase
// already validated it for the synchronous 400) and enforce the post-clone
// size cap. ctx carries the clone-stage timeout.
type SourceFetcher interface {
	Fetch(ctx context.Context, repoURL, branch string) (*domain.SourceCheckout, error)
}
