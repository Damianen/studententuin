package domain

import (
	"time"

	"github.com/google/uuid"
)

type PasswordCredential struct {
	UserId uuid.UUID `gorm:"type:uuid;primaryKey"`
	User User `gorm:"constraint:OnDelete:CASCADE"`
	PasswordHash string `gorm:"type:text;not null"`
	PasswordUpdatedAt time.Time `gorm:"not null;default:now()"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
