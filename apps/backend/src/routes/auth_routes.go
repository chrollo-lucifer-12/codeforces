package routes

import (
	"github.com/chrollo-lucifer-12/backend/src/controllers"
	"github.com/chrollo-lucifer-12/backend/src/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateAuthRoutes(rg *gin.RouterGroup, database *gorm.DB) {

	userService := services.NewUserService(database)
	sessionService := services.NewSessionService(database)
	tokenService := services.NewTokenService(database)

	auth := rg.Group("/auth")
	{
		auth.POST("/register", func(ctx *gin.Context) {
			controllers.Register(ctx, userService)
		})

		auth.POST("/login", func(ctx *gin.Context) {
			controllers.Login(ctx, userService, sessionService, tokenService)
		})
	}
}
