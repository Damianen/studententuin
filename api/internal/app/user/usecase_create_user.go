package user

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type CreateUser struct {
	userRepo ports.UserRepo
	passwordRepo ports.PasswordRepo
	clock ports.Clock
	passwordHasher ports.PasswordHasher
}

type UserInput struct {
	Email string
	Password string
	Name string
}

func (c *CreateUser) Execute(ctx context.Context, ui UserInput) error {
	now := c.clock.Now()

	email := strings.ToLower(strings.TrimSpace(ui.Email))
	if email == "" || ui.Password == "" {
		return errors.New("email and password required")
	}

	passwordHash, err := c.passwordHasher.Hash(ui.Password)
	if err != nil {
		return err
	}

	user := &domain.User{
		ID: uuid.New(),
		Email: email,
		DisplayName: ui.Name,
		Status: "active",
	}

	// The user row must exist before the credential row can reference it.
	err = c.userRepo.Create(user, ctx)
	if err != nil {
		return err
	}

	passwordCred := &domain.PasswordCredential{
		UserId: user.ID,
		PasswordHash: passwordHash,
		PasswordUpdatedAt: now,
	}

	return c.passwordRepo.Create(passwordCred, ctx)
}
