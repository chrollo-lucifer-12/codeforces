package routes

import (
	"github.com/chrollo-lucifer-12/backend/src/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	auth := rg.Group("/auth")
	{
		auth.POST("/signup", func(ctx *gin.Context) {
			handlers.SignupHandler(ctx, db)
		})
	}
}