package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleType string

const (
	RoleUser  RoleType = "user"
	RoleAdmin RoleType = "admin"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Role      string         `gorm:"type:role;default:'user'"`
	Username  string         `gorm:"uniqueIndex;size:100;not null"`
	Email     string         `gorm:"uniqueIndex;size:255;not null"`
	Password  string         `gorm:"size:255;not null"`
	Fullname  string         `gorm:"size:255;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
