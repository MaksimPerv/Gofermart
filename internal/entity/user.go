package entity

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserWithId struct {
	Id       int
	Login    string
	Password string
}
