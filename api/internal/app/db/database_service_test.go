package db

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"api/internal/mocks"
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func newDBMocks(t *testing.T) (*mocks.MockDatabaseRepo, *mocks.MockServerManagerClient, *Service) {
	t.Helper()
	ctrl := gomock.NewController(t)
	dr := mocks.NewMockDatabaseRepo(ctrl)
	sm := mocks.NewMockServerManagerClient(ctrl)
	clock := mocks.NewMockClock(ctrl)
	clock.EXPECT().Now().Return(time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)).AnyTimes()
	svc := NewService(Dependencies{DatabaseRepo: dr, ServerManager: sm, Clock: clock})
	svc.Poller.Interval = time.Millisecond
	return dr, sm, svc
}

func validInput(subdomainID string) DatabaseInput {
	return DatabaseInput{
		SubdomainID: subdomainID,
		Name:        "My Database",
		Version:     "16",
		Type:        domain.DatabaseTypePostgres,
		DbName:      "appdb",
	}
}

func TestCreateDatabase_Execute(t *testing.T) {
	subID := uuid.New().String()

	t.Run("creates the record and fires the provision with generated creds", func(t *testing.T) {
		dr, sm, svc := newDBMocks(t)

		var created *domain.Database
		dr.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *domain.Database, _ context.Context) error {
				created = db
				return nil
			})
		var gotSpec ports.DBProvisionSpec
		sm.EXPECT().ProvisionDatabase(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, dbID string, spec ports.DBProvisionSpec) (string, error) {
				gotSpec = spec
				return "prov-1", nil
			})
		// The poller starts polling; park it on a non-terminal status.
		sm.EXPECT().DatabaseStatus(gomock.Any(), gomock.Any()).
			Return(&ports.DBStatus{Provision: &ports.DBProvisionState{Status: "pulling"}}, nil).AnyTimes()

		if err := svc.Create.Execute(context.Background(), validInput(subID)); err != nil {
			t.Fatalf("Execute: %v", err)
		}

		if created == nil || created.Status != domain.DatabaseStatusProvisioning {
			t.Fatalf("created record = %+v", created)
		}
		if created.DbName == nil || *created.DbName != "appdb" || created.DbUser == nil || *created.DbUser != "app" {
			t.Errorf("record identifiers = %v/%v", created.DbName, created.DbUser)
		}
		if created.ConnectionString != nil {
			t.Error("connection string must not be set before the provision succeeds")
		}

		if gotSpec.Type != "postgres" || gotSpec.Version != "16" || gotSpec.DBName != "appdb" || gotSpec.DBUser != "app" {
			t.Errorf("spec = %+v", gotSpec)
		}
		if !regexp.MustCompile(`^[0-9a-f]{32}$`).MatchString(gotSpec.DBPassword) {
			t.Errorf("password %q is not 32 hex chars", gotSpec.DBPassword)
		}
		if gotSpec.MemoryLimit != dbMemoryDefault {
			t.Errorf("memory = %q, want the %s database default", gotSpec.MemoryLimit, dbMemoryDefault)
		}
	})

	t.Run("provision request failure marks the record failed, create still succeeds", func(t *testing.T) {
		dr, sm, svc := newDBMocks(t)
		dr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		sm.EXPECT().ProvisionDatabase(gomock.Any(), gomock.Any(), gomock.Any()).
			Return("", errors.New("manager unreachable"))
		failed := make(chan map[string]any, 1)
		dr.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				failed <- updates
				return nil
			})

		if err := svc.Create.Execute(context.Background(), validInput(subID)); err != nil {
			t.Fatalf("Execute: %v", err)
		}
		select {
		case updates := <-failed:
			if updates["status"] != domain.DatabaseStatusFailed {
				t.Errorf("updates = %+v, want failed status", updates)
			}
		case <-time.After(2 * time.Second):
			t.Error("record never marked failed")
		}
	})

	t.Run("invalid subdomain uuid", func(t *testing.T) {
		_, _, svc := newDBMocks(t)
		if err := svc.Create.Execute(context.Background(), validInput("not-a-uuid")); err == nil {
			t.Fatal("expected error for bad uuid")
		}
	})

	t.Run("repo error propagates", func(t *testing.T) {
		dr, _, svc := newDBMocks(t)
		dr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
		if err := svc.Create.Execute(context.Background(), validInput(subID)); err == nil {
			t.Fatal("expected repo error")
		}
	})
}

