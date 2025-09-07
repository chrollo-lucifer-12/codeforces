package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/chrollo-lucifer-12/backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleware(jwtSecret string, db *gorm.DB, jwtTTL time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tokenStr string

		authHeader := ctx.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
			}
		}

		refreshToken, _ := ctx.Cookie("refresh_token")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			if refreshToken == "" {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				ctx.Abort()
				return
			}

			var rt models.RefreshToken
			if err := db.Where("token = ?", refreshToken).First(&rt).Error; err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
				ctx.Abort()
				return
			}

			if rt.ExpiresAt.Before(time.Now()) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
				ctx.Abort()
				return
			}

			var session models.Session
			if err := db.Where("id = ?", rt.SessionID).First(&session).Error; err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
				ctx.Abort()
				return
			}

			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userId":    session.UserID,
				"sessionId": session.ID.String(),
				"aud":       "user",
				"exp":       time.Now().Add(jwtTTL).Unix(),
			})

			newTokenStr, err := newToken.SignedString([]byte(jwtSecret))
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
				ctx.Abort()
				return
			}

			ctx.Header("X-New-Token", newTokenStr)

			ctx.Set("userId", session.UserID)
			ctx.Set("session", session)

			ctx.Next()
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		ctx.Set("userId", claims["userId"])
		ctx.Set("sessionId", claims["sessionId"])
		ctx.Set("aud", claims["aud"])

		ctx.Next()
	}
}
