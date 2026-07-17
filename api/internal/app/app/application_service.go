package app

import "api/internal/app/ports"

type Dependencies struct {
	ApplicationRepo ports.ApplicationRepo
	DatabaseRepo    ports.DatabaseRepo
	DeploymentRepo  ports.DeploymentRepo
	ServerManager   ports.ServerManagerClient
	Clock           ports.Clock
}

type Service struct {
	Create          *CreateApplication
	Update          *UpdateApplication
	Delete          *DeleteApplication
	Get             *GetApplication
	Logs            *GetApplicationLogs
	StreamLogs      *StreamApplicationLogs
	Metrics         *GetApplicationMetrics
	Deploy          *DeployApplication
	GetDeployment   *GetApplicationDeployment
	ListDeployments *ListApplicationDeployments
	Start           *StartApplication
	Stop            *StopApplication
}

func NewService(d Dependencies) *Service {
	poller := NewDeploymentPoller(d.ApplicationRepo, d.DeploymentRepo, d.ServerManager, d.Clock)
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
		Metrics: &GetApplicationMetrics{
			applicationRepo: d.ApplicationRepo,
			serverManager:   d.ServerManager,
		},
		Deploy: &DeployApplication{
			applicationRepo: d.ApplicationRepo,
			databaseRepo:    d.DatabaseRepo,
			deploymentRepo:  d.DeploymentRepo,
			serverManager:   d.ServerManager,
			clock:           d.Clock,
			poller:          poller,
		},
		GetDeployment: &GetApplicationDeployment{
			applicationRepo: d.ApplicationRepo,
			serverManager:   d.ServerManager,
		},
		ListDeployments: &ListApplicationDeployments{
			applicationRepo: d.ApplicationRepo,
			deploymentRepo:  d.DeploymentRepo,
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
