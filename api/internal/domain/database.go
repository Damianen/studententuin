package domain

import (
	"time"

	"github.com/google/uuid"
)

type DatabaseType string

const (
	DatabaseTypePostgres DatabaseType = "postgres"
	DatabaseTypeMySQL    DatabaseType = "mysql"
	DatabaseTypeMongoDB  DatabaseType = "mongodb"
)

type DatabaseStatus string

const (
	DatabaseStatusProvisioning DatabaseStatus = "provisioning"
	DatabaseStatusRunning      DatabaseStatus = "running"
	DatabaseStatusStopped      DatabaseStatus = "stopped"
	DatabaseStatusFailed       DatabaseStatus = "failed"
)

type Database struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SubdomainID uuid.UUID `gorm:"type:uuid;not null;index"`
	Subdomain Subdomain
	Name    string       `gorm:"type:text;not null"`
	Type    DatabaseType `gorm:"type:text;not null"`
	Version string       `gorm:"type:text;not null"`
	Status  DatabaseStatus `gorm:"type:text;not null"`
	Host             *string `gorm:"type:text"`
	Port             *int    `gorm:"type:int"`
	DockerImage         *string `gorm:"type:text"`
	DockerContainerID   *string `gorm:"type:text"`
	DockerContainerName *string `gorm:"type:text"`
	Volumes []string `gorm:"type:jsonb"`
	MemoryLimit *string `gorm:"type:text"`
	CpuLimit    *string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
