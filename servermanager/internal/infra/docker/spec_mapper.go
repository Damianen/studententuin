package docker

import (
	"fmt"
	"sort"

	"servermanager/internal/domain"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

// maxRestartRetries bounds the on-failure restart policy so a crash-looping
// app can't burn CPU indefinitely (§3.3).
const maxRestartRetries = 3

// buildCreateOptions maps a validated ContainerSpec to the full Docker create
// payload, applying every §3.3 hardening default. Pure — no client, no I/O —
// so the security checklist's HostConfig assertions run without a daemon.
func buildCreateOptions(spec domain.ContainerSpec) (client.ContainerCreateOptions, error) {
	if !domain.ValidRuntime(spec.Runtime) {
		return client.ContainerCreateOptions{}, fmt.Errorf("runtime %q not on allowlist: %w", spec.Runtime, domain.ErrInvalid)
	}

	cfg := &container.Config{
		Image: spec.Image,
		Env:   envSlice(spec.Env),
		Cmd:   spec.Cmd,
	}
	if spec.Port > 0 {
		port, err := network.ParsePort(fmt.Sprintf("%d/tcp", spec.Port))
		if err != nil {
			return client.ContainerCreateOptions{}, fmt.Errorf("port %d: %w", spec.Port, domain.ErrInvalid)
		}
		cfg.ExposedPorts = network.PortSet{port: {}}
	}

	mounts := make([]mount.Mount, 0, len(spec.Volumes))
	for _, v := range spec.Volumes {
		name, target, err := domain.ParseVolumeSpec(v)
		if err != nil {
			return client.ContainerCreateOptions{}, fmt.Errorf("%w: %w", err, domain.ErrInvalid)
		}
		mounts = append(mounts, mount.Mount{Type: mount.TypeVolume, Source: name, Target: target})
	}

	netName := domain.AppNetworkName(spec.AppID)
	hostCfg := &container.HostConfig{
		Resources: container.Resources{
			Memory:     spec.MemoryBytes,
			MemorySwap: spec.MemoryBytes, // same value = swap disabled
			NanoCPUs:   spec.NanoCPUs,
			PidsLimit:  &spec.PidsLimit,
		},
		SecurityOpt:    []string{"no-new-privileges:true"},
		CapDrop:        []string{"ALL"},
		NetworkMode:    container.NetworkMode(netName),
		Mounts:         mounts,
		ReadonlyRootfs: spec.ReadonlyRootfs,
		Tmpfs:          map[string]string{"/tmp": "rw,noexec,nosuid,size=64m"},
		RestartPolicy: container.RestartPolicy{
			Name:              container.RestartPolicyOnFailure,
			MaximumRetryCount: maxRestartRetries,
		},
		Runtime: spec.Runtime,
		LogConfig: container.LogConfig{
			Type:   "json-file",
			Config: map[string]string{"max-size": "10m", "max-file": "3"},
		},
	}

	return client.ContainerCreateOptions{
		Name:       domain.AppContainerName(spec.AppID),
		Config:     cfg,
		HostConfig: hostCfg,
		NetworkingConfig: &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{netName: {}},
		},
	}, nil
}

// envSlice converts the env map to Docker's "K=V" form, sorted so the output
// is deterministic.
func envSlice(env map[string]string) []string {
	if len(env) == 0 {
		return nil
	}
	out := make([]string, 0, len(env))
	for k, v := range env {
		out = append(out, k+"="+v)
	}
	sort.Strings(out)
	return out
}
