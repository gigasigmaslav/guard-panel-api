package entity

import "time"

type Office struct {
	ID            int64
	Name          string
	Address       string
	CreatedByID   int64
	CreatedByName string
	CreatedAt     time.Time
	DeletedAt     *time.Time
}
