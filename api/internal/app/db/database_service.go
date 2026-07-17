package db

import "api/internal/app/ports"

type Dependencies struct {
	DatabaseRepo  ports.DatabaseRepo
	Clock         ports.Clock
	ServerManager ports.ServerManagerClient
}

type Service struct {
	Create     *CreateDatabase
	Update     *UpdateDatabase
	Delete     *DeleteDatabase
	Get        *GetDatabase
	Logs       *GetDatabaseLogs
	StreamLogs *StreamDatabaseLogs
	Metrics    *GetDatabaseMetrics
	Poller     *ProvisionPoller
}

func NewService(d Dependencies) *Service {
	poller := NewProvisionPoller(d.DatabaseRepo, d.ServerManager, d.Clock)
	return &Service{
		Create: &CreateDatabase{
			databaseRepo:  d.DatabaseRepo,
			serverManager: d.ServerManager,
			poller:        poller,
		},
		Update: &UpdateDatabase{
			databaseRepo: d.DatabaseRepo,
			clock:        d.Clock,
		},
		Delete: &DeleteDatabase{
			databaseRepo:  d.DatabaseRepo,
			serverManager: d.ServerManager,
		},
		Get: &GetDatabase{
			databaseRepo: d.DatabaseRepo,
		},
		Logs: &GetDatabaseLogs{
			databaseRepo:  d.DatabaseRepo,
			serverManager: d.ServerManager,
		},
		StreamLogs: &StreamDatabaseLogs{
			databaseRepo:  d.DatabaseRepo,
			serverManager: d.ServerManager,
		},
		Metrics: &GetDatabaseMetrics{
			databaseRepo:  d.DatabaseRepo,
			serverManager: d.ServerManager,
		},
		Poller: poller,
	}
}
