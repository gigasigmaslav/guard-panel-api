package server

import (
	"context"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
	"github.com/gigasigmaslav/guard-panel-api/pkg/api/v1/message"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateEmployee(
	ctx context.Context,
	req *message.CreateEmployeeRequest,
) (*message.CreatedResponse, error) {
	pos := mapEmployeePositionToEntity(req.GetPosition())

	emp := entity.Employee{
		FullName:    req.GetFullName(),
		Position:    pos,
		CreatedByID: req.GetCreatedById(),
	}

	id, err := s.dependencies.CreateEmployeeUseCase.Create(ctx, emp)
	if err != nil {
		return nil, err
	}

	return &message.CreatedResponse{Id: id}, nil
}

func (s *Server) UpdateEmployee(
	ctx context.Context,
	req *message.UpdateEmployeeRequest,
) (*emptypb.Empty, error) {
	emp := entity.Employee{
		ID:       req.GetId(),
		FullName: req.GetFullName(),
	}

	if req.Position != nil {
		emp.Position = mapEmployeePositionToEntity(req.GetPosition())
	}

	if err := s.dependencies.UpdateEmployeeUseCase.Update(ctx, emp); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteEmployee(
	ctx context.Context,
	req *message.DeleteByIDRequest,
) (*emptypb.Empty, error) {
	if err := s.dependencies.DeleteEmployeeUseCase.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) SearchEmployees(
	ctx context.Context,
	_ *emptypb.Empty,
) (*message.SearchEmployeesResponse, error) {
	emps, err := s.dependencies.SearchEmployeeUseCase.Search(ctx)
	if err != nil {
		return nil, err
	}

	out := &message.SearchEmployeesResponse{}
	for i := range emps {
		e := emps[i]
		out.Items = append(out.Items, &message.Employee{
			Id:       e.ID,
			FullName: e.FullName,
			Position: mapEmployeePositionToProto(e.Position),
			CreatedBy: &message.Lookup{
				Id:   e.CreatedByID,
				Name: e.CreatedByName,
			},
			CreatedAt: timestamppb.New(e.CreatedAt),
		})
	}

	return out, nil
}

func mapEmployeePositionToEntity(p message.EmployeePosition) entity.EmployeePosition {
	switch p {
	case message.EmployeePosition_EMPLOYEE_POSITION_UNSPECIFIED:
		return entity.EmployeePositionSec
	case message.EmployeePosition_SEC:
		return entity.EmployeePositionSec
	case message.EmployeePosition_SEC_HEAD:
		return entity.EmployeePositionSecHead
	default:
		return entity.EmployeePositionSec
	}
}

func mapEmployeePositionToProto(p entity.EmployeePosition) message.EmployeePosition {
	switch p {
	case entity.EmployeePositionUnspecified:
		return message.EmployeePosition_EMPLOYEE_POSITION_UNSPECIFIED
	case entity.EmployeePositionSec:
		return message.EmployeePosition_SEC
	case entity.EmployeePositionSecHead:
		return message.EmployeePosition_SEC_HEAD
	default:
		return message.EmployeePosition_EMPLOYEE_POSITION_UNSPECIFIED
	}
}
