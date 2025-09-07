package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleUser Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Role Role `gorm:"type:user_role;default:'user';not null"`
}
