package router

import (
	"github.com/Zapharaos/fihub-backend/internal/auth"
	"github.com/Zapharaos/fihub-backend/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"time"
)

// New Sets up the server's router
func New() *chi.Mux {
	// Create router
	r := chi.NewRouter()

	// Create auth
	a := auth.New(auth.CheckHeader)

	// Setup router
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Declare routes
	r.Route("/api/v1", func(r chi.Router) {

		// Health
		r.Get("/health", handlers.HealthCheckHandler)

		// Authentication
		r.Post("/auth/token", a.GetToken)

		// Protected routes
		r.Group(func(r chi.Router) {

			// TODO : use auth middleware

			// Users
			r.Route("/users", func(r chi.Router) {
				r.Post("/", handlers.CreateUser)
				r.Get("/{id}", handlers.GetUser)
			})
		})
	})

	// Return router
	return r
}
