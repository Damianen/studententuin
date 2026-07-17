// Package cgroup implements ports.MetricsSource with direct cgroup v2 reads
// (SERVERMANAGEMENT.md §6 phase 6). Host-level on purpose: the research
// showed `docker stats` undercounts non-runc runtimes because their overhead
// lives outside the container's cgroup view.
package cgroup

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

// containerIDRe validates the Docker container ID before it is joined into a
// filesystem path — hex only, never free text.
var containerIDRe = regexp.MustCompile(`^[0-9a-f]{12,64}$`)

// Source reads container samples from a cgroup v2 tree (the host's, mounted
// read-only into the manager container as SM_CGROUP_ROOT).
type Source struct {
	root string
}

var _ ports.MetricsSource = (*Source)(nil)

// New builds a source over root. A root without cgroup.controllers (not a
// cgroup2 mount, or the compose volume is missing) only warns: metrics are
// non-critical and must never block startup.
func New(root string) *Source {
	if _, err := os.Stat(filepath.Join(root, "cgroup.controllers")); err != nil {
		slog.Warn("cgroup v2 root not usable — metrics series will stay empty",
			"root", root, "error", err)
	}
	return &Source{root: root}
}

// Sample reads the container's cpu and memory files. Working set follows the
// cAdvisor convention: memory.current minus inactive_file, clamped at 0.
func (s *Source) Sample(_ context.Context, containerID string) (domain.StatsSample, error) {
	if !containerIDRe.MatchString(containerID) {
		return domain.StatsSample{}, fmt.Errorf("container id is not a hex id: %w", domain.ErrInvalid)
	}
	dir, err := s.containerDir(containerID)
	if err != nil {
		return domain.StatsSample{}, err
	}

	usage, err := readParsed(dir, "cpu.stat", parseCPUStat)
	if err != nil {
		return domain.StatsSample{}, err
	}
	cores, err := readParsed(dir, "cpu.max", parseCPUMax)
	if err != nil {
		return domain.StatsSample{}, err
	}
	current, err := readParsed(dir, "memory.current", parseUintValue)
	if err != nil {
		return domain.StatsSample{}, err
	}
	inactiveFile, err := readParsed(dir, "memory.stat", parseMemoryStat)
	if err != nil {
		return domain.StatsSample{}, err
	}

	var workingSet uint64
	if current > inactiveFile {
		workingSet = current - inactiveFile
	}
	return domain.StatsSample{
		CPUUsageUsec:       usage,
		CPULimitCores:      cores,
		MemWorkingSetBytes: workingSet,
	}, nil
}

// containerDir probes the systemd driver layout first (the confirmed layout
// of the target hosts), then plain cgroupfs. Neither existing means the
// container is gone or the root is wrong — ErrNotFound, the caller skips the
// tick.
func (s *Source) containerDir(containerID string) (string, error) {
	for _, dir := range []string{
		filepath.Join(s.root, "system.slice", "docker-"+containerID+".scope"),
		filepath.Join(s.root, "docker", containerID),
	} {
		if _, err := os.Stat(dir); err == nil {
			return dir, nil
		}
	}
	return "", fmt.Errorf("cgroup dir for container %s: %w", containerID[:12], domain.ErrNotFound)
}

// readParsed reads one cgroup file and applies its parser, prefixing errors
// with the file name so a broken layout is diagnosable from the warn log.
func readParsed[T any](dir, file string, parse func(string) (T, error)) (T, error) {
	var zero T
	// #nosec G304 -- dir is the config root plus a hex-validated container
	// id, file is a compile-time constant cgroup file name.
	data, err := os.ReadFile(filepath.Join(dir, file))
	if err != nil {
		return zero, fmt.Errorf("reading %s: %w", file, err)
	}
	v, err := parse(string(data))
	if err != nil {
		return zero, fmt.Errorf("%s: %w", file, err)
	}
	return v, nil
}
