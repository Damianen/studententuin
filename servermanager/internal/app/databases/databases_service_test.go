package databases

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"servermanager/internal/app/apps"
	"servermanager/internal/domain"
	"servermanager/internal/mocks"

	"go.uber.org/mock/gomock"
)

const (
	testDBID     = "6f1d2c3b-4a5e-4f60-8b7a-9c0d1e2f3a4b"
	testPassword = "0123456789abcdef0123456789abcdef" // gitleaks:allow — fixture, not a secret
)

func validProvisionInput() ProvisionInput {
	return ProvisionInput{
		DBID:       testDBID,
		Type:       "postgres",
		Version:    "16",
		DBName:     "appdb",
		DBUser:     "app",
		DBPassword: testPassword,
	}
}

func testLimits() apps.Limits {
	return apps.Limits{
		DefaultMemoryBytes: 256 * 1024 * 1024,
		MaxMemoryBytes:     1024 * 1024 * 1024,
		DefaultNanoCPUs:    500_000_000,
		MaxNanoCPUs:        2_000_000_000,
		DefaultPidsLimit:   256,
		DefaultRuntime:     domain.RuntimeRunc,
	}
}

func newTestService(t *testing.T) (*Service, *mocks.MockContainerRuntime) {
	t.Helper()
	ctrl := gomock.NewController(t)
	runtime := mocks.NewMockContainerRuntime(ctrl)
	svc := NewService(Dependencies{
		Runtime: runtime,
		Limits:  testLimits(),
		Clock:   testClock(),
		Budgets: Budgets{PullTimeout: time.Second, HealthBudget: 300 * time.Millisecond},
	})
	svc.Provision.pollInterval = 5 * time.Millisecond
	return svc, runtime
}

// waitTerminal polls until the latest job reaches running/failed — the
// provision pipeline is a detached goroutine by design.
func waitTerminal(t *testing.T, svc *Service, dbID string) domain.ProvisionJob {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		job, ok := svc.Jobs.Latest(dbID)
		if ok && job.Terminal() {
			return job
		}
		if time.Now().After(deadline) {
			t.Fatalf("job never terminal (latest: %+v, %v)", job, ok)
		}
		time.Sleep(2 * time.Millisecond)
	}
}

// expectNoContainer is the sync existence pre-check every provision starts with.
func expectNoContainer(runtime *mocks.MockContainerRuntime) {
	runtime.EXPECT().Inspect(gomock.Any(), domain.DBContainerName(testDBID)).
		Return(nil, fmt.Errorf("inspect: %w", domain.ErrNotFound))
}

func TestProvisionValidation(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*ProvisionInput)
	}{
		{"bad uuid", func(in *ProvisionInput) { in.DBID = "not-a-uuid" }},
		{"type not on allowlist", func(in *ProvisionInput) { in.Type = "mysql" }},
		{"version not on allowlist", func(in *ProvisionInput) { in.Version = "15" }},
		{"bad db name", func(in *ProvisionInput) { in.DBName = "Bad Name" }},
		{"bad db user", func(in *ProvisionInput) { in.DBUser = "app;drop" }},
		{"short password", func(in *ProvisionInput) { in.DBPassword = "short" }},
		{"runtime not on allowlist", func(in *ProvisionInput) { in.Runtime = "gvisor" }},
		{"memory over cap", func(in *ProvisionInput) { in.MemoryLimit = "8g" }},
		{"bad cpu", func(in *ProvisionInput) { in.CpuLimit = "lots" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, _ := newTestService(t)
			in := validProvisionInput()
			tc.mutate(&in)
			if _, err := svc.Provision.Execute(context.Background(), in); !errors.Is(err, domain.ErrInvalid) {
				t.Errorf("Execute = %v, want ErrInvalid", err)
			}
		})
	}
}

func TestProvisionConflictWhenContainerExists(t *testing.T) {
	svc, runtime := newTestService(t)
	runtime.EXPECT().Inspect(gomock.Any(), domain.DBContainerName(testDBID)).
		Return(&domain.ContainerState{Exists: true, Running: true}, nil)

	if _, err := svc.Provision.Execute(context.Background(), validProvisionInput()); !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("Execute = %v, want ErrConflict", err)
	}
}

