package app

import (
	"api/internal/app/ports"
	"context"
)

type GetApplicationDeployment struct {
	applicationRepo ports.ApplicationRepo
	serverManager   ports.ServerManagerClient
}

// Execute proxies the manager's deployment status after the same ownership
// check as every other application route — the frontend polls this for
// stage/error/build-log feedback.
func (g *GetApplicationDeployment) Execute(ctx context.Context, subdomainID, appID, deploymentID string) (*ports.DeploymentStatus, error) {
	application, err := g.applicationRepo.FindByID(appID, ctx)
	if err != nil {
		return nil, err
	}
	if application.SubdomainID.String() != subdomainID {
		return nil, ErrNotInSubdomain
	}

	status, err := g.serverManager.DeploymentStatus(ctx, deploymentID)
	if err != nil {
		return nil, err
	}
	if status.AppID != appID {
		// A deployment ID belonging to someone else's app: same not-found
		// treatment as the ownership check above.
		return nil, ErrNotInSubdomain
	}
	return status, nil
}
