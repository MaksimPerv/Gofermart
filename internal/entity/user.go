package entity

import "time"

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
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type UserBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
