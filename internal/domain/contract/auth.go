package contract

import (
	"errors"
	"time"
)

var (
	ErrAlreadyRegistered       = errors.New("user already registered")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrInvalidEmployeePosition = errors.New("employee position not allowed for sign up")
	ErrInvalidToken            = errors.New("invalid access token")
)

type AuthResult struct {
	AccessToken string
	ExpiresAt   time.Time
}
