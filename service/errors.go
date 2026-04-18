package service

import "errors"

var (
	ErrEmailTaken    = errors.New("email already registered")
	ErrInvalidConfig = errors.New("invalid credentials")
)
