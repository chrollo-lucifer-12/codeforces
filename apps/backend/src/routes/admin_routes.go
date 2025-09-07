package routes

import (
	"github.com/chrollo-lucifer-12/backend/src/controllers"
	"github.com/chrollo-lucifer-12/backend/src/middlewares"
	"github.com/chrollo-lucifer-12/backend/src/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateAdminRoutes(rg *gin.RouterGroup, database *gorm.DB) {

	contestService := services.NewContestService(database)

	admin := rg.Group("/admin", middlewares.AuthMiddleware(database), middlewares.AdminMiddleware())
	{
		admin.POST("/create-contest", func(ctx *gin.Context) {
			controllers.CreateContest(ctx, contestService)
		})
	}
}