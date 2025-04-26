package router

import (
	"github.com/Zapharaos/fihub-backend/cmd/api/app/auth"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/spf13/viper"
	"strings"
	"time"
)

// New Sets up the api's router
func New(config auth.Config) *chi.Mux {

	// Create router
	r := chi.NewRouter()

	// Create auth
	var a *auth.Auth
	if config.Security {
		a = auth.New(auth.CheckHeader, config)
	}

	// Setup router
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS
	if config.CORS {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   strings.Split(config.AllowOrigin, ","),
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
	}

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Setup handler utils
	handlers.ReplaceGlobals(handlers.NewUtils())

	// Declare routes
	apiBasePath := viper.GetString("API_BASE_PATH")
	r.Route(apiBasePath, func(r chi.Router) {

		// Health
		r.Get("/health", handlers.HealthCheckHandler)

		// Authentication
		r.Route("/auth", func(r chi.Router) {

			// Token
			if config.Security {
				r.Post("/token", a.GetToken)
			}

			// Password routes
			r.Route("/password", func(r chi.Router) {

				// Create password reset request
				requestLimit := viper.GetInt("OTP_MIDDLEWARE_REQUEST_LIMIT")
				requestLength := viper.GetDuration("OTP_MIDDLEWARE_REQUEST_LENGTH")
				if requestLength == 0 {
					requestLength = 24 * time.Hour
				}
				r.With(httprate.LimitByIP(requestLimit, requestLength)).Post("/", handlers.CreatePasswordResetRequest)

				// Input token and retrieve requestID using userID
				inputLimit := viper.GetInt("OTP_MIDDLEWARE_INPUT_LIMIT")
				inputLength := viper.GetDuration("OTP_MIDDLEWARE_INPUT_WINDOW")
				if inputLength == 0 {
					inputLength = 1 * time.Hour
				}
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

		// Apply auth middleware only if security is enabled
		if a != nil && a.Config.Security {
			r.Use(a.Middleware)
		}

		// Users
		r.Route("/users", func(r chi.Router) {
			r.Post("/", handlers.CreateUser)
			// r.Get("/", handlers.GetAllUsersWithRoles)

			// User's self profile : retrieving userID through context
			r.Route("/me", func(r chi.Router) {
				r.Get("/", handlers.GetUserSelf)
				r.Put("/", handlers.UpdateUserSelf)
				r.Delete("/", handlers.DeleteUserSelf)

				// User's password : retrieving userID through context
				r.Put("/password", handlers.UpdateUserPassword)
			})

			// User specific
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetUser)

				// Roles
				r.Route("/roles", func(r chi.Router) {
					r.Get("/", handlers.GetUserRoles)
					r.Put("/", handlers.SetUserRoles)
				})
			})

			// User's brokers : retrieving userID through context
			r.Route("/brokers", func(r chi.Router) {
				r.Post("/", handlers.CreateUserBroker)
				r.Get("/", handlers.ListUserBrokers)

				r.Delete("/{id}", handlers.DeleteUserBroker)
			})
		})

		// Roles
		r.Route("/roles", func(r chi.Router) {
			r.Post("/", handlers.CreateRole)
			r.Get("/", handlers.GetRoles)

			// Role specific
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetRole)
				r.Put("/", handlers.UpdateRole)
				r.Delete("/", handlers.DeleteRole)

				// Users
				r.Route("/users", func(r chi.Router) {
					/*r.Get("/", handlers.GetRoleUsers)
					r.Put("/", handlers.PutUsersRole)
					r.Delete("/", handlers.DeleteUsersRole)*/
				})

				// Permissions
				r.Route("/permissions", func(r chi.Router) {
					r.Get("/", handlers.GetRolePermissions)
					r.Put("/", handlers.SetRolePermissions)
				})
			})

		})

		// Permissions
		r.Route("/permissions", func(r chi.Router) {
			r.Post("/", handlers.CreatePermission)
			r.Get("/{id}", handlers.GetPermission)
			r.Put("/{id}", handlers.UpdatePermission)
			r.Delete("/{id}", handlers.DeletePermission)
			r.Get("/", handlers.ListPermissions)
		})

		// Brokers
		r.Route("/brokers", func(r chi.Router) {
			r.Post("/", handlers.CreateBroker)
			r.Get("/", handlers.ListBrokers)

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
			r.Get("/", handlers.ListTransactions)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetTransaction)
				r.Put("/", handlers.UpdateTransaction)
				r.Delete("/", handlers.DeleteTransaction)
			})
		})
	}
}
