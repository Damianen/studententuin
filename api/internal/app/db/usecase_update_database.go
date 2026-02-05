package db

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
)

type UpdateDatabse struct {
	databaseRepo ports.DatabaseRepo
	clock ports.Clock
}

type DatabaseUpdateInput struct {
	ID string
	Name *string
	Type *domain.DatabaseType
	Status *domain.DatabaseStatus
	ConnectionString *string
	Host *string
	Port *int
	DockerImage *string
	DockerContainerID *string
	DockerContainerName *string
	Volumes []string
	MemoryLimit *string
	CpuLimit *string
}

func (u *UpdateDatabse) Execute(ctx context.Context, di DatabaseUpdateInput) error {
	return nil
}
