package controllers

import (
	"os"

	"github.com/AnshVM/flashpoll-backend/db"
	"github.com/AnshVM/flashpoll-backend/models"
	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	var user models.User
	accessToken, err := getAccessToken(ctx)
	if err != nil {
		return
	}
	claims, err := parseToken(accessToken, []byte(os.Getenv("ACCESS_TOKENS_SECRET_KEY")))
	if err != nil {
		unauthorized(ctx)
		return
	}
	if err := db.DB.First(&user, claims.UserID).Error; err != nil {
		unauthorized(ctx)
		return
	}
	ctx.Set("user", &user)
	ctx.Next()
}
