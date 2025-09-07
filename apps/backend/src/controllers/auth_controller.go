package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/chrollo-lucifer-12/backend/src/models"
	"github.com/chrollo-lucifer-12/backend/src/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(ctx *gin.Context, userService *services.UserService) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request", "userId": nil})
		return
	}

	_, err := userService.FindUser(req.Username)
	if err == nil {
		ctx.JSON(400, gin.H{"error": "Username already exists", "userId": nil})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(500, gin.H{"error": err.Error(), "userId": nil})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password", "userId": nil})
		return
	}

	userCreated, err := userService.CreateUser(req.Username, string(hashedPassword))
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error(), "userId": nil})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "registration successful", "userId": userCreated.ID})
}

func Login(ctx *gin.Context, userService *services.UserService, sessionService *services.SessionService, tokenService *services.TokenService) {

	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request", "userId": nil})
		return
	}

	findUser, err := userService.FindUser(req.Username)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(req.Password))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Wrong password"})
	}

	session, err := sessionService.CreateSession(findUser.ID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Error creating session"})
	}

	jwt, err := tokenService.GenerateJWT(findUser.ID, session.ID, []string{"user"}, 15*time.Minute)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Error creating JWT"})
	}

	refreshToken, err := tokenService.CreateRefreshToken(session.ID, 7*24*time.Hour)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error creating refresh token"})
		return
	}

	ctx.SetCookie(
		"refresh_token",
		refreshToken.Token,
		int(7*24*time.Hour.Seconds()),
		"/",
		"",
		true,
		true,
	)

	ctx.JSON(200, gin.H{"userId": findUser.ID, "accessToken": jwt})
}

func Logout(ctx *gin.Context, sessionService *services.SessionService, tokenService *services.TokenService) {
	refreshTokenCookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(400, gin.H{"error": "No refresh token"})
		return
	}

	var rt models.RefreshToken
	if err := tokenService.DB.Where("token = ?", refreshTokenCookie).First(&rt).Error; err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid refresh token"})
		return
	}

	rt.Revoked = true
	if err := tokenService.DB.Save(&rt).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to revoke refresh token"})
		return
	}

	if err := sessionService.DeleteSession(rt.SessionID); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to revoke session"})
		return
	}

	ctx.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		true, 
		true, 
	)

	ctx.JSON(200, gin.H{"message": "Logged out successfully"})
}
