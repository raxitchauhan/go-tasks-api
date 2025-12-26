package main

import (
	"context"

	"go-tasks-api/internal/config"
	"go-tasks-api/internal/migrator"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("initiating migration")
	cfg := config.LoadConfig()

	// run migration
	if err := migrator.Run(context.Background(), cfg.DSN(), cfg.DatabaseMigrationTable, migrator.FS); err != nil {
		log.Fatal().Err(err).Msg("failed to apply migration")
	}

	log.Info().Msg("migration successfully")
}
