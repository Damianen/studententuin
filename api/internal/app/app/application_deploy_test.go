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
func newDeployMocks(t *testing.T) (*mocks.MockApplicationRepo, *mocks.MockServerManagerClient, *mocks.MockClock, *Service) {
	ar, dr, sm, clock, svc := newDeployMocksWithDB(t)
	dr.EXPECT().FindBySubdomainID(gomock.Any(), gomock.Any()).
		Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	return ar, sm, clock, svc
}

func newDeployMocksWithDB(t *testing.T) (*mocks.MockApplicationRepo, *mocks.MockDatabaseRepo, *mocks.MockServerManagerClient, *mocks.MockClock, *Service) {
	t.Helper()
	ctrl := gomock.NewController(t)
	ar := mocks.NewMockApplicationRepo(ctrl)
	dr := mocks.NewMockDatabaseRepo(ctrl)
	sm := mocks.NewMockServerManagerClient(ctrl)
	clock := mocks.NewMockClock(ctrl)
	clock.EXPECT().Now().Return(time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)).AnyTimes()
	svc := NewService(Dependencies{ApplicationRepo: ar, DatabaseRepo: dr, ServerManager: sm, Clock: clock})
	// Tests drive the poller synchronously through its own tests; deploy
	// tests keep it quick.
	svc.Deploy.poller.Interval = time.Millisecond
	return ar, dr, sm, clock, svc
}

func TestDeployApplication_Execute(t *testing.T) {
	appID := uuid.New()
	subID := uuid.New()

	t.Run("builds the spec from the record and returns the deployment id", func(t *testing.T) {
		ar, sm, _, svc := newDeployMocks(t)
		application := deployableApp(appID, subID)

		var gotSpec ports.DeploySpec
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(application, nil)
		sm.EXPECT().Deploy(gomock.Any(), appID.String(), gomock.Any()).DoAndReturn(
			func(_ context.Context, _ string, spec ports.DeploySpec) (string, error) {
				gotSpec = spec
				return "dep-123", nil
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
		select {
		case <-pending:
		case <-time.After(2 * time.Second):
			t.Error("application never marked pending")
		}
	})

	t.Run("defaults the port when none is configured", func(t *testing.T) {
		ar, sm, _, svc := newDeployMocks(t)
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
		ar, _, _, svc := newDeployMocks(t)
		application := deployableApp(appID, subID)
		application.RepositoryURL = nil
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(application, nil)

		if _, err := svc.Deploy.Execute(context.Background(), subID.String(), appID.String()); !errors.Is(err, ErrNotDeployable) {
			t.Errorf("err = %v, want ErrNotDeployable", err)
		}
	})

	t.Run("wrong subdomain is treated as not found", func(t *testing.T) {
		ar, _, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, uuid.New()), nil)

		if _, err := svc.Deploy.Execute(context.Background(), subID.String(), appID.String()); !errors.Is(err, ErrNotInSubdomain) {
			t.Errorf("err = %v, want ErrNotInSubdomain", err)
		}
	})

	t.Run("conflict from the manager passes through", func(t *testing.T) {
		ar, sm, _, svc := newDeployMocks(t)
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
		ar, dr, sm, _, svc := newDeployMocksWithDB(t)
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
		ar, dr, sm, _, svc := newDeployMocksWithDB(t)
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
		ar, dr, sm, _, svc := newDeployMocksWithDB(t)
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

	newPoller := func(t *testing.T) (*mocks.MockApplicationRepo, *mocks.MockServerManagerClient, *DeploymentPoller) {
		ctrl := gomock.NewController(t)
		ar := mocks.NewMockApplicationRepo(ctrl)
		sm := mocks.NewMockServerManagerClient(ctrl)
		clock := mocks.NewMockClock(ctrl)
		clock.EXPECT().Now().Return(time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)).AnyTimes()
		p := NewDeploymentPoller(ar, sm, clock)
		p.Interval = time.Millisecond
		return ar, sm, p
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

	t.Run("running records the container and deploy time", func(t *testing.T) {
		ar, sm, p := newPoller(t)
		gomock.InOrder(
			sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-1").
				Return(&ports.DeploymentStatus{Status: "building"}, nil),
			sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-1").
				Return(&ports.DeploymentStatus{
					Status: "running", Image: "stt-app-x:abcd1234",
					ContainerID: "cid-9", ContainerName: "stt-app-x",
				}, nil),
		)
		done := make(chan map[string]any, 1)
		ar.EXPECT().Update(appID, gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				done <- updates
				return nil
			})

		p.Watch(appID, "dep-1")
		updates := waitUpdate(t, done)
		if updates["status"] != domain.ApplicationStatusRunning ||
			updates["docker_image"] != "stt-app-x:abcd1234" ||
			updates["docker_container_id"] != "cid-9" ||
			updates["docker_container_name"] != "stt-app-x" ||
			updates["last_deployed_at"] == nil {
			t.Errorf("updates = %v", updates)
		}
	})

	t.Run("failed marks the application failed", func(t *testing.T) {
		ar, sm, p := newPoller(t)
		sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-2").
			Return(&ports.DeploymentStatus{Status: "failed", Error: "build failed"}, nil)
		done := make(chan map[string]any, 1)
		ar.EXPECT().Update(appID, gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				done <- updates
				return nil
			})

		p.Watch(appID, "dep-2")
		if updates := waitUpdate(t, done); updates["status"] != domain.ApplicationStatusFailed {
			t.Errorf("updates = %v", updates)
		}
	})

	t.Run("repeated 404 means the job is lost", func(t *testing.T) {
		ar, sm, p := newPoller(t)
		sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-3").
			Return(nil, fmt.Errorf("status: %w", ports.ErrDeploymentNotFound)).Times(maxNotFoundStrikes)
		done := make(chan map[string]any, 1)
		ar.EXPECT().Update(appID, gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				done <- updates
				return nil
			})

		p.Watch(appID, "dep-3")
		if updates := waitUpdate(t, done); updates["status"] != domain.ApplicationStatusFailed {
			t.Errorf("updates = %v", updates)
		}
	})

	t.Run("budget exhaustion fails the deployment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ar := mocks.NewMockApplicationRepo(ctrl)
		sm := mocks.NewMockServerManagerClient(ctrl)
		p := NewDeploymentPoller(ar, sm, &advancingClock{t: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)})
		p.Interval = time.Millisecond
		p.Budget = 5 * time.Minute // the clock advances a minute per Now()

		sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-4").
			Return(&ports.DeploymentStatus{Status: "building"}, nil).AnyTimes()
		done := make(chan map[string]any, 1)
		ar.EXPECT().Update(appID, gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				done <- updates
				return nil
			})

		p.Watch(appID, "dep-4")
		if updates := waitUpdate(t, done); updates["status"] != domain.ApplicationStatusFailed {
			t.Errorf("updates = %v", updates)
		}
	})
}

