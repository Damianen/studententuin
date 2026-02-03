package auth

import (
	"api/internal/app/ports"
	"context"
	"errors"
)

type LoginAuth struct {
	userRepo ports.UserRepo
	passwordRepo ports.PasswordRepo
	clock ports.Clock
	passwordHasher ports.PasswordHasher
	jwtTokenizer ports.JwtTokenizer
}

type LoginInput struct {
	Email string
	Password string
}

func (c *LoginAuth) Execute(ctx context.Context, li LoginInput) (string, error) {
	user, err := c.userRepo.FindByEmail(li.Email, ctx)
	if err != nil   {
		return "", err
	}

	passwordCred, err := c.passwordRepo.FindById(user.ID.String(), ctx)
	if err != nil {
		return "", err
	}

	same, err := c.passwordHasher.Compare(li.Password, passwordCred.PasswordHash)
	if err != nil {
		return "", err
	}
	if !same {
		return "", errors.New("invalid credentials")
	}

	token, err := c.jwtTokenizer.CreateToken(user.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}
