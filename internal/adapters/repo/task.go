package repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

func (q *queries) GetTaskByID(ctx context.Context, id int64) (entity.Task, error) {
	const query = `
		SELECT
			id,
			damage_amount,
			priority,
			status,
			start_date,
			end_date,
			executor_id,
			executor_name,
			office_id,
			office_name,
			created_by_id,
			created_by_name,
			created_at
		FROM guard.tasks
		WHERE id = $1
	`

	var out entity.Task
	var priority int32
	var status int32
	err := q.db.QueryRow(ctx, query, id).Scan(
		&out.ID,
		&out.DamageAmount,
		&priority,
		&status,
		&out.StartDate,
		&out.EndDate,
		&out.ExecutorID,
		&out.ExecutorName,
		&out.OfficeID,
		&out.OfficeName,
		&out.CreatedByID,
		&out.CreatedByName,
		&out.CreatedAt,
	)
	if err != nil {
		return entity.Task{}, fmt.Errorf("get task by ID storage error: %w", err)
	}

	out.Priority = entity.TaskPriority(priority)
	out.Status = entity.TaskStatus(status)

	return out, nil
}

func (q *queries) CreateTask(ctx context.Context, task entity.Task) (int64, error) {
	const query = `
		INSERT INTO guard.tasks (
			damage_amount,
			priority,
			status,
			end_date,
			executor_id,
			executor_name,
			office_id,
			office_name,
			created_by_id,
			created_by_name
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
		RETURNING id;
	`

	var out int64
	err := q.db.QueryRow(
		ctx,
		query,
		task.DamageAmount,
		int32(task.Priority),
		int32(task.Status),
		task.EndDate,
		task.ExecutorID,
		task.ExecutorName,
		task.OfficeID,
		task.OfficeName,
		task.CreatedByID,
		task.CreatedByName,
	).Scan(&out)
	if err != nil {
		return 0, fmt.Errorf("create task storage error: %w", err)
	}

	return out, nil
}

func (q *queries) UpdateTaskByID(ctx context.Context, task entity.Task) error {
	const query = `
		UPDATE guard.tasks
		SET
			damage_amount = $1,
			priority = $2,
			status = $3,
			end_date = $4,
			executor_id = $5,
			executor_name = $6
		WHERE id = $7
	`

	_, err := q.db.Exec(
		ctx,
		query,
		task.DamageAmount,
		int32(task.Priority),
		int32(task.Status),
		task.EndDate,
		task.ExecutorID,
		task.ExecutorName,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("update task storage error: %w", err)
	}

	return nil
}

func (q *queries) SearchTasks(
	ctx context.Context,
	req contract.SearchTasksRequest,
) ([]contract.TaskLookupItem, int32, error) {
	builder := sq.Select(
		"t.id",
		"t.damage_amount",
		"t.priority",
		"t.executor_id",
		"t.executor_name",
		"COALESCE(v.id, 0) AS violator_id",
		"COALESCE(v.full_name, '') AS violator_name",
		"COALESCE(v.type, 0) AS violator_type",
		"v.phone AS violator_phone",
		"t.office_id",
		"t.office_name",
		"t.created_by_id",
		"t.created_by_name",
		"t.created_at",
		"t.start_date",
		"t.end_date",
		"COUNT(*) OVER() as total_count",
	).
		From("guard.tasks t").
		LeftJoin(
			"LATERAL (" +
				"SELECT id, full_name, type, phone FROM guard.violators v WHERE v.task_id = t.id ORDER BY v.id LIMIT 1" +
				") v ON TRUE",
		).
		PlaceholderFormat(sq.Dollar)

	builder = applyTaskSearchFilters(builder, req)
	orderBy := buildTaskSearchOrderBy(req.Sorting)
	builder = applyTaskSearchPaginationAndOrder(builder, req, orderBy)

	dataSQL, dataArgs, err := builder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build search tasks query: %w", err)
	}

	rows, err := q.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute search tasks query: %w", err)
	}
	defer rows.Close()

	out := make([]contract.TaskLookupItem, 0)
	var total int32
	for rows.Next() {
		var cur contract.TaskLookupItem
		var priority int32
		var violatorType int32

		if err = rows.Scan(
			&cur.ID,
			&cur.DamageAmount,
			&priority,
			&cur.ExecutorID,
			&cur.ExecutorName,
			&cur.ViolatorID,
			&cur.ViolatorName,
			&violatorType,
			&cur.ViolatorPhone,
			&cur.OfficeID,
			&cur.OfficeName,
			&cur.CreatedByID,
			&cur.CreatedByName,
			&cur.CreatedAt,
			&cur.StartDate,
			&cur.EndDate,
			&total,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan task row: %w", err)
		}

		cur.Priority = entity.TaskPriority(priority)
		cur.ViolatorType = entity.ViolatorType(violatorType)
		out = append(out, cur)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating over task rows: %w", err)
	}

	return out, total, nil
}

