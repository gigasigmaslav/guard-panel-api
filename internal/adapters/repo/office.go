package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

func (q *queries) GetOfficeByID(ctx context.Context, id int64) (entity.Office, error) {
	const query = `
		SELECT
			id,
			name,
			address,
			created_by_id,
			created_by_name,
			created_at,
			deleted_at
		FROM staff.offices
		WHERE id = $1
			AND deleted_at IS NULL
	`

	var out entity.Office
	err := q.db.QueryRow(ctx, query, id).Scan(
		&out.ID,
		&out.Name,
		&out.Address,
		&out.CreatedByID,
		&out.CreatedByName,
		&out.CreatedAt,
		&out.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Office{}, entity.ErrOfficeNotFound
		}
		return entity.Office{}, fmt.Errorf("get office by id storage error: %w", err)
	}

	return out, nil
}

func (q *queries) GetOffices(ctx context.Context) ([]entity.Office, error) {
	const query = `
		SELECT
			id,
			name,
			address,
			created_by_id,
			created_by_name,
			created_at,
			deleted_at
		FROM staff.offices
		WHERE deleted_at IS NULL
		ORDER BY id DESC
	`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all offices query err: %w", err)
	}
	defer rows.Close()

	out := make([]entity.Office, 0)

	for rows.Next() {
		var cur entity.Office

		if err = rows.Scan(
			&cur.ID,
			&cur.Name,
			&cur.Address,
			&cur.CreatedByID,
			&cur.CreatedByName,
			&cur.CreatedAt,
			&cur.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("get all offices scan err: %w", err)
		}

		out = append(out, cur)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("get all offices rows err: %w", err)
	}

	return out, nil
}

func (q *queries) CreateOffice(ctx context.Context, office entity.Office) (int64, error) {
	const query = `
		INSERT INTO staff.offices (
			name,
			address,
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
		office.Name,
		office.Address,
		office.CreatedByID,
		office.CreatedByName,
	).Scan(&out)
	if err != nil {
		return 0, fmt.Errorf("create office storage error: %w", err)
	}

	return out, nil
}

func (q *queries) UpdateOfficeByID(ctx context.Context, office entity.Office) error {
	const query = `
		UPDATE staff.offices
		SET
			name = $1,
			address = $2
		WHERE id = $3
			AND deleted_at IS NULL
	`

	_, err := q.db.Exec(
		ctx,
		query,
		office.Name,
		office.Address,
		office.ID,
	)
	if err != nil {
		return fmt.Errorf("update office storage error: %w", err)
	}

	return nil
}

func (q *queries) DeleteOfficeByID(ctx context.Context, id int64) error {
	const query = `
		UPDATE staff.offices
		SET deleted_at = NOW()
		WHERE id = $1
			AND deleted_at IS NULL
	`

	_, err := q.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete office storage error: %w", err)
	}

	return nil
}
