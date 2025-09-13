package routes

import (
	"github.com/chrollo-lucifer-12/backend/src/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	users := rg.Group("/user")
	{
		users.GET("/me", middlewares.AuthMiddlware,func(ctx *gin.Context) {

		})
	}
}
