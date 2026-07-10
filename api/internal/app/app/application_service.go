package app

import "api/internal/app/ports"

type Dependencies struct {
	ApplicationRepo ports.ApplicationRepo
	ServerManager   ports.ServerManagerClient
	Clock           ports.Clock
}

type Service struct {
	Create        *CreateApplication
	Update        *UpdateApplication
	Delete        *DeleteApplication
	Get           *GetApplication
	Logs          *GetApplicationLogs
	StreamLogs    *StreamApplicationLogs
	Deploy        *DeployApplication
	GetDeployment *GetApplicationDeployment
	Start         *StartApplication
	Stop          *StopApplication
}

func NewService(d Dependencies) *Service {
	poller := NewDeploymentPoller(d.ApplicationRepo, d.ServerManager, d.Clock)
	return &Service{
		Create: &CreateApplication{
			applicationRepo: d.ApplicationRepo,
		},
		Update: &UpdateApplication{
			applicationRepo: d.ApplicationRepo,
			clock:           d.Clock,
		},
		Delete: &DeleteApplication{
			applicationRepo: d.ApplicationRepo,
			serverManager:   d.ServerManager,
		},
		Get: &GetApplication{
			applicationRepo: d.ApplicationRepo,
		},
		Logs: &GetApplicationLogs{
			applicationRepo: d.ApplicationRepo,
			serverManager:   d.ServerManager,
		},
		StreamLogs: &StreamApplicationLogs{
			applicationRepo: d.ApplicationRepo,
			serverManager:   d.ServerManager,
		},
		Deploy: &DeployApplication{
			applicationRepo: d.ApplicationRepo,
			serverManager:   d.ServerManager,
			clock:           d.Clock,
			poller:          poller,
		},
		GetDeployment: &GetApplicationDeployment{
			applicationRepo: d.ApplicationRepo,
			serverManager:   d.ServerManager,
		},
		Start: &StartApplication{
			applicationRepo: d.ApplicationRepo,
			serverManager:   d.ServerManager,
			clock:           d.Clock,
		},
		Stop: &StopApplication{
			applicationRepo: d.ApplicationRepo,
			serverManager:   d.ServerManager,
			clock:           d.Clock,
		},
	}
}
