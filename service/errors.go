package service

import "errors"

var (
	ErrEmailTaken         = errors.New("Email already registered")
	ErrUsernameTaken      = errors.New("Username already taken")
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrImageUpload        = errors.New("Failed to upload image")
	ErrClassification     = errors.New("Failed to classify batik")
)
