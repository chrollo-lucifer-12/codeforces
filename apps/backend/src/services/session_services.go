package services

import (
	"time"

	"github.com/chrollo-lucifer-12/backend/src/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionService struct {
	DB *gorm.DB
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{DB: db}
}

func (s *SessionService) CreateSession(userId uuid.UUID) (*models.Session,  error) {

	expiration := time.Now().Add(240 * time.Hour)

	session := &models.Session{
		UserID: userId,
		ExpiresAt: expiration,
	}

	if err := s.DB.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}
