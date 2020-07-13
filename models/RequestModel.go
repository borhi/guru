package models

type GetUserRequestModel struct {
	Id    uint64 `json:"id" validate:"required"`
	Token string `json:"token" validate:"required"`
}

type DepositRequestModel struct {
	UserId    uint64  `json:"user_id" validate:"required"`
	DepositId uint64  `json:"deposit_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,min=0"`
	Token     string  `json:"token" validate:"required"`
}

type TransactionRequestModel struct {
	UserId        uint64  `json:"user_id" validate:"required"`
	TransactionId uint64  `json:"transaction_id" validate:"required"`
	Type          string  `json:"type" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,min=0"`
	Token         string  `json:"token" validate:"required"`
}
