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
	return g.subdomainRepo.FindByID(subdomainID, ctx)
}
