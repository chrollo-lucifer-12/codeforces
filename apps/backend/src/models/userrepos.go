package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRepo struct {
	ID          uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()"`
	UserID      uuid.UUID `gorm:"not null;index"`
	ChallengeID uuid.UUID `gorm:"not null;index"`
	RepoURL     string    
	LastCommit  string    
	Status      string    
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
