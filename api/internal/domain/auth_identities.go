package domain

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuthIdentity struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User User
	Provider string `gorm:"type:text;not null;index"`
	ProviderUserID *string `gorm:"type:text"`
	Email *string `gorm:"type:text"`
	EmailVerified *bool
	LastLoginAt *time.Time
	Metadata datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
