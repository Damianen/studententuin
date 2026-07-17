package deployments

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

// startStageTimeout bounds the cutover (inspect/stop/remove/create/start);
// the health grace sleep comes on top of it.
const startStageTimeout = 60 * time.Second

// gcTimeout bounds the best-effort image sweep after a successful deploy.
const gcTimeout = 30 * time.Second

type Deploy struct {
	fetcher  ports.SourceFetcher
	builder  ports.ImageBuilder
	runtime  ports.ContainerRuntime
	images   ports.ImageStore
	runner   AppRunner
	jobs     *JobStore
	gitHosts []string
	budgets  Budgets
}

// DeployInput mirrors the §3.2 deploy body. Build- and StartCommand are
// shell strings by design — they run inside build/user containers via the
// generated image, never on the manager host; only /run takes argv arrays.
type DeployInput struct {
	AppID         string
	RepositoryURL string
	Branch        string
	Type          string // accepted, unused: nixpacks auto-detects (§3.4)
	BuildCommand  string
	StartCommand  string
	Env           map[string]string
	Port          int
	MemoryLimit   string
	CpuLimit      string
	Runtime       string
	DatabaseID    string
}

// Execute validates everything it can synchronously (the HTTP 400/409
// surface), registers the job, and detaches the pipeline goroutine. The
// request context is deliberately not inherited — the deploy outlives it.
func (u *Deploy) Execute(_ context.Context, in DeployInput) (string, error) {
	if err := domain.ValidateRepoURL(in.RepositoryURL, u.gitHosts); err != nil {
		return "", err
	}
	if err := domain.ValidateBranch(in.Branch); err != nil {
		return "", err
	}
	if err := domain.ValidateCommand("build command", in.BuildCommand); err != nil {
		return "", err
	}
	if err := domain.ValidateCommand("start command", in.StartCommand); err != nil {
		return "", err
	}

	deploymentID := uuid.NewString()
	imageRef := domain.AppImageRef(in.AppID, deploymentID)

	runInput := apps.RunInput{
		AppID:       in.AppID,
		Image:       imageRef,
		Env:         withPort(in.Env, in.Port),
		Port:        in.Port,
		MemoryLimit: in.MemoryLimit,
		CpuLimit:    in.CpuLimit,
		Runtime:     in.Runtime,
		DatabaseID:  in.DatabaseID,
		// No StartCommand: nixpacks bakes it into the image CMD.
	}
	if err := u.runner.Validate(runInput); err != nil {
		return "", err
	}

	if _, err := u.jobs.Create(deploymentID, in.AppID); err != nil {
		return "", err
	}

	// #nosec G118 -- deliberately detached from the request context: the
	// deploy is an async job (202) that must outlive the HTTP request; every
	// stage carries its own config-driven timeout.
	go u.run(deploymentID, in, runInput, imageRef)
	return deploymentID, nil
}

