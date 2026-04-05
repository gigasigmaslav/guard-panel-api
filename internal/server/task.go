package server

import (
	"context"
	"time"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
	"github.com/gigasigmaslav/guard-panel-api/pkg/api/v1/message"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateTask(
	ctx context.Context,
	req *message.CreateTaskRequest,
) (*message.CreatedResponse, error) {
	priority := mapTaskPriorityFromProto(req.GetPriority())

	task := entity.Task{
		DamageAmount: req.GetDamageAmount(),
		Priority:     priority,
		ExecutorID:   req.GetExecutorId(),
		CreatedByID:  req.GetCreatedById(),
	}

	id, err := s.dependencies.CreateTaskUseCase.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	return &message.CreatedResponse{Id: id}, nil
}

func (s *Server) GetTaskDetails(
	ctx context.Context,
	req *message.GetByIDRequest,
) (*message.GetTaskDetailsResponse, error) {
	dto, err := s.dependencies.GetTaskByIDUseCase.GetDetails(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	task := dto.Task
	resp := &message.GetTaskDetailsResponse{
		Task: &message.Task{
			Id:           task.ID,
			DamageAmount: task.DamageAmount,
			Priority:     mapTaskPriorityToProto(task.Priority),
			Status:       mapTaskStatusToProto(task.Status),
			Executor: &message.Lookup{
				Id:   task.ExecutorID,
				Name: task.ExecutorName,
			},
			Violator: &message.Lookup{
				Id:   task.ViolatorID,
				Name: task.ViolatorName,
			},
			Office: &message.Lookup{
				Id:   task.OfficeID,
				Name: task.OfficeName,
			},
			CreatedBy: &message.Lookup{
				Id:   task.CreatedByID,
				Name: task.CreatedByName,
			},
			CreatedAt: timestamppb.New(task.CreatedAt),
			StartDate: timestamppb.New(task.StartDate),
		},
	}
	if task.EndDate != nil {
		resp.Task.EndDate = timestamppb.New(*task.EndDate)
	}

	if dto.VUDDecision != nil {
		dec := dto.VUDDecision
		resp.Task.VudDecisions = &message.Task_TaskVudDecision{
			Id:   dec.ID,
			Kusp: dec.KUSP,
		}
		if dec.UD != nil {
			resp.Task.VudDecisions.Ud = dec.UD
		}
		if dec.CriminalCaseOpened != nil {
			resp.Task.VudDecisions.CriminalCaseOpened = dec.CriminalCaseOpened
		}
		if dec.Comment != nil {
			resp.Task.VudDecisions.Comment = dec.Comment
		}
	}

	for i := range dto.Refunds {
		r := dto.Refunds[i]
		resp.Task.Refunds = append(resp.Task.Refunds, &message.Task_TaskRefund{
			Id:      r.ID,
			Amount:  r.Amount,
			Comment: r.Comment,
		})
	}

	for i := range dto.Comments {
		c := dto.Comments[i]
		resp.Task.Comments = append(resp.Task.Comments, &message.Task_TaskComment{
			Id:      c.ID,
			Comment: c.Comment,
			CreatedBy: &message.Lookup{
				Id:   c.CreatedByID,
				Name: c.CreatedByName,
			},
			CreatedAt: timestamppb.New(c.CreatedAt),
		})
	}

	for i := range dto.HistoryEvents {
		h := dto.HistoryEvents[i]
		resp.Task.HistoryChanges = append(resp.Task.HistoryChanges, &message.Task_TaskHistoryChange{
			Event: mapHistoryEventToProto(h.Event),
			CreatedBy: &message.Lookup{
				Id:   h.CreatedByID,
				Name: h.CreatedByName,
			},
			CreatedAt: timestamppb.New(h.CreatedAt),
		})
	}

	return resp, nil
}

func (s *Server) UpdateTask(
	ctx context.Context,
	req *message.UpdateTaskRequest,
) (*emptypb.Empty, error) {
	var (
		priority entity.TaskPriority
		status   entity.TaskStatus
	)

	if req.Priority != nil {
		priority = mapTaskPriorityFromProto(req.GetPriority())
	}

	if req.Status != nil {
		status = contract.MapTaskStatusToContract(req.GetStatus())
	}

	var endDate *time.Time
	if req.EndDate != nil {
		t := req.GetEndDate().AsTime()
		endDate = &t
	}

	task := entity.Task{
		ID:           req.GetId(),
		DamageAmount: req.GetDamageAmount(),
		Priority:     priority,
		Status:       status,
		EndDate:      endDate,
		ExecutorID:   req.GetExecutorId(),
	}

	if err := s.dependencies.UpdateTaskUseCase.Update(ctx, task); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) SearchTasks(
	ctx context.Context,
	req *message.SearchTasksRequest,
) (*message.SearchTasksResponse, error) {
	searchReq := contract.SearchTasksRequest{
		Status:  contract.MapTaskStatusToContract(req.GetStatus()),
		Page:    req.GetPage(),
		PerPage: req.GetPerPage(),
		Filter:  contract.MapToSearchTasksFilter(req.GetFilter()),
		Sorting: contract.SearchTasksSorting{
			OrderDirection: contract.MapOrderDirectionToContract(req.GetSorting().GetOrderDirection()),
			OrderBy:        contract.MapOrderByToContract(req.GetSorting().GetOrderBy()),
		},
	}

	items, total, err := s.dependencies.SearchTasksUseCase.Search(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	resp := &message.SearchTasksResponse{
		Total: total,
	}

	for i := range items {
		it := items[i]
		msg := &message.SearchTasksResponse_TaskLookup{
			Id:           it.ID,
			DamageAmount: it.DamageAmount,
			Priority:     mapTaskPriorityToProto(it.Priority),
			Executor: &message.Lookup{
				Id:   it.ExecutorID,
				Name: it.ExecutorName,
			},
			Violator: &message.Lookup{
				Id:   it.ViolatorID,
				Name: it.ViolatorName,
			},
			Office: &message.Lookup{
				Id:   it.OfficeID,
				Name: it.OfficeName,
			},
			CreatedBy: &message.Lookup{
				Id:   it.CreatedByID,
				Name: it.CreatedByName,
			},
			CreatedAt: timestamppb.New(it.CreatedAt),
			StartDate: timestamppb.New(it.StartDate),
		}
		if it.EndDate != nil {
			msg.EndDate = timestamppb.New(*it.EndDate)
		}
		resp.Items = append(resp.Items, msg)
	}

	return resp, nil
}

func mapTaskStatusToProto(status entity.TaskStatus) message.TaskStatus {
	switch status {
	case entity.TaskStatusUnspecified:
		return message.TaskStatus_CASE_STATUS_UNSPECIFIED
	case entity.TaskStatusNew:
		return message.TaskStatus_NEW
	case entity.TaskStatusDocumentGathering:
		return message.TaskStatus_DOCUMENT_GATHERING
	case entity.TaskStatusWaitingForVUD:
		return message.TaskStatus_WAITING_FOR_VUD
	case entity.TaskStatusJudicialProceedings:
		return message.TaskStatus_JUDICIAL_PROCEEDINGS
	case entity.TaskStatusCompleted:
		return message.TaskStatus_COMPLETED
	default:
		return message.TaskStatus_CASE_STATUS_UNSPECIFIED
	}
}

func mapHistoryEventToProto(e entity.TaskHistoryChangeEvent) message.Task_TaskHistoryChange_TaskHistoryChangeEvent {
	switch e {
	case entity.TaskHistoryChangeEventCommentAdded:
		return message.Task_TaskHistoryChange_TASK_HISTORY_CHANGE_EVENT_COMMENT_ADDED
	case entity.TaskHistoryChangeEventRefundAdded:
		return message.Task_TaskHistoryChange_TASK_HISTORY_CHANGE_EVENT_REFUND_ADDED
	case entity.TaskHistoryChangeEventVUDDecisionAdded:
		return message.Task_TaskHistoryChange_TASK_HISTORY_CHANGE_EVENT_VUD_DECISION_ADDED
	default:
		return message.Task_TaskHistoryChange_TASK_HISTORY_CHANGE_EVENT_UNSPECIFIED
	}
}

func mapTaskPriorityToProto(priority entity.TaskPriority) message.TaskPriority {
	switch priority {
	case entity.TaskPriorityLow:
		return message.TaskPriority_LOW
	case entity.TaskPriorityHigh:
		return message.TaskPriority_HIGH
	default:
		return message.TaskPriority_TASK_PRIORITY_UNSPECIFIED
	}
}

func mapTaskPriorityFromProto(priority message.TaskPriority) entity.TaskPriority {
	switch priority {
	case message.TaskPriority_LOW:
		return entity.TaskPriorityLow
	case message.TaskPriority_HIGH:
		return entity.TaskPriorityHigh
	default:
		return entity.TaskPriorityLow
	}
}
