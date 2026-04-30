package dto

type LeaderboardResponse struct {
	Rank     int     `json:"rank"`
	Username string  `json:"username"`
	XP       float64 `json:"xp"`
}