func TestGetDatabase_Execute(t *testing.T) {
	dbID := uuid.New()
	subID := uuid.New()

	t.Run("success", func(t *testing.T) {
		dr, _, svc := newDBMocks(t)
		dr.EXPECT().FindByID(dbID.String(), gomock.Any()).
			Return(&domain.Database{ID: dbID, SubdomainID: subID}, nil)
		result, err := svc.Get.Execute(context.Background(), subID.String(), dbID.String())
		if err != nil || result.ID != dbID {
			t.Fatalf("result = %+v err = %v", result, err)
		}
	})

	t.Run("wrong subdomain is treated as not found", func(t *testing.T) {
		dr, _, svc := newDBMocks(t)
		dr.EXPECT().FindByID(dbID.String(), gomock.Any()).
			Return(&domain.Database{ID: dbID, SubdomainID: uuid.New()}, nil)
		if _, err := svc.Get.Execute(context.Background(), subID.String(), dbID.String()); !errors.Is(err, ErrNotInSubdomain) {
			t.Fatalf("err = %v, want ErrNotInSubdomain", err)
		}
	})
}

func TestUpdateDatabase_Execute(t *testing.T) {
	dbID := uuid.New()
	subID := uuid.New()
	name := "newname"

	t.Run("belongs check runs before the update", func(t *testing.T) {
		dr, _, svc := newDBMocks(t)
		dr.EXPECT().FindByID(dbID.String(), gomock.Any()).
			Return(&domain.Database{ID: dbID, SubdomainID: uuid.New()}, nil)
		err := svc.Update.Execute(context.Background(), DatabaseUpdateInput{ID: dbID.String(), SubdomainID: subID.String(), Name: &name})
		if !errors.Is(err, ErrNotInSubdomain) {
			t.Fatalf("err = %v, want ErrNotInSubdomain", err)
		}
	})

	t.Run("success with name", func(t *testing.T) {
		dr, _, svc := newDBMocks(t)
		dr.EXPECT().FindByID(dbID.String(), gomock.Any()).
			Return(&domain.Database{ID: dbID, SubdomainID: subID}, nil)
		dr.EXPECT().Update(dbID.String(), gomock.Any(), gomock.Any()).Return(nil)
		if err := svc.Update.Execute(context.Background(), DatabaseUpdateInput{ID: dbID.String(), SubdomainID: subID.String(), Name: &name}); err != nil {
			t.Fatalf("Execute: %v", err)
		}
	})
}

func TestDeleteDatabase_Execute(t *testing.T) {
	dbID := uuid.New()
	subID := uuid.New()
	record := func() *domain.Database { return &domain.Database{ID: dbID, SubdomainID: subID} }

	t.Run("removes on the hosting server before deleting the row", func(t *testing.T) {
		dr, sm, svc := newDBMocks(t)
		db := record()
		gomock.InOrder(
			dr.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(db, nil),
			sm.EXPECT().RemoveDatabase(gomock.Any(), dbID.String()).Return(nil),
			dr.EXPECT().Delete(db, gomock.Any()).Return(nil),
		)
		if err := svc.Delete.Execute(context.Background(), subID.String(), dbID.String()); err != nil {
			t.Fatalf("Execute: %v", err)
		}
	})

	t.Run("provision in flight blocks the delete", func(t *testing.T) {
		dr, sm, svc := newDBMocks(t)
		dr.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(record(), nil)
		sm.EXPECT().RemoveDatabase(gomock.Any(), dbID.String()).
			Return(fmt.Errorf("servermanager remove database: %w", ports.ErrProvisionInFlight))
		if err := svc.Delete.Execute(context.Background(), subID.String(), dbID.String()); !errors.Is(err, ports.ErrProvisionInFlight) {
			t.Fatalf("err = %v, want ErrProvisionInFlight", err)
		}
	})

	t.Run("wrong subdomain never reaches the manager", func(t *testing.T) {
		dr, _, svc := newDBMocks(t)
		dr.EXPECT().FindByID(dbID.String(), gomock.Any()).
			Return(&domain.Database{ID: dbID, SubdomainID: uuid.New()}, nil)
		if err := svc.Delete.Execute(context.Background(), subID.String(), dbID.String()); !errors.Is(err, ErrNotInSubdomain) {
			t.Fatalf("err = %v, want ErrNotInSubdomain", err)
		}
	})
}

