package entity

import (
	"fmt"
	"strings"
	"time"
)

type TaskPriority int32

const (
	TaskPriorityUnspecified TaskPriority = iota
	TaskPriorityLow
	TaskPriorityHigh
)

type TaskStatus int32

const (
	TaskStatusUnspecified TaskStatus = iota
	// TaskStatusNew - новое задача.
	TaskStatusNew
	// TaskStatusInProgress - в работе.
	TaskStatusInProgress
	// TaskStatusPendingVUD - ожидание ВУД.
	TaskStatusPendingVUD
	// TaskStatusInCourt - судебная стадия.
	TaskStatusInCourt
	// TaskStatusCompleted - завершено.
	TaskStatusCompleted
)

func MapTaskStatus(s string) (TaskStatus, error) {
	switch strings.ToLower(s) {
	case "unspecified":
		return TaskStatusUnspecified, nil
	case "new":
		return TaskStatusNew, nil
	case "in_progress":
		return TaskStatusInProgress, nil
	case "pending_vud":
		return TaskStatusPendingVUD, nil
	case "in_court":
		return TaskStatusInCourt, nil
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
