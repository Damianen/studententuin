package deployments

import (
	"errors"
	"fmt"
	"strings"
	"sync"
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

	if _, err := s.Create("d1", "app-1"); err != nil {
		t.Fatalf("first Create: %v", err)
	}
	if _, err := s.Create("d2", "app-1"); !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("second Create = %v, want ErrConflict", err)
	}
	// A different app is unaffected.
	if _, err := s.Create("d3", "app-2"); err != nil {
		t.Fatalf("Create for other app: %v", err)
	}

	// Terminal states release the slot.
	s.Fail("d1", "boom")
	if _, err := s.Create("d4", "app-1"); err != nil {
		t.Fatalf("Create after Fail: %v", err)
	}
	s.Succeed("d4", "cid", "name")
	if _, err := s.Create("d5", "app-1"); err != nil {
		t.Fatalf("Create after Succeed: %v", err)
	}
}

func TestJobStoreGetReturnsCopies(t *testing.T) {
	s := NewJobStore(testClock())
	if _, err := s.Create("d1", "app-1"); err != nil {
		t.Fatal(err)
	}

	job, err := s.Get("d1")
	if err != nil {
		t.Fatal(err)
	}
	job.Status = domain.DeploymentStatusFailed // must not write through

	fresh, _ := s.Get("d1")
	if fresh.Status != domain.DeploymentStatusQueued {
		t.Errorf("store job mutated through a Get copy: %s", fresh.Status)
	}

	if _, err := s.Get("unknown"); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Get(unknown) = %v, want ErrNotFound", err)
	}
}

func TestJobStoreLifecycleFields(t *testing.T) {
	s := NewJobStore(testClock())
	if _, err := s.Create("d1", "app-1"); err != nil {
		t.Fatal(err)
	}

	s.SetStatus("d1", domain.DeploymentStatusCloning)
	s.SetSource("d1", domain.SourceCheckout{
		CommitSHA:     "abc123",
		CommitMessage: "fix: water the plants",
		CommitAuthor:  "Gardener",
	}, "stt-app-x:abcd1234")
	s.Succeed("d1", "cid-9", "stt-app-x")

	job, _ := s.Get("d1")
	if job.Status != domain.DeploymentStatusRunning || job.CommitSHA != "abc123" ||
		job.Image != "stt-app-x:abcd1234" || job.ContainerID != "cid-9" || job.ContainerName != "stt-app-x" {
		t.Errorf("job = %+v", job)
	}
	if job.CommitMessage != "fix: water the plants" || job.CommitAuthor != "Gardener" {
		t.Errorf("commit metadata = %q by %q, want threaded through", job.CommitMessage, job.CommitAuthor)
	}

	s.Fail("d1", "late failure should still record") // terminal transition still writes
	if job, _ = s.Get("d1"); job.Error == "" {
		t.Errorf("Fail after Succeed recorded nothing: %+v", job)
	}
}

func TestJobStoreLogCap(t *testing.T) {
	s := NewJobStore(testClock())
	if _, err := s.Create("d1", "app-1"); err != nil {
		t.Fatal(err)
	}
	w := s.LogWriter("d1")

	line := strings.Repeat("x", 1023) + "\n"
	for written := 0; written < maxBuildLogBytes*2; written += len(line) {
		if _, err := w.Write([]byte(line)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := w.Write([]byte("THE END\n")); err != nil {
		t.Fatal(err)
	}

	job, _ := s.Get("d1")
	if len(job.BuildLog) > maxBuildLogBytes+len(truncationMarker) {
		t.Errorf("build log %d bytes, want <= %d", len(job.BuildLog), maxBuildLogBytes+len(truncationMarker))
	}
	if !strings.HasPrefix(job.BuildLog, truncationMarker) {
		t.Errorf("truncated log missing marker prefix: %q", job.BuildLog[:40])
	}
	if !strings.HasSuffix(job.BuildLog, "THE END\n") {
		t.Errorf("log tail lost, ends with %q", job.BuildLog[len(job.BuildLog)-20:])
	}
	// Cut lands on a line boundary: marker is followed by a full line.
	afterMarker := strings.TrimPrefix(job.BuildLog, truncationMarker)
	if !strings.HasPrefix(afterMarker, "x") && !strings.HasPrefix(afterMarker, "THE END") {
		t.Errorf("cut not at line boundary: %q", afterMarker[:20])
	}

	// Writes to unknown jobs are dropped, not panics.
	s.LogWriter("nope").Write([]byte("dropped")) //nolint:errcheck
}

func TestJobStoreEviction(t *testing.T) {
	s := NewJobStore(testClock())

	// One active job that must survive any eviction pressure.
	if _, err := s.Create("active", "app-active"); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < maxRetainedJobs+50; i++ {
		id := fmt.Sprintf("d%d", i)
		if _, err := s.Create(id, "app-"+id); err != nil {
			t.Fatal(err)
		}
		s.Fail(id, "done")
	}

	if n := len(s.jobs); n > maxRetainedJobs {
		t.Errorf("store holds %d jobs, want <= %d", n, maxRetainedJobs)
	}
	if _, err := s.Get("active"); err != nil {
		t.Errorf("active job evicted: %v", err)
	}
	if _, err := s.Get("d0"); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("oldest terminal job not evicted: %v", err)
	}
	newest := fmt.Sprintf("d%d", maxRetainedJobs+49)
	if _, err := s.Get(newest); err != nil {
		t.Errorf("newest job evicted: %v", err)
	}
}

// TestJobStoreConcurrency is the -race hammer: writers, readers, and status
// transitions on the same job.
func TestJobStoreConcurrency(t *testing.T) {
	s := NewJobStore(testClock())
	if _, err := s.Create("d1", "app-1"); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := s.LogWriter("d1")
			for j := 0; j < 500; j++ {
				w.Write([]byte("line from a concurrent writer\n")) //nolint:errcheck
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 500; j++ {
			if _, err := s.Get("d1"); err != nil {
				t.Error(err)
				return
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.SetStatus("d1", domain.DeploymentStatusCloning)
		s.SetSource("d1", domain.SourceCheckout{CommitSHA: "sha"}, "img")
		s.SetStatus("d1", domain.DeploymentStatusBuilding)
		s.SetStatus("d1", domain.DeploymentStatusStarting)
		s.Succeed("d1", "cid", "name")
	}()
	wg.Wait()

	if job, _ := s.Get("d1"); job.Status != domain.DeploymentStatusRunning {
		t.Errorf("final status = %s, want running", job.Status)
	}
}
