package user

import "api/internal/app/ports"

type Dependencies struct {
	UserRepo ports.UserRepo
	PasswordRepo ports.PasswordRepo
	Clock ports.Clock
	Hasher ports.PasswordHasher
}

type Service struct {
	Create *CreateUser
}

func NewService(d Dependencies) *Service {
	return &Service{
		Create: &CreateUser{userRepo: d.UserRepo,
			passwordRepo: d.PasswordRepo,
			clock: d.Clock,
			passwordHasher: d.Hasher, },
	}
}
