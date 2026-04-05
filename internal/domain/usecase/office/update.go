package office

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type UpdateOfficeUseCase struct {
	transactor contract.RepoTransactor
}

func NewUpdateOfficeUseCase(
	transactor contract.RepoTransactor,
) *UpdateOfficeUseCase {
	return &UpdateOfficeUseCase{
		transactor: transactor,
	}
}

func (uo *UpdateOfficeUseCase) Update(ctx context.Context, office entity.Office) error {
	return uo.transactor.InTx(ctx, func(tx contract.TxRepo) error {
		existingOffice, getErr := tx.GetOfficeByID(ctx, office.ID)
		if getErr != nil {
			return fmt.Errorf("get office by id: %w", getErr)
		}

		newOffice := entity.Office{
			ID:      existingOffice.ID,
			Name:    office.Name,
			Address: office.Address,
		}

		if err := tx.UpdateOfficeByID(ctx, newOffice); err != nil {
			return fmt.Errorf("update office: %w", err)
		}
		return nil
	})
}
