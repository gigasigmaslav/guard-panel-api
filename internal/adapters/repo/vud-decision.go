package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

func (q *queries) GetVUDDecisionByID(ctx context.Context, id int64) (entity.VUDDecision, error) {
	const query = `
		SELECT
			id,
			task_id,
			criminal_case_opened,
			comment,
			kusp,
			ud,
			created_by_id,
			created_by_name,
			created_at
		FROM guard.vud_decisions
		WHERE id = $1
	`

	var out entity.VUDDecision
	err := q.db.QueryRow(ctx, query, id).Scan(
		&out.ID,
		&out.TaskID,
		&out.CriminalCaseOpened,
		&out.Comment,
		&out.KUSP,
		&out.UD,
		&out.CreatedByID,
		&out.CreatedByName,
		&out.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.VUDDecision{}, entity.ErrVUDDecisionNotFound
		}
		return entity.VUDDecision{}, fmt.Errorf("get vud decision by id storage error: %w", err)
	}

	return out, nil
}

func (q *queries) GetVUDDecisionByTaskID(ctx context.Context, taskID int64) (entity.VUDDecision, error) {
	const query = `
		SELECT
			id,
			task_id,
			criminal_case_opened,
			comment,
			kusp,
			ud,
			created_by_id,
			created_by_name,
			created_at
		FROM guard.vud_decisions
		WHERE task_id = $1
	`

	var out entity.VUDDecision
	err := q.db.QueryRow(ctx, query, taskID).Scan(
		&out.ID,
		&out.TaskID,
		&out.CriminalCaseOpened,
		&out.Comment,
		&out.KUSP,
		&out.UD,
		&out.CreatedByID,
		&out.CreatedByName,
		&out.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.VUDDecision{}, entity.ErrVUDDecisionNotFound
		}
		return entity.VUDDecision{}, fmt.Errorf("get vud decision by task ID query err: %w", err)
	}

	return out, nil
}

func (q *queries) CreateVUDDecision(ctx context.Context, decision entity.VUDDecision) (int64, error) {
	const query = `
		INSERT INTO guard.vud_decisions (
			task_id,
			criminal_case_opened,
			comment,
			kusp,
			ud,
			created_by_id,
			created_by_name
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id;
	`

	var out int64
	err := q.db.QueryRow(
		ctx,
		query,
		decision.TaskID,
		decision.CriminalCaseOpened,
		decision.Comment,
		decision.KUSP,
		decision.UD,
		decision.CreatedByID,
		decision.CreatedByName,
	).Scan(&out)
	if err != nil {
		return 0, fmt.Errorf("create vud decision storage error: %w", err)
	}

	return out, nil
}

func (q *queries) UpdateVUDDecisionByID(ctx context.Context, decision entity.VUDDecision) error {
	const query = `
		UPDATE guard.vud_decisions
		SET
			criminal_case_opened = $1,
			comment = $2,
			kusp = $3,
			ud = $4
		WHERE id = $5
	`

	_, err := q.db.Exec(
		ctx,
		query,
		decision.CriminalCaseOpened,
		decision.Comment,
		decision.KUSP,
		decision.UD,
		decision.ID,
	)
	if err != nil {
		return fmt.Errorf("update vud decision storage error: %w", err)
	}

	return nil
}
