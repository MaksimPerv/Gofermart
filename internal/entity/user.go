package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserWithId struct {
	Id       int
	Login    string
	Password string
}

type DBUser struct {
	Number     string           `json:"number"`
	Status     string           `json:"status"`
	Accrual    *decimal.Decimal `json:"accrual,omitempty"`
	UploadedAt time.Time        `json:"uploaded_at"`
}

type UserBalance struct {
	Current   decimal.Decimal `json:"current"`
	Withdrawn decimal.Decimal `json:"withdrawn"`
}

type WithdrawRequest struct {
	Order string          `json:"order"`
	Sum   decimal.Decimal `json:"sum"`
}
