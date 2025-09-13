package routes

import (
	"log"
	"net/http"

	"github.com/chrollo-lucifer-12/backend/src/middlewares"
	"github.com/chrollo-lucifer-12/backend/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

	}
}
