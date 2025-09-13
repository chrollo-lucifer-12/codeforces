package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chrollo-lucifer-12/backend/src/config"
	"github.com/chrollo-lucifer-12/backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	users := rg.Group("/user")
	{
		users.GET("/me", func(ctx *gin.Context) {
			authHeader := ctx.GetHeader("Authorization")
			if authHeader == "" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
				return
			}

			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(config.LoadConfig().JwtSecret), nil
			})

			if err != nil || !token.Valid {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			}

			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
				return
			}

			userId, err := uuid.Parse(userIDStr)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id format"})
				return
			}

			var findUser models.User
			result := db.Where("id = ?", userId).First(&findUser)
			if result.Error != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{"user": findUser})
		})
	}
}
