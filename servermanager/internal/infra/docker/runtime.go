package docker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"servermanager/internal/domain"
	"servermanager/internal/ports"

	"github.com/moby/moby/client"
)

// stopTimeoutSeconds is the graceful SIGTERM window before the daemon kills
// the container.
const stopTimeoutSeconds = 10

// Runtime implements ports.ContainerRuntime against the Docker Engine API.
type Runtime struct {
	cli *client.Client
}

var _ ports.ContainerRuntime = (*Runtime)(nil)

// New connects to the Docker daemon (honouring DOCKER_HOST) and pings it, so
// a misconfigured socket fails at startup instead of per request.
func New(ctx context.Context) (*Runtime, error) {
	cli, err := client.New(client.FromEnv, client.WithUserAgent("studententuin-servermanager"))
	if err != nil {
		return nil, fmt.Errorf("docker client: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if _, err := cli.Ping(pingCtx, client.PingOptions{}); err != nil {
		if cerr := cli.Close(); cerr != nil {
			slog.Warn("closing docker client after failed ping", "error", cerr)
		}
		return nil, fmt.Errorf("docker ping: %w", err)
	}
	return &Runtime{cli: cli}, nil
}

func (r *Runtime) Close() error { return r.cli.Close() }

// Create ensures the per-app network exists and creates the hardened
// container attached to it. It never pulls images: Phase 1 images are
// hand-pushed by the operator, later phases build them locally.
func (r *Runtime) Create(ctx context.Context, spec domain.ContainerSpec) (string, error) {
	opts, err := buildCreateOptions(spec)
	if err != nil {
		return "", err
	}

	if err := r.ensureNetwork(ctx, spec.AppID); err != nil {
		return "", err
	}

	res, err := r.cli.ContainerCreate(ctx, opts)
	if err != nil {
		// The container can't be "not found" while being created — a
		// not-found here means the image is absent on the host.
		if errors.Is(mapErr("", err), domain.ErrNotFound) {
			return "", fmt.Errorf("create %s: image %q: %w", opts.Name, spec.Image, domain.ErrImageMissing)
		}
		return "", mapErr("create "+opts.Name, err)
	}
	return res.ID, nil
}

func (r *Runtime) Start(ctx context.Context, nameOrID string) error {
	_, err := r.cli.ContainerStart(ctx, nameOrID, client.ContainerStartOptions{})
	return mapErr("start "+nameOrID, err)
}

// Stop is idempotent: the daemon treats stopping an already-stopped
// container as success.
func (r *Runtime) Stop(ctx context.Context, nameOrID string) error {
	timeout := stopTimeoutSeconds
	_, err := r.cli.ContainerStop(ctx, nameOrID, client.ContainerStopOptions{Timeout: &timeout})
	return mapErr("stop "+nameOrID, err)
}

// Remove force-removes the container (and its anonymous volumes — named
// volumes are kept), then tears down its per-app networks. Image removal is
// deferred to phase 3 when the manager owns the images it builds.
func (r *Runtime) Remove(ctx context.Context, nameOrID string) error {
	res, err := r.cli.ContainerInspect(ctx, nameOrID, client.ContainerInspectOptions{})
	if err != nil {
		return mapErr("remove "+nameOrID, err)
	}

	networks := appNetworks(res)

	if _, err := r.cli.ContainerRemove(ctx, nameOrID, client.ContainerRemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	}); err != nil {
		return mapErr("remove "+nameOrID, err)
	}

	return r.removeNetworks(ctx, networks)
}

func (r *Runtime) Inspect(ctx context.Context, nameOrID string) (*domain.ContainerState, error) {
	res, err := r.cli.ContainerInspect(ctx, nameOrID, client.ContainerInspectOptions{})
	if err != nil {
		return nil, mapErr("inspect "+nameOrID, err)
	}

	state := &domain.ContainerState{Exists: true, RestartCount: res.Container.RestartCount}
	if s := res.Container.State; s != nil {
		state.Running = s.Running
		state.Status = string(s.Status)
		state.StartedAt = parseDockerTime(s.StartedAt)
	}
	return state, nil
}

// Logs lands in phase 2.
func (r *Runtime) Logs(ctx context.Context, nameOrID string, opts domain.LogOptions) ([]domain.LogLine, error) {
	return nil, fmt.Errorf("logs: %w", errors.ErrUnsupported)
}

// parseDockerTime parses the inspect API's string timestamps. Docker uses
// "0001-01-01T00:00:00Z" for never-started; that and garbage both map to the
// zero time.
func parseDockerTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// appNetworks returns the manager-owned networks a container is attached to.
func appNetworks(res client.ContainerInspectResult) []string {
	if res.Container.NetworkSettings == nil {
		return nil
	}
	var names []string
	for name := range res.Container.NetworkSettings.Networks {
		if strings.HasPrefix(name, "stt-net-") {
			names = append(names, name)
		}
	}
	return names
}