func TestDatabaseLogs_Execute(t *testing.T) {
	dbID := uuid.New()
	subID := uuid.New()

	t.Run("not provisioned reads as empty logs", func(t *testing.T) {
		dr, sm, svc := newDBMocks(t)
		dr.EXPECT().FindByID(dbID.String(), gomock.Any()).
			Return(&domain.Database{ID: dbID, SubdomainID: subID}, nil)
		sm.EXPECT().DatabaseLogs(gomock.Any(), dbID.String(), gomock.Any()).
			Return(nil, fmt.Errorf("servermanager logs: %w", ports.ErrDatabaseNotProvisioned))
		entries, err := svc.Logs.Execute(context.Background(), subID.String(), dbID.String(), ports.LogOptions{Tail: 10})
		if err != nil || len(entries) != 0 {
			t.Fatalf("entries = %v err = %v, want empty and nil", entries, err)
		}
	})

	t.Run("wrong subdomain is not found", func(t *testing.T) {
		dr, _, svc := newDBMocks(t)
		dr.EXPECT().FindByID(dbID.String(), gomock.Any()).
			Return(&domain.Database{ID: dbID, SubdomainID: uuid.New()}, nil)
		if _, err := svc.Logs.Execute(context.Background(), subID.String(), dbID.String(), ports.LogOptions{}); !errors.Is(err, ErrNotInSubdomain) {
			t.Fatalf("err = %v, want ErrNotInSubdomain", err)
		}
	})
}

func TestDatabaseMetrics_Execute(t *testing.T) {
	dbID := uuid.New()
	subID := uuid.New()

	t.Run("passes the manager envelope through", func(t *testing.T) {
		dr, sm, svc := newDBMocks(t)
		envelope := &ports.MetricsResponse{Range: "24h", Series: map[string][]ports.MetricPoint{
			"conn": {{Time: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC), Value: 3}},
			"qps":  {},
		}}
		dr.EXPECT().FindByID(dbID.String(), gomock.Any()).
			Return(&domain.Database{ID: dbID, SubdomainID: subID}, nil)
		sm.EXPECT().DatabaseMetrics(gomock.Any(), dbID.String(), "24h").Return(envelope, nil)

		got, err := svc.Metrics.Execute(context.Background(), subID.String(), dbID.String(), "24h")
		if err != nil || got != envelope {
			t.Fatalf("got = %+v, err = %v, want the envelope untouched", got, err)
		}
	})

	t.Run("not provisioned reads as an empty envelope", func(t *testing.T) {
		dr, sm, svc := newDBMocks(t)
		dr.EXPECT().FindByID(dbID.String(), gomock.Any()).
			Return(&domain.Database{ID: dbID, SubdomainID: subID}, nil)
		sm.EXPECT().DatabaseMetrics(gomock.Any(), dbID.String(), "24h").
			Return(nil, fmt.Errorf("servermanager metrics: %w", ports.ErrDatabaseNotProvisioned))

		got, err := svc.Metrics.Execute(context.Background(), subID.String(), dbID.String(), "24h")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil || got.Series == nil || len(got.Series) != 0 || got.Range != "24h" {
			t.Errorf("got = %+v, want an empty envelope with a non-nil series map", got)
		}
	})

	// The manager is never called for a database in someone else's subdomain
	// — no sm.EXPECT() proves it.
	t.Run("wrong subdomain is not found", func(t *testing.T) {
		dr, _, svc := newDBMocks(t)
		dr.EXPECT().FindByID(dbID.String(), gomock.Any()).
			Return(&domain.Database{ID: dbID, SubdomainID: uuid.New()}, nil)
		if _, err := svc.Metrics.Execute(context.Background(), subID.String(), dbID.String(), "24h"); !errors.Is(err, ErrNotInSubdomain) {
			t.Fatalf("err = %v, want ErrNotInSubdomain", err)
		}
	})
}

