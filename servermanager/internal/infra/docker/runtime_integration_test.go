//go:build integration

package docker

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"servermanager/internal/domain"

	"github.com/google/uuid"
	"github.com/moby/moby/client"
)

// fixtureImage must stay alive under CapDrop ALL + no-new-privileges +
// readonly rootfs; busybox sleep does.
const fixtureImage = "docker.io/library/busybox:1.37"

func integrationRuntime(t *testing.T) *Runtime {
	t.Helper()
	ctx := context.Background()

	r, err := New(ctx)
	if err != nil {
		t.Fatalf("docker daemon unreachable (integration tests need one): %v", err)
	}
	t.Cleanup(func() {
		if err := r.Close(); err != nil {
			t.Logf("closing runtime: %v", err)
		}
	})

	// The adapter never pulls images by design — the test provides the
	// "hand-pushed" fixture.
	resp, err := r.cli.ImagePull(ctx, fixtureImage, client.ImagePullOptions{})
	if err != nil {
		t.Fatalf("pulling fixture image: %v", err)
	}
	if _, err := io.Copy(io.Discard, resp); err != nil {
		t.Fatalf("draining image pull: %v", err)
	}
	if err := resp.Close(); err != nil {
		t.Fatalf("closing image pull: %v", err)
	}

	return r
}

func integrationSpec(appID string) domain.ContainerSpec {
	return domain.ContainerSpec{
		AppID:          appID,
		Image:          fixtureImage,
		Env:            map[string]string{"FIXTURE": "1"},
		Port:           3000,
		Cmd:            []string{"sleep", "3600"},
		MemoryBytes:    256 * 1024 * 1024,
		NanoCPUs:       500_000_000,
		PidsLimit:      256,
		Runtime:        domain.RuntimeRunc, // CI runners only have runc
		ReadonlyRootfs: true,
	}
}

func cleanupApp(t *testing.T, r *Runtime, appID string) {
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = r.cli.ContainerRemove(ctx, domain.AppContainerName(appID),
			client.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
		_, _ = r.cli.NetworkRemove(ctx, domain.AppNetworkName(appID), client.NetworkRemoveOptions{})
	})
}

