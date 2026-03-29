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
	MetadataJSON  *string
	CreatedByID   int64
	CreatedByName string
	CreatedAt     time.Time
}

type TaskFieldName string

const (
	FieldNameKUSPNumber TaskFieldName = "kusp_number"
	FieldNameUDNumber   TaskFieldName = "ud_number"
)

// TaskHistoryChangeUpdatedMetadata represents metadata for Task field update
//
// JSON tags are in lowerCamelTask to fit proto conventions.
type TaskHistoryChangeUpdatedMetadata struct {
	FieldName TaskFieldName `json:"fieldName"`
	OldValue  string        `json:"oldValue"`
	NewValue  string        `json:"newValue"`
}
