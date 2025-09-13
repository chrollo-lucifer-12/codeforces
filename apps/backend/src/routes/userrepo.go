package routes

import (
	"net/http"

	"github.com/chrollo-lucifer-12/backend/src/middlewares"
	"github.com/chrollo-lucifer-12/backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserRepoRoutes(rg *gin.RouterGroup, database *gorm.DB) {
	userRepos := rg.Group("/repos")
	{

		userRepos.POST("/fork", middlewares.AuthMiddleware(database), func(ctx *gin.Context) {
			userIDVal, exists := ctx.Get("user_id")
			if !exists {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
				return
			}

			userIDStr, ok := userIDVal.(string)
			if !ok {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
				return
			}

			// Parse userID string to UUID
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id UUID"})
				return
			}

			var body struct {
				ChallengeID string `json:"challengeId"`
				RepoURL     string `json:"repoURL"`
			}

			if err := ctx.ShouldBindJSON(&body); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Parse ChallengeID string to UUID
			challengeID, err := uuid.Parse(body.ChallengeID)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge_id UUID"})
				return
			}

			newUserRepo := models.UserRepo{
				UserID:      userID,
				ChallengeID: challengeID,
				RepoURL:     body.RepoURL,
			}

			if err := database.Create(&newUserRepo).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(http.StatusCreated, newUserRepo)
		})

		userRepos.GET("/:userRepoId", func(ctx *gin.Context) {
			userRepoID := ctx.Param("userRepoId")

			var userRepo models.UserRepo
			result := database.First(&userRepo, "id = ?", userRepoID)
			if result.Error != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "user repo not found"})
				return
			}

			ctx.JSON(http.StatusOK, userRepo)
		})

		userRepos.GET("/user/:userId", func(ctx *gin.Context) {
			userID := ctx.Param("userId")

			var userRepos []models.UserRepo
			result := database.Where("user_id = ?", userID).Find(&userRepos)
			if result.Error != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}

			ctx.JSON(http.StatusOK, userRepos)
		})
	}
}
