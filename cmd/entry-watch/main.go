package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/charmbracelet/log"
	_ "modernc.org/sqlite"

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
