package models

import "time"

const (
	TypeWin = "Win"
	TypeBet = "Bet"
)

type TransactionModel struct {
	Id            uint64    `json:"id"`
	UserId        uint64    `json:"user_id"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedAt     time.Time `json:"created_at"`
}
