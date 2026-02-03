package user

import (
	"api/internal/app/ports"
	"context"
	"strings"
)

type UpdateUser struct {
	userRepo ports.UserRepo
	clock ports.Clock
}

type UserUpdateInput struct {
	ID string
	Email *string
	Name *string
}

func (u *UpdateUser) Execute(ctx context.Context, ui UserUpdateInput) error {
	now := u.clock.Now()
	user, err := u.userRepo.FindByID(ui.ID, ctx)
	if err != nil {
		return err
	}

	if ui.Email != nil {
		user.Email = strings.ToLower(strings.TrimSpace(*ui.Email))
	}

	if ui.Name != nil {
		user.DisplayName = *ui.Name
	}

	user.UpdatedAt = now

	return u.userRepo.Update(user, ctx)
}
