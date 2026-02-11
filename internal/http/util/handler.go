package util

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/templates/common"
)

func Handler(
	logger *slog.Logger,
	h func(w http.ResponseWriter, r *http.Request) error,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			HandleError(w, r, logger, err)
		}
	})
}

// Add here special handling for different error types
func HandleError(
	w http.ResponseWriter, r *http.Request, logger *slog.Logger, err error,
) {
	if e, ok := errorAs[*ErrorWithCode](err); ok {
		errorModal(w, r, e.Error(), e.code)
	} else if e, ok := errorAs[*entry.ForbiddenError](err); ok {
		errorModal(w, r, e.Error(), http.StatusForbidden)
	} else if e, ok := errorAs[*entry.UnauthorizedError](err); ok {
		errorModal(w, r, e.Error(), http.StatusUnauthorized)
	} else if e, ok := errorAs[entry.UserSafeError](err); ok {
		errorModal(w, r, e.Error(), http.StatusBadRequest)
	} else {
		// TODO: Get request ID from context and expose it to the user
		logger.LogAttrs(
			r.Context(),
			slog.LevelError,
			"internal server error",
			slog.String("error", err.Error()),
		)
		errorModal(w, r, "internal server error", http.StatusInternalServerError)
	}
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

func errorModal(w http.ResponseWriter, r *http.Request, msg string, code int) {
	w.Header().Set("HX-Retarget", "#error-modal")
	w.WriteHeader(code)
	_ = common.ErrorModal(msg).Render(r.Context(), w)
}
