package routes

import (
	"github.com/chrollo-lucifer-12/backend/src/controllers"
	"github.com/chrollo-lucifer-12/backend/src/middlewares"
	"github.com/chrollo-lucifer-12/backend/src/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// get active and finished contest - with pagination 🟢
// get a specific contest returns all challenges too - param  🟢
// get a specific contest and challenge - challengeID 🟢
// submit a challenge - challengeId
// create contest, challenge - only admins

func CreateContestRoutes(rg *gin.RouterGroup, database *gorm.DB) {
	contestService := services.NewContestService(database)

	contest := rg.Group("/contest", middlewares.AuthMiddleware(database))
	{
		contest.GET("/active", func(ctx *gin.Context) {
			controllers.GetContests(ctx, contestService, "active")
		})

		contest.GET("/inactive", func(ctx *gin.Context) {
			controllers.GetContests(ctx, contestService, "inactive")
		})

		contest.GET("/:contestId", func(ctx *gin.Context) {
			controllers.GetContest(ctx, contestService)
		})

		contest.GET("/:contestId/challenge/:challengeId", func(ctx *gin.Context) {
			controllers.GetChallenge(ctx, contestService)
		})
	}
}
