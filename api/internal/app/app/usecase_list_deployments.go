package app

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
)

type ListApplicationDeployments struct {
	applicationRepo ports.ApplicationRepo
	deploymentRepo  ports.DeploymentRepo
}

// Execute returns the app's persisted deploy history, newest first, after
// the same ownership validation as every other application route.
func (l *ListApplicationDeployments) Execute(ctx context.Context, subdomainID, appID string, limit int) ([]domain.Deployment, error) {
	application, err := l.applicationRepo.FindByID(appID, ctx)
	if err != nil {
		return nil, err
	}
	if application.SubdomainID.String() != subdomainID {
		return nil, ErrNotInSubdomain
	}
	return l.deploymentRepo.ListByApplication(appID, limit, ctx)
}
