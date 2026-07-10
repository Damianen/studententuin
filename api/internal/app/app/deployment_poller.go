package app

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	defaultPollInterval = 2 * time.Second
	// defaultPollBudget comfortably covers the manager's clone+build budgets
	// (2m + 10m) plus slack.
	defaultPollBudget = 15 * time.Minute
	// maxNotFoundStrikes: consecutive 404s before the deployment is declared
	// lost (the manager's in-memory job store forgets jobs on restart).
	maxNotFoundStrikes = 3
)

// DeploymentPoller watches in-flight deployments and writes the outcome onto
// the application row: the api owns all Status transitions (§4). Polling
// state is in-memory — an api restart drops the watch, and the deployment
// endpoint remains the frontend's source of truth for progress.
type DeploymentPoller struct {
	applicationRepo ports.ApplicationRepo
	serverManager   ports.ServerManagerClient
	clock           ports.Clock

	// Interval and Budget are overridable for tests.
	Interval time.Duration
	Budget   time.Duration
}

func NewDeploymentPoller(repo ports.ApplicationRepo, sm ports.ServerManagerClient, clock ports.Clock) *DeploymentPoller {
	return &DeploymentPoller{
		applicationRepo: repo,
		serverManager:   sm,
		clock:           clock,
		Interval:        defaultPollInterval,
		Budget:          defaultPollBudget,
	}
}

// Watch follows one deployment to its terminal state in the background.
func (p *DeploymentPoller) Watch(appID, deploymentID string) {
	go p.poll(appID, deploymentID)
}

func (p *DeploymentPoller) poll(appID, deploymentID string) {
	ctx := context.Background()
	deadline := p.clock.Now().Add(p.Budget)
	strikes := 0

	for {
		time.Sleep(p.Interval)

		status, err := p.serverManager.DeploymentStatus(ctx, deploymentID)
		switch {
		case errors.Is(err, ports.ErrDeploymentNotFound):
			if strikes++; strikes >= maxNotFoundStrikes {
				p.markFailed(ctx, appID)
				return
			}
		case err != nil:
			// Transient manager trouble: keep polling within the budget.
		case status.Status == "running":
			p.markRunning(ctx, appID, status)
			return
		case status.Status == "failed":
			p.markFailed(ctx, appID)
			return
		default:
			strikes = 0
		}

		if p.clock.Now().After(deadline) {
			p.markFailed(ctx, appID)
			return
		}
	}
}

func (p *DeploymentPoller) markRunning(ctx context.Context, appID string, status *ports.DeploymentStatus) {
	now := p.clock.Now()
	updates := map[string]any{
		"status":           domain.ApplicationStatusRunning,
		"last_deployed_at": now,
		"updated_at":       now,
	}
	if status.Image != "" {
		updates["docker_image"] = status.Image
	}
	if status.ContainerID != "" {
		updates["docker_container_id"] = status.ContainerID
	}
	if status.ContainerName != "" {
		updates["docker_container_name"] = status.ContainerName
	}
	if err := p.applicationRepo.Update(appID, updates, ctx); err != nil {
		fmt.Println("recording successful deployment:", err.Error())
	}
}

func (p *DeploymentPoller) markFailed(ctx context.Context, appID string) {
	if err := p.applicationRepo.Update(appID, map[string]any{
		"status":     domain.ApplicationStatusFailed,
		"updated_at": p.clock.Now(),
	}, ctx); err != nil {
		fmt.Println("recording failed deployment:", err.Error())
	}
}
