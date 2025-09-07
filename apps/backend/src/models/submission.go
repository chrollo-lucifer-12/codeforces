package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ContestID   uuid.UUID `gorm:"type:uuid;not null;index"`
	ChallengeID uuid.UUID `gorm:"type:uuid;not null;index"`
	Points      int
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
