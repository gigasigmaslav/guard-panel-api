package task

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type sTaskRepo interface {
	SearchTasks(
		ctx context.Context,
		req contract.SearchTasksRequest,
	) ([]contract.TaskLookupItem, int32, error)
}

type sViolatorRepo interface {
	GetViolatorByTaskID(ctx context.Context, id int64) (entity.Violator, error)
}

type SearchTasksUseCase struct {
	taskRepo     sTaskRepo
	violatorRepo sViolatorRepo
}

func NewSearchTasksUseCase(
	taskRepo sTaskRepo,
	violatorRepo sViolatorRepo,
) *SearchTasksUseCase {
	return &SearchTasksUseCase{
		taskRepo:     taskRepo,
		violatorRepo: violatorRepo,
	}
}

func (s *SearchTasksUseCase) Search(
	ctx context.Context,
	req contract.SearchTasksRequest,
) ([]contract.TaskLookupItem, int32, error) {
	items, total, err := s.taskRepo.SearchTasks(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("search tasks: %w", err)
	}

	for i := range items {
		violator, violatorErr := s.violatorRepo.GetViolatorByTaskID(ctx, items[i].ID)
		if violatorErr != nil {
			return nil, 0, fmt.Errorf("get violator by task id: %w", violatorErr)
		}
		items[i].ViolatorID = violator.ID
		items[i].ViolatorName = violator.FullName
	}

	return items, total, nil
}
