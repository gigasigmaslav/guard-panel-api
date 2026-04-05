package entity

import (
	"errors"
	"time"
)

var ErrOfficeNotFound = errors.New("office not found")

type Office struct {
	ID            int64
	Name          string
	Address       string
	CreatedByID   int64
	CreatedByName string
	CreatedAt     time.Time
	DeletedAt     *time.Time
}
