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
	return g.userRepo.FindByID(userID, ctx)
}
