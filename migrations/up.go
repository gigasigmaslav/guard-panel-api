package migrations

import (
	"context"
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

const lockKey = 228

//go:embed postgres
var sqlDir embed.FS

func Up(_ context.Context, db *sql.DB) error {
	goose.SetTableName(goose.DefaultTablename)
	goose.SetBaseFS(sqlDir)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if _, err := db.Exec("SELECT pg_advisory_lock($1)", lockKey); err != nil {
		return err
	}

	defer func() {
		if _, err := db.Exec("SELECT pg_advisory_unlock($1)", lockKey); err != nil {
			log.Error().Err(err).Msg("got err unlocking DB")
		}
	}()

	return goose.Up(db, "postgres", goose.WithAllowMissing())
}
