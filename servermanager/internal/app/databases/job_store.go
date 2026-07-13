package databases

import (
	"fmt"
	"sync"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

// maxRetainedJobs bounds store memory; terminal jobs of other databases are
// evicted oldest-first. Far above anything a single hosting server sees.
const maxRetainedJobs = 200

// JobStore is the in-memory provision registry for the single-instance MVP,
// keyed by database ID: the api polls the database, not the job, so only the
// latest job per database matters. A non-terminal latest job doubles as the
// single-flight lock. Jobs vanish on restart; container truth stays
// recoverable via the status endpoint.
type JobStore struct {
	mu    sync.Mutex
	jobs  map[string]*domain.ProvisionJob // dbID → latest job
	order []string                        // dbID insertion order, for eviction
	clock ports.Clock
}

func NewJobStore(clock ports.Clock) *JobStore {
	return &JobStore{
		jobs:  map[string]*domain.ProvisionJob{},
		clock: clock,
	}
}

// Create registers a queued job for the database; a provision already in
// flight returns ErrConflict. A terminal previous job is replaced.
func (s *JobStore) Create(id, dbID string) (domain.ProvisionJob, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if prev, ok := s.jobs[dbID]; ok && !prev.Terminal() {
		return domain.ProvisionJob{}, fmt.Errorf(
			"provision %s already in flight for this database: %w", prev.ID, domain.ErrConflict)
	}

	now := s.clock.Now()
	job := &domain.ProvisionJob{
		ID:        id,
		DBID:      dbID,
		Status:    domain.ProvisionStatusQueued,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, existed := s.jobs[dbID]; !existed {
		s.order = append(s.order, dbID)
	}
	s.jobs[dbID] = job
	s.evictLocked()
	return *job, nil
}

// Latest returns a copy of the database's most recent job, if any.
func (s *JobStore) Latest(dbID string) (domain.ProvisionJob, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[dbID]
	if !ok {
		return domain.ProvisionJob{}, false
	}
	return *job, true
}

func (s *JobStore) SetStatus(dbID string, status domain.ProvisionStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, ok := s.jobs[dbID]; ok {
		job.Status = status
		job.UpdatedAt = s.clock.Now()
	}
}

// Fail marks the latest job terminal with a user-facing reason.
func (s *JobStore) Fail(dbID, reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, ok := s.jobs[dbID]; ok {
		job.Status = domain.ProvisionStatusFailed
		job.Error = reason
		job.UpdatedAt = s.clock.Now()
	}
}

// Succeed marks the latest job running and records the started container.
func (s *JobStore) Succeed(dbID, image, containerID, containerName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, ok := s.jobs[dbID]; ok {
		job.Status = domain.ProvisionStatusRunning
		job.Image = image
		job.ContainerID = containerID
		job.ContainerName = containerName
		job.UpdatedAt = s.clock.Now()
	}
}

// Drop forgets the database's job — used by delete so a removed database
// doesn't keep reporting a stale provision state.
func (s *JobStore) Drop(dbID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, dbID)
}

// evictLocked removes the oldest terminal jobs over the retention cap.
// In-flight jobs are never evicted — they are the single-flight lock.
func (s *JobStore) evictLocked() {
	over := len(s.jobs) - maxRetainedJobs
	if over <= 0 {
		return
	}
	kept := s.order[:0]
	for _, dbID := range s.order {
		job, ok := s.jobs[dbID]
		if !ok {
			continue
		}
		if over > 0 && job.Terminal() {
			delete(s.jobs, dbID)
			over--
			continue
		}
		kept = append(kept, dbID)
	}
	s.order = kept
}
