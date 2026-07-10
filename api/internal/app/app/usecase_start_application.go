package app

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
	"fmt"
)

type StartApplication struct {
	applicationRepo ports.ApplicationRepo
	serverManager   ports.ServerManagerClient
	clock           ports.Clock
}

// Execute starts the app's existing container and records the new status.
// ErrAppNotDeployed propagates when there is no container yet.
func (s *StartApplication) Execute(ctx context.Context, subdomainID, appID string) error {
	application, err := s.applicationRepo.FindByID(appID, ctx)
	if err != nil {
		return err
	}
	if application.SubdomainID.String() != subdomainID {
		return ErrNotInSubdomain
	}

	if err := s.serverManager.Start(ctx, appID); err != nil {
		return err
	}

	if err := s.applicationRepo.Update(appID, map[string]any{
		"status":     domain.ApplicationStatusRunning,
		"updated_at": s.clock.Now(),
	}, ctx); err != nil {
		fmt.Println("recording started application:", err.Error())
	}
	return nil
}
