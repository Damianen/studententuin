package deployments

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

const (
	// maxBuildLogBytes keeps the LAST chunk of the build log — the tail is
	// where the failure lives.
	maxBuildLogBytes = 64 * 1024
	// maxRetainedJobs bounds store memory; only terminal jobs are evicted,
	// oldest first (~13MiB worst case with full logs).
	maxRetainedJobs = 200

	truncationMarker = "…[log truncated]\n"
)

// JobStore is the in-memory deployment registry for the single-instance MVP.
// Jobs vanish on restart; container truth stays recoverable via the apps
// status endpoint.
type JobStore struct {
	mu     sync.Mutex
	jobs   map[string]*domain.DeploymentJob
	active map[string]string // appID → in-flight deployment ID (single-flight)
	order  []string          // creation order, for eviction
	clock  ports.Clock
}

func NewJobStore(clock ports.Clock) *JobStore {
	return &JobStore{
		jobs:   map[string]*domain.DeploymentJob{},
		active: map[string]string{},
		clock:  clock,
	}
}

// Create registers a queued job and claims the app's single deploy slot;
// a deploy already in flight for the app returns ErrConflict.
func (s *JobStore) Create(id, appID string) (domain.DeploymentJob, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if inFlight, ok := s.active[appID]; ok {
		return domain.DeploymentJob{}, fmt.Errorf(
			"deployment %s already in flight for this application: %w", inFlight, domain.ErrConflict)
	}

	now := s.clock.Now()
	job := &domain.DeploymentJob{
		ID:        id,
		AppID:     appID,
		Status:    domain.DeploymentStatusQueued,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.jobs[id] = job
	s.active[appID] = id
	s.order = append(s.order, id)
	s.evictLocked()
	return *job, nil
}

// Get returns a copy — callers never share memory with the running job.
func (s *JobStore) Get(id string) (domain.DeploymentJob, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		return domain.DeploymentJob{}, fmt.Errorf("deployment %s: %w", id, domain.ErrNotFound)
	}
	return *job, nil
}

func (s *JobStore) SetStatus(id string, status domain.DeploymentStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, ok := s.jobs[id]; ok {
		job.Status = status
		job.UpdatedAt = s.clock.Now()
	}
}

// SetSource records what the clone produced (sha plus display metadata) and
// which image will be built.
func (s *JobStore) SetSource(id string, checkout domain.SourceCheckout, image string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, ok := s.jobs[id]; ok {
		job.CommitSHA = checkout.CommitSHA
		job.CommitMessage = checkout.CommitMessage
		job.CommitAuthor = checkout.CommitAuthor
		job.Image = image
		job.UpdatedAt = s.clock.Now()
	}
}

// Fail marks the job terminal with a user-facing reason and frees the app's
// deploy slot.
func (s *JobStore) Fail(id, reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, ok := s.jobs[id]; ok {
		job.Status = domain.DeploymentStatusFailed
		job.Error = reason
		job.UpdatedAt = s.clock.Now()
		delete(s.active, job.AppID)
	}
}

// Succeed marks the job running, records the started container, and frees
// the app's deploy slot.
func (s *JobStore) Succeed(id, containerID, containerName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, ok := s.jobs[id]; ok {
		job.Status = domain.DeploymentStatusRunning
		job.ContainerID = containerID
		job.ContainerName = containerName
		job.UpdatedAt = s.clock.Now()
		delete(s.active, job.AppID)
	}
}

// AppendLog appends build output, keeping only the last maxBuildLogBytes
// (cut at a line boundary, with a truncation marker).
func (s *JobStore) AppendLog(id string, p []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		return
	}
	job.BuildLog += string(p)
	if len(job.BuildLog) > maxBuildLogBytes {
		cut := job.BuildLog[len(job.BuildLog)-maxBuildLogBytes:]
		if nl := strings.IndexByte(cut, '\n'); nl >= 0 && nl+1 < len(cut) {
			cut = cut[nl+1:]
		}
		job.BuildLog = truncationMarker + cut
	}
}

// LogWriter adapts AppendLog to io.Writer for the builder's output; safe for
// concurrent writers because AppendLog locks.
func (s *JobStore) LogWriter(id string) io.Writer { return jobLogWriter{store: s, id: id} }

type jobLogWriter struct {
	store *JobStore
	id    string
}

func (w jobLogWriter) Write(p []byte) (int, error) {
	w.store.AppendLog(w.id, p)
	return len(p), nil
}

// evictLocked drops the oldest terminal jobs over the retention cap. Active
// jobs are never evicted — their single-flight slot must stay accounted.
func (s *JobStore) evictLocked() {
	over := len(s.jobs) - maxRetainedJobs
	if over <= 0 {
		return
	}
	kept := s.order[:0]
	for _, id := range s.order {
		job, ok := s.jobs[id]
		if !ok {
			continue
		}
		if over > 0 && job.Terminal() {
			delete(s.jobs, id)
			over--
			continue
		}
		kept = append(kept, id)
	}
	s.order = kept
}
