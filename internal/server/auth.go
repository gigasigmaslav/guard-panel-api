package server

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
	"github.com/gigasigmaslav/guard-panel-api/internal/pkg/grpc"
	"github.com/gigasigmaslav/guard-panel-api/pkg/api/v1/message"
)

func (s *Server) SignUp(ctx context.Context, req *message.SignUpRequest) (*message.AuthTokensResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := s.dependencies.SignUp(ctx, req.GetEmployeeId(), req.GetPassword())
	if err != nil {
		return nil, mapAuthError(err)
	}

	return &message.AuthTokensResponse{
		AccessToken: res.AccessToken,
		ExpiresAt:   timestamppb.New(res.ExpiresAt),
	}, nil
}

func (s *Server) SignIn(ctx context.Context, req *message.SignInRequest) (*message.AuthTokensResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := s.dependencies.SignIn(ctx, req.GetEmployeeId(), req.GetPassword())
	if err != nil {
		return nil, mapAuthError(err)
	}

	return &message.AuthTokensResponse{
		AccessToken: res.AccessToken,
		ExpiresAt:   timestamppb.New(res.ExpiresAt),
	}, nil
}

func (s *Server) WhoAmI(ctx context.Context, _ *emptypb.Empty) (*message.WhoAmIResponse, error) {
	employeeID := grpc.GetEmployeeIDFronCtx(ctx)
	if employeeID == 0 {
		return nil, status.Error(codes.Unauthenticated, "employee_id not found in context")
	}

	prof, err := s.dependencies.WhoAmIUseCase.WhoAmI(ctx, employeeID)
	if err != nil {
		return nil, mapAuthError(err)
	}

	return &message.WhoAmIResponse{
		EmployeeId: prof.EmployeeID,
		FullName:   prof.FullName,
		Position:   mapEmployeePositionToProto(prof.Position),
	}, nil
}

func mapAuthError(err error) error {
	switch {
	case errors.Is(err, contract.ErrAlreadyRegistered):
		return status.Error(codes.AlreadyExists, contract.ErrAlreadyRegistered.Error())
	case errors.Is(err, contract.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, contract.ErrInvalidCredentials.Error())
	case errors.Is(err, contract.ErrInvalidToken):
		return status.Error(codes.Unauthenticated, contract.ErrInvalidToken.Error())
	case errors.Is(err, contract.ErrInvalidEmployeePosition):
		return status.Error(codes.FailedPrecondition, contract.ErrInvalidEmployeePosition.Error())
	case errors.Is(err, entity.ErrEmployeeNotFound):
		return status.Error(codes.NotFound, entity.ErrEmployeeNotFound.Error())
	default:
		return err
	}
}
