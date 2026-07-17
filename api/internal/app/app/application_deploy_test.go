package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"api/internal/app/ports"
	"api/internal/domain"
	"api/internal/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func strPtr(s string) *string { return &s }

func deployableApp(id, subdomainID uuid.UUID) *domain.Application {
	return &domain.Application{
		ID:            id,
		SubdomainID:   subdomainID,
		RepositoryURL: strPtr("https://github.com/user/repo"),
		Branch:        strPtr("main"),
		Type:          domain.ApplicationTypeNodejs,
		StartCommand:  strPtr("node index.js"),
		EnvironmentVariables: map[string]string{
			"NODE_ENV": "production",
		},
		Ports:       []int{8081},
		MemoryLimit: strPtr("256m"),
		Status:      domain.ApplicationStatusStopped,
	}
}

// newDeployMocks wires a subdomain without a database — the common case.
func newDeployMocks(t *testing.T) (*mocks.MockApplicationRepo, *mocks.MockDeploymentRepo, *mocks.MockServerManagerClient, *mocks.MockClock, *Service) {
	ar, dr, depr, sm, clock, svc := newDeployMocksWithDB(t)
	dr.EXPECT().FindBySubdomainID(gomock.Any(), gomock.Any()).
		Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	return ar, depr, sm, clock, svc
}

func newDeployMocksWithDB(t *testing.T) (*mocks.MockApplicationRepo, *mocks.MockDatabaseRepo, *mocks.MockDeploymentRepo, *mocks.MockServerManagerClient, *mocks.MockClock, *Service) {
	t.Helper()
	ctrl := gomock.NewController(t)
	ar := mocks.NewMockApplicationRepo(ctrl)
	dr := mocks.NewMockDatabaseRepo(ctrl)
	depr := mocks.NewMockDeploymentRepo(ctrl)
	sm := mocks.NewMockServerManagerClient(ctrl)
	clock := mocks.NewMockClock(ctrl)
	clock.EXPECT().Now().Return(time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)).AnyTimes()
	svc := NewService(Dependencies{ApplicationRepo: ar, DatabaseRepo: dr, DeploymentRepo: depr, ServerManager: sm, Clock: clock})
	// Tests drive the poller synchronously through its own tests; deploy
	// tests keep it quick.
	svc.Deploy.poller.Interval = time.Millisecond
	return ar, dr, depr, sm, clock, svc
}

