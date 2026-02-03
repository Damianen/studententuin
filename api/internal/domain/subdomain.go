package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subdomain struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primairyKey"`
	Name string `gorm:"type:text"`
	FullDomain string `gorm:"type:text;uniqueIndex"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User User
	Application *Application `gorm:"constraint:OnDelete:CASCADE"`
	Database *Database `gorm:"constraint:OnDelete:CASCADE"`
	IsActive bool `gorm:"type:bool"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
