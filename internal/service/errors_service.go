package service

import "errors"

//auth

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid login or password")
)

//order

var (
	ErrOrderAlreadyUploaded      = errors.New("order already uploaded by this user")
	ErrOrderBelongsToAnotherUser = errors.New("order already uploaded by another user")
)

var (
	ErrInsufficientPoints = errors.New("insufficient points")
)
