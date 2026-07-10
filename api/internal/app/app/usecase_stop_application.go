package app

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
	"fmt"
)

type StopApplication struct {
	applicationRepo ports.ApplicationRepo
	serverManager   ports.ServerManagerClient
	clock           ports.Clock
}

// Execute stops the app's container and records the new status.
// ErrAppNotDeployed propagates when there is no container.
func (s *StopApplication) Execute(ctx context.Context, subdomainID, appID string) error {
	application, err := s.applicationRepo.FindByID(appID, ctx)
	if err != nil {
		return err
	}
	if application.SubdomainID.String() != subdomainID {
		return ErrNotInSubdomain
	}

	if err := s.serverManager.Stop(ctx, appID); err != nil {
		return err
	}

	if err := s.applicationRepo.Update(appID, map[string]any{
		"status":     domain.ApplicationStatusStopped,
		"updated_at": s.clock.Now(),
	}, ctx); err != nil {
		fmt.Println("recording stopped application:", err.Error())
	}
	return nil
}
