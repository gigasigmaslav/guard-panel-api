package office

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
)

type DeleteOfficeUseCase struct {
	transactor contract.RepoTransactor
}

func NewDeleteOfficeUseCase(
	transactor contract.RepoTransactor,
) *DeleteOfficeUseCase {
	return &DeleteOfficeUseCase{
		transactor: transactor,
	}
}

func (do *DeleteOfficeUseCase) Delete(ctx context.Context, id int64) error {
	return do.transactor.InTx(ctx, func(tx contract.TxRepo) error {
		if _, err := tx.GetOfficeByID(ctx, id); err != nil {
			return fmt.Errorf("get office by id: %w", err)
		}
		if err := tx.DeleteOfficeByID(ctx, id); err != nil {
			return fmt.Errorf("delete office by id: %w", err)
		}
		return nil
	})
}
