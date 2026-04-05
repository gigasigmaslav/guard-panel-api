package repo

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

func (q *queries) GetRefundsByTaskID(ctx context.Context, id int64) ([]entity.Refund, error) {
	const query = `
		SELECT
			id,
			task_id,
			amount,
			comment,
			created_by_id,
			created_by_name,
			created_at
		FROM guard.refunds
		WHERE task_id = $1
		ORDER BY created_at
	`

	rows, err := q.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("get refunds by task ID query err: %w", err)
	}
	defer rows.Close()

	out := make([]entity.Refund, 0)

	for rows.Next() {
		var cur entity.Refund

		if err = rows.Scan(
			&cur.ID,
			&cur.TaskID,
			&cur.Amount,
			&cur.Comment,
			&cur.CreatedByID,
			&cur.CreatedByName,
			&cur.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("get refunds by task ID scan err: %w", err)
		}

		out = append(out, cur)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("get refunds by task ID rows err: %w", err)
	}

	return out, nil
}

func (q *queries) CreateRefund(ctx context.Context, refund entity.Refund) (int64, error) {
	const query = `
		INSERT INTO guard.refunds (
			task_id,
			amount,
			comment,
			created_by_id,
			created_by_name
		)
		VALUES (
			$1, $2, $3, $4, $5
		)
		RETURNING id;
	`

	var out int64
	err := q.db.QueryRow(
		ctx,
		query,
		refund.TaskID,
		refund.Amount,
		refund.Comment,
		refund.CreatedByID,
		refund.CreatedByName,
	).Scan(&out)
	if err != nil {
		return 0, fmt.Errorf("create refund storage error: %w", err)
	}

	return out, nil
}
