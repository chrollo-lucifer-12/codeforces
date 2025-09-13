package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	ID         uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()"`
	UserRepoID uuid.UUID `gorm:"not null;index"`
	CommitSHA  string    `gorm:"not null;index"`
	Status     string
	Score      float64
	CreatedAt  time.Time
}
