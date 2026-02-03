package user

import (
	"api/internal/app/ports"
	"context"
)

type DeleteUser struct {
	userRepo ports.UserRepo
}

func (d *DeleteUser) Execute(ctx context.Context, userID string) error {
	user, err := d.userRepo.FindByID(userID, ctx)
	if err != nil {
		return err
	}

	return d.userRepo.Delete(user, ctx)
}
