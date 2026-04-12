package server

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/pkg/api/v1/message"
)

func (s *Server) GetTaskDashboardKPI(
	ctx context.Context,
	req *message.GetTaskDashboardKPIRequest,
) (*message.GetTaskDashboardKPIResponse, error) {
	from, to := time.Time{}, time.Time{}
	if ts := req.GetPeriodFrom(); ts != nil {
		from = ts.AsTime()
	}
	if ts := req.GetPeriodTo(); ts != nil {
		to = ts.AsTime()
	}

	kpi, top, err := s.dependencies.GetAnalyticsUseCase.Get(ctx, from, to)
	if err != nil {
		if errors.Is(err, contract.ErrInvalidAnalyticsPeriod) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &message.GetTaskDashboardKPIResponse{
		ActiveTasksCount:        kpi.ActiveTasksCount,
		CompletedInPeriod:       kpi.CompletedInPeriod,
		CreatedInPeriod:         kpi.CreatedInPeriod,
		CompletedToCreatedRatio: kpi.CompletedToCreatedRatio,
	}

	resp.TopExecutor = &message.GetTaskDashboardKPIResponse_TopExecutorByCompleted{
		ExecutorId:     top.ExecutorID,
		ExecutorName:   top.ExecutorName,
		CompletedTasks: top.Completed,
	}

	return resp, nil
}
