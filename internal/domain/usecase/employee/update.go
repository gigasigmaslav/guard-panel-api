package employee

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type UpdateEmployeeUseCase struct {
	transactor contract.RepoTransactor
}

func NewUpdateEmployeeUseCase(
	transactor contract.RepoTransactor,
) *UpdateEmployeeUseCase {
	return &UpdateEmployeeUseCase{
		transactor: transactor,
	}
}

func (ue *UpdateEmployeeUseCase) Update(ctx context.Context, emp entity.Employee) error {
	return ue.transactor.InTx(ctx, func(tx contract.TxRepo) error {
		existingEmp, getErr := tx.GetEmployeeByID(ctx, emp.ID)
		if getErr != nil {
			return fmt.Errorf("get employee by id: %w", getErr)
		}

		updatedEmp := entity.Employee{
			ID:       existingEmp.ID,
			FullName: emp.FullName,
			Position: emp.Position,
		}

		if updErr := tx.UpdateEmployeeByID(ctx, updatedEmp); updErr != nil {
			return fmt.Errorf("update employee: %w", updErr)
		}

		return nil
	})
}