// run is the async pipeline. Stage ordering is the §3.4 guarantee: the old
// container is not touched until the new image is fully built, so failed
// clones and builds leave the running app running.
func (u *Deploy) run(id string, in DeployInput, runInput apps.RunInput, imageRef string) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("deploy panicked", "deployment_id", id, "panic", r, "stack", string(debug.Stack()))
			u.jobs.Fail(id, "internal error during deploy")
		}
	}()
	sink := u.jobs.LogWriter(id)

	// Clone.
	u.jobs.SetStatus(id, domain.DeploymentStatusCloning)
	cloneCtx, cancelClone := context.WithTimeout(context.Background(), u.budgets.CloneTimeout)
	checkout, err := u.fetcher.Fetch(cloneCtx, in.RepositoryURL, in.Branch)
	cancelClone()
	if err != nil {
		u.jobs.Fail(id, err.Error())
		return
	}
	defer checkout.Cleanup()
	u.jobs.SetSource(id, *checkout, imageRef)

	// Build.
	u.jobs.SetStatus(id, domain.DeploymentStatusBuilding)
	buildCtx, cancelBuild := context.WithTimeout(context.Background(), u.budgets.BuildTimeout)
	err = u.builder.Build(buildCtx, domain.BuildInput{
		Dir:          checkout.Dir,
		ImageRef:     imageRef,
		Env:          in.Env,
		BuildCommand: in.BuildCommand,
		StartCommand: in.StartCommand,
		Port:         in.Port,
	}, sink)
	timedOut := errors.Is(buildCtx.Err(), context.DeadlineExceeded)
	cancelBuild()
	if err != nil {
		if timedOut {
			err = fmt.Errorf("build timed out after %s", u.budgets.BuildTimeout)
		}
		u.jobs.Fail(id, err.Error())
		return
	}

	// Cutover: replace the container in place (network survives).
	u.jobs.SetStatus(id, domain.DeploymentStatusStarting)
	ctx, cancel := context.WithTimeout(context.Background(), startStageTimeout+u.budgets.HealthGrace)
	defer cancel()

	name := domain.AppContainerName(in.AppID)
	if _, err := u.runtime.Inspect(ctx, name); err == nil {
		if err := u.runtime.Stop(ctx, name); err != nil {
			u.jobs.Fail(id, "could not stop the previous container: "+err.Error())
			return
		}
		if err := u.runtime.RemoveContainer(ctx, name); err != nil {
			u.jobs.Fail(id, "could not remove the previous container: "+err.Error())
			return
		}
	} else if !errors.Is(err, domain.ErrNotFound) {
		u.jobs.Fail(id, "could not inspect the previous container: "+err.Error())
		return
	}

	containerID, err := u.runner.Execute(ctx, runInput)
	if err != nil {
		u.jobs.Fail(id, "could not start the new container: "+err.Error())
		return
	}

	// Health check (MVP): still running, no restarts, after a short grace —
	// catches images that exit immediately (bad start command). The failed
	// container is kept so its logs stay available for debugging.
	select {
	case <-time.After(u.budgets.HealthGrace):
	case <-ctx.Done():
	}
	state, err := u.runtime.Inspect(ctx, containerID)
	switch {
	case err != nil:
		u.jobs.Fail(id, "could not verify the new container: "+err.Error())
		return
	case !state.Running || state.RestartCount > 0:
		u.jobs.Fail(id, fmt.Sprintf(
			"app exited during startup (status %s, restarts %d) — check the runtime logs", state.Status, state.RestartCount))
		return
	}

	u.jobs.Succeed(id, containerID, name)
	u.gcImages(in.AppID, imageRef)
}

// gcImages keeps only the image this deploy built (§3.4 cleanup). Best
// effort: an in-use or already-gone image is logged and skipped. Listing all
// tags instead of remembering "the previous one" also sweeps strays from
// earlier failed cutovers, and survives manager restarts.
func (u *Deploy) gcImages(appID, keepRef string) {
	ctx, cancel := context.WithTimeout(context.Background(), gcTimeout)
	defer cancel()

	refs, err := u.images.ListAppImages(ctx, appID)
	if err != nil {
		slog.Warn("image gc: listing app images", "app_id", appID, "error", err)
		return
	}
	for _, ref := range refs {
		if ref == keepRef {
			continue
		}
		if err := u.images.RemoveImage(ctx, ref); err != nil {
			slog.Warn("image gc: removing image", "ref", ref, "error", err)
		}
	}
}

// withPort mirrors the builder's PORT injection for the runtime env, so the
// container and the baked image agree on the port the app listens on.
func withPort(env map[string]string, port int) map[string]string {
	if port <= 0 {
		return env
	}
	if _, ok := env["PORT"]; ok {
		return env
	}
	out := make(map[string]string, len(env)+1)
	for k, v := range env {
		out[k] = v
	}
	out["PORT"] = fmt.Sprintf("%d", port)
	return out
}
