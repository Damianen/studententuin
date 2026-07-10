// Package git implements ports.SourceFetcher with go-git — pure Go, so no
// git binary and no argv injection surface (§3.4).
package git

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"servermanager/internal/domain"
	"servermanager/internal/ports"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type Fetcher struct {
	hosts    []string
	maxBytes int64
}

var _ ports.SourceFetcher = (*Fetcher)(nil)

func NewFetcher(hosts []string, maxBytes int64) *Fetcher {
	return &Fetcher{hosts: hosts, maxBytes: maxBytes}
}

// Fetch validates the URL (defense in depth — the deploy usecase already
// rejected bad input synchronously), clones depth-1 into a throwaway temp
// dir, and enforces the size cap. Every error path removes the dir.
func (f *Fetcher) Fetch(ctx context.Context, repoURL, branch string) (*domain.SourceCheckout, error) {
	if err := domain.ValidateRepoURL(repoURL, f.hosts); err != nil {
		return nil, err
	}
	if err := domain.ValidateBranch(branch); err != nil {
		return nil, err
	}
	return f.fetchInto(ctx, repoURL, branch)
}

// fetchInto is the post-validation clone path. Integration tests enter here
// with a local fixture path; the https/allowlist gate stays covered by the
// domain unit tests.
func (f *Fetcher) fetchInto(ctx context.Context, repoURL, branch string) (*domain.SourceCheckout, error) {
	dir, err := os.MkdirTemp("", "stt-clone-*")
	if err != nil {
		return nil, fmt.Errorf("creating clone dir: %w", err)
	}
	cleanup := domain.OnceCleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			slog.Warn("removing clone dir", "dir", dir, "error", err)
		}
	})

	opts := &gogit.CloneOptions{
		URL:          repoURL,
		Depth:        1,
		SingleBranch: true,
		Tags:         gogit.NoTags,
	}
	if branch != "" {
		opts.ReferenceName = plumbing.NewBranchReferenceName(branch)
	}

	repo, err := gogit.PlainCloneContext(ctx, dir, false, opts)
	if err != nil {
		cleanup()
		return nil, mapCloneErr(ctx, err, branch)
	}

	if err := checkTreeSize(dir, f.maxBytes); err != nil {
		cleanup()
		return nil, err
	}

	head, err := repo.Head()
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("reading clone HEAD: %w", err)
	}

	return &domain.SourceCheckout{
		Dir:       dir,
		CommitSHA: head.Hash().String(),
		Cleanup:   cleanup,
	}, nil
}

// mapCloneErr turns go-git failures into the user-facing messages that end
// up on the deployment job.
func mapCloneErr(ctx context.Context, err error, branch string) error {
	var noRef gogit.NoMatchingRefSpecError
	switch {
	case ctx.Err() != nil:
		return errors.New("clone timed out — repository too large or host unreachable")
	case errors.Is(err, transport.ErrAuthenticationRequired), errors.Is(err, transport.ErrAuthorizationFailed):
		return errors.New("repository not found or private (only public repositories are supported)")
	case errors.Is(err, transport.ErrRepositoryNotFound):
		return errors.New("repository not found")
	case errors.Is(err, plumbing.ErrReferenceNotFound), errors.As(err, &noRef):
		return fmt.Errorf("branch %q not found in the repository", branch)
	default:
		return fmt.Errorf("clone failed: %w", err)
	}
}

// checkTreeSize sums regular files (packfiles under .git included — they are
// the disk cost) and rejects trees over the cap.
func checkTreeSize(dir string, maxBytes int64) error {
	var total int64
	err := filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if total += info.Size(); total > maxBytes {
			return fmt.Errorf("repository exceeds the %d MiB size limit", maxBytes/(1024*1024))
		}
		return nil
	})
	return err
}
