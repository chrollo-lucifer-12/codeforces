package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		CreateAuthRoutes(api, db)
	}

	return r
}