package db

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
	"errors"
)

type UpdateDatabase struct {
	databaseRepo ports.DatabaseRepo
	clock ports.Clock
}

type DatabaseUpdateInput struct {
	ID string
	Name *string
	Type *domain.DatabaseType
	Status *domain.DatabaseStatus
	Version *string
	ConnectionString *string
	DockerImage *string
	DockerContainerID *string
	DockerContainerName *string
	Volumes *[]string
	MemoryLimit *string
	CpuLimit *string
}

func (u *UpdateDatabase) Execute(ctx context.Context, di DatabaseUpdateInput) error {
	now := u.clock.Now()
	updates := map[string]any{
		"updated_at": now,
	}

	err := ports.SetIfNotNil(updates, "name", di.Name, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("name cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "version", di.Version, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("version cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "type", di.Type, func(s domain.DatabaseType) (any, error) {
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "status", di.Status, func(s domain.DatabaseStatus) (any, error) {
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "connection_string", di.ConnectionString, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("connection string cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "docker_image", di.DockerImage, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("docker image cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "docker_container_id", di.DockerContainerID, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("docker container id cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "docker_container_name", di.DockerContainerName, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("docker container name cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "memory_limit", di.MemoryLimit, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("memory limit cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "cpu_limit", di.CpuLimit, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("cpu limit cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "volumes", di.Volumes, func(s []string) (any, error) {
		return s, nil
	})
	if err != nil { return err }
	if len(updates) == 1 {
		return errors.New("no fields to update!")
	}

	return u.databaseRepo.Update(di.ID, updates, ctx)
}
