package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/sessions"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

type wrappedWritter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWritter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func CanonicalLoggerMiddleware(logger *slog.Logger, app *entry.App, store sessions.Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &wrappedWritter{w, http.StatusOK}

		// Attempt to authenticate user from session-backed uid.
		ctx := r.Context()
		if app != nil && store != nil {
			sess, err := store.Get(r, "entrywatch_session")
			if err == nil && sess != nil {
				if v, ok := sess.Values["uid"]; ok {
					// robustly parse numeric value
					var uid int64
					_, err := fmt.Sscanf(fmt.Sprintf("%v", v), "%d", &uid)
					if err == nil && uid > 0 {
						if su, err := app.UserGetByID(ctx, uid); err == nil && su != nil && su.Enabled {
							ctx = entry.WithUser(ctx, &entry.User{
								ID:            su.ID,
								CondominiumID: su.CondominiumID,
								Role:          su.Role,
								Enabled:       su.Enabled,
							})
						} else {
							// clear invalid or disabled session uid to avoid stale auth
							delete(sess.Values, "uid")
							// Set Secure flag based on TLS so local dev HTTP works
							secure := r.TLS != nil
							sess.Options = &sessions.Options{
								Path:     "/",
								HttpOnly: true,
								MaxAge:   60 * 60 * 12,
								Secure:   secure,
								SameSite: http.SameSiteStrictMode,
							}
							_ = sess.Save(r, w)
						}
					}

				}

			}
		}

		next.ServeHTTP(ww, r.WithContext(ctx))

		attrs := []slog.Attr{
			slog.String("url", r.URL.String()),
			slog.String("method", r.Method),
			slog.Int("status_code", ww.statusCode),
			slog.Duration("duration", time.Since(start)),
		}

		// If we have an authenticated user in the context, add its ID and role
		if u := entry.UserFromCtx(ctx); u != nil {
			attrs = append(attrs, slog.Int64("user_id", u.ID))
			attrs = append(attrs, slog.String("role", string(u.Role)))
		}

		logger.LogAttrs(
			ctx,
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
