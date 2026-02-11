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
	} else if _, ok := errorAs[*ErrorCodeOnly](err); ok {
		// Handler already rendered the complete response (status + body)
		// No further action needed - error serves only to signal no modal should be shown
		return
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

// ErrorWithCode represents an error with a message and HTTP status code.
// The message will be displayed in an error modal.
type ErrorWithCode struct {
	msg  string
	code int
}

// NewErrorWithCode creates an error with a message that will be shown in the error modal.
func NewErrorWithCode(msg string, code int) *ErrorWithCode {
	return &ErrorWithCode{
		msg:  msg,
		code: code,
	}
}

func (e *ErrorWithCode) Error() string {
	return e.msg
}

// ErrorCodeOnly represents an error with only an HTTP status code.
// Use this when the handler has already rendered the error message in another way.
type ErrorCodeOnly struct {
	code int
}

// NewErrorCodeOnly creates an error with only a status code.
// The handler is responsible for rendering the error message.
func NewErrorCodeOnly(code int) *ErrorCodeOnly {
	return &ErrorCodeOnly{
		code: code,
	}
}

func (e *ErrorCodeOnly) Error() string {
	return http.StatusText(e.code)
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
