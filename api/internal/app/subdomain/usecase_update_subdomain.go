package subdomain

import (
	"api/internal/app/ports"
	"context"
)

type UpdateSubdomain struct {
	subdomainRepo ports.SubdomainRepo
	clock ports.Clock
}

type SubdomainUpdateInput struct {
	ID string
	Name *string
	FullDomain *string
}

func (u *UpdateSubdomain) Execute(ctx context.Context, si SubdomainUpdateInput) error {
	now := u.clock.Now()
	subdomain, err := u.subdomainRepo.FindByID(si.ID, ctx)
	if err != nil {
		return err
	}

	if si.FullDomain != nil {
		subdomain.FullDomain = *si.FullDomain
	}

	if si.Name != nil {
		subdomain.Name = *si.Name
	}

	subdomain.UpdatedAt = now

	return u.subdomainRepo.Update(subdomain, ctx)
}
