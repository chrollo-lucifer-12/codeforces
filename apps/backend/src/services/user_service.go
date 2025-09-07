package services

import (
	"github.com/chrollo-lucifer-12/backend/src/models"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) FindUser(username string) (*models.User, error) {
	var exisitingUser models.User
	err := s.DB.First(&exisitingUser, "username = ?", username).Error

	return &exisitingUser, err
}

func (s *UserService) CreateUser(username, password string) (*models.User, error) {
	user := &models.User{
		Username: username,
		Password: string(password),
	}

	if err := s.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
