package employee

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type seEmployeeRepo interface {
	GetEmployees(ctx context.Context) ([]entity.Employee, error)
}

type SearchEmployeeUseCase struct {
	employeeRepo seEmployeeRepo
}

func NewSearchEmployeeUseCase(
	employeeRepo seEmployeeRepo,
) *SearchEmployeeUseCase {
	return &SearchEmployeeUseCase{
		employeeRepo: employeeRepo,
	}
}

func (se *SearchEmployeeUseCase) Search(ctx context.Context) ([]entity.Employee, error) {
	emps, err := se.employeeRepo.GetEmployees(ctx)
	if err != nil {
		return nil, fmt.Errorf("search employees: %w", err)
	}
	return emps, nil
}