func TestStartStopApplication_Execute(t *testing.T) {
	appID := uuid.New()
	subID := uuid.New()

	t.Run("start records running", func(t *testing.T) {
		ar, sm, _, svc := newDeployMocks(t)
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
		ar, sm, _, svc := newDeployMocks(t)
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
		ar, sm, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, subID), nil)
		sm.EXPECT().Stop(gomock.Any(), appID.String()).
			Return(fmt.Errorf("servermanager stop: %w", ports.ErrAppNotDeployed))

		if err := svc.Stop.Execute(context.Background(), subID.String(), appID.String()); !errors.Is(err, ports.ErrAppNotDeployed) {
			t.Errorf("err = %v, want ErrAppNotDeployed", err)
		}
	})

	t.Run("ownership check", func(t *testing.T) {
		ar, _, _, svc := newDeployMocks(t)
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
		ar, sm, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, subID), nil)
		sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-1").
			Return(&ports.DeploymentStatus{ID: "dep-1", AppID: appID.String(), Status: "building"}, nil)

		status, err := svc.GetDeployment.Execute(context.Background(), subID.String(), appID.String(), "dep-1")
		if err != nil || status.Status != "building" {
			t.Errorf("status = %+v, err = %v", status, err)
		}
	})

	t.Run("a deployment of another app is not found", func(t *testing.T) {
		ar, sm, _, svc := newDeployMocks(t)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployableApp(appID, subID), nil)
		sm.EXPECT().DeploymentStatus(gomock.Any(), "dep-1").
			Return(&ports.DeploymentStatus{ID: "dep-1", AppID: uuid.NewString(), Status: "building"}, nil)

		if _, err := svc.GetDeployment.Execute(context.Background(), subID.String(), appID.String(), "dep-1"); !errors.Is(err, ErrNotInSubdomain) {
			t.Errorf("err = %v, want ErrNotInSubdomain", err)
		}
	})
}