func TestIntegrationLifecycle(t *testing.T) {
	r := integrationRuntime(t)
	ctx := context.Background()
	appID := uuid.NewString()
	spec := integrationSpec(appID)
	cleanupApp(t, r, appID)

	containerID, err := r.Create(ctx, spec)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if containerID == "" {
		t.Fatal("Create returned empty container ID")
	}

	// Daemon-side twin of the spec_mapper unit test: assert the hardening
	// actually landed in the created container.
	raw, err := r.cli.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{})
	if err != nil {
		t.Fatalf("raw inspect: %v", err)
	}
	hc := raw.Container.HostConfig
	if hc.Memory != spec.MemoryBytes || hc.MemorySwap != spec.MemoryBytes {
		t.Errorf("Memory/MemorySwap = %d/%d, want %d/%d", hc.Memory, hc.MemorySwap, spec.MemoryBytes, spec.MemoryBytes)
	}
	if hc.NanoCPUs != spec.NanoCPUs {
		t.Errorf("NanoCPUs = %d, want %d", hc.NanoCPUs, spec.NanoCPUs)
	}
	if hc.PidsLimit == nil || *hc.PidsLimit != spec.PidsLimit {
		t.Errorf("PidsLimit = %v, want %d", hc.PidsLimit, spec.PidsLimit)
	}
	if len(hc.SecurityOpt) != 1 || hc.SecurityOpt[0] != "no-new-privileges:true" {
		t.Errorf("SecurityOpt = %v", hc.SecurityOpt)
	}
	// The daemon may normalize capability casing.
	if len(hc.CapDrop) != 1 || !strings.EqualFold(hc.CapDrop[0], "ALL") {
		t.Errorf("CapDrop = %v, want [ALL]", hc.CapDrop)
	}
	if len(hc.CapAdd) != 0 {
		t.Errorf("CapAdd = %v, want empty", hc.CapAdd)
	}
	if hc.Runtime != domain.RuntimeRunc {
		t.Errorf("Runtime = %q, want runc", hc.Runtime)
	}
	if want := domain.AppNetworkName(appID); string(hc.NetworkMode) != want {
		t.Errorf("NetworkMode = %q, want %q", hc.NetworkMode, want)
	}
	if hc.RestartPolicy.Name != "on-failure" || hc.RestartPolicy.MaximumRetryCount != maxRestartRetries {
		t.Errorf("RestartPolicy = %+v", hc.RestartPolicy)
	}
	if hc.LogConfig.Type != "json-file" || hc.LogConfig.Config["max-size"] != "10m" || hc.LogConfig.Config["max-file"] != "3" {
		t.Errorf("LogConfig = %+v", hc.LogConfig)
	}
	if !hc.ReadonlyRootfs {
		t.Error("ReadonlyRootfs = false, want true")
	}
	if hc.Tmpfs["/tmp"] == "" {
		t.Errorf("Tmpfs = %v, want /tmp entry", hc.Tmpfs)
	}
	if len(hc.Binds) != 0 {
		t.Errorf("Binds = %v, want empty", hc.Binds)
	}
	if len(hc.PortBindings) != 0 {
		t.Errorf("PortBindings = %v, want empty", hc.PortBindings)
	}

	if err := r.Start(ctx, containerID); err != nil {
		t.Fatalf("Start: %v", err)
	}
	state, err := r.Inspect(ctx, domain.AppContainerName(appID))
	if err != nil {
		t.Fatalf("Inspect after start: %v", err)
	}
	if !state.Exists || !state.Running {
		t.Fatalf("state after start = %+v, want exists+running", state)
	}
	if state.StartedAt.IsZero() {
		t.Error("StartedAt is zero after start")
	}

	if err := r.Stop(ctx, containerID); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	state, err = r.Inspect(ctx, containerID)
	if err != nil {
		t.Fatalf("Inspect after stop: %v", err)
	}
	if state.Running {
		t.Errorf("state after stop = %+v, want not running", state)
	}

	// Stop is idempotent.
	if err := r.Stop(ctx, containerID); err != nil {
		t.Errorf("second Stop: %v, want nil", err)
	}

	if err := r.Remove(ctx, containerID); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if _, err := r.Inspect(ctx, domain.AppContainerName(appID)); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Inspect after remove = %v, want ErrNotFound", err)
	}
	if _, err := r.cli.NetworkInspect(ctx, domain.AppNetworkName(appID), client.NetworkInspectOptions{}); err == nil {
		t.Error("app network still exists after Remove")
	}
	if err := r.Remove(ctx, domain.AppContainerName(appID)); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("second Remove = %v, want ErrNotFound", err)
	}

	// Name and network are free again: same appID can be re-created.
	if _, err := r.Create(ctx, spec); err != nil {
		t.Fatalf("re-Create after remove: %v", err)
	}
}

func TestIntegrationInspectUnknown(t *testing.T) {
	r := integrationRuntime(t)
	if _, err := r.Inspect(context.Background(), domain.AppContainerName(uuid.NewString())); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Inspect = %v, want ErrNotFound", err)
	}
}

func TestIntegrationCreateImageMissing(t *testing.T) {
	r := integrationRuntime(t)
	appID := uuid.NewString()
	cleanupApp(t, r, appID)

	spec := integrationSpec(appID)
	spec.Image = "docker.io/library/studententuin-does-not-exist:" + uuid.NewString()[:8]
	if _, err := r.Create(context.Background(), spec); !errors.Is(err, domain.ErrImageMissing) {
		t.Errorf("Create = %v, want ErrImageMissing", err)
	}
}

func TestIntegrationCreateConflict(t *testing.T) {
	r := integrationRuntime(t)
	ctx := context.Background()
	appID := uuid.NewString()
	spec := integrationSpec(appID)
	cleanupApp(t, r, appID)

	if _, err := r.Create(ctx, spec); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := r.Create(ctx, spec); !errors.Is(err, domain.ErrConflict) {
		t.Errorf("second Create = %v, want ErrConflict", err)
	}
}

func TestIntegrationStopTimeout(t *testing.T) {
	// Guards the graceful-stop window: busybox sleep ignores SIGTERM, so a
	// stop should take roughly stopTimeoutSeconds, not hang forever.
	r := integrationRuntime(t)
	ctx := context.Background()
	appID := uuid.NewString()
	spec := integrationSpec(appID)
	cleanupApp(t, r, appID)

	cid, err := r.Create(ctx, spec)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := r.Start(ctx, cid); err != nil {
		t.Fatalf("Start: %v", err)
	}

	start := time.Now()
	if err := r.Stop(ctx, cid); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	if elapsed := time.Since(start); elapsed > (stopTimeoutSeconds+10)*time.Second {
		t.Errorf("Stop took %v, want ~%ds", elapsed, stopTimeoutSeconds)
	}
}
