package vuddecision

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type UpdateVUDDecisionUseCase struct {
	transactor contract.RepoTransactor
}

func NewUpdateVUDDecisionUseCase(
	transactor contract.RepoTransactor,
) *UpdateVUDDecisionUseCase {
	return &UpdateVUDDecisionUseCase{
		transactor: transactor,
	}
}

func (uu *UpdateVUDDecisionUseCase) Update(ctx context.Context, dec entity.VUDDecision) error {
	return uu.transactor.InTx(ctx, func(tx contract.TxRepo) error {
		existingDec, err := tx.GetVUDDecisionByID(ctx, dec.ID)
		if err != nil {
			return fmt.Errorf("get vud decision by id: %w", err)
		}

		newDecision := entity.VUDDecision{
			ID:                 existingDec.ID,
			CriminalCaseOpened: dec.CriminalCaseOpened,
			Comment:            dec.Comment,
			KUSP:               dec.KUSP,
			UD:                 dec.UD,
		}

		if err = tx.UpdateVUDDecisionByID(ctx, newDecision); err != nil {
			return fmt.Errorf("update vud decision: %w", err)
		}
		return nil
	})
}
