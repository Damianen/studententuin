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
	deploymentRepo  ports.DeploymentRepo
	serverManager   ports.ServerManagerClient
	clock           ports.Clock

	// Interval and Budget are overridable for tests.
	Interval time.Duration
	Budget   time.Duration
}

func NewDeploymentPoller(repo ports.ApplicationRepo, deployments ports.DeploymentRepo, sm ports.ServerManagerClient, clock ports.Clock) *DeploymentPoller {
	return &DeploymentPoller{
		applicationRepo: repo,
		deploymentRepo:  deployments,
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
				p.markFailed(ctx, appID, deploymentID, "deployment lost (servermanager restarted)", nil)
				return
			}
		case err != nil:
			// Transient manager trouble: keep polling within the budget.
		case status.Status == "running":
			p.markRunning(ctx, appID, deploymentID, status)
			return
		case status.Status == "failed":
			reason := status.Error
			if reason == "" {
				reason = "deployment failed"
			}
			p.markFailed(ctx, appID, deploymentID, reason, status)
			return
		default:
			strikes = 0
		}

		if p.clock.Now().After(deadline) {
			p.markFailed(ctx, appID, deploymentID, "deployment timed out", nil)
			return
		}
	}
}

func (p *DeploymentPoller) markRunning(ctx context.Context, appID, deploymentID string, status *ports.DeploymentStatus) {
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

	record := map[string]any{
		"status":      domain.DeploymentStatusSucceeded,
		"finished_at": now,
		"updated_at":  now,
	}
	copyCommitMeta(record, status)
	p.finalize(ctx, deploymentID, record)
}

// markFailed finalizes a failed deploy; status is nil when the manager never
// reported one (lost job, poll budget exhausted).
func (p *DeploymentPoller) markFailed(ctx context.Context, appID, deploymentID, reason string, status *ports.DeploymentStatus) {
	now := p.clock.Now()
	if err := p.applicationRepo.Update(appID, map[string]any{
		"status":     domain.ApplicationStatusFailed,
		"updated_at": now,
	}, ctx); err != nil {
		fmt.Println("recording failed deployment:", err.Error())
	}

	record := map[string]any{
		"status":      domain.DeploymentStatusFailed,
		"error":       reason,
		"finished_at": now,
		"updated_at":  now,
	}
	copyCommitMeta(record, status)
	p.finalize(ctx, deploymentID, record)
}

// copyCommitMeta copies the commit fields the manager reports into a history
// update. They exist from the clone stage onward, so a failed build still
// names the commit that broke it.
func copyCommitMeta(record map[string]any, status *ports.DeploymentStatus) {
	if status == nil {
		return
	}
	if status.CommitSHA != "" {
		record["commit_sha"] = status.CommitSHA
	}
	if status.CommitMessage != "" {
		record["commit_message"] = status.CommitMessage
	}
	if status.CommitAuthor != "" {
		record["commit_author"] = status.CommitAuthor
	}
}

// finalize writes the history row's terminal state, keyed by the manager
// deployment ID. History is best-effort: repo trouble never breaks the
// status transition.
func (p *DeploymentPoller) finalize(ctx context.Context, deploymentID string, updates map[string]any) {
	if err := p.deploymentRepo.UpdateByManagerID(deploymentID, updates, ctx); err != nil {
		fmt.Println("recording deployment history:", err.Error())
	}
}
