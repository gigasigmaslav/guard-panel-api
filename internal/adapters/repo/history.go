package repo

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

func (q *queries) CreateTaskHistoryChange(ctx context.Context, change entity.TaskHistoryChange) (int64, error) {
	const query = `
		INSERT INTO guard.history_changes (
			task_id,
			event,
			created_by_id,
			created_by_name
		)
		VALUES (
			$1, $2, $3, $4
		)
		RETURNING id;
	`

	var out int64
	err := q.db.QueryRow(
		ctx,
		query,
		change.TaskID,
		int32(change.Event),
		change.CreatedByID,
		change.CreatedByName,
	).Scan(&out)
	if err != nil {
		return 0, fmt.Errorf("create task history change storage error: %w", err)
	}

	return out, nil
}

func (q *queries) GetTaskHistoryChangesByTaskID(ctx context.Context, id int64) ([]entity.TaskHistoryChange, error) {
	const query = `
		SELECT
			id,
			task_id,
			event,
			created_by_id,
			created_by_name,
			created_at
		FROM guard.history_changes
		WHERE task_id = $1
		ORDER BY created_at
	`

	rows, err := q.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("get task history changes by task ID query err: %w", err)
	}
	defer rows.Close()

	out := make([]entity.TaskHistoryChange, 0)

	for rows.Next() {
		var cur entity.TaskHistoryChange
		var event int32

		if err = rows.Scan(
			&cur.ID,
			&cur.TaskID,
			&event,
			&cur.CreatedByID,
			&cur.CreatedByName,
			&cur.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("get task history changes by task ID scan err: %w", err)
		}

		cur.Event = entity.TaskHistoryChangeEvent(event)
		out = append(out, cur)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("get task history changes by task ID rows err: %w", err)
	}

	return out, nil
}
