package server

import (
	"go-tasks-api/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter sets up the router with all routes and middleware
func NewRouter(a *handler.Task) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	// tasks routes
	router.Route("/api/v1/tasks", func(r chi.Router) {
		r.Post("/", a.Create)
		r.Get("/", a.List)
		r.Get("/{id}", a.Get)
		r.Put("/{id}", a.Update)
		r.Delete("/{id}", a.Delete)
	})

	return router
}
