package server

import (
	"context"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
	"github.com/gigasigmaslav/guard-panel-api/pkg/api/v1/message"
)

func (s *Server) CreateRefund(
	ctx context.Context,
	req *message.CreateRefundRequest,
) (*message.CreatedResponse, error) {
	ref := entity.Refund{
		TaskID:      req.GetTaskId(),
		Amount:      req.GetAmount(),
		Comment:     req.GetComment(),
		CreatedByID: req.GetCreatedById(),
	}

	id, err := s.dependencies.CreateRefundUseCase.Create(ctx, ref)
	if err != nil {
		return nil, err
	}

	return &message.CreatedResponse{Id: id}, nil
}
