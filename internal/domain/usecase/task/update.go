package task

import (
	"context"
	"fmt"
	"time"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type UpdateTaskUseCase struct {
	transactor contract.RepoTransactor
}

func NewUpdateTaskUseCase(
	transactor contract.RepoTransactor,
) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{
		transactor: transactor,
	}
}

func (ut *UpdateTaskUseCase) Update(ctx context.Context, dto UpdateTaskDTO) error {
	return ut.transactor.InTx(ctx, func(tx contract.TxRepo) error {
		existingTask, err := tx.GetTaskByID(ctx, dto.ID)
		if err != nil {
			return fmt.Errorf("get task by id: %w", err)
		}

		newTask, err := ut.mergeTaskUpdates(ctx, tx, existingTask, dto)
		if err != nil {
			return err
		}

		if err = tx.UpdateTaskByID(ctx, newTask); err != nil {
			return fmt.Errorf("update task: %w", err)
		}
		return nil
	})
}

func (ut *UpdateTaskUseCase) mergeTaskUpdates(
	ctx context.Context,
	tx contract.TxRepo,
	existingTask entity.Task,
	dto UpdateTaskDTO,
) (entity.Task, error) {
	newTask := existingTask

	if dto.DamageAmount != nil {
		newTask.DamageAmount = *dto.DamageAmount
	}
	if dto.Priority != nil {
		newTask.Priority = *dto.Priority
	}
	if dto.Status != nil {
		newTask.Status = *dto.Status
	}
	if dto.EndDate != nil {
		newTask.EndDate = dto.EndDate
	}
	if dto.ExecutorID != nil {
		newTask.ExecutorID = *dto.ExecutorID

		executor, err := tx.GetEmployeeByID(ctx, *dto.ExecutorID)
		if err != nil {
			return entity.Task{}, fmt.Errorf("get executor employee by id: %w", err)
		}

		newTask.ExecutorName = executor.FullName
	}

	return newTask, nil
}

type UpdateTaskDTO struct {
	ID           int64
	DamageAmount *int64
	Priority     *entity.TaskPriority
	Status       *entity.TaskStatus
	EndDate      *time.Time
	ExecutorID   *int64
}
