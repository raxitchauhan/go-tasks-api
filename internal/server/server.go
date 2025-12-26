package server

import (
	"net/http"

	"go-tasks-api/internal/handler"
)

// NewServer creates and configures a new HTTP server
func NewServer(a *handler.Task) *http.Server {
	r := NewRouter(a)

	return &http.Server{
		Addr:    ":3000",
		Handler: r,
	}
}
