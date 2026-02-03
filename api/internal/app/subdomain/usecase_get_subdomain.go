package subdomain

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
)


type GetSubdomain struct {
	subdomainRepo ports.SubdomainRepo
}

func (g *GetSubdomain) Execute(ctx context.Context, subdomainID string) (*domain.Subdomain, error) {
	subdomain, err := g.subdomainRepo.FindByID(subdomainID, ctx)
	if err != nil {
		return nil, err
	}

	return subdomain, nil
}
