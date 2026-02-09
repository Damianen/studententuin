package domain

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationStatus string
type ApplicationType string

const (
	ApplicationStatusPending   ApplicationStatus = "pending"
	ApplicationStatusRunning   ApplicationStatus = "running"
	ApplicationStatusStopped   ApplicationStatus = "stopped"
	ApplicationStatusFailed    ApplicationStatus = "failed"
)

const (
	ApplicationTypeNodejs ApplicationType = "Nodejs"
)

type Application struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SubdomainID uuid.UUID `gorm:"type:uuid;not null;index"`
	Subdomain Subdomain
	Name *string `gorm:"type:text"`
	RepositoryURL *string `gorm:"type:text"`
	Branch        *string `gorm:"type:text"`
	Type     ApplicationType `gorm:"type:text;not null"`
	BuildCommand *string `gorm:"type:text"`
	StartCommand *string `gorm:"type:text"`
	EnvironmentVariables map[string]string `gorm:"type:jsonb"`
	DockerImage        *string `gorm:"type:text"`
	DockerContainerID  *string `gorm:"type:text"`
	DockerContainerName *string `gorm:"type:text"`
	Ports   []int  `gorm:"type:jsonb"`
	Volumes []string `gorm:"type:jsonb"`
	MemoryLimit *string `gorm:"type:text"`
	CpuLimit    *string `gorm:"type:text"`
	Status ApplicationStatus `gorm:"type:text;not null"`
	LastDeployedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
