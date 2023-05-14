package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AnshVM/flashpoll-backend/db"
	"github.com/AnshVM/flashpoll-backend/models"
	"github.com/AnshVM/flashpoll-backend/types"
	"github.com/AnshVM/flashpoll-backend/ws"
	"github.com/gin-gonic/gin"
)

type OptionResponse = types.OptionResponse
type CreatePollRequest = types.CreatePollRequest
type GetPollResponse = types.GetPollResponse
type UpdatePollResponse = types.UpdatePollResponse

func newPoll(title string, userID uint, optionNames []string) *models.Poll {
	var options []models.Option
	for _, option := range optionNames {
		options = append(options, models.Option{Name: option, Count: 0})
	}
	return &models.Poll{
		Title:   title,
		User:    userID,
		Options: options,
	}
}

func getPollData(ctx *gin.Context, pollID uint, dest *GetPollResponse) error {
	var poll models.Poll
	err := db.FindById(uint(pollID), &poll)
	if err != nil {
		invalidRequestBody(ctx)
		return err
	}
	fmt.Printf("%+v", poll)

	var options []models.Option
	db.DB.Model(&poll).Association("Options").Find(&options)
	fmt.Printf("%+v", options)

	var resOptions []OptionResponse
	var totalVotes uint = 0

	for _, v := range options {
		resOptions = append(resOptions, OptionResponse{Name: v.Name, Votes: v.Count, ID: v.ID})
		totalVotes = totalVotes + v.Count
	}

	for i, v := range resOptions {
		if v.Votes == 0 && totalVotes == 0 {
			resOptions[i].VotesPercent = 0
			continue
		}
		resOptions[i].VotesPercent = (float32(v.Votes) / float32(totalVotes)) * 100
	}

	*dest = GetPollResponse{
		Title:      poll.Title,
		Options:    resOptions,
		TotalVotes: totalVotes,
		ID:         pollID,
	}
	return nil
}

func CreatePoll(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)

	var req CreatePollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		invalidRequestBody(ctx)
		return
	}

	poll := newPoll(req.Title, user.ID, req.OptionNames)
	err := db.DB.Create(poll).Error

	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": poll.ID})

}

func GetPollById(ctx *gin.Context) {

	user := ctx.MustGet("user").(*models.User)
	pollID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		invalidRequestBody(ctx)
		return
	}

	var poll GetPollResponse
	err = getPollData(ctx, uint(pollID), &poll)
	if err != nil {
		return
	}
	var option models.Option
	getUserVoteForPoll(user, uint(pollID), &option)

	if option.ID != 0 {
		poll.UserVote = OptionResponse{ID: option.ID, Name: option.Name, Votes: option.Count}
	}

	ctx.JSON(http.StatusOK, poll)
}

type SubmitVoteReq struct {
	OptionID uint `json:"optionID"`
}

func SubmitVote(ctx *gin.Context, hub *ws.Hub) {

	user := ctx.MustGet("user").(*models.User)
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

	var votedOption models.Option
	getUserVoteForPoll(user, option.Poll, &votedOption)
	if votedOption.ID != 0 {
		ctx.JSON(http.StatusConflict, "ALREADY_VOTED")
		return
	}

	db.DB.Model(&user).Association("Votes").Append(&option)
	option.Count = option.Count + 1
	if err := db.DB.Save(&option).Error; err != nil {
		fmt.Println(err)
		return
	}

	var pollData GetPollResponse
	getPollData(ctx, option.Poll, &pollData)
	wsUpdate := UpdatePollResponse{
		Options:    pollData.Options,
		TotalVotes: pollData.TotalVotes,
		ID:         pollData.ID,
	}
	hub.Broadcast <- wsUpdate

	ctx.JSON(http.StatusOK, "VOTE_SUBMITTED")

}

func getUserVoteForPoll(user *models.User, pollID uint, dest *models.Option) {
	db.DB.Model(&user).Association("Votes").Find(dest, "poll = ?", pollID)
}
