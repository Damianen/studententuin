package ports

import (
	"context"

	"servermanager/internal/domain"
)

// MetricsSource reads one host-level resource sample for a container. The
// MVP implementation reads the cgroup v2 tree directly (§6 phase 6 — never
// `docker stats`, which undercounts non-runc runtimes); cAdvisor can swap in
// behind this port later. A container without a resolvable cgroup maps to
// domain.ErrNotFound.
type MetricsSource interface {
	Sample(ctx context.Context, containerID string) (domain.StatsSample, error)
}
