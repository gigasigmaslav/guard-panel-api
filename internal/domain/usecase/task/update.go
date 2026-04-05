package task

import (
	"context"
	"fmt"

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

func (uu *UpdateTaskUseCase) Update(ctx context.Context, task entity.Task) error {
	return uu.transactor.InTx(ctx, func(tx contract.TxRepo) error {
		existingTask, err := tx.GetTaskByID(ctx, task.ID)
		if err != nil {
			return fmt.Errorf("get task by id: %w", err)
		}

		newTask := entity.Task{
			ID:           existingTask.ID,
			DamageAmount: task.DamageAmount,
			Priority:     task.Priority,
			Status:       task.Status,
			EndDate:      task.EndDate,
			ExecutorID:   task.ExecutorID,
			ExecutorName: task.ExecutorName,
		}

		if err = tx.UpdateTaskByID(ctx, newTask); err != nil {
			return fmt.Errorf("update task: %w", err)
		}
		return nil
	})
}
