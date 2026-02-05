package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/Polo123456789/assert"
	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/auth"
	"github.com/gorilla/sessions"
)

func NewServer(
	address string,
	port int,
	app *entry.App,
	logger *slog.Logger,
	session sessions.Store,
	userStore auth.UserStore,
) *http.Server {
	assert.NotEquals(address, "")
	assert.MoreThan(port, 0)
	assert.NotNil(logger)

	mux := http.NewServeMux()

	setupRoutes(
		mux,
		app,
		logger,
		session,
		userStore,
	)

	// Global middlewares
	var handler http.Handler = mux
	handler = CanonicalLoggerMiddleware(logger, session, handler)
	handler = RecoverMiddleware(logger, handler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: handler,
	}

	return server
}

func RunServer(
	globalCtx context.Context,
	cancel context.CancelFunc,
	server *http.Server,
	logger *slog.Logger,
) {
	assert.NotNil(server)
	assert.NotNil(logger)

	go func() {
		logger.Info("Starting server", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
			cancel()
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-globalCtx.Done()

		logger.Info("Shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Failed to shutdown server", "error", err)
		}
	}()

	wg.Wait()
}
