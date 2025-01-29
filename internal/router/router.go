package router

import (
	"github.com/Zapharaos/fihub-backend/internal/auth"
	"github.com/Zapharaos/fihub-backend/internal/handlers"
	"github.com/Zapharaos/fihub-backend/pkg/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
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
		r.Route("/auth", func(r chi.Router) {

			// Token
			r.Post("/token", a.GetToken)

			// Password routes
			r.Route("/password", func(r chi.Router) {

				// Create password reset request
				requestLimit := env.GetInt("OTP_MIDDLEWARE_REQUEST_LIMIT", 3)
				requestLength := env.GetDuration("OTP_MIDDLEWARE_REQUEST_LENGTH", 24*time.Hour)
				r.With(httprate.LimitByIP(requestLimit, requestLength)).Post("/", handlers.CreatePasswordResetRequest)

				// Input token and retrieve requestID using userID
				inputLimit := env.GetInt("OTP_MIDDLEWARE_INPUT_LIMIT", 5)
				inputLength := env.GetDuration("OTP_MIDDLEWARE_INPUT_WINDOW", 1*time.Hour)
				r.With(httprate.LimitByIP(inputLimit, inputLength)).Get("/{id}/{token}", handlers.GetPasswordResetRequestID)

				// Reset password using userID and requestID
				r.Put("/{id}/{request_id}", handlers.ResetPassword)
			})
		})

		// Protected routes
		r.Group(buildProtectedRoutes(a))
	})

	// Return router
	return r
}

func buildProtectedRoutes(a *auth.Auth) func(r chi.Router) {
	return func(r chi.Router) {

		r.Use(a.Middleware)

		// Users
		r.Route("/users", func(r chi.Router) {
			r.Post("/", handlers.CreateUser)

			// User's self profile : retrieving userID through context
			r.Route("/me", func(r chi.Router) {
				r.Get("/", handlers.GetUserSelf)
				r.Put("/", handlers.UpdateUserSelf)
				r.Delete("/", handlers.DeleteUserSelf)

				// User's password : retrieving userID through context
				r.Put("/password", handlers.ChangeUserPassword)
			})

			// User's brokers : retrieving userID through context
			r.Route("/brokers", func(r chi.Router) {
				r.Post("/", handlers.CreateUserBroker)
				r.Get("/", handlers.GetUserBrokers)

				r.Delete("/{id}", handlers.DeleteUserBroker)
			})
		})

		// Brokers
		r.Route("/brokers", func(r chi.Router) {
			r.Post("/", handlers.CreateBroker)
			r.Get("/", handlers.GetBrokers)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetBroker)
				r.Put("/", handlers.UpdateBroker)
				r.Delete("/", handlers.DeleteBroker)

				// Image
				r.Route("/image", func(r chi.Router) {
					r.Post("/", handlers.CreateBrokerImage)

					r.Route("/{image_id}", func(r chi.Router) {
						r.Get("/", handlers.GetBrokerImage)
						r.Put("/", handlers.UpdateBrokerImage)
						r.Delete("/", handlers.DeleteBrokerImage)
					})
				})
			})

		})

		// Transactions : retrieving userID through context
		r.Route("/transactions", func(r chi.Router) {
			r.Post("/", handlers.CreateTransaction)
			r.Get("/", handlers.GetTransactions)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetTransaction)
				r.Put("/", handlers.UpdateTransaction)
				r.Delete("/", handlers.DeleteTransaction)
			})
		})
	}
}
