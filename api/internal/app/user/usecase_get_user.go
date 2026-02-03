package user

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
)

type GetUser struct {
	userRepo ports.UserRepo
}

func (g *GetUser) Execute(ctx context.Context, userID string) (*domain.User, error) {
	user, err := g.userRepo.FindByID(userID, ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}
