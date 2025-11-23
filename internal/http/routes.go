package http

import (
	"log/slog"
	"net/http"
)

func setupRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	})
}