func TestDeployApplication_Execute(t *testing.T) {
	appID := uuid.New()
	subID := uuid.New()

	t.Run("builds the spec from the record and returns the deployment id", func(t *testing.T) {
		ar, depr, sm, _, svc := newDeployMocks(t)
		application := deployableApp(appID, subID)

		var gotSpec ports.DeploySpec
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(application, nil)
		sm.EXPECT().Deploy(gomock.Any(), appID.String(), gomock.Any()).DoAndReturn(
			func(_ context.Context, _ string, spec ports.DeploySpec) (string, error) {
				gotSpec = spec
				return "dep-123", nil
			})
		var gotDeployment *domain.Deployment
		depr.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(deployment *domain.Deployment, _ context.Context) error {
				gotDeployment = deployment
				return nil
			})
		pending := make(chan struct{})
		ar.EXPECT().Update(appID.String(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				if updates["status"] == domain.ApplicationStatusPending {
					close(pending)
				}
				return nil
			})
		// The poller starts polling; park it on a non-terminal status.
		sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-123").
			Return(&ports.DeploymentStatus{ID: "dep-123", Status: "building"}, nil).AnyTimes()

		id, err := svc.Deploy.Execute(context.Background(), subID.String(), appID.String())
		if err != nil {
			t.Fatalf("Execute: %v", err)
		}
		if id != "dep-123" {
			t.Errorf("deployment id = %q", id)
		}
		if gotSpec.RepositoryURL != "https://github.com/user/repo" || gotSpec.Branch != "main" ||
			gotSpec.StartCommand != "node index.js" || gotSpec.Port != 8081 ||
			gotSpec.MemoryLimit != "256m" || gotSpec.Env["NODE_ENV"] != "production" ||
			gotSpec.Type != "Nodejs" {
			t.Errorf("spec = %+v", gotSpec)
		}
		if gotDeployment == nil || gotDeployment.ApplicationID != appID ||
			gotDeployment.ManagerDeploymentID != "dep-123" ||
			gotDeployment.Status != domain.DeploymentStatusInFlight ||
			gotDeployment.Branch == nil || *gotDeployment.Branch != "main" ||
			gotDeployment.StartedAt.IsZero() {
			t.Errorf("history row = %+v", gotDeployment)
		}
		select {
		case <-pending:
		case <-time.After(2 * time.Second):
			t.Error("application never marked pending")
		}
	})

	t.Run("a history write failure never aborts the deploy", func(t *testing.T) {
		ar, depr, sm, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, subID), nil)
		sm.EXPECT().Deploy(gomock.Any(), appID.String(), gomock.Any()).Return("dep-124", nil)
		depr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("history table gone"))
		ar.EXPECT().Update(appID.String(), gomock.Any(), gomock.Any()).Return(nil)
		sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-124").
			Return(&ports.DeploymentStatus{ID: "dep-124", Status: "building"}, nil).AnyTimes()

		id, err := svc.Deploy.Execute(context.Background(), subID.String(), appID.String())
		if err != nil || id != "dep-124" {
			t.Fatalf("id = %q, err = %v — the deploy must survive a history failure", id, err)
		}
	})

	t.Run("defaults the port when none is configured", func(t *testing.T) {
		ar, _, sm, _, svc := newDeployMocks(t)
		application := deployableApp(appID, subID)
		application.Ports = nil

		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(application, nil)
		var gotPort int
		sm.EXPECT().Deploy(gomock.Any(), appID.String(), gomock.Any()).DoAndReturn(
			func(_ context.Context, _ string, spec ports.DeploySpec) (string, error) {
				gotPort = spec.Port
				return "", errors.New("stop here")
			})

		if _, err := svc.Deploy.Execute(context.Background(), subID.String(), appID.String()); err == nil {
			t.Fatal("expected the injected error")
		}
		if gotPort != defaultAppPort {
			t.Errorf("port = %d, want the %d default", gotPort, defaultAppPort)
		}
	})

	t.Run("missing repository url is not deployable", func(t *testing.T) {
		ar, _, _, _, svc := newDeployMocks(t)
		application := deployableApp(appID, subID)
		application.RepositoryURL = nil
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(application, nil)

		if _, err := svc.Deploy.Execute(context.Background(), subID.String(), appID.String()); !errors.Is(err, ErrNotDeployable) {
			t.Errorf("err = %v, want ErrNotDeployable", err)
		}
	})

	t.Run("wrong subdomain is treated as not found", func(t *testing.T) {
		ar, _, _, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, uuid.New()), nil)

		if _, err := svc.Deploy.Execute(context.Background(), subID.String(), appID.String()); !errors.Is(err, ErrNotInSubdomain) {
			t.Errorf("err = %v, want ErrNotInSubdomain", err)
		}
	})

	t.Run("conflict from the manager passes through", func(t *testing.T) {
		ar, _, sm, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, subID), nil)
		sm.EXPECT().Deploy(gomock.Any(), appID.String(), gomock.Any()).
			Return("", fmt.Errorf("servermanager deploy: %w", ports.ErrDeployInFlight))

		if _, err := svc.Deploy.Execute(context.Background(), subID.String(), appID.String()); !errors.Is(err, ports.ErrDeployInFlight) {
			t.Errorf("err = %v, want ErrDeployInFlight", err)
		}
	})
}

