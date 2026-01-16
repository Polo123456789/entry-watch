package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"database/sql"
	"github.com/charmbracelet/log"
	_ "modernc.org/sqlite"

	"github.com/Polo123456789/entry-watch/db"
	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http"
)

// Set in config, you set that
const DEBUG = true

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt, os.Kill,
	)

	var logger *slog.Logger
	if DEBUG {
		logger = slog.New(log.NewWithOptions(os.Stderr, log.Options{
			Level:           log.DebugLevel,
			ReportTimestamp: false,
			ReportCaller:    true,
			CallerOffset:    0,
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}

	// Ensure data directory exists and open sqlite DB for prod wiring
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		logger.Error("failed to create data dir", "error", err)
		os.Exit(1)
	}

	// Allow overriding DSN via environment for flexibility in prod/tests
	dsn := os.Getenv("ENTRYWATCH_DB_DSN")
	if dsn == "" {
		dsn = "./data/db.sqlite3"
	}

	// Open DB and run migrations. Keep MemStore for now until sqlstore implemented.
	var sqlDB *sql.DB
	{
		var err error
		sqlDB, err = sql.Open("sqlite", dsn)
		if err != nil {
			logger.Error("failed to open sqlite db", "error", err, "dsn", dsn)
			os.Exit(1)
		}
		// Enable foreign keys pragma
		if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
			logger.Info("failed to set PRAGMA foreign_keys", "error", err)
		}
		// Run embedded migrations
		if err := db.AutoMigrate(sqlDB, logger); err != nil {
			logger.Error("migrations failed", "error", err)
			os.Exit(1)
		}
	}

	// NOTE: sqlstore not implemented yet; keep using MemStore until sqlc/sqlstore is added.
	store := entry.NewMemStore()
	app := entry.NewApp(logger, store)

	// Ensure a bootstrap superadmin exists in the store
	if err := app.BootstrapSuperadmin(ctx); err != nil {
		logger.Error("bootstrap superadmin failed", "error", err)
		os.Exit(1)
	}

	server := http.NewServer(
		"0.0.0.0",
		8080,
		app,
		logger,
	)

	http.RunServer(ctx, cancel, server, logger)
}
