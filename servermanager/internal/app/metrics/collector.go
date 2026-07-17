package metrics

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

type Dependencies struct {
	Runtime   ports.ContainerRuntime
	Source    ports.MetricsSource
	Exec      ports.ContainerExec
	Store     *Store
	Clock     ports.Clock
	Interval  time.Duration
	Retention time.Duration
}

// Collector is the manager's background sampling loop: one goroutine, one
// tick per interval, one Record per running managed container. All state
// (baselines, warn-once sets) is owned by that goroutine — no locks.
type Collector struct {
	runtime   ports.ContainerRuntime
	source    ports.MetricsSource
	store     *Store
	clock     ports.Clock
	interval  time.Duration
	retention time.Duration

	// cpuBaselines is keyed by container ID on purpose: a cutover's new
	// container must not inherit the old counter, and a restart resets it.
	cpuBaselines map[string]cpuBaseline
	// sampleWarned keeps the per-container cgroup failure to one warning per
	// failure streak instead of one per tick.
	sampleWarned map[string]bool

	db *dbSampler
}

type cpuBaseline struct {
	usageUsec uint64
	at        time.Time
}

func NewCollector(d Dependencies) *Collector {
	return &Collector{
		runtime:      d.Runtime,
		source:       d.Source,
		store:        d.Store,
		clock:        d.Clock,
		interval:     d.Interval,
		retention:    d.Retention,
		cpuBaselines: map[string]cpuBaseline{},
		sampleWarned: map[string]bool{},
		db:           newDBSampler(d.Exec),
	}
}

// Run ticks until ctx is canceled. The first tick is immediate so mem/db
// series don't wait a full interval after startup.
func (c *Collector) Run(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	c.tick(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.tick(ctx)
		}
	}
}

// tick samples every running managed container. Per-container failures warn
// and continue — one broken container must never starve the others' series.
func (c *Collector) tick(ctx context.Context) {
	// Bound the whole pass: a stalled daemon call (list, inspect) must never
	// wedge the collector across ticks.
	ctx, cancel := context.WithTimeout(ctx, c.interval)
	defer cancel()

	now := c.clock.Now()
	containers, err := c.runtime.ListManagedContainers(ctx)
	if err != nil {
		slog.Warn("metrics: listing managed containers", "error", err)
		return
	}

	seen := make(map[string]bool, len(containers))
	for _, ctr := range containers {
		seen[ctr.ID] = true
		values := map[string]float64{}

		sample, err := c.source.Sample(ctx, ctr.ID)
		switch {
		case err != nil:
			if !c.sampleWarned[ctr.ID] {
				slog.Warn("metrics: sampling container", "container", ctr.Name, "error", err)
				c.sampleWarned[ctr.ID] = true
			}
		default:
			delete(c.sampleWarned, ctr.ID)
			if cpu, ok := c.cpuPercent(ctr.ID, sample, now); ok {
				values["cpu"] = cpu
			}
			values["mem"] = float64(sample.MemWorkingSetBytes) / (1024 * 1024)
		}

		// DB vitals come from postgres itself and don't depend on the cgroup
		// sample succeeding.
		if ctr.Kind == domain.KindDB {
			c.db.sample(ctx, ctr, now, values)
		}

		if len(values) > 0 {
			c.store.Record(ctr.Name, now, values)
		}
	}

	c.prune(seen)
	c.store.DropStale(now.Add(-c.retention))
}

// cpuPercent turns the cumulative counter into a percentage of the
// container's cpu limit, clamped to [0,100]. The first sighting of an ID and
// counter resets (restart, cutover reuse of a recycled ID) record a fresh
// baseline and skip the point.
func (c *Collector) cpuPercent(containerID string, sample domain.StatsSample, now time.Time) (float64, bool) {
	prev, ok := c.cpuBaselines[containerID]
	c.cpuBaselines[containerID] = cpuBaseline{usageUsec: sample.CPUUsageUsec, at: now}
	if !ok || sample.CPUUsageUsec < prev.usageUsec || !now.After(prev.at) {
		return 0, false
	}

	cores := sample.CPULimitCores
	if cores <= 0 {
		// cpu.max "max": every managed container gets NanoCPUs (§3.3), but
		// stay correct for hand-made ones — percentage of the whole host.
		cores = float64(runtime.NumCPU())
	}

	wallUsec := now.Sub(prev.at).Seconds() * 1e6
	pct := float64(sample.CPUUsageUsec-prev.usageUsec) / wallUsec / cores * 100
	return min(max(pct, 0), 100), true
}

// prune forgets per-ID state for containers not seen this tick.
func (c *Collector) prune(seen map[string]bool) {
	for id := range c.cpuBaselines {
		if !seen[id] {
			delete(c.cpuBaselines, id)
		}
	}
	for id := range c.sampleWarned {
		if !seen[id] {
			delete(c.sampleWarned, id)
		}
	}
	c.db.prune(seen)
}
