package entity

type ViolatorType int32

const (
	ViolatorTypeUnspecified ViolatorType = iota
	ViolatorTypeEmployee
	ViolatorTypeClient
)

type Violator struct {
	ID          int64
	TaskID      int64
	Type        ViolatorType
	FullName    string
	PhoneNumber *string
}
