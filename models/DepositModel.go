package models

import "time"

type DepositModel struct {
	Id            uint64    `json:"id"`
	UserId        uint64    `json:"user_id"`
	Amount        float64   `json:"amount"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedAt     time.Time `json:"created_at"`
}
