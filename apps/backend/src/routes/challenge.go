package routes

import (
	"crypto/rand"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/chrollo-lucifer-12/backend/src/middlewares"
	"github.com/chrollo-lucifer-12/backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateChallengeInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	RepoURL     string `json:"repoURL"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		sb.WriteByte(charset[num.Int64()])
	}
	return sb.String()
}

func ChallengeRoutes(c *gin.RouterGroup, database *gorm.DB) {
	challenge := c.Group("/challenges")
	{
		challenge.GET("/all", middlewares.AuthMiddleware(database), func(ctx *gin.Context) {
			var challenges []models.Challenge
			database.Find(&challenges)

			log.Println(challenges)

			ctx.JSON(http.StatusFound, gin.H{"challenges": challenges})
		})

		challenge.GET("/:slug", middlewares.AuthMiddleware(database), func(ctx *gin.Context) {
			slug := ctx.Param("slug")

			var challenge models.Challenge

			if err := database.Where("slug = ?", slug).First(&challenge).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					ctx.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
					return
				}
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}

			ctx.JSON(http.StatusOK, challenge)
		})

		challenge.POST("/create", middlewares.AuthMiddleware(database), func(ctx *gin.Context) {

			var body CreateChallengeInput

			if err := ctx.ShouldBindJSON(&body); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			userIDVal, exists := ctx.Get("user_id")
			if !exists {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
				return
			}

			userID, ok := userIDVal.(string)
			if !ok {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
				return
			}

			userIdUUId, err := uuid.Parse(userID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var findUser models.User
			result := database.Where("id = ?", userIdUUId).First(&findUser)
			if result.Error != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "no user found"})
				return
			}

			if findUser.Role != "admin" {
				ctx.JSON(400, gin.H{"error": "user is not admin"})
				return
			}

			newChallenge := models.Challenge{
				Title:       body.Title,
				Difficulty:  models.DifficultyType(body.Difficulty),
				Description: body.Description,
				RepoURL:     body.RepoURL,
				Slug:        RandomString(6),
			}

			if err := database.Create(&newChallenge).Error; err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(200, gin.H{"challenge": newChallenge})
		})

	}
}
