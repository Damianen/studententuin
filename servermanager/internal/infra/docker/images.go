package docker

import (
	"context"

	"servermanager/internal/domain"
	"servermanager/internal/ports"

	"github.com/moby/moby/client"
)

var _ ports.ImageStore = (*Runtime)(nil)

// ListAppImages returns every repo:tag the manager built for the app, found
// by the deterministic repository name — never by user input.
func (r *Runtime) ListAppImages(ctx context.Context, appID string) ([]string, error) {
	repo := domain.AppImageRepo(appID)
	res, err := r.cli.ImageList(ctx, client.ImageListOptions{
		Filters: client.Filters{}.Add("reference", repo+":*"),
	})
	if err != nil {
		return nil, mapErr("list images "+repo, err)
	}

	var refs []string
	for _, summary := range res.Items {
		refs = append(refs, summary.RepoTags...)
	}
	return refs, nil
}

// RemoveImage untags one ref; child layers of legacy builds are pruned with
// it. In-use images surface as ErrConflict, missing ones as ErrNotFound.
func (r *Runtime) RemoveImage(ctx context.Context, ref string) error {
	_, err := r.cli.ImageRemove(ctx, ref, client.ImageRemoveOptions{PruneChildren: true})
	return mapErr("remove image "+ref, err)
}
