package employee

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type ceEmployeeRepo interface {
	GetEmployeeByID(ctx context.Context, id int64) (entity.Employee, error)
	CreateEmployee(ctx context.Context, employee entity.Employee) (int64, error)
}

type CreateEmployeeUseCase struct {
	employeeRepo ceEmployeeRepo
}

func NewCreateEmployeeUseCase(
	employeeRepo ceEmployeeRepo,
) *CreateEmployeeUseCase {
	return &CreateEmployeeUseCase{
		employeeRepo: employeeRepo,
	}
}

func (ce *CreateEmployeeUseCase) Create(ctx context.Context, emp entity.Employee) (int64, error) {
	creator, err := ce.employeeRepo.GetEmployeeByID(ctx, emp.CreatedByID)
	if err != nil {
		return 0, fmt.Errorf("get employee by id: %w", err)
	}

	emp.CreatedByName = creator.FullName

	id, err := ce.employeeRepo.CreateEmployee(ctx, emp)
	if err != nil {
		return 0, fmt.Errorf("create employee: %w", err)
	}
	return id, nil
}
