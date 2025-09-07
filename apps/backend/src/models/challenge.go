package models

import "github.com/google/uuid"

type Challenge struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title       string
	NotionDocId string
	MaxPoints   int
	Contests    []Contest `gorm:"many2many:contest_models;"`
}