func TestDeployApplication_DatabaseLink(t *testing.T) {
	appID := uuid.New()
	subID := uuid.New()
	dbID := uuid.New()
	connStr := "postgres://app:secret@stt-db-" + dbID.String() + ":5432/appdb?sslmode=disable"

	linkedDatabase := func() *domain.Database {
		return &domain.Database{ID: dbID, SubdomainID: subID, ConnectionString: strPtr(connStr)}
	}

	// stopDeploy makes the manager call fail so the test ends synchronously
	// after capturing the spec.
	captureSpec := func(sm *mocks.MockServerManagerClient, gotSpec *ports.DeploySpec) {
		sm.EXPECT().Deploy(gomock.Any(), appID.String(), gomock.Any()).DoAndReturn(
			func(_ context.Context, _ string, spec ports.DeploySpec) (string, error) {
				*gotSpec = spec
				return "", errors.New("stop here")
			})
	}

	t.Run("injects DATABASE_URL and the database id", func(t *testing.T) {
		ar, dr, _, sm, _, svc := newDeployMocksWithDB(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, subID), nil)
		dr.EXPECT().FindBySubdomainID(subID.String(), gomock.Any()).Return(linkedDatabase(), nil)

		var gotSpec ports.DeploySpec
		captureSpec(sm, &gotSpec)
		_, _ = svc.Deploy.Execute(context.Background(), subID.String(), appID.String())

		if gotSpec.DatabaseID != dbID.String() {
			t.Errorf("DatabaseID = %q, want %s", gotSpec.DatabaseID, dbID)
		}
		if gotSpec.Env["DATABASE_URL"] != connStr {
			t.Errorf("DATABASE_URL = %q", gotSpec.Env["DATABASE_URL"])
		}
		if gotSpec.Env["NODE_ENV"] != "production" {
			t.Errorf("existing env lost: %+v", gotSpec.Env)
		}
	})

	t.Run("a user-defined DATABASE_URL wins", func(t *testing.T) {
		ar, dr, _, sm, _, svc := newDeployMocksWithDB(t)
		application := deployableApp(appID, subID)
		application.EnvironmentVariables["DATABASE_URL"] = "postgres://users-own"
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(application, nil)
		dr.EXPECT().FindBySubdomainID(subID.String(), gomock.Any()).Return(linkedDatabase(), nil)

		var gotSpec ports.DeploySpec
		captureSpec(sm, &gotSpec)
		_, _ = svc.Deploy.Execute(context.Background(), subID.String(), appID.String())

		if gotSpec.Env["DATABASE_URL"] != "postgres://users-own" {
			t.Errorf("DATABASE_URL = %q, want the user's value", gotSpec.Env["DATABASE_URL"])
		}
		// The network link still happens: the id rides along regardless.
		if gotSpec.DatabaseID != dbID.String() {
			t.Errorf("DatabaseID = %q, want %s", gotSpec.DatabaseID, dbID)
		}
	})

	t.Run("a database without a connection string deploys unlinked", func(t *testing.T) {
		ar, dr, _, sm, _, svc := newDeployMocksWithDB(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, subID), nil)
		database := linkedDatabase()
		database.ConnectionString = nil // still provisioning, or failed
		dr.EXPECT().FindBySubdomainID(subID.String(), gomock.Any()).Return(database, nil)

		var gotSpec ports.DeploySpec
		captureSpec(sm, &gotSpec)
		_, _ = svc.Deploy.Execute(context.Background(), subID.String(), appID.String())

		if gotSpec.DatabaseID != "" {
			t.Errorf("DatabaseID = %q, want empty for an unprovisioned database", gotSpec.DatabaseID)
		}
		if _, ok := gotSpec.Env["DATABASE_URL"]; ok {
			t.Error("DATABASE_URL injected for an unprovisioned database")
		}
	})
}

// advancingClock moves one minute per Now() call, so poll budgets expire
// quickly in tests.
type advancingClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *advancingClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(time.Minute)
	return c.t
}

