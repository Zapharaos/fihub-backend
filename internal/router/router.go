package router

import (
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/auth"
	"github.com/Zapharaos/fihub-backend/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
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
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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
		r.Group(buildProtectedRoutes(a))
	})

	// Return router
	return r
}

func printRoutes(r chi.Router) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
}

func buildProtectedRoutes(a *auth.Auth) func(r chi.Router) {
	return func(r chi.Router) {

		r.Use(a.Middleware)

		// Users
		r.Route("/users", func(r chi.Router) {
			r.Post("/", handlers.CreateUser)
			r.Get("/me", handlers.GetUserSelf)

			// User's brokers : retrieving userID through context
			r.Route("/brokers", func(r chi.Router) {
				r.Post("/", handlers.CreateUserBroker)
				r.Get("/", handlers.GetUserBrokers)

				r.Delete("/{id}", handlers.DeleteUserBroker)
			})
		})

		// TODO : scan utils

		// Brokers
		r.Route("/brokers", func(r chi.Router) {
			r.Post("/", handlers.CreateBroker)
			r.Get("/", handlers.GetBrokers)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetBroker)
				r.Put("/", handlers.UpdateBroker)
				r.Delete("/", handlers.DeleteBroker)
			})
		})

		// Transactions : retrieving userID through context
		r.Route("/transactions", func(r chi.Router) {
			r.Post("/", handlers.CreateTransaction)
			r.Get("/", handlers.GetTransactions)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetTransaction)
				r.Delete("/", handlers.DeleteTransaction)
			})
		})
	}
}
