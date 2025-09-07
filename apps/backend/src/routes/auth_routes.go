package routes

import (
	"github.com/chrollo-lucifer-12/backend/src/controllers"
	"github.com/chrollo-lucifer-12/backend/src/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateAuthRoutes(rg *gin.RouterGroup, database *gorm.DB) {

	userService := services.NewUserService(database)

	auth := rg.Group("/auth")
	{
		auth.POST("/register", func(ctx *gin.Context) {
			controllers.Register(ctx, userService)
		})
	}
}
