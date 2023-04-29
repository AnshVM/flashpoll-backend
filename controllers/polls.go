package controllers

import (
	"net/http"
	"os"

	"strconv"

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
		return
	}

	claims, err := parseToken(accessToken, []byte(os.Getenv("ACCESS_TOKENS_SECRET_KEY")))

	if err != nil {
		unauthorized(ctx)
		return
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

type OptionRes struct {
	Name  string `json:"name"`
	Votes uint   `json:"votes"`
	ID    uint   `json:"id"`
}

type GetPollResponse struct {
	Title   string      `json:"title"`
	Options []OptionRes `json:"options"`
}

func GetPollById(ctx *gin.Context) {
	pollID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		invalidRequestBody(ctx)
		return
	}

	var poll models.Poll
	db.FindById(uint(pollID), &poll)

	var options []models.Option
	db.DB.Model(&poll).Association("Options").Find(&options)

	var resOptions []OptionRes

	for _, v := range options {
		resOptions = append(resOptions, OptionRes{Name: v.Name, Votes: v.Count, ID: v.ID})
	}

	res := GetPollResponse{
		Title:   poll.Title,
		Options: resOptions,
	}

	ctx.JSON(http.StatusOK, res)
}

type SubmitVoteReq struct {
	OptionID uint `json:"optionID"`
}

func SubmitVote(ctx *gin.Context) {

	accessToken, err := getAccessToken(ctx)

	if err != nil {
		unauthorized(ctx)
		return
	}

	claims, err := parseToken(accessToken, []byte(os.Getenv("ACCESS_TOKENS_SECRET_KEY")))

	if err != nil {
		unauthorized(ctx)
		return
	}

	var req SubmitVoteReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		invalidRequestBody(ctx)
		return
	}

	var option models.Option
	if err := db.FindById(req.OptionID, &option); err != nil {
		invalidRequestBody(ctx)
		return
	}

	var user models.User
	if err := db.FindById(claims.UserID, &user); err != nil {
		invalidRequestBody(ctx)
		return
	}
	db.DB.Model(&user).Association("Votes").Append(&option)

	option.Count = option.Count + 1
	if err := db.DB.Save(&option).Error; err != nil {
		return
	}

	ctx.JSON(http.StatusOK, "VOTE_SUBMITTED")

}
