package subdomain

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
)

type GetAllSubdomains struct {
	subdomainRepo ports.SubdomainRepo
}

func (g *GetAllSubdomains) Execute(ctx context.Context, userID string) ([]domain.Subdomain, error) {
	return g.subdomainRepo.FindAllByUserID(userID, ctx)
}
