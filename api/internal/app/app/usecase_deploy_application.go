package app

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"context"
	"errors"
	"fmt"
)

// ErrNotDeployable means the application record is missing something a
// deploy needs (repository URL); the wrapped message says what.
var ErrNotDeployable = errors.New("application is not deployable")

// defaultAppPort is used when the record has no ports configured — the PaaS
// convention nixpacks apps listen on (PORT env var).
const defaultAppPort = 3000

type DeployApplication struct {
	applicationRepo ports.ApplicationRepo
	databaseRepo    ports.DatabaseRepo
	deploymentRepo  ports.DeploymentRepo
	serverManager   ports.ServerManagerClient
	clock           ports.Clock
	poller          *DeploymentPoller
}

// Execute queues a deploy on the servermanager from the stored application
// record, marks the app pending, and starts the background poller that will
// flip it to running/failed. Returns the manager's deployment ID for the
// frontend to poll.
func (d *DeployApplication) Execute(ctx context.Context, subdomainID, appID string) (string, error) {
	application, err := d.applicationRepo.FindByID(appID, ctx)
	if err != nil {
		return "", err
	}
	if application.SubdomainID.String() != subdomainID {
		return "", ErrNotInSubdomain
	}

	if application.RepositoryURL == nil || *application.RepositoryURL == "" {
		return "", fmt.Errorf("%w: set a repository URL first", ErrNotDeployable)
	}

	spec := deploySpec(application)
	d.linkDatabase(ctx, subdomainID, &spec)

	deploymentID, err := d.serverManager.Deploy(ctx, appID, spec)
	if err != nil {
		return "", err
	}

	// History is best-effort: a write failure never aborts a deploy the
	// manager already accepted.
	if err := d.deploymentRepo.Create(&domain.Deployment{
		ApplicationID:       application.ID,
		ManagerDeploymentID: deploymentID,
		Status:              domain.DeploymentStatusInFlight,
		Branch:              application.Branch,
		StartedAt:           d.clock.Now(),
	}, ctx); err != nil {
		fmt.Println("recording deployment history:", err.Error())
	}

	if err := d.applicationRepo.Update(appID, map[string]any{
		"status":     domain.ApplicationStatusPending,
		"updated_at": d.clock.Now(),
	}, ctx); err != nil {
		// The job is already queued — the poller will still record the
		// outcome; only the intermediate "pending" write failed.
		fmt.Println("marking application pending:", err.Error())
	}

	d.poller.Watch(appID, deploymentID)
	return deploymentID, nil
}

// deploySpec maps the Application record onto the §3.2 deploy body.
func deploySpec(application *domain.Application) ports.DeploySpec {
	return ports.DeploySpec{
		RepositoryURL: *application.RepositoryURL,
		Branch:        deref(application.Branch),
		Type:          string(application.Type),
		BuildCommand:  deref(application.BuildCommand),
		StartCommand:  deref(application.StartCommand),
		Env:           application.EnvironmentVariables,
		Port:          deployPort(application),
		MemoryLimit:   deref(application.MemoryLimit),
		CpuLimit:      deref(application.CpuLimit),
	}
}

// linkDatabase injects the subdomain's provisioned database into the deploy:
// DATABASE_URL in the env (a user-defined DATABASE_URL always wins) and the
// database id so the manager attaches the container to the database's
// private network. A subdomain without a running database deploys unlinked —
// that's the common case, not an error.
func (d *DeployApplication) linkDatabase(ctx context.Context, subdomainID string, spec *ports.DeploySpec) {
	database, err := d.databaseRepo.FindBySubdomainID(subdomainID, ctx)
	if err != nil || database.ConnectionString == nil || *database.ConnectionString == "" {
		return
	}

	spec.DatabaseID = database.ID.String()
	if _, exists := spec.Env["DATABASE_URL"]; exists {
		return
	}
	env := make(map[string]string, len(spec.Env)+1)
	for k, v := range spec.Env {
		env[k] = v
	}
	env["DATABASE_URL"] = *database.ConnectionString
	spec.Env = env
}

// deployPort picks the app's configured port, falling back to the
// conventional default when none is stored yet.
func deployPort(application *domain.Application) int {
	if len(application.Ports) > 0 && application.Ports[0] > 0 {
		return application.Ports[0]
	}
	return defaultAppPort
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
