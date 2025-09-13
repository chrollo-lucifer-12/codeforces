package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api/v1")
	UserRoutes(api, db)
	AuthRoutes(api, db)
}
