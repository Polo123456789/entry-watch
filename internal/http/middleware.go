package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Polo123456789/entry-watch/internal/app"
)

type wrappedWritter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWritter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func CanonicalLoggerMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &wrappedWritter{w, http.StatusOK}

		// TODO: Use actual auth instead of the mock one
		ctx := app.WithUser(r.Context(), &app.User{
			ID:            1,
			CondominiumID: 1,
			Role:          app.RoleSuperAdmin,
			Enabled:       true,
		})

		next.ServeHTTP(ww, r.WithContext(ctx))

		attrs := []slog.Attr{
			slog.String("url", r.URL.String()),
			slog.String("method", r.Method),
			slog.Int("status_code", ww.statusCode),
			slog.Duration("duration", time.Since(start)),
		}

		//

		/*
			You might want to add more stuff in here, like ips, the user that
			made the request, or the request ID if you have one.

			if u, ok := CurrentUser(r); ok {
				attrs = append(attrs, slog.Int64("user_id", u.ID))
				// Maybe add it to the context here?
			}

			if ip := r.Header.Get("X-Real-IP"); ip != "" {
				attrs = append(attrs, slog.String("ip", ip))
			} else {
				attrs = append(attrs, slog.String("ip", r.RemoteAddr))
			}
		*/

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
