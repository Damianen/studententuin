package app

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"

	"github.com/google/uuid"
)

type CreateApplication struct {
	applicationRepo ports.ApplicationRepo
	subdomainRepo ports.SubdomainRepo
}

type ApplicationInput struct {
	SubdomainID string
	Name string
	RepoUrl *string
	Branch *string
	BuildCommand *string
	StartCommand *string
	Type domain.ApplicationType
	Status domain.ApplicationStatus
	MemoryLimit *string
	CpuLimit *string
}

func (c *CreateApplication) Execute(ctx context.Context, ai ApplicationInput) error {
	id, err := uuid.Parse(ai.SubdomainID)
	if err != nil {
		return err
	}

	application := &domain.Application{
		ID: uuid.New(),
		SubdomainID: id,
		Name: &ai.Name,
		RepositoryURL: ai.RepoUrl,
		Branch: ai.Branch,
		BuildCommand: ai.BuildCommand,
		StartCommand: ai.StartCommand,
		Type: ai.Type,
		Status: ai.Status,
		MemoryLimit: ai.MemoryLimit,
		CpuLimit: ai.CpuLimit,
	}

	return c.applicationRepo.Create(application, ctx)
}
