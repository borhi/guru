package models

type StatisticModel struct {
	Id           uint64  `bson:"_id"`
	DepositCount int     `bson:"deposit_count"`
	DepositSum   float64 `bson:"deposit_sum"`
	BetCount     int     `bson:"bet_count"`
	BetSum       float64 `bson:"bet_sum"`
	WinCount     int     `bson:"win_count"`
	WinSum       float64 `bson:"win_sum"`
}
