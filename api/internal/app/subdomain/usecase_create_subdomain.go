package subdomain

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"

	"github.com/google/uuid"
)

type CreateSubdomain struct {
	subdomainRepo ports.SubdomainRepo
}

type SubdomainInput struct {
	Name string
	FullDomain string
	UserID string
}

func (c *CreateSubdomain) Execute(ctx context.Context, si SubdomainInput) error {
	id, err := uuid.Parse(si.UserID)
	if err != nil {
		return err
	}

	subdomain := &domain.Subdomain{
		ID: uuid.New(),
		Name: si.Name,
		FullDomain: si.FullDomain,
		UserID: id,
		IsActive: true,
	}

	return c.subdomainRepo.Create(subdomain, ctx)
}
