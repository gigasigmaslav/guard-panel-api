package entity

import (
	"errors"
	"time"
)

var ErrVUDDecisionNotFound = errors.New("vud decision not found")

type VUDDecision struct {
	ID                 int64
	TaskID             int64
	CriminalCaseOpened *bool
	Comment            *string
	KUSP               string
	UD                 *string
	CreatedByID        int64
	CreatedByName      string
	CreatedAt          time.Time
}
