package auth

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type waiEmployeeRepo interface {
	GetEmployeeByID(ctx context.Context, id int64) (entity.Employee, error)
}

type WhoAmIUseCase struct {
	employeeRepo waiEmployeeRepo
}

func NewWhoAmIUseCase(
	employeeRepo waiEmployeeRepo,
) *WhoAmIUseCase {
	return &WhoAmIUseCase{
		employeeRepo: employeeRepo,
	}
}

func (wai *WhoAmIUseCase) WhoAmI(ctx context.Context, employeeID int64) (ProfileDTO, error) {
	emp, err := wai.employeeRepo.GetEmployeeByID(ctx, employeeID)
	if err != nil {
		return ProfileDTO{}, fmt.Errorf("whoami get employee: %w", err)
	}

	return ProfileDTO{
		EmployeeID: emp.ID,
		FullName:   emp.FullName,
		Position:   emp.Position,
	}, nil
}

type ProfileDTO struct {
	EmployeeID int64
	FullName   string
	Position   entity.EmployeePosition
}
