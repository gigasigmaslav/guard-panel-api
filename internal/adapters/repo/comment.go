package repo

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

func (q *queries) GetCommentsByTaskID(ctx context.Context, id int64) ([]entity.Comment, error) {
	const query = `
        SELECT 
            id,
            comment,
            created_by_id,
            created_by_name,
            created_at
        FROM guard.comments
        WHERE 
			task_id = $1
			AND deleted_at IS NULL
        ORDER BY created_at
    `

	rows, err := q.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("get comments by task ID query err: %w", err)
	}
	defer rows.Close()

	out := make([]entity.Comment, 0)

	for rows.Next() {
		var cur entity.Comment

		if err = rows.Scan(
			&cur.ID,
			&cur.Comment,
			&cur.CreatedByID,
			&cur.CreatedByName,
			&cur.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("get comments by task ID scan err: %w", err)
		}

		out = append(out, cur)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("get comments by task ID rows err: %w", err)
	}

	return out, nil
}

func (q *queries) CreateComment(ctx context.Context, com entity.Comment) (int64, error) {
	const query = `
		INSERT INTO guard.comments (
			task_id,
			comment,
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
		com.TaskID,
		com.Comment,
		com.CreatedByID,
		com.CreatedByName,
	).Scan(&out)
	if err != nil {
		return 0, fmt.Errorf("create comment storage error: %w", err)
	}

	return out, nil
}

func (q *queries) DeleteCommentByID(ctx context.Context, id int64) error {
	const query = `
		UPDATE guard.comments
		SET deleted_at = NOW()
		WHERE id = $1
	`

	_, err := q.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete comment storage error: %w", err)
	}

	return nil
}
