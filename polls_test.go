package main

import (
	"fmt"
	"testing"

	"github.com/AnshVM/flashpoll-backend/db"
	"github.com/AnshVM/flashpoll-backend/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type OptionRes struct {
	Name  string `json:"name"`
	Votes uint   `json:"votes"`
	ID    uint   `json:"id"`
}

type GetPollResponse struct {
	Title      string      `json:"title"`
	Options    []OptionRes `json:"options"`
	TotalVotes uint        `json:"totalVotes"`
}

func login(r *gin.Engine, email string, password string) string {

	payload := map[string]any{
		"email":    email,
		"password": password,
	}
	response := request(r, "POST", "/login", payload, nil)

	decoded := decodeJSON[map[string]string](response)

	return decoded["accessToken"]

}

func getPollById(r *gin.Engine, pollID uint) GetPollResponse {
	response := request(r, "GET", fmt.Sprintf("/poll/%d", pollID), nil, nil)
	return decodeJSON[GetPollResponse](response)
}

func submitVote(r *gin.Engine, optionID uint, accessToken string) {
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)}
	payload := map[string]any{"optionID": optionID}

	request(r, "POST", "/poll/submit", payload, headers)
}

func TestVoteSubmission(t *testing.T) {

	POLL_ID := 8

	godotenv.Load()

	err := db.CreateConnection()

	if err != nil {
		panic(err)
	}

	r := router.SetupRouter()
	accessToken := login(r, "test@mail.com", "test")

	pollData := getPollById(r, uint(POLL_ID))

	votesBefore := pollData.Options[0].Votes
	submitVote(r, pollData.Options[0].ID, accessToken)

	pollData = getPollById(r, uint(POLL_ID))
	votesAfter := pollData.Options[0].Votes

	assert.Equal(t, 1, votesAfter-votesBefore)

}
