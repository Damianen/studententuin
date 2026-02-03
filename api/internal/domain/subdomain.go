package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subdomain struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primairyKey"`
	Name string
	FullDomain string
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User User
	applications
	IsActive bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