func buildTaskSearchOrderBy(sorting contract.SearchTasksSorting) string {
	const defaultOrderBy = "t.created_at"

	var orderBy string
	switch sorting.OrderBy {
	case contract.OrderByDamageAmountKopecks:
		orderBy = "t.damage_amount"
	case contract.OrderByCreatedAt, contract.OrderByUnspecified:
		orderBy = defaultOrderBy
	default:
		orderBy = defaultOrderBy
	}

	if sorting.OrderDirection == contract.OrderDirectionAsc {
		return orderBy + " ASC"
	}

	return orderBy + " DESC"
}

func applyTaskSearchPaginationAndOrder(
	builder sq.SelectBuilder,
	req contract.SearchTasksRequest,
	orderBy string,
) sq.SelectBuilder {
	limit := req.PerPage
	offset := (req.Page - 1) * req.PerPage

	if limit > 0 {
		builder = builder.Limit(uint64(limit))
	}
	if offset >= 0 {
		builder = builder.Offset(uint64(offset))
	}

	return builder.OrderBy(orderBy)
}

func applyTaskSearchFilters(
	builder sq.SelectBuilder,
	req contract.SearchTasksRequest,
) sq.SelectBuilder {
	if req.Status != entity.TaskStatusUnspecified {
		builder = builder.Where(sq.Eq{"t.status": int32(req.Status)})
	}

	if req.Filter == nil {
		return builder
	}

	if req.Filter.ID != nil {
		builder = builder.Where(sq.Eq{"t.id": *req.Filter.ID})
	}

	if req.Filter.Priority != nil {
		builder = builder.Where(sq.Eq{"t.priority": int32(*req.Filter.Priority)})
	}

	if req.Filter.ExecutorID != nil {
		builder = builder.Where(sq.Eq{"t.executor_id": *req.Filter.ExecutorID})
	}

	if req.Filter.OfficeID != nil {
		builder = builder.Where(sq.Eq{"t.office_id": *req.Filter.OfficeID})
	}

	if req.Filter.ViolatorType != nil {
		builder = builder.Where(
			`EXISTS (
				SELECT 1
				FROM guard.violators v2
				WHERE v2.task_id = t.id
					AND v2.type = ?
			)`,
			int32(*req.Filter.ViolatorType),
		)
	}

	if req.Filter.KUSP != nil {
		builder = builder.Where(
			`EXISTS (
				SELECT 1
				FROM guard.vud_decisions vd2
				WHERE vd2.task_id = t.id
					AND vd2.kusp = ?
			)`,
			*req.Filter.KUSP,
		)
	}

	if req.Filter.UD != nil {
		builder = builder.Where(
			`EXISTS (
				SELECT 1
				FROM guard.vud_decisions vd2
				WHERE vd2.task_id = t.id
					AND vd2.ud = ?
			)`,
			*req.Filter.UD,
		)
	}

	return builder
}
