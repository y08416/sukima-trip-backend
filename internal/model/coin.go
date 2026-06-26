package model

type CoinResponse struct {
	Balance int `json:"balance"`
}

type CoinTodayResponse struct {
	EarnedToday int `json:"earned_today"`
}
