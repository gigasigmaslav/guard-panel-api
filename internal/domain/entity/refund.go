package entity

import "time"

type Refund struct {
	ID            int64
	TaskID        int64
	Amount        int64
	Comment       string
	CreatedByID   int64
	CreatedByName string
	CreatedAt     time.Time
}
