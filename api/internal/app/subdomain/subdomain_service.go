package subdomain

import "api/internal/app/ports"

type Dependencies struct {
	SubdomainRepo ports.SubdomainRepo
	UserRepo ports.UserRepo
	clock ports.Clock
}

type Service struct {
	Create *CreateSubdomain
	Update *UpdateSubdomain
	Delete *DeleteSubdomain
	Get *GetSubdomain
	CheckUser *CheckUser
}

func NewService(d Dependencies) *Service {
	return &Service{
		Create: &CreateSubdomain{
			subdomainRepo: d.SubdomainRepo,
		},
		Update: &UpdateSubdomain{
			subdomainRepo: d.SubdomainRepo,
			clock: d.clock,
		},
		Delete: &DeleteSubdomain{
			subdomainRepo: d.SubdomainRepo,
		},
		Get: &GetSubdomain{
			subdomainRepo: d.SubdomainRepo,
		},
		CheckUser: &CheckUser{
			subdomainRepo: d.SubdomainRepo,
			userRepo: d.UserRepo,
		},
	}
}