func TestProvisionPoller(t *testing.T) {
	dbID := uuid.New().String()
	creds := ProvisionCreds{User: "app", Password: "0123456789abcdef0123456789abcdef", DBName: "appdb"} // gitleaks:allow — fixture, not a secret

	newPoller := func(t *testing.T) (*mocks.MockDatabaseRepo, *mocks.MockServerManagerClient, *ProvisionPoller) {
		t.Helper()
		ctrl := gomock.NewController(t)
		dr := mocks.NewMockDatabaseRepo(ctrl)
		sm := mocks.NewMockServerManagerClient(ctrl)
		clock := mocks.NewMockClock(ctrl)
		clock.EXPECT().Now().Return(time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)).AnyTimes()
		p := NewProvisionPoller(dr, sm, clock)
		p.Interval = time.Millisecond
		return dr, sm, p
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

	t.Run("provision running writes status and connection string", func(t *testing.T) {
		dr, sm, p := newPoller(t)
		sm.EXPECT().DatabaseStatus(gomock.Any(), dbID).Return(&ports.DBStatus{
			Exists: true, Running: true, Health: "healthy",
			Provision: &ports.DBProvisionState{
				ID: "prov-1", Status: "running", Image: "postgres:16",
				ContainerID: "cid-1", ContainerName: "stt-db-" + dbID,
			},
		}, nil)

		done := make(chan map[string]any, 1)
		dr.EXPECT().Update(dbID, gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				done <- updates
				return nil
			})

		p.Watch(dbID, creds)
		updates := waitUpdate(t, done)
		if updates["status"] != domain.DatabaseStatusRunning {
			t.Errorf("status = %v", updates["status"])
		}
		want := "postgres://app:0123456789abcdef0123456789abcdef@stt-db-" + dbID + ":5432/appdb?sslmode=disable"
		if updates["connection_string"] != want {
			t.Errorf("connection_string = %v, want %s", updates["connection_string"], want)
		}
		if updates["docker_image"] != "postgres:16" || updates["docker_container_name"] != "stt-db-"+dbID {
			t.Errorf("docker fields = %+v", updates)
		}
	})

	t.Run("provision failed marks the record failed", func(t *testing.T) {
		dr, sm, p := newPoller(t)
		sm.EXPECT().DatabaseStatus(gomock.Any(), dbID).Return(&ports.DBStatus{
			Provision: &ports.DBProvisionState{ID: "prov-2", Status: "failed", Error: "pull failed"},
		}, nil)
		done := make(chan map[string]any, 1)
		dr.EXPECT().Update(dbID, gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				done <- updates
				return nil
			})

		p.Watch(dbID, creds)
		if updates := waitUpdate(t, done); updates["status"] != domain.DatabaseStatusFailed {
			t.Errorf("status = %v, want failed", updates["status"])
		}
	})

	t.Run("lost provision strikes out as failed", func(t *testing.T) {
		dr, sm, p := newPoller(t)
		// Manager restarted mid-provision: no job, no container.
		sm.EXPECT().DatabaseStatus(gomock.Any(), dbID).
			Return(&ports.DBStatus{Exists: false}, nil).Times(maxLostStrikes)
		done := make(chan map[string]any, 1)
		dr.EXPECT().Update(dbID, gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				done <- updates
				return nil
			})

		p.Watch(dbID, creds)
		if updates := waitUpdate(t, done); updates["status"] != domain.DatabaseStatusFailed {
			t.Errorf("status = %v, want failed", updates["status"])
		}
	})

	t.Run("healthy container without a job still counts as success", func(t *testing.T) {
		dr, sm, p := newPoller(t)
		// Manager restarted after the container came up: Docker state is truth.
		sm.EXPECT().DatabaseStatus(gomock.Any(), dbID).
			Return(&ports.DBStatus{Exists: true, Running: true, Health: "healthy"}, nil)
		done := make(chan map[string]any, 1)
		dr.EXPECT().Update(dbID, gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ string, updates map[string]any, _ context.Context) error {
				done <- updates
				return nil
			})

		p.Watch(dbID, creds)
		updates := waitUpdate(t, done)
		if updates["status"] != domain.DatabaseStatusRunning {
			t.Errorf("status = %v, want running", updates["status"])
		}
		if updates["docker_container_name"] != "stt-db-"+dbID {
			t.Errorf("container name = %v (derived fallback expected)", updates["docker_container_name"])
		}
	})
}
