package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Querier
	Transactor
}

type Querier interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	CopyFrom(
		ctx context.Context,
		tableName pgx.Identifier,
		columnNames []string,
		rowSrc pgx.CopyFromSource,
	) (int64, error)
}

type Transactor interface {
	InTx(context.Context, func(tx pgx.Tx) error) error
}

type Storage struct {
	tx Transactor
	*queries
}

func NewStorage(db DB) *Storage {
	return &Storage{
		tx:      db,
		queries: newQueries(db),
	}
}

func (s *Storage) InTx(ctx context.Context, f func(repo any) error) error {
	return s.tx.InTx(ctx, func(tx pgx.Tx) error {
		return f(newQueries(tx))
	})
}

type queries struct {
	db Querier
}

func newQueries(querier Querier) *queries {
	return &queries{
		db: querier,
	}
}
