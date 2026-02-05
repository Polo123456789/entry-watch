package util

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
		wantBody string
	}{
		{
			name: "ErrorWithCode",
			builder: func(r *http.Request, logger *slog.Logger) (*http.Request, http.Handler) {
				return r, Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
					return NewErrorWithCode("teapot", 418)
				})
			},
			wantCode: 418,
			wantBody: "teapot",
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
			wantBody: "user is disabled",
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
			wantBody: "user not authenticated",
		},
		{
			name: "UserSafeError",
			builder: func(r *http.Request, logger *slog.Logger) (*http.Request, http.Handler) {
				return r, Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
					return entry.NewUserSafeError("bad input")
				})
			},
			wantCode: http.StatusBadRequest,
			wantBody: "bad input",
		},
		{
			name: "Internal",
			builder: func(r *http.Request, logger *slog.Logger) (*http.Request, http.Handler) {
				return r, Handler(logger, func(w http.ResponseWriter, r *http.Request) error {
					return errors.New("boom")
				})
			},
			wantCode: http.StatusInternalServerError,
			wantBody: "internal server error",
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

			body := strings.TrimSpace(rec.Body.String())
			if body != tc.wantBody {
				t.Fatalf("%s: body = %q; want %q", tc.name, body, tc.wantBody)
			}
		})
	}
}
