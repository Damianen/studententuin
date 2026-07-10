package app

import "api/internal/app/ports"

type Dependencies struct {
	ApplicationRepo ports.ApplicationRepo
	ServerManager ports.ServerManagerClient
	Clock ports.Clock
}

type Service struct {
	Create *CreateApplication
	Update *UpdateApplication
	Delete *DeleteApplication
	Get *GetApplication
	Logs *GetApplicationLogs
	StreamLogs *StreamApplicationLogs
}

func NewService(d Dependencies) *Service {
	return &Service{
		Create: &CreateApplication{
			applicationRepo: d.ApplicationRepo,
		},
		Update: &UpdateApplication{
			applicationRepo: d.ApplicationRepo,
			clock: d.Clock,
		},
		Delete: &DeleteApplication{
			applicationRepo: d.ApplicationRepo,
		},
		Get: &GetApplication{
			applicationRepo: d.ApplicationRepo,
		},
		Logs: &GetApplicationLogs{
			applicationRepo: d.ApplicationRepo,
			serverManager: d.ServerManager,
		},
		StreamLogs: &StreamApplicationLogs{
			applicationRepo: d.ApplicationRepo,
			serverManager: d.ServerManager,
		},
	}
}
