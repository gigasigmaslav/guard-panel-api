package task

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type ctEmployeeRepo interface {
	GetEmployeeByID(ctx context.Context, id int64) (entity.Employee, error)
}

type ctOfficeRepo interface {
	GetOfficeByID(ctx context.Context, id int64) (entity.Office, error)
}

type CreateTaskUseCase struct {
	employeeRepo ctEmployeeRepo
	officeRepo   ctOfficeRepo
	transactor   contract.RepoTransactor
}

func NewCreateTaskUseCase(
	employeeRepo ctEmployeeRepo,
	officeRepo ctOfficeRepo,
	transactor contract.RepoTransactor,
) *CreateTaskUseCase {
	return &CreateTaskUseCase{
		employeeRepo: employeeRepo,
		officeRepo:   officeRepo,
		transactor:   transactor,
	}
}

func (ct *CreateTaskUseCase) Create(
	ctx context.Context,
	dto CreateTaskDTO,
) (int64, error) {
	creator, err := ct.employeeRepo.GetEmployeeByID(ctx, dto.CreatedByID)
	if err != nil {
		return 0, fmt.Errorf("get created_by employee by id: %w", err)
	}

	executor, err := ct.employeeRepo.GetEmployeeByID(ctx, dto.ExecutorID)
	if err != nil {
		return 0, fmt.Errorf("get executor employee by id: %w", err)
	}

	office, err := ct.officeRepo.GetOfficeByID(ctx, dto.OfficeID)
	if err != nil {
		return 0, fmt.Errorf("get office by id: %w", err)
	}

	task := entity.Task{
		DamageAmount:  dto.DamageAmount,
		Priority:      dto.Priority,
		Status:        entity.TaskStatusNew,
		OfficeID:      dto.OfficeID,
		OfficeName:    office.Name,
		CreatedByID:   dto.CreatedByID,
		CreatedByName: creator.FullName,
		ExecutorID:    dto.ExecutorID,
		ExecutorName:  executor.FullName,
	}

	violator := entity.Violator{
		Type:        dto.ViolatorType,
		FullName:    dto.ViolatorFullName,
		PhoneNumber: dto.ViolatorPhoneNumber,
	}

	var out int64
	if txErr := ct.transactor.InTx(ctx, func(tx contract.TxRepo) error {
		id, createErr := tx.CreateTask(ctx, task)
		if createErr != nil {
			return fmt.Errorf("create task: %w", createErr)
		}

		violator.TaskID = id
		if _, err = tx.CreateViolator(ctx, violator); err != nil {
			return fmt.Errorf("create violator: %w", err)
		}

		out = id
		return nil
	}); txErr != nil {
		return 0, txErr
	}

	return out, nil
}

type CreateTaskDTO struct {
	DamageAmount        int64
	Priority            entity.TaskPriority
	ExecutorID          int64
	OfficeID            int64
	CreatedByID         int64
	ViolatorType        entity.ViolatorType
	ViolatorFullName    string
	ViolatorPhoneNumber *string
}
