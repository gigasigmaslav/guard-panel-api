package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type DB struct {
	*pgxpool.Pool
}

func New(ctx context.Context, cfg Config) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("db.New: %w", err)
	}

	c, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("db.New: %w", err)
	}

	if err = c.Ping(ctx); err != nil {
		return nil, err
	}

	return &DB{c}, nil
}

func (d *DB) InTx(ctx context.Context, f func(pgx.Tx) error) (err error) {
	tx, err := d.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Ctx(ctx).Panic().
				Any("cause", p).
				Msg("tx finished with panic")
			_ = tx.Rollback(ctx)

			panic(p)
		} else if err != nil {
			log.Ctx(ctx).Error().
				Err(err).
				Msg("tx finished with err")
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = f(tx)

	return err
}

type Config struct {
	DBName   string
	HostPort string
	Username string
	Password string
}

func (c Config) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable&timezone=UTC",
		c.Username,
		c.Password,
		c.HostPort,
		c.DBName,
	)
}
