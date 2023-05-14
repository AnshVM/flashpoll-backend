package types

type GetPollResponse struct {
	Title      string           `json:"title"`
	Options    []OptionResponse `json:"options"`
	TotalVotes uint             `json:"totalVotes"`
	UserVote   OptionResponse   `json:"userVote"`
	ID         uint             `json:"id"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type OptionResponse struct {
	Name         string  `json:"name"`
	Votes        uint    `json:"votes"`
	ID           uint    `json:"id"`
	VotesPercent float32 `json:"votesPercent"`
}
type UpdatePollResponse struct {
	Options    []OptionResponse `json:"options"`
	TotalVotes uint             `json:"totalVotes"`
	ID         uint             `json:"id"`
}
