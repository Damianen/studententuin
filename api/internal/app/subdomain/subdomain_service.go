package subdomain

import "api/internal/app/ports"

type Dependencies struct {
	SubdomainRepo   ports.SubdomainRepo
	UserRepo        ports.UserRepo
	ApplicationRepo ports.ApplicationRepo
	DatabaseRepo    ports.DatabaseRepo
	ServerManager   ports.ServerManagerClient
	Clock           ports.Clock
}

type Service struct {
	Create    *CreateSubdomain
	Update    *UpdateSubdomain
	Delete    *DeleteSubdomain
	Get       *GetSubdomain
	GetAll    *GetAllSubdomains
	CheckUser *CheckUser
}

func NewService(d Dependencies) *Service {
	return &Service{
		Create: &CreateSubdomain{
			subdomainRepo: d.SubdomainRepo,
		},
		Update: &UpdateSubdomain{
			subdomainRepo: d.SubdomainRepo,
			clock:         d.Clock,
		},
		Delete: &DeleteSubdomain{
			subdomainRepo:   d.SubdomainRepo,
			applicationRepo: d.ApplicationRepo,
			databaseRepo:    d.DatabaseRepo,
			serverManager:   d.ServerManager,
		},
		Get: &GetSubdomain{
			subdomainRepo: d.SubdomainRepo,
		},
		GetAll: &GetAllSubdomains{
			subdomainRepo: d.SubdomainRepo,
		},
		CheckUser: &CheckUser{
			subdomainRepo: d.SubdomainRepo,
			userRepo:      d.UserRepo,
		},
	}
}
