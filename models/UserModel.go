package models

const (
	StatusNew      = "New"
	StatusModified = "Modified"
)

type UserModel struct {
	Id      uint64  `json:"id" validate:"required"`
	Balance float64 `json:"balance" validate:"min=0"`
	Token   string  `json:"token" validate:"required"`
	Status  string
}
