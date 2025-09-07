package services

import (
	"time"

	"github.com/chrollo-lucifer-12/backend/src/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenService struct {
	DB *gorm.DB
}

var jwtSecret = []byte("axqiDDn?y|eWEAV")

func NewTokenService(db *gorm.DB) *TokenService {
	return &TokenService{DB: db}
}

func (t *TokenService) GenerateJWT(userId uuid.UUID, sessionId uuid.UUID, role string, ttl time.Duration) (string, error) {
	claims := models.Claims{
		UserID:    userId,
		SessionID: sessionId,
		Audience:  []string{role}, // optional DB persistence
	}

	t.DB.Create(&claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    claims.UserID,
		"sessionId": claims.SessionID,
		"aud":       role, // store as string for easier AdminMiddleware check
		"exp":       time.Now().Add(ttl).Unix(),
		"iat":       time.Now().Unix(),
		"jti":       uuid.New().String(),
	})

	return token.SignedString(jwtSecret)
}

func (t *TokenService) CreateRefreshToken(sessionId uuid.UUID, ttl time.Duration) (*models.RefreshToken, error) {

	if err := t.DB.Where("session_id = ?", sessionId).Delete(&models.RefreshToken{}).Error; err != nil {
		return nil, err
	}

	token := uuid.New().String()
	refreshToken := &models.RefreshToken{
		SessionID: sessionId,
		Token:     token,
		ExpiresAt: time.Now().Add(ttl),
	}

	if err := t.DB.Create(refreshToken).Error; err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (t *TokenService) RevokeRefreshToken(token string) error {
	return t.DB.Model(&models.RefreshToken{}).Where("token = ?", token).Update("revoked", true).Error
}
