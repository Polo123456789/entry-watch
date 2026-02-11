package util

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

func TestHandleErrorResponses(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name string
		// builder may return a modified request and the handler to exercise
		builder  func(*http.Request, *slog.Logger) (*http.Request, http.Handler)
		wantCode int
	}{
		{
			name: "ErrorWithCode",
			builder: func(r *http.Request, logger *slog.Logger) (*http.Request, http.Handler) {
				return r, Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
					return NewErrorWithCode("teapot", 418)
				})
			},
			wantCode: 418,
		},
		{
			name: "Forbidden",
			builder: func(r *http.Request, logger *slog.Logger) (*http.Request, http.Handler) {
				// create a disabled user so RequireRole returns ForbiddenError
				user := &entry.User{Enabled: false, Role: entry.RoleUser}
				req2 := r.WithContext(entry.WithUser(r.Context(), user))
				return req2, Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
					_, err := entry.RequireRole(r.Context(), entry.RoleUser)
					return err
				})
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Unauthorized",
			builder: func(r *http.Request, logger *slog.Logger) (*http.Request, http.Handler) {
				// no user in context -> RequireRole returns UnauthorizedError
				return r, Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
					_, err := entry.RequireRole(r.Context(), entry.RoleUser)
					return err
				})
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "UserSafeError",
			builder: func(r *http.Request, logger *slog.Logger) (*http.Request, http.Handler) {
				return r, Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
					return entry.NewUserSafeError("bad input")
				})
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Internal",
			builder: func(r *http.Request, logger *slog.Logger) (*http.Request, http.Handler) {
				return r, Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
					return errors.New("boom")
				})
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "ErrorCodeOnly",
			builder: func(r *http.Request, logger *slog.Logger) (*http.Request, http.Handler) {
				return r, Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
					// Handler is responsible for rendering the error page and status code
					// ErrorCodeOnly just signals that no error modal should be shown
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte("custom error page"))
					return NewErrorCodeOnly(http.StatusBadRequest)
				})
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			req2, h := tc.builder(req, logger)
			if req2 == nil {
				req2 = req
			}

			h.ServeHTTP(rec, req2)

			if rec.Code != tc.wantCode {
				t.Fatalf("%s: status = %d; want %d", tc.name, rec.Code, tc.wantCode)
			}
		})
	}
}
