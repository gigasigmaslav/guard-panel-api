package refund

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type CreateRefundUseCase struct {
	transactor contract.RepoTransactor
}

func NewCreateRefundUseCase(
	transactor contract.RepoTransactor,
) *CreateRefundUseCase {
	return &CreateRefundUseCase{
		transactor: transactor,
	}
}

func (cu *CreateRefundUseCase) Create(ctx context.Context, ref entity.Refund) (int64, error) {
	var out int64
	if err := cu.transactor.InTx(ctx, func(tx contract.TxRepo) error {
		emp, err := tx.GetEmployeeByID(ctx, ref.CreatedByID)
		if err != nil {
			return fmt.Errorf("get employee by id: %w", err)
		}
		ref.CreatedByName = emp.FullName

		id, err := tx.CreateRefund(ctx, ref)
		if err != nil {
			return fmt.Errorf("create refund: %w", err)
		}

		_, err = tx.CreateTaskHistoryChange(ctx, entity.TaskHistoryChange{
			TaskID:        ref.TaskID,
			Event:         entity.TaskHistoryChangeEventRefundAdded,
			CreatedByID:   ref.CreatedByID,
			CreatedByName: ref.CreatedByName,
		})
		if err != nil {
			return fmt.Errorf("create task history change (refund added): %w", err)
		}

		out = id
		return nil
	}); err != nil {
		return 0, err
	}
	return out, nil
}
