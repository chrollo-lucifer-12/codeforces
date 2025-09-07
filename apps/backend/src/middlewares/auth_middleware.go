package middlewares

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chrollo-lucifer-12/backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var jwtSecret = []byte("axqiDDn?y|eWEAV")

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("AuthMiddleware: started")

		authHeader := ctx.GetHeader("Authorization")
		log.Println("Authorization header:", authHeader)

		var tokenStr string
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
				log.Println("JWT token found")
			}
		}

		refreshToken, _ := ctx.Cookie("refresh_token")
		log.Println("Refresh token from cookie:", refreshToken)

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Println("JWT invalid or expired:", err)

			if refreshToken == "" {
				log.Println("No refresh token, aborting")
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				ctx.Abort()
				return
			}

			var rt models.RefreshToken
			if err := db.Where("token = ?", refreshToken).First(&rt).Error; err != nil {
				log.Println("Refresh token not found:", err)
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
				ctx.Abort()
				return
			}

			if rt.ExpiresAt.Before(time.Now()) {
				log.Println("Refresh token expired")
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
				ctx.Abort()
				return
			}

			var session models.Session
			if err := db.Preload("User").Where("id = ?", rt.SessionID).First(&session).Error; err != nil {
				log.Println("Session not found:", err)
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
				ctx.Abort()
				return
			}

			log.Println("Generating new JWT for user:", session.UserID, "role:", session.User.Role)

			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userId":    session.UserID,
				"sessionId": session.ID.String(),
				"aud":       session.User.Role,
				"exp":       time.Now().Add(15 * time.Minute).Unix(),
				"iat":       time.Now().Unix(),
				"jti":       uuid.New().String(),
			})

			newTokenStr, err := newToken.SignedString(jwtSecret)
			if err != nil {
				log.Println("Failed to sign new JWT:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				ctx.Abort()
				return
			}

			ctx.Header("X-New-Token", newTokenStr)
			ctx.Set("userId", session.UserID)
			ctx.Set("session", session)
			ctx.Set("aud", session.User.Role)

			ctx.Next()
			return
		}

		// Token is valid
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Invalid token claims")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		log.Println("JWT valid. Claims:", claims)

		ctx.Set("userId", claims["userId"])
		ctx.Set("sessionId", claims["sessionId"])
		ctx.Set("aud", claims["aud"])

		ctx.Next()
		log.Println("AuthMiddleware: finished")
	}
}
