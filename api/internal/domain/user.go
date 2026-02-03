package domain

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email string `gorm:"type:text;uniqueIndex"`
	EmailVerifiedAt *time.Time
	DisplayName string
	Status string `gorm:"type:text;not null;default:'active'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	AuthIdentities []AuthIdentity `gorm:"constraint:OnDelete:CASCADE"`
	PasswordCred *PasswordCredential `gorm:"constraint:OnDelete:CASCADE"`
}
