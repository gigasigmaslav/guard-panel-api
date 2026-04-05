package vuddecision

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type CreateVUDDecisionUseCase struct {
	transactor contract.RepoTransactor
}

func NewCreateVUDDecisionUseCase(
	transactor contract.RepoTransactor,
) *CreateVUDDecisionUseCase {
	return &CreateVUDDecisionUseCase{
		transactor: transactor,
	}
}

func (cu *CreateVUDDecisionUseCase) Create(ctx context.Context, dec entity.VUDDecision) (int64, error) {
	var out int64
	if err := cu.transactor.InTx(ctx, func(tx contract.TxRepo) error {
		emp, err := tx.GetEmployeeByID(ctx, dec.CreatedByID)
		if err != nil {
			return fmt.Errorf("get employee by id: %w", err)
		}
		dec.CreatedByName = emp.FullName

		id, err := tx.CreateVUDDecision(ctx, dec)
		if err != nil {
			return fmt.Errorf("create vud decision: %w", err)
		}

		_, err = tx.CreateTaskHistoryChange(ctx, entity.TaskHistoryChange{
			TaskID:        dec.TaskID,
			Event:         entity.TaskHistoryChangeEventVUDDecisionAdded,
			CreatedByID:   dec.CreatedByID,
			CreatedByName: dec.CreatedByName,
		})
		if err != nil {
			return fmt.Errorf("create task history change (vud decision added): %w", err)
		}

		out = id
		return nil
	}); err != nil {
		return 0, err
	}
	return out, nil
}
