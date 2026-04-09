package employee

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
)

type DeleteEmployeeUseCase struct {
	transactor contract.RepoTransactor
}

func NewDeleteEmployeeUseCase(
	transactor contract.RepoTransactor,
) *DeleteEmployeeUseCase {
	return &DeleteEmployeeUseCase{
		transactor: transactor,
	}
}

func (de *DeleteEmployeeUseCase) Delete(ctx context.Context, id int64) error {
	return de.transactor.InTx(ctx, func(tx contract.TxRepo) error {
		_, getErr := tx.GetEmployeeByID(ctx, id)
		if getErr != nil {
			return fmt.Errorf("get employee by id: %w", getErr)
		}

		if delErr := tx.DeleteEmployeeByID(ctx, id); delErr != nil {
			return fmt.Errorf("delete employee by id: %w", delErr)
		}

		return nil
	})
}
