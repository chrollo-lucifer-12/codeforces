package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DifficultyType string

const (
	DifficultyEasy   DifficultyType = "easy"
	DifficultyMedium DifficultyType = "medium"
	DifficultyHard   DifficultyType = "hard"
)

type Challenge struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string         `gorm:"not null"`
	Slug        string         `gorm:"uniqueIndex;not null"`
	Difficulty  DifficultyType `gorm:"type:difficulty;default:'easy'"`
	Description string         `gorm:"type:text"`
	RepoURL     string
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
