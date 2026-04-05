package task

import (
	"context"
	"errors"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type gtTaskRepo interface {
	GetTaskByID(ctx context.Context, id int64) (entity.Task, error)
}

type gtViolatorRepo interface {
	GetViolatorByTaskID(ctx context.Context, id int64) (entity.Violator, error)
}

type gtRefundRepo interface {
	GetRefundsByTaskID(ctx context.Context, id int64) ([]entity.Refund, error)
}

type gtVudDecisionRepo interface {
	GetVUDDecisionByTaskID(ctx context.Context, taskID int64) (entity.VUDDecision, error)
}

type gtCommentRepo interface {
	GetCommentsByTaskID(ctx context.Context, id int64) ([]entity.Comment, error)
}

type gtHistoryChangesRepo interface {
	GetTaskHistoryChangesByTaskID(ctx context.Context, id int64) ([]entity.TaskHistoryChange, error)
}

type GetTaskByIDUseCase struct {
	taskRepo           gtTaskRepo
	violatorRepo       gtViolatorRepo
	refundRepo         gtRefundRepo
	vudDecisionRepo    gtVudDecisionRepo
	commentRepo        gtCommentRepo
	historyChangesRepo gtHistoryChangesRepo
}

func NewGetTaskByIDUseCase(
	taskRepo gtTaskRepo,
	violatorRepo gtViolatorRepo,
	refundRepo gtRefundRepo,
	vudDecisionRepo gtVudDecisionRepo,
	commentRepo gtCommentRepo,
	historyChangesRepo gtHistoryChangesRepo,
) *GetTaskByIDUseCase {
	return &GetTaskByIDUseCase{
		taskRepo:           taskRepo,
		violatorRepo:       violatorRepo,
		refundRepo:         refundRepo,
		vudDecisionRepo:    vudDecisionRepo,
		commentRepo:        commentRepo,
		historyChangesRepo: historyChangesRepo,
	}
}

func (gt *GetTaskByIDUseCase) GetDetails(ctx context.Context, id int64) (contract.TaskDetailsDTO, error) {
	task, err := gt.taskRepo.GetTaskByID(ctx, id)
	if err != nil {
		return contract.TaskDetailsDTO{}, fmt.Errorf("get task by id: %w", err)
	}

	violator, err := gt.violatorRepo.GetViolatorByTaskID(ctx, id)
	if err != nil {
		return contract.TaskDetailsDTO{}, fmt.Errorf("get violator by task id: %w", err)
	}

	refunds, err := gt.refundRepo.GetRefundsByTaskID(ctx, id)
	if err != nil {
		return contract.TaskDetailsDTO{}, fmt.Errorf("get refunds by task id: %w", err)
	}

	decision, decErr := gt.vudDecisionRepo.GetVUDDecisionByTaskID(ctx, id)
	if decErr != nil && !errors.Is(decErr, entity.ErrVUDDecisionNotFound) {
		return contract.TaskDetailsDTO{}, fmt.Errorf("get vud decision by task id: %w", decErr)
	}

	comments, err := gt.commentRepo.GetCommentsByTaskID(ctx, id)
	if err != nil {
		return contract.TaskDetailsDTO{}, fmt.Errorf("get comments by task id: %w", err)
	}

	history, err := gt.historyChangesRepo.GetTaskHistoryChangesByTaskID(ctx, id)
	if err != nil {
		return contract.TaskDetailsDTO{}, fmt.Errorf("get task history changes by task id: %w", err)
	}

	taskItem := contract.TaskLookupItem{
		ID:            task.ID,
		Status:        task.Status,
		DamageAmount:  task.DamageAmount,
		Priority:      task.Priority,
		ExecutorID:    task.ExecutorID,
		ExecutorName:  task.ExecutorName,
		ViolatorID:    violator.ID,
		ViolatorName:  violator.FullName,
		OfficeID:      task.OfficeID,
		OfficeName:    task.OfficeName,
		CreatedByID:   task.CreatedByID,
		CreatedByName: task.CreatedByName,
		CreatedAt:     task.CreatedAt,
		StartDate:     task.StartDate,
		EndDate:       task.EndDate,
	}

	return contract.TaskDetailsDTO{
		Task:          taskItem,
		Refunds:       refunds,
		VUDDecision:   mapVUDDecisionPtr(decision),
		Comments:      comments,
		HistoryEvents: history,
	}, nil
}

func mapVUDDecisionPtr(dec entity.VUDDecision) *entity.VUDDecision {
	if dec.ID == 0 {
		return nil
	}
	return &dec
}
