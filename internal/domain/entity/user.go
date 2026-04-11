package entity

import (
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID           int64
	EmployeeID   int64
	PasswordHash string
	CreatedAt    time.Time
	DeletedAt    *time.Time
}
