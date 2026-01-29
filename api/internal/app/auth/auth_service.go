package auth

import "api/internal/app/ports"


type Dependencies struct {
	UserRepo ports.UserRepo
	PasswordRepo ports.PasswordRepo
	Clock ports.Clock
	Hasher ports.PasswordHasher
	JwtTokenizer ports.JwtTokenizer
}

type Service struct {
	Login *LoginAuth
}

func NewService(d Dependencies) *Service {
	return &Service{
		Login: &LoginAuth{userRepo: d.UserRepo,
			passwordRepo: d.PasswordRepo,
			clock: d.Clock,
			passwordHasher: d.Hasher,
			jwtTokenizer: d.JwtTokenizer,},
	}
}
