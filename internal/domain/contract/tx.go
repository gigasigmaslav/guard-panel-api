package contract

import (
	"context"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type RepoTransactor interface {
	InTx(ctx context.Context, f func(tx TxRepo) error) error
}

type TxRepo interface {
	// employee
	GetEmployeeByID(ctx context.Context, id int64) (entity.Employee, error)
	UpdateEmployeeByID(ctx context.Context, emp entity.Employee) error
	DeleteEmployeeByID(ctx context.Context, id int64) error

	// office
	GetOfficeByID(ctx context.Context, id int64) (entity.Office, error)
	UpdateOfficeByID(ctx context.Context, office entity.Office) error
	DeleteOfficeByID(ctx context.Context, id int64) error

	// task
	CreateTask(ctx context.Context, task entity.Task) (int64, error)
	GetTaskByID(ctx context.Context, id int64) (entity.Task, error)
	UpdateTaskByID(ctx context.Context, task entity.Task) error

	// violator
	CreateViolator(ctx context.Context, violator entity.Violator) (int64, error)

	// vud decision
	GetVUDDecisionByID(ctx context.Context, id int64) (entity.VUDDecision, error)
	CreateVUDDecision(ctx context.Context, decision entity.VUDDecision) (int64, error)
	UpdateVUDDecisionByID(ctx context.Context, decision entity.VUDDecision) error

	// refund
	CreateRefund(ctx context.Context, refund entity.Refund) (int64, error)

	// comment
	CreateComment(ctx context.Context, com entity.Comment) (int64, error)

	// history
	CreateTaskHistoryChange(ctx context.Context, change entity.TaskHistoryChange) (int64, error)

	// user
	GetUserByEmployeeID(ctx context.Context, employeeID int64) (entity.User, error)
	UserExistsByEmployeeID(ctx context.Context, employeeID int64) (bool, error)
	CreateUser(ctx context.Context, u entity.User) (int64, error)
}
