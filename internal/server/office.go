package server

import (
	"context"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
	grpcctx "github.com/gigasigmaslav/guard-panel-api/internal/pkg/grpc"
	"github.com/gigasigmaslav/guard-panel-api/pkg/api/v1/message"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateOffice(
	ctx context.Context,
	req *message.CreateOfficeRequest,
) (*message.CreatedResponse, error) {
	office := entity.Office{
		Name:        req.GetName(),
		Address:     req.GetAddress(),
		CreatedByID: grpcctx.GetEmployeeIDFromCtx(ctx),
	}

	id, err := s.dependencies.CreateOfficeUseCase.Create(ctx, office)
	if err != nil {
		return nil, err
	}

	return &message.CreatedResponse{Id: id}, nil
}

func (s *Server) UpdateOffice(
	ctx context.Context,
	req *message.UpdateOfficeRequest,
) (*emptypb.Empty, error) {
	office := entity.Office{
		ID:      req.GetId(),
		Name:    req.GetName(),
		Address: req.GetAddress(),
	}

	if err := s.dependencies.UpdateOfficeUseCase.Update(ctx, office); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteOffice(
	ctx context.Context,
	req *message.DeleteByIDRequest,
) (*emptypb.Empty, error) {
	if err := s.dependencies.DeleteOfficeUseCase.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) SearchOffices(
	ctx context.Context,
	_ *emptypb.Empty,
) (*message.SearchOfficesResponse, error) {
	offices, err := s.dependencies.SearchOfficeUseCase.Search(ctx)
	if err != nil {
		return nil, err
	}

	resp := &message.SearchOfficesResponse{}
	for i := range offices {
		o := offices[i]
		resp.Items = append(resp.Items, &message.Office{
			Id:      o.ID,
			Name:    o.Name,
			Address: o.Address,
			CreatedBy: &message.Lookup{
				Id:   o.CreatedByID,
				Name: o.CreatedByName,
			},
			CreatedAt: timestamppb.New(o.CreatedAt),
		})
	}

	return resp, nil
}
