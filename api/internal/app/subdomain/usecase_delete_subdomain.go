package subdomain

import (
	"api/internal/app/ports"
	"context"
)

type DeleteSubdomain struct {
	subdomainRepo ports.SubdomainRepo
}

func (d *DeleteSubdomain) Execute(ctx context.Context, subdomainID string) error {
	subdomain, err := d.subdomainRepo.FindByID(subdomainID, ctx)
	if err != nil {
		return err
	}

	return d.subdomainRepo.Delete(subdomain, ctx)
}
