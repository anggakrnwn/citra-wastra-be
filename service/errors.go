package service

import "errors"

var (
	ErrEmailTaken    = errors.New("Email already registered")
	ErrInvalidConfig = errors.New("Invalid credentials")
)
