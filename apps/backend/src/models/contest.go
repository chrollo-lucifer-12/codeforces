package models

import (
	"time"

	"github.com/google/uuid"
)

type Contest struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title      string
	StartTime  time.Time
	Challenges []Challenge `gorm:"many2many:contest_models;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