func TestDeploymentPoller(t *testing.T) {
	appID := uuid.New().String()

	newPoller := func(t *testing.T) (*mocks.MockApplicationRepo, *mocks.MockDeploymentRepo, *mocks.MockServerManagerClient, *DeploymentPoller) {
		ctrl := gomock.NewController(t)
		ar := mocks.NewMockApplicationRepo(ctrl)
		depr := mocks.NewMockDeploymentRepo(ctrl)
		sm := mocks.NewMockServerManagerClient(ctrl)
		clock := mocks.NewMockClock(ctrl)
		clock.EXPECT().Now().Return(time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)).AnyTimes()
		p := NewDeploymentPoller(ar, depr, sm, clock)
		p.Interval = time.Millisecond
		return ar, depr, sm, p
	}

	waitUpdate := func(t *testing.T, done <-chan map[string]any) map[string]any {
		t.Helper()
		select {
		case updates := <-done:
			return updates
		case <-time.After(2 * time.Second):
			t.Fatal("poller never wrote the outcome")
			return nil
		}
	}

	// captureAppUpdate/captureHistory funnel the poller's two terminal writes
	// into channels; history is finalized after the app row, so waiting on
	// both keeps the mock controller ahead of the background goroutine.
	captureAppUpdate := func(ar *mocks.MockApplicationRepo) <-chan map[string]any {
		done := make(chan map[string]any, 1)
		ar.EXPECT().Update(appID, gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				done <- updates
				return nil
			})
		return done
	}
	captureHistory := func(depr *mocks.MockDeploymentRepo, deploymentID string) <-chan map[string]any {
		done := make(chan map[string]any, 1)
		depr.EXPECT().UpdateByManagerID(deploymentID, gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				done <- updates
				return nil
			})
		return done
	}

	t.Run("running records the container and deploy time", func(t *testing.T) {
		ar, depr, sm, p := newPoller(t)
		gomock.InOrder(
			sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-1").
				Return(&ports.DeploymentStatus{Status: "building"}, nil),
			sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-1").
				Return(&ports.DeploymentStatus{
					Status: "running", Image: "stt-app-x:abcd1234",
					ContainerID: "cid-9", ContainerName: "stt-app-x",
					CommitSHA: "abcd1234ef", CommitMessage: "feat: metrics", CommitAuthor: "Dev",
				}, nil),
		)
		done := captureAppUpdate(ar)
		history := captureHistory(depr, "dep-1")

		p.Watch(appID, "dep-1")
		updates := waitUpdate(t, done)
		if updates["status"] != domain.ApplicationStatusRunning ||
			updates["docker_image"] != "stt-app-x:abcd1234" ||
			updates["docker_container_id"] != "cid-9" ||
			updates["docker_container_name"] != "stt-app-x" ||
			updates["last_deployed_at"] == nil {
			t.Errorf("updates = %v", updates)
		}
		record := waitUpdate(t, history)
		if record["status"] != domain.DeploymentStatusSucceeded ||
			record["finished_at"] == nil ||
			record["commit_sha"] != "abcd1234ef" ||
			record["commit_message"] != "feat: metrics" ||
			record["commit_author"] != "Dev" {
			t.Errorf("history = %v", record)
		}
	})

	t.Run("failed marks the application failed", func(t *testing.T) {
		ar, depr, sm, p := newPoller(t)
		sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-2").
			Return(&ports.DeploymentStatus{
				Status: "failed", Error: "build failed",
				CommitSHA: "abcd1234ef", CommitMessage: "feat: metrics", CommitAuthor: "Dev",
			}, nil)
		done := captureAppUpdate(ar)
		history := captureHistory(depr, "dep-2")

		p.Watch(appID, "dep-2")
		if updates := waitUpdate(t, done); updates["status"] != domain.ApplicationStatusFailed {
			t.Errorf("updates = %v", updates)
		}
		record := waitUpdate(t, history)
		if record["status"] != domain.DeploymentStatusFailed || record["error"] != "build failed" ||
			record["finished_at"] == nil {
			t.Errorf("history = %v", record)
		}
		// A failed build happens after the clone — the history row must still
		// name the commit that broke it.
		if record["commit_sha"] != "abcd1234ef" ||
			record["commit_message"] != "feat: metrics" ||
			record["commit_author"] != "Dev" {
			t.Errorf("history lost commit metadata: %v", record)
		}
	})

	t.Run("failed without a manager error gets the fallback reason", func(t *testing.T) {
		ar, depr, sm, p := newPoller(t)
		sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-2b").
			Return(&ports.DeploymentStatus{Status: "failed"}, nil)
		done := captureAppUpdate(ar)
		history := captureHistory(depr, "dep-2b")

		p.Watch(appID, "dep-2b")
		waitUpdate(t, done)
		if record := waitUpdate(t, history); record["error"] != "deployment failed" {
			t.Errorf("history = %v", record)
		}
	})

	t.Run("repeated 404 means the job is lost", func(t *testing.T) {
		ar, depr, sm, p := newPoller(t)
		sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-3").
			Return(nil, fmt.Errorf("status: %w", ports.ErrDeploymentNotFound)).Times(maxNotFoundStrikes)
		done := captureAppUpdate(ar)
		history := captureHistory(depr, "dep-3")

		p.Watch(appID, "dep-3")
		if updates := waitUpdate(t, done); updates["status"] != domain.ApplicationStatusFailed {
			t.Errorf("updates = %v", updates)
		}
		record := waitUpdate(t, history)
		if record["status"] != domain.DeploymentStatusFailed ||
			record["error"] != "deployment lost (servermanager restarted)" {
			t.Errorf("history = %v", record)
		}
	})

	t.Run("budget exhaustion fails the deployment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ar := mocks.NewMockApplicationRepo(ctrl)
		depr := mocks.NewMockDeploymentRepo(ctrl)
		sm := mocks.NewMockServerManagerClient(ctrl)
		p := NewDeploymentPoller(ar, depr, sm, &advancingClock{t: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)})
		p.Interval = time.Millisecond
		p.Budget = 5 * time.Minute // the clock advances a minute per Now()

		sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-4").
			Return(&ports.DeploymentStatus{Status: "building"}, nil).AnyTimes()
		done := captureAppUpdate(ar)
		history := captureHistory(depr, "dep-4")

		p.Watch(appID, "dep-4")
		if updates := waitUpdate(t, done); updates["status"] != domain.ApplicationStatusFailed {
			t.Errorf("updates = %v", updates)
		}
		record := waitUpdate(t, history)
		if record["status"] != domain.DeploymentStatusFailed || record["error"] != "deployment timed out" {
			t.Errorf("history = %v", record)
		}
	})
}

