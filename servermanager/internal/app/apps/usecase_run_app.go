package apps

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"servermanager/internal/domain"
	"servermanager/internal/ports"

	"github.com/distribution/reference"
)

// minMemoryBytes is the Docker daemon's memory-limit floor (6 MiB); reject
// below it up front instead of surfacing a daemon error.
const minMemoryBytes = 6 * 1024 * 1024

var envKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type RunApp struct {
	runtime ports.ContainerRuntime
	limits  Limits
}

// RunInput mirrors the §3.2 deploy body minus the git/build fields, so the
// phase 3 deploy job ends in this same path.
type RunInput struct {
	AppID          string
	Image          string
	Env            map[string]string
	Port           int
	StartCommand   []string
	MemoryLimit    string
	CpuLimit       string
	Runtime        string
	Volumes        []string
	ReadonlyRootfs bool
}

// Execute validates and defaults the input against the server-side caps,
// creates the hardened container, and starts it. A container that fails to
// start is removed again — there is no reconciler yet to catch orphans.
func (u *RunApp) Execute(ctx context.Context, in RunInput) (string, error) {
	spec, err := u.buildSpec(in)
	if err != nil {
		return "", err
	}

	containerID, err := u.runtime.Create(ctx, *spec)
	if err != nil {
		return "", err
	}

	if err := u.runtime.Start(ctx, containerID); err != nil {
		if rmErr := u.runtime.Remove(ctx, containerID); rmErr != nil {
			slog.Error("removing container after failed start", "container_id", containerID, "error", rmErr)
		}
		return "", err
	}
	return containerID, nil
}

func (u *RunApp) buildSpec(in RunInput) (*domain.ContainerSpec, error) {
	if _, err := reference.ParseNormalizedNamed(in.Image); err != nil {
		return nil, fmt.Errorf("image %q: %w", in.Image, domain.ErrInvalid)
	}
	if in.Port < 0 || in.Port > 65535 {
		return nil, fmt.Errorf("port %d out of range: %w", in.Port, domain.ErrInvalid)
	}
	for key := range in.Env {
		// Values are free text and are never logged (§7); keys must be sane.
		if !envKeyRe.MatchString(key) {
			return nil, fmt.Errorf("env key %q: %w", key, domain.ErrInvalid)
		}
	}
	for _, v := range in.Volumes {
		if _, _, err := domain.ParseVolumeSpec(v); err != nil {
			return nil, fmt.Errorf("%w: %w", err, domain.ErrInvalid)
		}
	}

	runtime := in.Runtime
	if runtime == "" {
		runtime = u.limits.DefaultRuntime
	}
	if !domain.ValidRuntime(runtime) {
		return nil, fmt.Errorf("runtime %q not on allowlist: %w", runtime, domain.ErrInvalid)
	}

	memory := u.limits.DefaultMemoryBytes
	if in.MemoryLimit != "" {
		var err error
		if memory, err = domain.ParseMemoryLimit(in.MemoryLimit); err != nil {
			return nil, fmt.Errorf("%w: %w", err, domain.ErrInvalid)
		}
	}
	if memory < minMemoryBytes {
		return nil, fmt.Errorf("memory limit below the %d-byte minimum: %w", int64(minMemoryBytes), domain.ErrInvalid)
	}
	// Reject over-cap requests instead of clamping: a silent clamp would
	// hide api/manager configuration drift.
	if memory > u.limits.MaxMemoryBytes {
		return nil, fmt.Errorf("memory limit exceeds server maximum: %w", domain.ErrInvalid)
	}

	cpus := u.limits.DefaultNanoCPUs
	if in.CpuLimit != "" {
		var err error
		if cpus, err = domain.ParseCPULimit(in.CpuLimit); err != nil {
			return nil, fmt.Errorf("%w: %w", err, domain.ErrInvalid)
		}
	}
	if cpus > u.limits.MaxNanoCPUs {
		return nil, fmt.Errorf("cpu limit exceeds server maximum: %w", domain.ErrInvalid)
	}

	return &domain.ContainerSpec{
		AppID:          in.AppID,
		Image:          in.Image,
		Env:            in.Env,
		Port:           in.Port,
		Cmd:            in.StartCommand,
		MemoryBytes:    memory,
		NanoCPUs:       cpus,
		PidsLimit:      u.limits.DefaultPidsLimit, // never from the request
		Runtime:        runtime,
		Volumes:        in.Volumes,
		ReadonlyRootfs: in.ReadonlyRootfs,
	}, nil
}
