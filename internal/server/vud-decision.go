package server

import (
	"context"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
	"github.com/gigasigmaslav/guard-panel-api/pkg/api/v1/message"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateVudDecision(
	ctx context.Context,
	req *message.CreateVudDecisionRequest,
) (*message.CreatedResponse, error) {
	var (
		udPtr   *string
		ccPtr   *bool
		comment *string
	)

	if u := req.GetUd(); u != "" {
		udPtr = &u
	}

	if req.CriminalCaseOpened != nil {
		v := req.GetCriminalCaseOpened()
		ccPtr = &v
	}

	if c := req.GetComment(); c != "" {
		comment = &c
	}

	dec := entity.VUDDecision{
		TaskID:             req.GetTaskId(),
		CriminalCaseOpened: ccPtr,
		Comment:            comment,
		KUSP:               req.GetKusp(),
		UD:                 udPtr,
		CreatedByID:        req.GetCreatedById(),
	}

	id, err := s.dependencies.CreateVUDDecisionUseCase.Create(ctx, dec)
	if err != nil {
		return nil, err
	}

	return &message.CreatedResponse{Id: id}, nil
}

func (s *Server) UpdateVudDecision(
	ctx context.Context,
	req *message.UpdateVudDecisionRequest,
) (*emptypb.Empty, error) {
	var (
		udPtr   *string
		ccPtr   *bool
		comment *string
		kusp    string
	)

	if req.Ud != nil {
		u := req.GetUd()
		udPtr = &u
	}

	if req.CriminalCaseOpened != nil {
		v := req.GetCriminalCaseOpened()
		ccPtr = &v
	}

	if req.Comment != nil {
		c := req.GetComment()
		comment = &c
	}

	if req.Kusp != nil {
		kusp = req.GetKusp()
	}

	dec := entity.VUDDecision{
		ID:                 req.GetId(),
		CriminalCaseOpened: ccPtr,
		Comment:            comment,
		KUSP:               kusp,
		UD:                 udPtr,
	}

	if err := s.dependencies.UpdateVUDDecisionUseCase.Update(ctx, dec); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