func TestListApplicationDeployments_Execute(t *testing.T) {
	appID := uuid.New()
	subID := uuid.New()

	t.Run("passes the limit through and returns the rows newest first", func(t *testing.T) {
		ar, depr, _, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, subID), nil)
		rows := []domain.Deployment{
			{ID: uuid.New(), ApplicationID: appID, Status: domain.DeploymentStatusSucceeded},
			{ID: uuid.New(), ApplicationID: appID, Status: domain.DeploymentStatusFailed},
		}
		depr.EXPECT().ListByApplication(appID.String(), 20, gomock.Any()).Return(rows, nil)

		got, err := svc.ListDeployments.Execute(context.Background(), subID.String(), appID.String(), 20)
		if err != nil {
			t.Fatalf("Execute: %v", err)
		}
		if len(got) != 2 || got[0].ID != rows[0].ID || got[1].ID != rows[1].ID {
			t.Errorf("deployments = %+v", got)
		}
	})

	// The repo is never queried for an app in someone else's subdomain — no
	// depr.EXPECT() proves it.
	t.Run("wrong subdomain is treated as not found", func(t *testing.T) {
		ar, _, _, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, uuid.New()), nil)

		if _, err := svc.ListDeployments.Execute(context.Background(), subID.String(), appID.String(), 20); !errors.Is(err, ErrNotInSubdomain) {
			t.Errorf("err = %v, want ErrNotInSubdomain", err)
		}
	})

	t.Run("app not found propagates", func(t *testing.T) {
		ar, _, _, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		if _, err := svc.ListDeployments.Execute(context.Background(), subID.String(), appID.String(), 20); !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("err = %v, want ErrRecordNotFound", err)
		}
	})
}

