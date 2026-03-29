package entity

import "time"

type EmployeePosition int32

const (
	EmployeePositionSec EmployeePosition = iota
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
