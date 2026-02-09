package app

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
	"errors"
)

type UpdateApplication struct {
	applicationRepo ports.ApplicationRepo
	clock ports.Clock
}

type ApplicationUpdateInput struct {
	ID string
	Name *string
	RepoUrl *string
	Branch *string
	BuildCommand *string
	StartCommand *string
	Type *domain.ApplicationType
	Status *domain.ApplicationStatus
	MemoryLimit *string
	CpuLimit *string
	Ports *[]int
	Volumes *[]string
	EnvVariables *map[string]string
	DockerImage *string
	DockerContainerID *string
	DockerContainerName *string
}

func (u *UpdateApplication) Execute(ctx context.Context, ai ApplicationUpdateInput) error {
	now := u.clock.Now()
	updates := map[string]any {
		"updated_at": now,
	}

	err := ports.SetIfNotNil(updates, "name", ai.Name, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("name cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "repository_url", ai.RepoUrl, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("repo url cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "branch", ai.Branch, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("branch cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "build_command", ai.BuildCommand, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("repo url cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "start_command", ai.StartCommand, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("repo url cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "memory_limit", ai.MemoryLimit, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("repo url cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "cpu_limit", ai.CpuLimit, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("repo url cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "docker_image", ai.DockerImage, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("docker image cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "docker_container_id", ai.DockerContainerID, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("docker container id cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "docker_container_name", ai.DockerContainerName, func(s string) (any, error) {
		if s == "" {
			return nil, errors.New("docker container name cannot be empty")
		}
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "type", ai.Type, func(s domain.ApplicationType) (any, error) {
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "status", ai.Status, func(s domain.ApplicationStatus) (any, error) {
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "volumes", ai.Volumes, func(s []string) (any, error) {
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "ports", ai.Ports, func(s []int) (any, error) {
		return s, nil
	})
	if err != nil { return err }
	err = ports.SetIfNotNil(updates, "enviroment_variables", ai.EnvVariables, func(s map[string]string) (any, error) {
		return s, nil
	})
	if err != nil { return err }

	if len(updates) == 1 {
		return errors.New("no fields to update!")
	}

	return u.applicationRepo.Update(ai.ID, updates, ctx)
}
