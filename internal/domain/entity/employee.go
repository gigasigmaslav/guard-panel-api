package entity

import (
	"errors"
	"time"
)

var ErrEmployeeNotFound = errors.New("employee not found")

type EmployeePosition int32

const (
	EmployeePositionUnspecified EmployeePosition = iota
	EmployeePositionSec
	EmployeePositionSecHead
)

type Employee struct {
	ID            int64
	FullName      string
	Position      EmployeePosition
	CreatedByID   int64
	CreatedByName string
	CreatedAt     time.Time
	DeletedAt     *time.Time
}
