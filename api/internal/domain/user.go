package domain

import (
	"time"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email string `gorm:"type:citex;uniqueIndex"`
	EmailVerifiedAt *time.Time
	DisplayName *string
	Status string `gorm:"type:text;not null;default:'active'"`
	Created_at time.Time
	Updated_at time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	AuthIdentities []AuthIdentity `gorm:"contraint:OnDelete:CASCADE"`
	PasswordCred *PasswordCredential `gorm:"contraint:OnDelete:CASCADE"`
}
