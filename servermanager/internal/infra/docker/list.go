package docker

import (
	"context"

	"servermanager/internal/domain"

	"github.com/moby/moby/client"
)

// ListManagedContainers lists the running containers named by the manager's
// scheme. The daemon's name filter is substring-based, so ParseManagedName
// re-validates every hit and drops lookalikes.
func (r *Runtime) ListManagedContainers(ctx context.Context) ([]domain.ManagedContainer, error) {
	res, err := r.cli.ContainerList(ctx, client.ContainerListOptions{
		All:     false, // running only — stopped containers have nothing to sample
		Filters: client.Filters{}.Add("name", "stt-app-", "stt-db-"),
	})
	if err != nil {
		return nil, mapErr("list managed containers", err)
	}

	var out []domain.ManagedContainer
	for _, summary := range res.Items {
		for _, name := range summary.Names {
			managed, ok := domain.ParseManagedName(name)
			if !ok {
				continue
			}
			managed.ID = summary.ID
			out = append(out, managed)
			break
		}
	}
	return out, nil
}
