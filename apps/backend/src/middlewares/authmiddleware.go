package middlewares

import "github.com/gin-gonic/gin"

func AuthMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		
		
		
		c.Next()
	}
}
