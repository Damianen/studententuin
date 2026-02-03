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
	Update *UpdateUser
	Delete *DeleteUser
	Get *GetUser
}

func NewService(d Dependencies) *Service {
	return &Service{
		Create: &CreateUser{
			userRepo: d.UserRepo,
			passwordRepo: d.PasswordRepo,
			clock: d.Clock,
			passwordHasher: d.Hasher, },
		Update: &UpdateUser{
			userRepo: d.UserRepo,
			clock: d.Clock,
		},
		Delete: &DeleteUser{
			userRepo: d.UserRepo,
		},
		Get: &GetUser{
			userRepo: d.UserRepo,
		},
	}
}
