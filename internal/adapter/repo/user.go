package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

func (q *queries) GetUserByEmployeeID(ctx context.Context, employeeID int64) (entity.User, error) {
	const query = `
		SELECT
			id,
			employee_id,
			password_hash,
			created_at
		FROM auth.users
		WHERE employee_id = $1
			AND deleted_at IS NULL
	`

	var out entity.User

	err := q.db.QueryRow(ctx, query, employeeID).Scan(
		&out.ID,
		&out.EmployeeID,
		&out.PasswordHash,
		&out.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, entity.ErrUserNotFound
		}

		return entity.User{}, fmt.Errorf("get user by employee id storage error: %w", err)
	}

	return out, nil
}

func (q *queries) UserExistsByEmployeeID(ctx context.Context, employeeID int64) (bool, error) {
	const query = `
		SELECT EXISTS(
			SELECT 1 FROM auth.users 
			WHERE employee_id = $1 AND deleted_at IS NULL
		)
	`

	var exists bool
	err := q.db.QueryRow(ctx, query, employeeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("user exists by employee id storage error: %w", err)
	}

	return exists, nil
}

func (q *queries) CreateUser(ctx context.Context, u entity.User) (int64, error) {
	const query = `
		INSERT INTO auth.users (
			employee_id,
			password_hash
		)
		VALUES (
			$1, $2
		)
		RETURNING id
	`

	var id int64

	err := q.db.QueryRow(ctx, query, u.EmployeeID, u.PasswordHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create user storage error: %w", err)
	}

	return id, nil
}
