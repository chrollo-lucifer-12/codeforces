package db

import (
	"fmt"

	"github.com/chrollo-lucifer-12/backend/src/config"
	"github.com/chrollo-lucifer-12/backend/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Database() (*gorm.DB, error) {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {

		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.UserRepo{}, &models.Challenge{}, &models.Submission{})
	if err != nil {
		return nil, err
	}

	return db, err
}
