package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/sessions"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/auth"
)

type wrappedWritter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWritter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func CanonicalLoggerMiddleware(
	logger *slog.Logger,
	session sessions.Store,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		authUser, userOk := auth.CurrentUser(session, r)
		if userOk {
			user := authUser.ToEntryUser()
			ctx := entry.WithUser(r.Context(), user)
			r = r.WithContext(ctx)
		}

		ww := &wrappedWritter{w, http.StatusOK}

		next.ServeHTTP(ww, r)

		attrs := []slog.Attr{
			slog.String("url", r.URL.String()),
			slog.String("method", r.Method),
			slog.Int("status_code", ww.statusCode),
			slog.Duration("duration", time.Since(start)),
		}

		if userOk {
			attrs = append(
				attrs,
				slog.Int64("user_id", authUser.ID),
				slog.String("user_role", string(authUser.Role)),
				slog.Int64("user_condo", authUser.CondominiumID),
			)
		} else {
			attrs = append(attrs, slog.String("user_id", "anonymous"))
		}

		logger.LogAttrs(
			r.Context(),
			slog.LevelInfo,
			"canonical-log",
			attrs...,
		)
	})
}

func RecoverMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.LogAttrs(
					r.Context(),
					slog.LevelError,
					"Panic Recovered",
					slog.String("path", r.URL.String()),
					slog.String("method", r.Method),
					slog.String("error", fmt.Sprintf("%v", err)),
				)
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
