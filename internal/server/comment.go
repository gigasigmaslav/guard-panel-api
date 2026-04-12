package server

import (
	"context"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
	grpcctx "github.com/gigasigmaslav/guard-panel-api/internal/pkg/grpc"
	"github.com/gigasigmaslav/guard-panel-api/pkg/api/v1/message"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateComment(
	ctx context.Context,
	req *message.CreateCommentRequest,
) (*message.CreatedResponse, error) {
	comment := entity.Comment{
		TaskID:      req.GetTaskId(),
		Comment:     req.GetComment(),
		CreatedByID: grpcctx.GetEmployeeIDFromCtx(ctx),
	}

	id, err := s.dependencies.CreateCommentUseCase.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &message.CreatedResponse{Id: id}, nil
}

func (s *Server) DeleteComment(
	ctx context.Context,
	req *message.DeleteByIDRequest,
) (*emptypb.Empty, error) {
	if err := s.dependencies.DeleteCommentUseCase.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
