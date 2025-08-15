package service

import "errors"

var (
	ErrUserExists = errors.New("user already exists")
)
