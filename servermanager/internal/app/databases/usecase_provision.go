package databases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"servermanager/internal/app/apps"
	"servermanager/internal/domain"
	"servermanager/internal/ports"

	"github.com/google/uuid"
)

const (
	// dbPort is the in-container postgres port; never published to the host —
	// the linked app reaches it over the database's private network.
	dbPort = 5432

	// dbShmSizeBytes raises /dev/shm from the 64m daemon default: postgres
	// parallel queries fail with "could not resize shared memory segment"
	// otherwise. Counts against the memory cgroup, so the api sends databases
	// a higher memory_limit than the app default.
	dbShmSizeBytes = 256 * 1024 * 1024

	// dbUser is the in-container OS user. Named (not numeric): the uid
	// differs between image variants, the name does not. Non-root from the
	// first process means CapDrop ALL costs nothing — the entrypoint skips
	// its chown/su branch entirely. If an image version ever breaks this,
	// the fallback is root + CapAdd [CHOWN,SETUID,SETGID,FOWNER,DAC_OVERRIDE].
	dbOSUser = "postgres"

	// startStageTimeout bounds create+start; the health wait comes on top.
	startStageTimeout = 60 * time.Second

	// healthPollInterval is how often the wait loop re-inspects the container.
	healthPollInterval = time.Second
)

type ProvisionDatabase struct {
	runtime ports.ContainerRuntime
	limits  apps.Limits
	jobs    *JobStore
	budgets Budgets
	// pollInterval is how often the wait loop re-inspects; a field so tests
	// can shrink it below their health budget.
	pollInterval time.Duration
}

// ProvisionInput carries everything the api decided: the image is resolved
// from type+version against the server-side allowlist, and the credentials
// are api-generated. They pass through to the container env and are never
// logged or retained.
type ProvisionInput struct {
	DBID        string
	Type        string
	Version     string
	DBName      string
	DBUser      string
	DBPassword  string
	MemoryLimit string
	CpuLimit    string
	Runtime     string
}

// Execute validates everything it can synchronously (the HTTP 400/409
// surface), registers the job, and detaches the provision goroutine. The
// request context is deliberately not inherited — the job outlives it.
func (u *ProvisionDatabase) Execute(ctx context.Context, in ProvisionInput) (string, error) {
	spec, img, err := u.buildSpec(in)
	if err != nil {
		return "", err
	}

	// Mirror /run's 409: an existing container means DELETE first. This also
	// catches containers orphaned by a manager restart mid-provision.
	name := domain.DBContainerName(in.DBID)
	if _, err := u.runtime.Inspect(ctx, name); err == nil {
		return "", fmt.Errorf("container %s already exists — delete the database first: %w", name, domain.ErrConflict)
	} else if !errors.Is(err, domain.ErrNotFound) {
		return "", err
	}

	provisionID := uuid.NewString()
	if _, err := u.jobs.Create(provisionID, in.DBID); err != nil {
		return "", err
	}

	// #nosec G118 -- deliberately detached from the request context: the
	// provision is an async job (202) that must outlive the HTTP request;
	// every stage carries its own config-driven timeout.
	go u.run(in.DBID, img.Ref, name, *spec)
	return provisionID, nil
}

// run is the async pipeline: pull → create → start → wait healthy. Failures
// keep whatever was created — the database delete sweep is the recovery path,
// and a crashed container's logs stay available for debugging.
func (u *ProvisionDatabase) run(dbID, imageRef, containerName string, spec domain.ContainerSpec) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("provision panicked", "db_id", dbID, "panic", r, "stack", string(debug.Stack()))
			u.jobs.Fail(dbID, "internal error during provision")
		}
	}()

	// Pull.
	u.jobs.SetStatus(dbID, domain.ProvisionStatusPulling)
	pullCtx, cancelPull := context.WithTimeout(context.Background(), u.budgets.PullTimeout)
	err := u.runtime.PullImage(pullCtx, imageRef)
	cancelPull()
	if err != nil {
		u.jobs.Fail(dbID, "could not pull "+imageRef+": "+err.Error())
		return
	}

	// Create + start.
	u.jobs.SetStatus(dbID, domain.ProvisionStatusStarting)
	startCtx, cancelStart := context.WithTimeout(context.Background(), startStageTimeout)
	containerID, err := u.runtime.Create(startCtx, spec)
	if err == nil {
		err = u.runtime.Start(startCtx, containerID)
	}
	cancelStart()
	if err != nil {
		u.jobs.Fail(dbID, "could not start the database container: "+err.Error())
		return
	}

	// Wait until the healthcheck reports healthy. The check forces TCP, so
	// this only passes once the real server (not the entrypoint's socket-only
	// init server) accepts connections.
	waitCtx, cancelWait := context.WithTimeout(context.Background(), u.budgets.HealthBudget)
	defer cancelWait()
	for {
		state, err := u.runtime.Inspect(waitCtx, containerID)
		switch {
		case err != nil:
			u.jobs.Fail(dbID, "could not verify the database container: "+err.Error())
			return
		case !state.Running || state.RestartCount > 0:
			u.jobs.Fail(dbID, fmt.Sprintf(
				"database exited during startup (status %s, restarts %d) — check the runtime logs", state.Status, state.RestartCount))
			return
		case state.Health == "healthy":
			u.jobs.Succeed(dbID, imageRef, containerID, containerName)
			return
		}

		select {
		case <-waitCtx.Done():
			u.jobs.Fail(dbID, fmt.Sprintf("database did not become healthy within %s", u.budgets.HealthBudget))
			return
		case <-time.After(u.pollInterval):
		}
	}
}

