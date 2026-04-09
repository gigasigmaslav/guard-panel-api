package entity

import "time"

type Comment struct {
	ID            int64
	TaskID        int64
	Comment       string
	CreatedByID   int64
	CreatedByName string
	CreatedAt     time.Time
}
