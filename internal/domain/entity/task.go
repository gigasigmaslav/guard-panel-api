package entity

import (
	"fmt"
	"strings"
	"time"
)

type TaskPriority int32

const (
	TaskPriorityLow TaskPriority = iota
	TaskPriorityHigh
)

type TaskStatus int32

const (
	TaskStatusUnspecified TaskStatus = iota
	// TaskStatusNew - новое задача.
	TaskStatusNew
	// TaskStatusDocumentGathering - сбор документов.
	TaskStatusDocumentGathering
	// TaskStatusWaitingForVUD - ожидание ВУД.
	TaskStatusWaitingForVUD
	// TaskStatusJudicialProceedings - судебное разбирательство.
	TaskStatusJudicialProceedings
	// TaskStatusCompleted - завершено.
	TaskStatusCompleted
)

func MapTaskStatus(s string) (TaskStatus, error) {
	switch strings.ToLower(s) {
	case "unspecified":
		return TaskStatusUnspecified, nil
	case "new":
		return TaskStatusNew, nil
	case "document_gathering":
		return TaskStatusDocumentGathering, nil
	case "waiting_for_vud":
		return TaskStatusWaitingForVUD, nil
	case "judicial_proceedings":
		return TaskStatusJudicialProceedings, nil
	case "completed":
		return TaskStatusCompleted, nil
	default:
		return TaskStatusUnspecified, fmt.Errorf("invalid TaskStatus: %s", s)
	}
}

type Task struct {
	ID            int64
	DamageAmount  int64
	BarcodesCount int32
	Priority      TaskPriority
	Status        TaskStatus
	StartDate     time.Time
	EndDate       *time.Time
	ExecutorID    int64
	ExecutorName  string
	OfficeID      int64
	OfficeName    string
	CreatedByID   int64
	CreatedByName string
	CreatedAt     time.Time
}
