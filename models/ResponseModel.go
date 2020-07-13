package models

type ErrorResponseModel struct {
	Error string `json:"error"`
}

type GetUserResponseModel struct {
	Id           uint64  `json:"id"`
	Balance      float64 `json:"balance"`
	DepositCount int     `json:"deposit_count"`
	DepositSum   float64 `json:"deposit_sum"`
	BetCount     int     `json:"bet_count"`
	BetSum       float64 `json:"bet_sum"`
	WinCount     int     `json:"win_count"`
	WinSum       float64 `json:"win_sum"`
}

type TransactionResponseModel struct {
	Error   string  `json:"error"`
	Balance float64 `json:"balance"`
}
