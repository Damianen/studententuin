package user

import (
	"api/internal/app/ports"
	"context"
	"errors"
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
	updates := map[string]any{
		"updated_at": now,
	}

	err := ports.SetIfNotNil(updates, "email", ui.Email, func(s string) (any, error) {
	    s = strings.ToLower(strings.TrimSpace(s))
	    if s == "" {
	        return nil, errors.New("email cannot be empty")
	    }
	    return s, nil
	})
	if err != nil { return err }

	err = ports.SetIfNotNil(updates, "display_name", ui.Name, func(s string) (any, error) {
	    if s == "" {
	        return nil, errors.New("name cannot be empty")
	    }
	    return s, nil
	})
	if err != nil { return err }

	if len(updates) == 1 {
		return errors.New("no fields to update")
	}

	return u.userRepo.Update(ui.ID, updates, ctx)
}
