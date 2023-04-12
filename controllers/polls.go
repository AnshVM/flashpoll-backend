package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AnshVM/flashpoll-backend/db"
	"github.com/AnshVM/flashpoll-backend/models"
	"github.com/gin-gonic/gin"
)

type CreatePollRequest struct {
	Title       string   `json:"title"`
	OptionNames []string `json:"options"`
}

func CreatePoll(ctx *gin.Context) {
	accessToken, err := getAccessToken(ctx)

	if err != nil {
		unauthorized(ctx)
	}

	claims, err := parseToken(accessToken, []byte(os.Getenv("ACCESS_TOKENS_SECRET_KEY")))

	if err != nil {
		unauthorized(ctx)
	}

	var req CreatePollRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		invalidRequestBody(ctx)
		return
	}

	var options []models.Option

	for _, option := range req.OptionNames {
		options = append(options, models.Option{Name: option, Count: 0})
	}

	fmt.Printf("%+v\n", options)

	poll := models.Poll{
		Title:   req.Title,
		User:    claims.UserID,
		Options: options,
	}

	err = db.DB.Create(&poll).Error

	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": poll.ID})

}
