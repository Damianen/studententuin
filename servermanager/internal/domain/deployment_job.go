package domain

import "time"

type DeploymentStatus string

const (
	DeploymentStatusQueued   DeploymentStatus = "queued"
	DeploymentStatusCloning  DeploymentStatus = "cloning"
	DeploymentStatusBuilding DeploymentStatus = "building"
	DeploymentStatusStarting DeploymentStatus = "starting"
	DeploymentStatusRunning  DeploymentStatus = "running"
	DeploymentStatusFailed   DeploymentStatus = "failed"
)

// DeploymentJob tracks an async deploy (clone → build → start). Jobs live in
// memory for the MVP and are reconstructable from Docker state after a restart.
type DeploymentJob struct {
	ID            string
	AppID         string
	Status        DeploymentStatus
	Error         string
	BuildLog      string
	CommitSHA     string
	CommitMessage string
	CommitAuthor  string
	Image         string
	ContainerID   string
	ContainerName string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Terminal reports whether the job reached an end state.
func (j *DeploymentJob) Terminal() bool {
	return j.Status == DeploymentStatusRunning || j.Status == DeploymentStatusFailed
}
