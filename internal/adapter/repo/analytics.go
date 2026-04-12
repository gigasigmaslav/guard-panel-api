package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
	"github.com/jackc/pgx/v5"
)

func (q *queries) GetTaskDashboardKPI(
	ctx context.Context,
	from time.Time,
	to time.Time,
) (contract.TaskDashboardKPI, error) {
	completedStatus := int32(entity.TaskStatusCompleted)

	const query = `
		WITH b AS (
			SELECT $1::timestamptz AS ts_from, $2::timestamptz AS ts_to
		),
		period_created AS (
			SELECT COUNT(*)::bigint AS n
			FROM guard.tasks t
			CROSS JOIN b
			WHERE t.created_at >= b.ts_from
	  			AND t.created_at <= b.ts_to
		),
		period_completed AS (
			SELECT COUNT(*)::bigint AS n
			FROM guard.tasks t
			CROSS JOIN b
			WHERE t.status = $3
	  			AND t.end_date IS NOT NULL
	  			AND t.end_date >= b.ts_from
	  			AND t.end_date <= b.ts_to
		)
		SELECT
		(SELECT COUNT(*)::bigint FROM guard.tasks WHERE status <> $3) AS active_tasks_count,
		(SELECT n FROM period_completed) AS completed_in_period,
		(SELECT n FROM period_created) AS created_in_period,
		CASE
		WHEN (SELECT n FROM period_created) = 0 THEN NULL::double precision
		ELSE (SELECT n FROM period_completed)::double precision
			/ (SELECT n FROM period_created)::double precision
		END AS completed_to_created_ratio
	`

	var out contract.TaskDashboardKPI
	var ratio sql.NullFloat64
	err := q.db.QueryRow(ctx, query, from, to, completedStatus).Scan(
		&out.ActiveTasksCount,
		&out.CompletedInPeriod,
		&out.CreatedInPeriod,
		&ratio,
	)
	if err != nil {
		return contract.TaskDashboardKPI{}, fmt.Errorf("get task dashboard kpi: %w", err)
	}
	if ratio.Valid {
		v := ratio.Float64
		out.CompletedToCreatedRatio = &v
	}
	return out, nil
}

func (q *queries) GetTopExecutorByCompletedInPeriod(
	ctx context.Context,
	from time.Time,
	to time.Time,
) (contract.TopExecutorByCompleted, error) {
	completedStatus := int32(entity.TaskStatusCompleted)

	const query = `
		SELECT
			executor_id,
			executor_name,
			COUNT(*)::bigint AS completed
		FROM guard.tasks
		WHERE status = $3
  			AND end_date IS NOT NULL
  			AND end_date >= $1
  			AND end_date <= $2
		GROUP BY executor_id, executor_name
		ORDER BY completed DESC, executor_id ASC
		LIMIT 1
	`

	var out contract.TopExecutorByCompleted
	err := q.db.QueryRow(ctx, query, from, to, completedStatus).Scan(
		&out.ExecutorID,
		&out.ExecutorName,
		&out.Completed,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return contract.TopExecutorByCompleted{}, nil
		}
		return contract.TopExecutorByCompleted{}, fmt.Errorf("get top executor by completed in period: %w", err)
	}
	return out, nil
}
