package domain

import (
	"time"

	"github.com/google/uuid"
)

type DeploymentStatus string

const (
	DeploymentStatusInFlight  DeploymentStatus = "in_flight"
	DeploymentStatusSucceeded DeploymentStatus = "succeeded"
	DeploymentStatusFailed    DeploymentStatus = "failed"
)

// Deployment is one persisted deploy attempt — the history behind the
// dashboard's deployments list. Duration is derived from the timestamps,
// never stored. Commit fields are pointers: a deploy that fails before the
// clone has no commit.
type Deployment struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ApplicationID uuid.UUID `gorm:"type:uuid;not null;index"`
	// ManagerDeploymentID keys the row to the manager job the poller follows,
	// so concurrent polls never finalize each other's rows.
	ManagerDeploymentID string           `gorm:"type:text;not null;uniqueIndex"`
	Status              DeploymentStatus `gorm:"type:text;not null"`
	Branch              *string          `gorm:"type:text"`
	CommitSHA           *string          `gorm:"type:text"`
	CommitMessage       *string          `gorm:"type:text"`
	CommitAuthor        *string          `gorm:"type:text"`
	Error               *string          `gorm:"type:text"`
	StartedAt           time.Time
	FinishedAt          *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
