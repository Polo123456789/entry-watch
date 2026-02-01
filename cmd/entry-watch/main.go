package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/charmbracelet/log"
	"github.com/gorilla/sessions"
	_ "modernc.org/sqlite"

	"github.com/Polo123456789/entry-watch/internal/entry"
	apphttp "github.com/Polo123456789/entry-watch/internal/http"
	"github.com/Polo123456789/entry-watch/internal/http/auth"
	"github.com/Polo123456789/entry-watch/internal/sqlc"
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

	db, err := sql.Open("sqlite", "./entry-watch.db")
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database", "error", err)
		}
	}()

	store := sqlc.NewStore(db)
	app := entry.NewApp(logger, store)

	userStore := sqlc.NewUserStore(db)

	if err := auth.EnsureSuperAdminExists(userStore, logger); err != nil {
		logger.Error("Failed to ensure superadmin exists", "error", err)
		os.Exit(1)
	}

	sessionKey := os.Getenv("SESSION_KEY")
	if len(sessionKey) < 32 {
		logger.Error("SESSION_KEY environment variable must be set with at least 32 characters")
		os.Exit(1)
	}

	secureCookies := !DEBUG || os.Getenv("FORCE_SECURE_COOKIES") == "true"

	sessionStore := sessions.NewCookieStore([]byte(sessionKey))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 12,
		HttpOnly: true,
		Secure:   secureCookies,
		SameSite: http.SameSiteLaxMode,
	}

	logger.Info("Session store configured",
		"secure", secureCookies,
		"same_site", "Lax",
	)

	server := apphttp.NewServer(
		"0.0.0.0",
		8080,
		app,
		logger,
		sessionStore,
		userStore,
	)

	apphttp.RunServer(ctx, cancel, server, logger)
}
