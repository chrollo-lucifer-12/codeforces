package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/chrollo-lucifer-12/backend/src/config"
	"github.com/chrollo-lucifer-12/backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginHandler(ctx *gin.Context, db *gorm.DB) {
	var body LoginHandlerInput

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var findUser models.User

	if err := db.Where("username = ?", body.Username).Find(&findUser).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "username not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(body.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  findUser.ID.String(),
		"username": findUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	cfg := config.LoadConfig()
	log.Println(cfg.JwtSecret)

	tokenString, err := token.SignedString([]byte(cfg.JwtSecret))
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Login successful",
		"accessToken": tokenString,
	})

}
