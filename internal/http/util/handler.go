package util

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

func Handler(
	logger *slog.Logger,
	h func(w http.ResponseWriter, r *http.Request) error,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)

		if err == nil {
			return
		} else if e, ok := errorAs[*ErrorWithCode](err); ok {
			http.Error(w, e.msg, e.code)
		} else if e, ok := errorAs[entry.ForbiddenError](err); ok {
			http.Error(w, e.Error(), http.StatusForbidden)
		} else if e, ok := errorAs[entry.UnauthorizedError](err); ok {
			http.Error(w, e.Error(), http.StatusUnauthorized)
		} else {
			// TODO: Get request ID from context
			logger.LogAttrs(
				r.Context(),
				slog.LevelError,
				"internal server error",
				slog.String("error", err.Error()),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
}

type ErrorWithCode struct {
	msg  string
	code int
}

func NewErrorWithCode(msg string, code int) *ErrorWithCode {
	return &ErrorWithCode{
		msg:  msg,
		code: code,
	}
}

func (e *ErrorWithCode) Error() string {
	return e.msg
}

func errorAs[T any](err error) (T, bool) {
	var target T
	if errors.As(err, &target) {
		return target, true
	}
	var zero T
	return zero, false
}
