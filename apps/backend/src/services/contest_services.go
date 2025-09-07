package services

import (
	"errors"
	"time"

	"github.com/chrollo-lucifer-12/backend/src/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContestService struct {
	DB *gorm.DB
}

func NewContestService(db *gorm.DB) *ContestService {
	return &ContestService{DB: db}
}

func (c *ContestService) GetContests(status string, limit, offset int) ([]models.Contest, error) {
	var contests []models.Contest

	query := c.DB.Select("title", "start_time").Where("status = ?", status).Limit(limit).Offset(offset)

	if err := query.Find(&contests).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("contest not found")
		}
		return nil, err
	}

	return contests, nil
}

func (c *ContestService) GetContest(contestId uuid.UUID) (*models.Contest, error) {
	var contest models.Contest

	if err := c.DB.Preload("Challenges").First(&contest, "id = ?", contestId).Error; err != nil {
		return nil, err
	}

	return &contest, nil
}

func (c *ContestService) GetChallenge(challengeId, contestId uuid.UUID) (*models.Challenge, error) {
	var challenge models.Challenge

	if err := c.DB.
		Joins("JOIN contest_models ON contest_models.challenge_id = challenges.id").
		Where("contest_models.contest_id = ? AND challenges.id = ?", contestId, challengeId).
		First(&challenge).Error; err != nil {
		return nil, err
	}

	return &challenge, nil
}

func (c *ContestService) SubmitChallenge (challengeId uuid.UUID) {
 // todo
}

func (c *ContestService) CreateContest (title string, startTime time.Time) (*models.Contest, error) {
	contest := &models.Contest {
		Title: title,
		StartTime: startTime,
		Status:  "active",
	}

	if err := c.DB.Create(contest).Error; err != nil {
		return nil, err
	}

	return contest, nil
}