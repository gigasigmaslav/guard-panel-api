package task

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type ctTaskRepo interface {
	GetEmployeeByID(ctx context.Context, id int64) (entity.Employee, error)
	CreateTask(ctx context.Context, task entity.Task) (int64, error)
}

type CreateTaskUseCase struct {
	taskRepo ctTaskRepo
}

func NewCreateTaskUseCase(
	taskRepo ctTaskRepo,
) *CreateTaskUseCase {
	return &CreateTaskUseCase{
		taskRepo: taskRepo,
	}
}

func (cu *CreateTaskUseCase) Create(ctx context.Context, t entity.Task) (int64, error) {
	creator, err := cu.taskRepo.GetEmployeeByID(ctx, t.CreatedByID)
	if err != nil {
		return 0, fmt.Errorf("get created_by employee by id: %w", err)
	}
	executor, err := cu.taskRepo.GetEmployeeByID(ctx, t.ExecutorID)
	if err != nil {
		return 0, fmt.Errorf("get executor employee by id: %w", err)
	}

	t.CreatedByName = creator.FullName
	t.ExecutorName = executor.FullName

	id, err := cu.taskRepo.CreateTask(ctx, t)
	if err != nil {
		return 0, fmt.Errorf("create task: %w", err)
	}
	return id, nil
}
