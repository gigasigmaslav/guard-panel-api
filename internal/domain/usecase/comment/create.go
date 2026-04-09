package comment

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type CreateCommentUseCase struct {
	transactor contract.RepoTransactor
}

func NewCreateCommentUseCase(
	transactor contract.RepoTransactor,
) *CreateCommentUseCase {
	return &CreateCommentUseCase{
		transactor: transactor,
	}
}

func (cc *CreateCommentUseCase) Create(ctx context.Context, com entity.Comment) (int64, error) {
	var out int64
	if err := cc.transactor.InTx(ctx, func(tx contract.TxRepo) error {
		emp, err := tx.GetEmployeeByID(ctx, com.CreatedByID)
		if err != nil {
			return fmt.Errorf("get employee by id: %w", err)
		}

		com.CreatedByName = emp.FullName

		id, err := tx.CreateComment(ctx, com)
		if err != nil {
			return fmt.Errorf("create comment: %w", err)
		}

		_, err = tx.CreateTaskHistoryChange(ctx, entity.TaskHistoryChange{
			TaskID:        com.TaskID,
			Event:         entity.TaskHistoryChangeEventCommentAdded,
			CreatedByID:   com.CreatedByID,
			CreatedByName: com.CreatedByName,
		})
		if err != nil {
			return fmt.Errorf("create task history change (comment added): %w", err)
		}

		out = id
		return nil
	}); err != nil {
		return 0, err
	}
	return out, nil
}
