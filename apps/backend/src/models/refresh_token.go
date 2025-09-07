package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	SessionID uuid.UUID `gorm:"type:uuid;not null;index"`
	Session   Session   `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"default:false;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
