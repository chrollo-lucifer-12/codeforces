package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("AdminMiddleware: started")

		audience, exists := ctx.Get("aud")
		if !exists {
			log.Println("No 'aud' value set in context")
			ctx.JSON(400, gin.H{"error": "not an admin"})
			ctx.Abort()
			return
		}

		log.Println("Audience from context:", audience)

		if audience != "admin" {
			log.Println("Access denied: user is not an admin")
			ctx.JSON(400, gin.H{"error": "not an admin"})
			ctx.Abort()
			return
		}

		log.Println("Access granted: user is admin")
		ctx.Next()
		log.Println("AdminMiddleware: finished")
	}
}
