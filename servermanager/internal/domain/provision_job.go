package domain

import "time"

type ProvisionStatus string

const (
	ProvisionStatusQueued   ProvisionStatus = "queued"
	ProvisionStatusPulling  ProvisionStatus = "pulling"
	ProvisionStatusStarting ProvisionStatus = "starting"
	ProvisionStatusRunning  ProvisionStatus = "running"
	ProvisionStatusFailed   ProvisionStatus = "failed"
)

// ProvisionJob tracks an async database provision (pull → create → start →
// wait healthy). Deliberately has no credential fields: the POSTGRES_* values
// pass through to the container env and are never retained or logged.
type ProvisionJob struct {
	ID            string
	DBID          string
	Status        ProvisionStatus
	Error         string
	Image         string
	ContainerID   string
	ContainerName string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Terminal reports whether the job reached an end state.
func (j *ProvisionJob) Terminal() bool {
	return j.Status == ProvisionStatusRunning || j.Status == ProvisionStatusFailed
}
