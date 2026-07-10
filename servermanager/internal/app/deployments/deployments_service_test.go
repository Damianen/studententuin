package deployments

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"servermanager/internal/app/apps"
	"servermanager/internal/domain"
	"servermanager/internal/mocks"

	"go.uber.org/mock/gomock"
)

const testAppID = "8b9f5e9e-9a3a-4b7e-9a59-1d2f3a4b5c6d"

// fakeRunner stands in for apps.RunApp: the spec-building rules have their
// own tests, the state machine only needs the outcomes.
type fakeRunner struct {
	mu          sync.Mutex
	validateErr error
	execErr     error
	execCalls   int
	lastInput   apps.RunInput
}

func (f *fakeRunner) Validate(apps.RunInput) error { return f.validateErr }

func (f *fakeRunner) Execute(_ context.Context, in apps.RunInput) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.execCalls++
	f.lastInput = in
	if f.execErr != nil {
		return "", f.execErr
	}
	return "cid-new", nil
}

type deps struct {
	fetcher *mocks.MockSourceFetcher
	builder *mocks.MockImageBuilder
	runtime *mocks.MockContainerRuntime
	images  *mocks.MockImageStore
	runner  *fakeRunner
	svc     *Service
}

func newTestService(t *testing.T) deps {
	t.Helper()
	ctrl := gomock.NewController(t)
	d := deps{
		fetcher: mocks.NewMockSourceFetcher(ctrl),
		builder: mocks.NewMockImageBuilder(ctrl),
		runtime: mocks.NewMockContainerRuntime(ctrl),
		images:  mocks.NewMockImageStore(ctrl),
		runner:  &fakeRunner{},
	}
	d.svc = NewService(Dependencies{
		Fetcher:  d.fetcher,
		Builder:  d.builder,
		Runtime:  d.runtime,
		Images:   d.images,
		Runner:   d.runner,
		Clock:    testClock(),
		GitHosts: []string{"github.com"},
		Budgets: Budgets{
			CloneTimeout: time.Second,
			BuildTimeout: time.Second,
			HealthGrace:  time.Millisecond,
		},
	})
	return d
}

func validInput() DeployInput {
	return DeployInput{
		AppID:         testAppID,
		RepositoryURL: "https://github.com/user/repo",
		Branch:        "main",
		Port:          3000,
	}
}

func checkoutFixture(t *testing.T) (*domain.SourceCheckout, *atomic.Int32) {
	t.Helper()
	var cleanups atomic.Int32
	return &domain.SourceCheckout{
		Dir:       t.TempDir(),
		CommitSHA: "sha-1",
		Cleanup:   func() { cleanups.Add(1) },
	}, &cleanups
}

// waitCleanup polls: the deferred checkout cleanup runs moments after the
// job turns terminal.
func waitCleanup(t *testing.T, cleanups *atomic.Int32) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for cleanups.Load() == 0 {
		if time.Now().After(deadline) {
			t.Fatal("checkout never cleaned up")
		}
		time.Sleep(2 * time.Millisecond)
	}
}

// waitTerminal polls until the job reaches running/failed — the pipeline is
// a detached goroutine by design.
func waitTerminal(t *testing.T, svc *Service, id string) domain.DeploymentJob {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		job, err := svc.Get.Execute(id)
		if err != nil {
			t.Fatalf("Get(%s): %v", id, err)
		}
		if job.Terminal() {
			return job
		}
		if time.Now().After(deadline) {
			t.Fatalf("deployment stuck in %s", job.Status)
		}
		time.Sleep(2 * time.Millisecond)
	}
}

