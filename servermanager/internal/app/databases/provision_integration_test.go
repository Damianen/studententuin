//go:build integration

package databases

import (
	"context"
	"strings"
	"testing"
	"time"

	"servermanager/internal/app/apps"
	"servermanager/internal/domain"
	clockinfra "servermanager/internal/infra/clock"
	"servermanager/internal/infra/docker"

	"github.com/google/uuid"
	"github.com/moby/moby/client"
)

const integrationPassword = "0123456789abcdef0123456789abcdef" // gitleaks:allow — fixture, not a secret

// TestIntegrationProvisionPipeline is the phase-5 acceptance against a real
// daemon: provision postgres:16 (proving User postgres + CapDrop ALL + fresh
// named volume start), reach it from a linked "app" container over the db
// network, survive app removal, and sweep everything on delete.
func TestIntegrationProvisionPipeline(t *testing.T) {
	ctx := context.Background()

	runtime, err := docker.New(ctx)
	if err != nil {
		t.Fatalf("docker daemon unreachable (integration tests need one): %v", err)
	}
	t.Cleanup(func() { runtime.Close() })

	// Raw client for read-only daemon-side assertions.
	raw, err := client.New(client.FromEnv)
	if err != nil {
		t.Fatalf("raw docker client: %v", err)
	}
	t.Cleanup(func() { raw.Close() })

	dbID := uuid.NewString()
	appID := uuid.NewString()
	limits := apps.Limits{
		DefaultMemoryBytes: 256 * 1024 * 1024,
		MaxMemoryBytes:     1024 * 1024 * 1024,
		DefaultNanoCPUs:    500_000_000,
		MaxNanoCPUs:        2_000_000_000,
		DefaultPidsLimit:   256,
		DefaultRuntime:     domain.RuntimeRunc,
	}
	svc := NewService(Dependencies{
		Runtime: runtime,
		Limits:  limits,
		Clock:   clockinfra.System{},
		Budgets: Budgets{PullTimeout: 5 * time.Minute, HealthBudget: 90 * time.Second},
	})
	appsSvc := apps.NewService(apps.Dependencies{Runtime: runtime, Limits: limits})

	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_ = appsSvc.Remove.Execute(cleanupCtx, appID)
		_ = svc.Delete.Execute(cleanupCtx, dbID)
	})

	provision := func(t *testing.T) domain.ProvisionJob {
		t.Helper()
		if _, err := svc.Provision.Execute(ctx, ProvisionInput{
			DBID:       dbID,
			Type:       "postgres",
			Version:    "16",
			DBName:     "appdb",
			DBUser:     "app",
			DBPassword: integrationPassword,
			// 512m: the 256m shm carve-out counts against the memory cgroup.
			MemoryLimit: "512m",
		}); err != nil {
			t.Fatalf("Provision.Execute: %v", err)
		}
		deadline := time.Now().Add(7 * time.Minute)
		for {
			job, ok := svc.Jobs.Latest(dbID)
			if ok && job.Terminal() {
				if job.Status != domain.ProvisionStatusRunning {
					t.Fatalf("provision failed: %s", job.Error)
				}
				return job
			}
			if time.Now().After(deadline) {
				t.Fatalf("provision stuck (job %+v)", job)
			}
			time.Sleep(time.Second)
		}
	}

	job := provision(t)
	if job.ContainerName != domain.DBContainerName(dbID) || job.Image != "postgres:16" {
		t.Errorf("job = %+v", job)
	}

	status, err := svc.Status.Execute(ctx, dbID)
	if err != nil || !status.State.Running || status.State.Health != "healthy" {
		t.Fatalf("status = %+v err = %v, want running+healthy", status, err)
	}

	// Daemon-side twin of the db spec-mapper test: the hardening plus the
	// postgres-specific settings actually landed.
	res, err := raw.ContainerInspect(ctx, job.ContainerName, client.ContainerInspectOptions{})
	if err != nil {
		t.Fatalf("raw inspect: %v", err)
	}
	cfg, hc := res.Container.Config, res.Container.HostConfig
	if cfg.User != "postgres" {
		t.Errorf("User = %q, want postgres", cfg.User)
	}
	if cfg.Healthcheck == nil || cfg.Healthcheck.Test[1] != "pg_isready" {
		t.Errorf("Healthcheck = %+v", cfg.Healthcheck)
	}
	if len(hc.CapDrop) != 1 || !strings.EqualFold(hc.CapDrop[0], "ALL") {
		t.Errorf("CapDrop = %v, want [ALL]", hc.CapDrop)
	}
	if len(hc.SecurityOpt) != 1 || hc.SecurityOpt[0] != "no-new-privileges:true" {
		t.Errorf("SecurityOpt = %v", hc.SecurityOpt)
	}
	if hc.ShmSize != dbShmSizeBytes {
		t.Errorf("ShmSize = %d, want %d", hc.ShmSize, dbShmSizeBytes)
	}
	if want := domain.DBNetworkName(dbID); string(hc.NetworkMode) != want {
		t.Errorf("NetworkMode = %q, want %q", hc.NetworkMode, want)
	}
	if len(hc.Binds) != 0 {
		t.Errorf("Binds = %v, want empty", hc.Binds)
	}
	foundVolume := false
	for _, m := range res.Container.Mounts {
		if m.Name == domain.DBVolumeName(dbID) {
			foundVolume = true
		}
	}
	if !foundVolume {
		t.Errorf("mounts %+v missing named volume %s", res.Container.Mounts, domain.DBVolumeName(dbID))
	}

	// Reachability (= the phase-5 acceptance "its app can reach it"): a
	// linked app container resolves stt-db-<id> over the db network and gets
	// a postgres response. The client reuses the already-pulled postgres
	// image for pg_isready.
	probe := "pg_isready -h " + domain.DBContainerName(dbID) + " -p 5432 -U app -t 5" +
		" && echo DB-REACHABLE || echo DB-UNREACHABLE; sleep 120"
	if _, err := appsSvc.Run.Execute(ctx, apps.RunInput{
		AppID:        appID,
		Image:        "postgres:16",
		StartCommand: []string{"sh", "-c", probe},
		DatabaseID:   dbID,
	}); err != nil {
		t.Fatalf("running linked app container: %v", err)
	}

	verdict := ""
	deadline := time.Now().Add(30 * time.Second)
	for verdict == "" {
		lines, err := appsSvc.Logs.Execute(ctx, appID, domain.LogOptions{Tail: 10})
		if err != nil {
			t.Fatalf("app logs: %v", err)
		}
		for _, line := range lines {
			if strings.Contains(line.Message, "DB-") {
				verdict = line.Message
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("no reachability verdict in app logs: %+v", lines)
		}
		time.Sleep(500 * time.Millisecond)
	}
	if !strings.Contains(verdict, "DB-REACHABLE") {
		t.Fatalf("verdict = %q, want DB-REACHABLE", verdict)
	}

	// Removing the app must not touch the database: its network has the
	// deliberately different stt-dbnet- prefix the app sweep ignores.
	if err := appsSvc.Remove.Execute(ctx, appID); err != nil {
		t.Fatalf("removing linked app: %v", err)
	}
	if _, err := raw.NetworkInspect(ctx, domain.DBNetworkName(dbID), client.NetworkInspectOptions{}); err != nil {
		t.Fatalf("db network gone after app removal: %v", err)
	}
	if _, err := raw.NetworkInspect(ctx, domain.AppNetworkName(appID), client.NetworkInspectOptions{}); err == nil {
		t.Error("app network still exists after app removal")
	}
	if status, err := svc.Status.Execute(ctx, dbID); err != nil || !status.State.Running {
		t.Fatalf("db not running after app removal: %+v err = %v", status, err)
	}

	// Delete sweeps container, network, and volume.
	if err := svc.Delete.Execute(ctx, dbID); err != nil {
		t.Fatalf("Delete.Execute: %v", err)
	}
	if status, err := svc.Status.Execute(ctx, dbID); err != nil || status.State.Exists {
		t.Fatalf("db still exists after delete: %+v err = %v", status, err)
	}
	if _, err := raw.NetworkInspect(ctx, domain.DBNetworkName(dbID), client.NetworkInspectOptions{}); err == nil {
		t.Error("db network still exists after delete")
	}
	if _, err := raw.VolumeInspect(ctx, domain.DBVolumeName(dbID), client.VolumeInspectOptions{}); err == nil {
		t.Error("db volume still exists after delete")
	}

	// The name, network, and volume are free again: same dbID re-provisions.
	if job := provision(t); job.Status != domain.ProvisionStatusRunning {
		t.Fatalf("re-provision = %+v", job)
	}
}

// TestIntegrationDeleteIdempotent covers the sweep on a database that never
// existed — every step must tolerate "already gone".
func TestIntegrationDeleteIdempotent(t *testing.T) {
	runtime, err := docker.New(context.Background())
	if err != nil {
		t.Fatalf("docker daemon unreachable: %v", err)
	}
	t.Cleanup(func() { runtime.Close() })

	svc := NewService(Dependencies{
		Runtime: runtime,
		Limits:  apps.Limits{DefaultMemoryBytes: 256 * 1024 * 1024, MaxMemoryBytes: 1024 * 1024 * 1024, DefaultNanoCPUs: 500_000_000, MaxNanoCPUs: 2_000_000_000, DefaultPidsLimit: 256, DefaultRuntime: domain.RuntimeRunc},
		Clock:   clockinfra.System{},
		Budgets: Budgets{PullTimeout: time.Minute, HealthBudget: time.Minute},
	})

	if err := svc.Delete.Execute(context.Background(), uuid.NewString()); err != nil {
		t.Fatalf("Delete of unknown database = %v, want nil", err)
	}
}
