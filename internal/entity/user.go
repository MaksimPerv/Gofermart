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
	Number      string    `json:"number"`
	Status      string    `json:"status"`
	Accrual     *float64  `json:"accrual,omitempty"`
	Uploaded_at time.Time `json:"uploaded_at"`
}
