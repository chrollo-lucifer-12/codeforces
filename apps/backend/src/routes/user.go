package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	users := rg.Group("/user")
	{
		users.GET("/")
	}
}
