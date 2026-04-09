package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

func (q *queries) GetEmployeeByID(ctx context.Context, id int64) (entity.Employee, error) {
	const query = `
		SELECT
			id,
			full_name,
			position,
			created_by_id,
			created_by_name,
			created_at
		FROM staff.employees
		WHERE id = $1
			AND deleted_at IS NULL
	`

	var out entity.Employee
	var position int32

	err := q.db.QueryRow(ctx, query, id).Scan(
		&out.ID,
		&out.FullName,
		&position,
		&out.CreatedByID,
		&out.CreatedByName,
		&out.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Employee{}, entity.ErrEmployeeNotFound
		}
		return entity.Employee{}, fmt.Errorf("get employee by id storage error: %w", err)
	}

	out.Position = entity.EmployeePosition(position)

	return out, nil
}

func (q *queries) GetEmployees(ctx context.Context) ([]entity.Employee, error) {
	const query = `
		SELECT
			id,
			full_name,
			position,
			created_by_id,
			created_by_name,
			created_at
		FROM staff.employees
		WHERE deleted_at IS NULL
		ORDER BY id DESC
	`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all employees query err: %w", err)
	}
	defer rows.Close()

	out := make([]entity.Employee, 0)

	for rows.Next() {
		var cur entity.Employee
		var position int32

		if err = rows.Scan(
			&cur.ID,
			&cur.FullName,
			&position,
			&cur.CreatedByID,
			&cur.CreatedByName,
			&cur.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("get all employees scan err: %w", err)
		}

		cur.Position = entity.EmployeePosition(position)

		out = append(out, cur)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("get all employees rows err: %w", err)
	}

	return out, nil
}

func (q *queries) CreateEmployee(ctx context.Context, emp entity.Employee) (int64, error) {
	const query = `
		INSERT INTO staff.employees (
			full_name,
			position,
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
		emp.FullName,
		int32(emp.Position),
		emp.CreatedByID,
		emp.CreatedByName,
	).Scan(&out)
	if err != nil {
		return 0, fmt.Errorf("create employee storage error: %w", err)
	}

	return out, nil
}

func (q *queries) UpdateEmployeeByID(ctx context.Context, emp entity.Employee) error {
	const query = `
		UPDATE staff.employees
		SET
			full_name = $1,
			position = $2
		WHERE id = $3
			AND deleted_at IS NULL
	`

	_, err := q.db.Exec(
		ctx,
		query,
		emp.FullName,
		int32(emp.Position),
		emp.ID,
	)
	if err != nil {
		return fmt.Errorf("update employee storage error: %w", err)
	}

	return nil
}

func (q *queries) DeleteEmployeeByID(ctx context.Context, id int64) error {
	const query = `
		UPDATE staff.employees
		SET deleted_at = NOW()
		WHERE id = $1
			AND deleted_at IS NULL
	`

	_, err := q.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete employee storage error: %w", err)
	}

	return nil
}
