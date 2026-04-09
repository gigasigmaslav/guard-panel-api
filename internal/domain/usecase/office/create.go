package office

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type coOfficeRepo interface {
	GetEmployeeByID(ctx context.Context, id int64) (entity.Employee, error)
	CreateOffice(ctx context.Context, office entity.Office) (int64, error)
}

type CreateOfficeUseCase struct {
	officeRepo coOfficeRepo
}

func NewCreateOfficeUseCase(
	officeRepo coOfficeRepo,
) *CreateOfficeUseCase {
	return &CreateOfficeUseCase{
		officeRepo: officeRepo,
	}
}

func (co *CreateOfficeUseCase) Create(ctx context.Context, office entity.Office) (int64, error) {
	creator, err := co.officeRepo.GetEmployeeByID(ctx, office.CreatedByID)
	if err != nil {
		return 0, fmt.Errorf("get employee by id: %w", err)
	}

	office.CreatedByName = creator.FullName

	id, err := co.officeRepo.CreateOffice(ctx, office)
	if err != nil {
		return 0, fmt.Errorf("create office: %w", err)
	}
	return id, nil
}
