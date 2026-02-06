package db

import "api/internal/app/ports"

type Dependencies struct {
	DatabaseRepo ports.DatabaseRepo
	Clock ports.Clock
}

type Service struct {
	Create *CreateDatabase
	Update *UpdateDatabase
	Delete *DeleteDatabase
	Get *GetDatabase
}

func NewService(d Dependencies)  *Service {
	return &Service{
		Create: &CreateDatabase{
			databaseRepo: d.DatabaseRepo,
		},
		Update: &UpdateDatabase{
			databaseRepo: d.DatabaseRepo,
			clock: d.Clock,
		},
		Delete: &DeleteDatabase{
			databaseRepo: d.DatabaseRepo,
		},
		Get: &GetDatabase{
			databaseRepo: d.DatabaseRepo,
		},
	}
}