func TestDeployHappyPath(t *testing.T) {
	d := newTestService(t)
	checkout, cleanups := checkoutFixture(t)
	oldRef := domain.AppImageRepo(testAppID) + ":old00000"

	var imageRef string
	d.fetcher.EXPECT().Fetch(gomock.Any(), "https://github.com/user/repo", "main").Return(checkout, nil)
	d.builder.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, in domain.BuildInput, sink io.Writer) error {
			imageRef = in.ImageRef
			if in.Dir != checkout.Dir {
				t.Errorf("Build dir = %q, want the checkout dir", in.Dir)
			}
			_, _ = io.WriteString(sink, "step 1 ok\n")
			return nil
		})
	// Old container exists → stop + remove (network survives), then health OK.
	name := domain.AppContainerName(testAppID)
	d.runtime.EXPECT().Inspect(gomock.Any(), name).Return(&domain.ContainerState{Exists: true, Running: true}, nil)
	d.runtime.EXPECT().Stop(gomock.Any(), name).Return(nil)
	d.runtime.EXPECT().RemoveContainer(gomock.Any(), name).Return(nil)
	d.runtime.EXPECT().Inspect(gomock.Any(), "cid-new").Return(&domain.ContainerState{Exists: true, Running: true}, nil)
	// GC removes everything but the fresh image.
	removed := make(chan string, 1)
	d.images.EXPECT().ListAppImages(gomock.Any(), testAppID).DoAndReturn(
		func(context.Context, string) ([]string, error) { return []string{oldRef, imageRef}, nil })
	d.images.EXPECT().RemoveImage(gomock.Any(), oldRef).DoAndReturn(
		func(_ context.Context, ref string) error { removed <- ref; return nil })

	id, err := d.svc.Deploy.Execute(context.Background(), validInput())
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	job := waitTerminal(t, d.svc, id)
	if job.Status != domain.DeploymentStatusRunning {
		t.Fatalf("status = %s (error %q), want running", job.Status, job.Error)
	}
	if job.CommitSHA != "sha-1" || job.Image != imageRef || job.ContainerID != "cid-new" || job.ContainerName != name {
		t.Errorf("job fields = %+v", job)
	}
	if !strings.Contains(job.BuildLog, "step 1 ok") {
		t.Errorf("build log missing builder output: %q", job.BuildLog)
	}
	if want := domain.AppImageRef(testAppID, id); imageRef != want {
		t.Errorf("image ref = %q, want %q", imageRef, want)
	}
	if d.runner.lastInput.Env["PORT"] != "3000" {
		t.Errorf("runtime env PORT = %q, want 3000", d.runner.lastInput.Env["PORT"])
	}
	select {
	case <-removed:
	case <-time.After(5 * time.Second):
		t.Fatal("old image never GC'd")
	}
	waitCleanup(t, cleanups)
}

func TestDeployFirstDeploySkipsCutover(t *testing.T) {
	d := newTestService(t)
	checkout, _ := checkoutFixture(t)

	d.fetcher.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any()).Return(checkout, nil)
	d.builder.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	// No old container: only the two inspects (old → not found, new → healthy).
	d.runtime.EXPECT().Inspect(gomock.Any(), domain.AppContainerName(testAppID)).Return(nil, domain.ErrNotFound)
	d.runtime.EXPECT().Inspect(gomock.Any(), "cid-new").Return(&domain.ContainerState{Exists: true, Running: true}, nil)
	d.images.EXPECT().ListAppImages(gomock.Any(), testAppID).Return(nil, nil)

	id, err := d.svc.Deploy.Execute(context.Background(), validInput())
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if job := waitTerminal(t, d.svc, id); job.Status != domain.DeploymentStatusRunning {
		t.Fatalf("status = %s (error %q), want running", job.Status, job.Error)
	}
}

func TestDeploySyncValidation(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*DeployInput)
	}{
		{"bad url", func(in *DeployInput) { in.RepositoryURL = "http://github.com/u/r" }},
		{"bad branch", func(in *DeployInput) { in.Branch = "bad branch" }},
		{"control chars in build command", func(in *DeployInput) { in.BuildCommand = "a\nb" }},
		{"control chars in start command", func(in *DeployInput) { in.StartCommand = "a\x00b" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := newTestService(t) // no mock expectations: nothing may be touched
			in := validInput()
			tc.mutate(&in)
			if _, err := d.svc.Deploy.Execute(context.Background(), in); !errors.Is(err, domain.ErrInvalid) {
				t.Errorf("Execute = %v, want ErrInvalid", err)
			}
		})
	}

	t.Run("runner validation", func(t *testing.T) {
		d := newTestService(t)
		d.runner.validateErr = fmt.Errorf("memory too big: %w", domain.ErrInvalid)
		if _, err := d.svc.Deploy.Execute(context.Background(), validInput()); !errors.Is(err, domain.ErrInvalid) {
			t.Errorf("Execute = %v, want the runner's ErrInvalid", err)
		}
	})
}

func TestDeployFetchFailureLeavesAppAlone(t *testing.T) {
	d := newTestService(t)
	// Builder/runtime/images have no expectations: any call fails the test —
	// that IS the "old container untouched" guarantee.
	d.fetcher.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("repository not found"))

	id, err := d.svc.Deploy.Execute(context.Background(), validInput())
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	job := waitTerminal(t, d.svc, id)
	if job.Status != domain.DeploymentStatusFailed || job.Error != "repository not found" {
		t.Errorf("job = %+v", job)
	}
}

func TestDeployBuildFailureLeavesAppAlone(t *testing.T) {
	d := newTestService(t)
	checkout, cleanups := checkoutFixture(t)

	d.fetcher.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any()).Return(checkout, nil)
	d.builder.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ domain.BuildInput, sink io.Writer) error {
			_, _ = io.WriteString(sink, "npm ERR! boom\n")
			return errors.New("build failed: exit 1")
		})

	id, err := d.svc.Deploy.Execute(context.Background(), validInput())
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	job := waitTerminal(t, d.svc, id)
	if job.Status != domain.DeploymentStatusFailed || !strings.Contains(job.Error, "build failed") {
		t.Errorf("job = %+v", job)
	}
	if !strings.Contains(job.BuildLog, "npm ERR! boom") {
		t.Errorf("build log lost on failure: %q", job.BuildLog)
	}
	waitCleanup(t, cleanups)
}