func TestStartStopApplication_Execute(t *testing.T) {
	appID := uuid.New()
	subID := uuid.New()

	t.Run("start records running", func(t *testing.T) {
		ar, _, sm, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, subID), nil)
		sm.EXPECT().Start(gomock.Any(), appID.String()).Return(nil)
		ar.EXPECT().Update(appID.String(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				if updates["status"] != domain.ApplicationStatusRunning {
					t.Errorf("updates = %v", updates)
				}
				return nil
			})

		if err := svc.Start.Execute(context.Background(), subID.String(), appID.String()); err != nil {
			t.Fatalf("Start: %v", err)
		}
	})

	t.Run("stop records stopped", func(t *testing.T) {
		ar, _, sm, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, subID), nil)
		sm.EXPECT().Stop(gomock.Any(), appID.String()).Return(nil)
		ar.EXPECT().Update(appID.String(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				if updates["status"] != domain.ApplicationStatusStopped {
					t.Errorf("updates = %v", updates)
				}
				return nil
			})

		if err := svc.Stop.Execute(context.Background(), subID.String(), appID.String()); err != nil {
			t.Fatalf("Stop: %v", err)
		}
	})

	t.Run("not deployed passes through without a status write", func(t *testing.T) {
		ar, _, sm, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, subID), nil)
		sm.EXPECT().Stop(gomock.Any(), appID.String()).
			Return(fmt.Errorf("servermanager stop: %w", ports.ErrAppNotDeployed))

		if err := svc.Stop.Execute(context.Background(), subID.String(), appID.String()); !errors.Is(err, ports.ErrAppNotDeployed) {
			t.Errorf("err = %v, want ErrAppNotDeployed", err)
		}
	})

	t.Run("ownership check", func(t *testing.T) {
		ar, _, _, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, uuid.New()), nil)

		if err := svc.Start.Execute(context.Background(), subID.String(), appID.String()); !errors.Is(err, ErrNotInSubdomain) {
			t.Errorf("err = %v, want ErrNotInSubdomain", err)
		}
	})
}

func TestGetApplicationDeployment_Execute(t *testing.T) {
	appID := uuid.New()
	subID := uuid.New()

	t.Run("returns the manager's status", func(t *testing.T) {
		ar, _, sm, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, subID), nil)
		sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-1").
			Return(&ports.DeploymentStatus{ID: "dep-1", AppID: appID.String(), Status: "building"}, nil)

		status, err := svc.GetDeployment.Execute(context.Background(), subID.String(), appID.String(), "dep-1")
		if err != nil || status.Status != "building" {
			t.Errorf("status = %+v, err = %v", status, err)
		}
	})

	t.Run("a deployment of another app is not found", func(t *testing.T) {
		ar, _, sm, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, subID), nil)
		sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-1").
			Return(&ports.DeploymentStatus{ID: "dep-1", AppID: uuid.NewString(), Status: "building"}, nil)

		if _, err := svc.GetDeployment.Execute(context.Background(), subID.String(), appID.String(), "dep-1"); !errors.Is(err, ErrNotInSubdomain) {
			t.Errorf("err = %v, want ErrNotInSubdomain", err)
		}
	})
}
