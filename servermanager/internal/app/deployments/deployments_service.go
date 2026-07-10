// Package deployments orchestrates the async deploy pipeline (§3.4):
// clone → build → cutover, tracked in an in-memory job store the api polls.
package deployments

import (
	"context"
	"time"

	"servermanager/internal/app/apps"
	"servermanager/internal/ports"
)

// AppRunner is the slice of apps.RunApp the deploy job needs — an interface
// so the state-machine tests mock the run step instead of re-testing spec
// building.
type AppRunner interface {
	Validate(in apps.RunInput) error
	Execute(ctx context.Context, in apps.RunInput) (string, error)
}

// Budgets are the per-stage time limits from config.
type Budgets struct {
	CloneTimeout time.Duration
	BuildTimeout time.Duration
	HealthGrace  time.Duration
}

type Dependencies struct {
	Fetcher  ports.SourceFetcher
	Builder  ports.ImageBuilder
	Runtime  ports.ContainerRuntime
	Images   ports.ImageStore
	Runner   AppRunner
	Clock    ports.Clock
	GitHosts []string
	Budgets  Budgets
}

type Service struct {
	Deploy *Deploy
	Get    *GetDeployment
	Jobs   *JobStore
}

func NewService(d Dependencies) *Service {
	jobs := NewJobStore(d.Clock)
	return &Service{
		Deploy: &Deploy{
			fetcher:  d.Fetcher,
			builder:  d.Builder,
			runtime:  d.Runtime,
			images:   d.Images,
			runner:   d.Runner,
			jobs:     jobs,
			gitHosts: d.GitHosts,
			budgets:  d.Budgets,
		},
		Get:  &GetDeployment{jobs: jobs},
		Jobs: jobs,
	}
}
