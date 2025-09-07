package controllers

import (
	"errors"
	"net/http"

	"github.com/chrollo-lucifer-12/backend/src/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterRequest struct {
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
