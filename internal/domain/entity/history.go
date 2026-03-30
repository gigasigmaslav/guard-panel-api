package entity

import "time"

type TaskHistoryChangeEvent int32

const (
	TaskHistoryChangeEventCommentAdded = iota
	TaskHistoryChangeEventRefundAdded
	TaskHistoryChangeEventVUDDecisionAdded
)

type TaskHistoryChange struct {
	ID            int64
	TaskID        int64
	Event         TaskHistoryChangeEvent
	CreatedByID   int64
	CreatedByName string
	CreatedAt     time.Time
}
