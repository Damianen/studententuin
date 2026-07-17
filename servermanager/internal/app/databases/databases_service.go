// Package databases orchestrates database provisioning (§6 phase 5):
// pull → create → start → wait-healthy, tracked in an in-memory job store the
// api polls via the database status endpoint.
package databases

import (
	"time"

	"servermanager/internal/app/apps"
	"servermanager/internal/app/metrics"
	"servermanager/internal/ports"
)

// Budgets are the per-stage time limits from config.
type Budgets struct {
	PullTimeout  time.Duration // first pull of an official image can be slow
	HealthBudget time.Duration // start → healthy, covers initdb
}

type Dependencies struct {
	Runtime ports.ContainerRuntime
	// Metrics is the shared collector-fed store; nil only in tests that never
	// hit the metrics endpoint.
	Metrics *metrics.Store
	Limits  apps.Limits
	Clock   ports.Clock
	Budgets Budgets
}

type Service struct {
	Provision *ProvisionDatabase
	Status    *GetDatabaseStatus
	Delete    *DeleteDatabase
	Logs      *GetDatabaseLogs
	Follow    *FollowDatabaseLogs
	Metrics   *GetDatabaseMetrics
	Jobs      *JobStore
}

func NewService(d Dependencies) *Service {
	jobs := NewJobStore(d.Clock)
	return &Service{
		Provision: &ProvisionDatabase{
			runtime:      d.Runtime,
			limits:       d.Limits,
			jobs:         jobs,
			budgets:      d.Budgets,
			pollInterval: healthPollInterval,
		},
		Status:  &GetDatabaseStatus{runtime: d.Runtime, jobs: jobs},
		Delete:  &DeleteDatabase{runtime: d.Runtime, jobs: jobs},
		Logs:    &GetDatabaseLogs{runtime: d.Runtime},
		Follow:  &FollowDatabaseLogs{runtime: d.Runtime},
		Metrics: &GetDatabaseMetrics{store: d.Metrics},
		Jobs:    jobs,
	}
}
