package subdomain

import (
	"api/internal/app/ports"
	"context"
)

type CheckUser struct {
	userRepo ports.UserRepo
	subdomainRepo ports.SubdomainRepo
}

func (c *CheckUser) Execute(ctx context.Context, userID string, subdomainID string) (bool, error) {
	subdomain, err := c.subdomainRepo.FindByID(subdomainID, ctx)
	if err != nil {
		return false, err
	}

	return subdomain.UserID.String() == userID, nil
}
