package repo

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

func (q *queries) GetViolatorByTaskID(ctx context.Context, id int64) (entity.Violator, error) {
	const query = `
		SELECT
			id,
			task_id,
			type,
			full_name,
			phone
		FROM guard.violators
		WHERE task_id = $1
		ORDER BY id
	`

	var cur entity.Violator

	var violatorType int32
	err := q.db.QueryRow(ctx, query, id).Scan(
		&cur.ID,
		&cur.TaskID,
		&violatorType,
		&cur.FullName,
		&cur.PhoneNumber,
	)
	if err != nil {
		return entity.Violator{}, fmt.Errorf("get violator by task ID query err: %w", err)
	}

	cur.Type = entity.ViolatorType(violatorType)

	return cur, nil
}

func (q *queries) CreateViolator(ctx context.Context, violator entity.Violator) (int64, error) {
	const query = `
		INSERT INTO guard.violators (
			task_id,
			type,
			full_name,
			phone
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
		violator.TaskID,
		int64(violator.Type),
		violator.FullName,
		violator.PhoneNumber,
	).Scan(&out)
	if err != nil {
		return 0, fmt.Errorf("create violator storage error: %w", err)
	}

	return out, nil
}
