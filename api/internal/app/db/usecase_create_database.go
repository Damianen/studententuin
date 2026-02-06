package db

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"

	"github.com/google/uuid"
)

type CreateDatabase struct {
	databaseRepo ports.DatabaseRepo
	subdomainRepo ports.SubdomainRepo
}

type DatabaseInput struct {
	SubdomainID string
	Name string
	Version string
	Type domain.DatabaseType
	Status domain.DatabaseStatus
	DockerImage *string
	Volumes []string
	MemoryLimit *string
	CpuLimit *string
}

func (c *CreateDatabase) Execute(ctx context.Context, di DatabaseInput) error {
	id, err := uuid.Parse(di.SubdomainID)
	if err != nil {
		return err
	}

	database := &domain.Database{
		ID: uuid.New(),
		SubdomainID: id,
		Name: di.Name,
		Version: di.Version,
		Type: di.Type,
		Status: di.Status,
		DockerImage: di.DockerImage,
		Volumes: di.Volumes,
		MemoryLimit: di.MemoryLimit,
		CpuLimit: di.CpuLimit,
	}

	return c.databaseRepo.Create(database, ctx)
}