func (u *ProvisionDatabase) buildSpec(in ProvisionInput) (*domain.ContainerSpec, domain.DBImage, error) {
	if _, err := uuid.Parse(in.DBID); err != nil {
		return nil, domain.DBImage{}, fmt.Errorf("id must be a UUID: %w", domain.ErrInvalid)
	}
	img, err := domain.DBImageFor(in.Type, in.Version)
	if err != nil {
		return nil, domain.DBImage{}, err
	}
	if err := domain.ValidateDBIdent("db_name", in.DBName); err != nil {
		return nil, domain.DBImage{}, err
	}
	if err := domain.ValidateDBIdent("db_user", in.DBUser); err != nil {
		return nil, domain.DBImage{}, err
	}
	if err := domain.ValidateDBPassword(in.DBPassword); err != nil {
		return nil, domain.DBImage{}, err
	}

	runtime := in.Runtime
	if runtime == "" {
		runtime = u.limits.DefaultRuntime
	}
	if !domain.ValidRuntime(runtime) {
		return nil, domain.DBImage{}, fmt.Errorf("runtime %q not on allowlist: %w", runtime, domain.ErrInvalid)
	}

	memory, cpus, err := resolveLimits(u.limits, in.MemoryLimit, in.CpuLimit)
	if err != nil {
		return nil, domain.DBImage{}, err
	}

	return &domain.ContainerSpec{
		AppID: in.DBID,
		Kind:  domain.KindDB,
		Image: img.Ref,
		Env: map[string]string{
			"POSTGRES_USER":     in.DBUser,
			"POSTGRES_PASSWORD": in.DBPassword,
			"POSTGRES_DB":       in.DBName,
		},
		Port: dbPort,
		User: dbOSUser,
		Healthcheck: &domain.Healthcheck{
			// -h forces TCP: without it pg_isready probes the unix socket and
			// reports ready during initdb's socket-only temp server.
			Test:     []string{"CMD", "pg_isready", "-h", "127.0.0.1", "-U", in.DBUser, "-d", in.DBName},
			Interval: 2 * time.Second,
			Timeout:  3 * time.Second,
			Retries:  30,
		},
		MemoryBytes:  memory,
		NanoCPUs:     cpus,
		PidsLimit:    u.limits.DefaultPidsLimit, // never from the request
		ShmSizeBytes: dbShmSizeBytes,
		Runtime:      runtime,
		Volumes:      []string{domain.DBVolumeName(in.DBID) + ":" + img.DataDir},
		// Postgres writes PGDATA config and /var/run/postgresql — no
		// read-only rootfs for databases.
		ReadonlyRootfs: false,
	}, img, nil
}

// minMemoryBytes is the Docker daemon's memory-limit floor (6 MiB) — same
// rule as the app run path.
const minMemoryBytes = 6 * 1024 * 1024

// resolveLimits applies the same default/cap policy as the app run path
// (usecase_run_app.go) to the database limits.
func resolveLimits(limits apps.Limits, memoryLimit, cpuLimit string) (memory, cpus int64, err error) {
	memory = limits.DefaultMemoryBytes
	if memoryLimit != "" {
		if memory, err = domain.ParseMemoryLimit(memoryLimit); err != nil {
			return 0, 0, fmt.Errorf("%w: %w", err, domain.ErrInvalid)
		}
	}
	if memory < minMemoryBytes {
		return 0, 0, fmt.Errorf("memory limit below the %d-byte minimum: %w", int64(minMemoryBytes), domain.ErrInvalid)
	}
	if memory > limits.MaxMemoryBytes {
		return 0, 0, fmt.Errorf("memory limit exceeds server maximum: %w", domain.ErrInvalid)
	}

	cpus = limits.DefaultNanoCPUs
	if cpuLimit != "" {
		if cpus, err = domain.ParseCPULimit(cpuLimit); err != nil {
			return 0, 0, fmt.Errorf("%w: %w", err, domain.ErrInvalid)
		}
	}
	if cpus > limits.MaxNanoCPUs {
		return 0, 0, fmt.Errorf("cpu limit exceeds server maximum: %w", domain.ErrInvalid)
	}
	return memory, cpus, nil
}
