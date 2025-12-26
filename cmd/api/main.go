package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-tasks-api/internal/config"
	"go-tasks-api/internal/database"
	"go-tasks-api/internal/handler"
	"go-tasks-api/internal/repository"
	"go-tasks-api/internal/server"

	"github.com/rs/zerolog/log"
)

type Service struct {
	taskHandler *handler.Task
}

func main() {
	// create a root context that can be cancelled on system interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	NewService(ctx).Run(ctx)
}

// NewService initializes the service with its dependencies
func NewService(ctx context.Context) *Service {
	cfg := config.LoadConfig()

	db, err := database.NewConnection(cfg.DSN(), cfg.DatabaseMaxOpenConns)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("failed to establish database connection: %w", err))
	}

	if err := database.MigrationVersionCheck(ctx, db, cfg.DatabaseMigrationTable, cfg.DatabaseMinVersion); err != nil {
		log.Fatal().Err(fmt.Errorf("failed while checking database migration version: %w", err))
	}

	taskRepo := repository.NewTaskRepo(db)
	taskHandler := handler.NewTaskHandler(taskRepo)

	return &Service{
		taskHandler,
	}
}

// Run starts the service
func (s *Service) Run(ctx context.Context) {
	webServer := server.NewServer(s.taskHandler)
	go func() {
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err)
		}
	}()

	defer func() {
		// new context for shutdown timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := webServer.Shutdown(shutdownCtx); err != nil {
			log.Fatal().Err(err).Msg("server shutdown failed")
		}
	}()

	<-ctx.Done()
}
