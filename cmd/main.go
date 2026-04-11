package main

import (
	"context"

	gokit "github.com/cripplemymind9/go-utils/go-kit"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/gigasigmaslav/guard-panel-api/internal/app"
	"github.com/gigasigmaslav/guard-panel-api/internal/config"
)

func main() {
	cfg, err := config.Get(viper.New())
	if err != nil {
		log.Fatal().Err(err).Msg("get config")
	}

	runner := gokit.NewRunner()

	app, err := app.New(context.Background(), runner, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("get new app")
	}

	if runErr := runner.Run(app); runErr != nil {
		log.Fatal().Err(runErr).Msg("run app")
	}
}