func TestProvisionSuccess(t *testing.T) {
	svc, runtime := newTestService(t)
	expectNoContainer(runtime)

	var spec domain.ContainerSpec
	runtime.EXPECT().PullImage(gomock.Any(), "postgres:16").Return(nil)
	runtime.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, s domain.ContainerSpec) (string, error) {
			spec = s
			return "cid-db", nil
		})
	runtime.EXPECT().Start(gomock.Any(), "cid-db").Return(nil)
	// First inspect: still starting; second: healthy.
	gomock.InOrder(
		runtime.EXPECT().Inspect(gomock.Any(), "cid-db").
			Return(&domain.ContainerState{Exists: true, Running: true, Health: "starting"}, nil),
		runtime.EXPECT().Inspect(gomock.Any(), "cid-db").
			Return(&domain.ContainerState{Exists: true, Running: true, Health: "healthy"}, nil).
			AnyTimes(),
	)

	provisionID, err := svc.Provision.Execute(context.Background(), validProvisionInput())
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	job := waitTerminal(t, svc, testDBID)

	if job.ID != provisionID || job.Status != domain.ProvisionStatusRunning {
		t.Fatalf("job = %+v, want running %s", job, provisionID)
	}
	if job.Image != "postgres:16" || job.ContainerID != "cid-db" || job.ContainerName != domain.DBContainerName(testDBID) {
		t.Fatalf("job container fields = %+v", job)
	}

	// The spec the runtime received is the §3.3-hardened db spec.
	if spec.Kind != domain.KindDB || spec.AppID != testDBID {
		t.Errorf("spec kind/id = %q/%q", spec.Kind, spec.AppID)
	}
	if spec.User != "postgres" || spec.ReadonlyRootfs || spec.ShmSizeBytes != dbShmSizeBytes {
		t.Errorf("spec user/rootfs/shm = %q/%v/%d", spec.User, spec.ReadonlyRootfs, spec.ShmSizeBytes)
	}
	if want := []string{domain.DBVolumeName(testDBID) + ":/var/lib/postgresql/data"}; len(spec.Volumes) != 1 || spec.Volumes[0] != want[0] {
		t.Errorf("spec volumes = %v, want %v", spec.Volumes, want)
	}
	if spec.Env["POSTGRES_USER"] != "app" || spec.Env["POSTGRES_DB"] != "appdb" || spec.Env["POSTGRES_PASSWORD"] != testPassword {
		t.Errorf("spec env = %v", spec.Env)
	}
	hc := spec.Healthcheck
	if hc == nil || hc.Test[0] != "CMD" || hc.Test[1] != "pg_isready" {
		t.Fatalf("healthcheck = %+v", hc)
	}
	// The probe must force TCP — the socket-only init server would false-positive.
	if !contains(hc.Test, "-h") || !contains(hc.Test, "127.0.0.1") {
		t.Errorf("healthcheck not TCP-forced: %v", hc.Test)
	}

	// The job (what the status endpoint exposes) must never carry the password.
	if strings.Contains(fmt.Sprintf("%+v", job), testPassword) {
		t.Error("provision job leaks the password")
	}
}

func TestProvisionPullFailure(t *testing.T) {
	svc, runtime := newTestService(t)
	expectNoContainer(runtime)
	runtime.EXPECT().PullImage(gomock.Any(), "postgres:16").Return(errors.New("registry unreachable"))

	if _, err := svc.Provision.Execute(context.Background(), validProvisionInput()); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	job := waitTerminal(t, svc, testDBID)
	if job.Status != domain.ProvisionStatusFailed || !strings.Contains(job.Error, "could not pull") {
		t.Fatalf("job = %+v, want pull failure", job)
	}
}

func TestProvisionStartFailure(t *testing.T) {
	svc, runtime := newTestService(t)
	expectNoContainer(runtime)
	runtime.EXPECT().PullImage(gomock.Any(), "postgres:16").Return(nil)
	runtime.EXPECT().Create(gomock.Any(), gomock.Any()).Return("cid-db", nil)
	runtime.EXPECT().Start(gomock.Any(), "cid-db").Return(errors.New("boom"))

	if _, err := svc.Provision.Execute(context.Background(), validProvisionInput()); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	job := waitTerminal(t, svc, testDBID)
	if job.Status != domain.ProvisionStatusFailed || !strings.Contains(job.Error, "could not start") {
		t.Fatalf("job = %+v, want start failure", job)
	}
}

func TestProvisionCrashLoopFailsFast(t *testing.T) {
	svc, runtime := newTestService(t)
	expectNoContainer(runtime)
	runtime.EXPECT().PullImage(gomock.Any(), "postgres:16").Return(nil)
	runtime.EXPECT().Create(gomock.Any(), gomock.Any()).Return("cid-db", nil)
	runtime.EXPECT().Start(gomock.Any(), "cid-db").Return(nil)
	runtime.EXPECT().Inspect(gomock.Any(), "cid-db").
		Return(&domain.ContainerState{Exists: true, Running: false, Status: "exited", RestartCount: 3}, nil)

	if _, err := svc.Provision.Execute(context.Background(), validProvisionInput()); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	job := waitTerminal(t, svc, testDBID)
	if job.Status != domain.ProvisionStatusFailed || !strings.Contains(job.Error, "exited during startup") {
		t.Fatalf("job = %+v, want startup-exit failure", job)
	}
}

