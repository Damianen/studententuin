package databases

import (
	"errors"
	"testing"
	"time"

	"servermanager/internal/domain"
)

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

func testClock() fixedClock {
	return fixedClock{t: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)}
}

func TestJobStoreSingleFlight(t *testing.T) {
	s := NewJobStore(testClock())

	if _, err := s.Create("p1", "db-1"); err != nil {
		t.Fatalf("first Create: %v", err)
	}
	if _, err := s.Create("p2", "db-1"); !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("second Create = %v, want ErrConflict", err)
	}
	// A different database is unaffected.
	if _, err := s.Create("p3", "db-2"); err != nil {
		t.Fatalf("Create for other db: %v", err)
	}

	// A terminal job frees the slot and is replaced by the next provision.
	s.Fail("db-1", "pull failed")
	if _, err := s.Create("p4", "db-1"); err != nil {
		t.Fatalf("Create after terminal: %v", err)
	}
	job, ok := s.Latest("db-1")
	if !ok || job.ID != "p4" || job.Status != domain.ProvisionStatusQueued {
		t.Fatalf("Latest = %+v, %v; want p4 queued", job, ok)
	}
}

func TestJobStoreLifecycle(t *testing.T) {
	s := NewJobStore(testClock())
	if _, err := s.Create("p1", "db-1"); err != nil {
		t.Fatal(err)
	}

	s.SetStatus("db-1", domain.ProvisionStatusPulling)
	if job, _ := s.Latest("db-1"); job.Status != domain.ProvisionStatusPulling {
		t.Fatalf("Status = %s, want pulling", job.Status)
	}

	s.Succeed("db-1", "postgres:16", "cid-1", "stt-db-db-1")
	job, _ := s.Latest("db-1")
	if job.Status != domain.ProvisionStatusRunning || job.Image != "postgres:16" ||
		job.ContainerID != "cid-1" || job.ContainerName != "stt-db-db-1" {
		t.Fatalf("after Succeed: %+v", job)
	}

	s.Drop("db-1")
	if _, ok := s.Latest("db-1"); ok {
		t.Fatal("Latest after Drop should report no job")
	}
}
