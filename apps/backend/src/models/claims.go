package models

import (
	"github.com/google/uuid"
)

type Claims struct {
	UserID    uuid.UUID `json:"userId"`
	SessionID uuid.UUID `json:"sessionId"`
	Audience  []string  `gorm:"type:jsonb"`
}