func TestDeployCutoverFailures(t *testing.T) {
	name := domain.AppContainerName(testAppID)
	running := &domain.ContainerState{Exists: true, Running: true}

	t.Run("stop old fails", func(t *testing.T) {
		d := newTestService(t)
		checkout, _ := checkoutFixture(t)
		d.fetcher.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any()).Return(checkout, nil)
		d.builder.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		d.runtime.EXPECT().Inspect(gomock.Any(), name).Return(running, nil)
		d.runtime.EXPECT().Stop(gomock.Any(), name).Return(errors.New("daemon says no"))

		id, _ := d.svc.Deploy.Execute(context.Background(), validInput())
		job := waitTerminal(t, d.svc, id)
		if job.Status != domain.DeploymentStatusFailed || !strings.Contains(job.Error, "stop the previous") {
			t.Errorf("job = %+v", job)
		}
		if d.runner.execCalls != 0 {
			t.Error("runner executed after failed stop")
		}
	})

	t.Run("start new fails", func(t *testing.T) {
		d := newTestService(t)
		checkout, _ := checkoutFixture(t)
		d.fetcher.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any()).Return(checkout, nil)
		d.builder.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		d.runtime.EXPECT().Inspect(gomock.Any(), name).Return(nil, domain.ErrNotFound)
		d.runner.execErr = errors.New("image entrypoint broken")

		id, _ := d.svc.Deploy.Execute(context.Background(), validInput())
		job := waitTerminal(t, d.svc, id)
		if job.Status != domain.DeploymentStatusFailed || !strings.Contains(job.Error, "start the new container") {
			t.Errorf("job = %+v", job)
		}
	})

	t.Run("health check fails", func(t *testing.T) {
		for _, state := range []*domain.ContainerState{
			{Exists: true, Running: false, Status: "exited"},
			{Exists: true, Running: true, RestartCount: 2, Status: "running"},
		} {
			d := newTestService(t)
			checkout, _ := checkoutFixture(t)
			d.fetcher.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any()).Return(checkout, nil)
			d.builder.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			d.runtime.EXPECT().Inspect(gomock.Any(), name).Return(nil, domain.ErrNotFound)
			d.runtime.EXPECT().Inspect(gomock.Any(), "cid-new").Return(state, nil)

			id, _ := d.svc.Deploy.Execute(context.Background(), validInput())
			job := waitTerminal(t, d.svc, id)
			if job.Status != domain.DeploymentStatusFailed || !strings.Contains(job.Error, "exited during startup") {
				t.Errorf("state %+v: job = %+v", state, job)
			}
		}
	})
}

func TestDeployConflictWhileInFlight(t *testing.T) {
	d := newTestService(t)
	release := make(chan struct{})
	checkout, _ := checkoutFixture(t)

	d.fetcher.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(context.Context, string, string) (*domain.SourceCheckout, error) {
			<-release
			return checkout, errors.New("released")
		})

	id1, err := d.svc.Deploy.Execute(context.Background(), validInput())
	if err != nil {
		t.Fatalf("first Execute: %v", err)
	}
	if _, err := d.svc.Deploy.Execute(context.Background(), validInput()); !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("second Execute = %v, want ErrConflict", err)
	}

	close(release)
	waitTerminal(t, d.svc, id1)

	// The slot frees once the first deploy is terminal. Wait for the new
	// job's terminal state so the async Fetch demonstrably happened.
	d.fetcher.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("nope"))
	id3, err := d.svc.Deploy.Execute(context.Background(), validInput())
	if err != nil {
		t.Fatalf("Execute after terminal = %v, want accepted", err)
	}
	waitTerminal(t, d.svc, id3)
}

func TestDeployPanicBecomesFailure(t *testing.T) {
	d := newTestService(t)
	checkout, _ := checkoutFixture(t)
	d.fetcher.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any()).Return(checkout, nil)
	d.builder.EXPECT().Build(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(context.Context, domain.BuildInput, io.Writer) error { panic("builder bug") })

	id, err := d.svc.Deploy.Execute(context.Background(), validInput())
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	job := waitTerminal(t, d.svc, id)
	if job.Status != domain.DeploymentStatusFailed || job.Error != "internal error during deploy" {
		t.Errorf("job = %+v", job)
	}

	// The single-flight slot is released even on panic.
	d.fetcher.EXPECT().Fetch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("nope"))
	id2, err := d.svc.Deploy.Execute(context.Background(), validInput())
	if err != nil {
		t.Fatalf("Execute after panic = %v, want accepted", err)
	}
	waitTerminal(t, d.svc, id2)
}
