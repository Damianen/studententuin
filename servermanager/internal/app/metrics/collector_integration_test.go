//go:build integration

package metrics

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"servermanager/internal/domain"
	"servermanager/internal/infra/cgroup"
	clockinfra "servermanager/internal/infra/clock"
	"servermanager/internal/infra/docker"

	"github.com/google/uuid"
	"github.com/moby/moby/client"
)

const hostCgroupRoot = "/sys/fs/cgroup"

// TestIntegrationCollectorBusyContainer runs the real pipeline — docker list,
// cgroup sample, ring store — against a busy busybox and expects nonzero
// cpu and mem within a few short ticks.
func TestIntegrationCollectorBusyContainer(t *testing.T) {
	if _, err := os.ReadFile(hostCgroupRoot + "/cgroup.controllers"); err != nil {
		t.Skipf("cgroup v2 root not readable (%v) — collector integration needs it", err)
	}

	ctx := context.Background()
	runtime, err := docker.New(ctx)
	if err != nil {
		t.Fatalf("docker daemon unreachable (integration tests need one): %v", err)
	}
	t.Cleanup(func() { runtime.Close() })

	const fixtureImage = "docker.io/library/busybox:1.37"
	appID := uuid.NewString()
	spec := domain.ContainerSpec{
		AppID:       appID,
		Image:       fixtureImage,
		Cmd:         []string{"sh", "-c", "while :; do :; done"}, // hot loop = visible cpu
		MemoryBytes: 256 * 1024 * 1024,
		NanoCPUs:    500_000_000,
		PidsLimit:   64,
		Runtime:     domain.RuntimeRunc,
	}

	raw, err := client.New(client.FromEnv)
	if err != nil {
		t.Fatalf("raw docker client: %v", err)
	}
	t.Cleanup(func() { raw.Close() })
	pull, err := raw.ImagePull(ctx, fixtureImage, client.ImagePullOptions{})
	if err != nil {
		t.Fatalf("pulling fixture image: %v", err)
	}
	if _, err := io.Copy(io.Discard, pull); err != nil {
		t.Fatalf("draining image pull: %v", err)
	}
	pull.Close()

	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_, _ = raw.ContainerRemove(cleanupCtx, domain.AppContainerName(appID),
			client.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
		_, _ = raw.NetworkRemove(cleanupCtx, domain.AppNetworkName(appID), client.NetworkRemoveOptions{})
	})

	cid, err := runtime.Create(ctx, spec)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := runtime.Start(ctx, cid); err != nil {
		t.Fatalf("Start: %v", err)
	}

	store := NewStore(time.Second, time.Hour, clockinfra.System{})
	collector := NewCollector(Dependencies{
		Runtime:   runtime,
		Source:    cgroup.New(hostCgroupRoot),
		Exec:      runtime,
		Store:     store,
		Clock:     clockinfra.System{},
		Interval:  time.Second,
		Retention: time.Hour,
	})

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		collector.Run(runCtx)
	}()

	name := domain.AppContainerName(appID)
	deadline := time.Now().Add(30 * time.Second)
	for {
		series := store.Query(name, AppSeriesKeys, time.Hour)
		var cpu, mem float64
		for _, p := range series["cpu"] {
			cpu = max(cpu, p.Value)
		}
		for _, p := range series["mem"] {
			mem = max(mem, p.Value)
		}
		if cpu > 0 && mem > 0 {
			if cpu > 100 {
				t.Errorf("cpu = %v, want clamped <= 100", cpu)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("no nonzero cpu/mem within deadline: %+v", series)
		}
		time.Sleep(500 * time.Millisecond)
	}

	cancel()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("collector did not stop on context cancel")
	}
}