func TestProvisionHealthTimeout(t *testing.T) {
	svc, runtime := newTestService(t)
	expectNoContainer(runtime)
	runtime.EXPECT().PullImage(gomock.Any(), "postgres:16").Return(nil)
	runtime.EXPECT().Create(gomock.Any(), gomock.Any()).Return("cid-db", nil)
	runtime.EXPECT().Start(gomock.Any(), "cid-db").Return(nil)
	// Healthy never arrives.
	runtime.EXPECT().Inspect(gomock.Any(), "cid-db").
		Return(&domain.ContainerState{Exists: true, Running: true, Health: "starting"}, nil).AnyTimes()

	if _, err := svc.Provision.Execute(context.Background(), validProvisionInput()); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	job := waitTerminal(t, svc, testDBID)
	if job.Status != domain.ProvisionStatusFailed || !strings.Contains(job.Error, "did not become healthy") {
		t.Fatalf("job = %+v, want health timeout", job)
	}
}

func TestGetStatus(t *testing.T) {
	svc, runtime := newTestService(t)
	runtime.EXPECT().Inspect(gomock.Any(), domain.DBContainerName(testDBID)).
		Return(nil, fmt.Errorf("inspect: %w", domain.ErrNotFound))

	status, err := svc.Status.Execute(context.Background(), testDBID)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if status.State.Exists || status.Provision != nil {
		t.Fatalf("status = %+v, want not-exists without provision", status)
	}
}

func TestDeleteSweepsEverything(t *testing.T) {
	svc, runtime := newTestService(t)
	name := domain.DBContainerName(testDBID)
	gomock.InOrder(
		runtime.EXPECT().Stop(gomock.Any(), name).Return(nil),
		runtime.EXPECT().RemoveContainer(gomock.Any(), name).Return(nil),
		runtime.EXPECT().RemoveNetwork(gomock.Any(), domain.DBNetworkName(testDBID)).Return(nil),
		runtime.EXPECT().RemoveVolume(gomock.Any(), domain.DBVolumeName(testDBID)).Return(nil),
	)

	if err := svc.Delete.Execute(context.Background(), testDBID); err != nil {
		t.Fatalf("Execute: %v", err)
	}
}

func TestDeleteIdempotentWhenGone(t *testing.T) {
	svc, runtime := newTestService(t)
	notFound := fmt.Errorf("gone: %w", domain.ErrNotFound)
	runtime.EXPECT().Stop(gomock.Any(), gomock.Any()).Return(notFound)
	runtime.EXPECT().RemoveContainer(gomock.Any(), gomock.Any()).Return(notFound)
	runtime.EXPECT().RemoveNetwork(gomock.Any(), gomock.Any()).Return(notFound)
	runtime.EXPECT().RemoveVolume(gomock.Any(), gomock.Any()).Return(notFound)

	if err := svc.Delete.Execute(context.Background(), testDBID); err != nil {
		t.Fatalf("Execute = %v, want nil for an already-gone database", err)
	}
}

func TestDeleteInUseVolumePropagates(t *testing.T) {
	svc, runtime := newTestService(t)
	runtime.EXPECT().Stop(gomock.Any(), gomock.Any()).Return(nil)
	runtime.EXPECT().RemoveContainer(gomock.Any(), gomock.Any()).Return(nil)
	runtime.EXPECT().RemoveNetwork(gomock.Any(), gomock.Any()).Return(nil)
	runtime.EXPECT().RemoveVolume(gomock.Any(), gomock.Any()).Return(fmt.Errorf("in use: %w", domain.ErrConflict))

	if err := svc.Delete.Execute(context.Background(), testDBID); !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("Execute = %v, want ErrConflict", err)
	}
}

func TestDeleteBlockedWhileProvisionInFlight(t *testing.T) {
	svc, _ := newTestService(t)
	if _, err := svc.Jobs.Create("p1", testDBID); err != nil {
		t.Fatal(err)
	}
	if err := svc.Delete.Execute(context.Background(), testDBID); !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("Execute = %v, want ErrConflict while provision in flight", err)
	}
}

func contains(ss []string, want string) bool {
	for _, s := range ss {
		if s == want {
			return true
		}
	}
	return false
}
