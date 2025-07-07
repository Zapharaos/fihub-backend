package router

import (
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/middleware"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/server"
	"github.com/go-chi/chi/v5"
	mw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/spf13/viper"
	"strings"
	"time"
)

// New Sets up the api's router
func New(config server.Config) *chi.Mux {

	// Create router
	r := chi.NewRouter()

	// Setup router
	r.Use(mw.RequestID)
	r.Use(mw.RealIP)
	r.Use(mw.Logger)
	r.Use(mw.Recoverer)

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
	r.Use(mw.Timeout(60 * time.Second))

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
			r.Post("/token", handlers.GetToken)

			// User registration
			r.Route("/register", func(r chi.Router) {

				// Retrieve config values for OTP request : rate limiting
				otpRequestRateLimit := viper.GetInt("OTP_REQUEST_RATE_LIMIT")
				otpRequestRateWindow := viper.GetDuration("OTP_REQUEST_RATE_WINDOW")
				if otpRequestRateWindow == 0 {
					otpRequestRateWindow = 24 * time.Hour
				}

				// Retrieve config values for OTP input : attempts limiting
				otpInputAttemptLimit := viper.GetInt("OTP_INPUT_ATTEMPT_LIMIT")
				otpInputAttemptWindow := viper.GetDuration("OTP_INPUT_ATTEMPT_WINDOW")
				if otpInputAttemptWindow == 0 {
					otpInputAttemptWindow = 1 * time.Hour
				}

				r.With(httprate.LimitByIP(otpRequestRateLimit, otpRequestRateWindow)).Post("/otp", handlers.GenerateSignupOTP)
				r.With(httprate.LimitByIP(otpInputAttemptLimit, otpInputAttemptWindow)).Post("/otp/validate", handlers.ValidateSignupOTP)
				r.Post("/", handlers.RegisterUser)
			})

			// Password routes
			r.Route("/password", func(r chi.Router) {

				// Forgot password flow
				r.Route("/reset", func(r chi.Router) {

					// Retrieve config values for OTP request : rate limiting
					otpRequestRateLimit := viper.GetInt("OTP_REQUEST_RATE_LIMIT")
					otpRequestRateWindow := viper.GetDuration("OTP_REQUEST_RATE_WINDOW")
					if otpRequestRateWindow == 0 {
						otpRequestRateWindow = 24 * time.Hour
					}

					// Retrieve config values for OTP input : attempts limiting
					otpInputAttemptLimit := viper.GetInt("OTP_INPUT_ATTEMPT_LIMIT")
					otpInputAttemptWindow := viper.GetDuration("OTP_INPUT_ATTEMPT_WINDOW")
					if otpInputAttemptWindow == 0 {
						otpInputAttemptWindow = 1 * time.Hour
					}

					r.With(httprate.LimitByIP(otpRequestRateLimit, otpRequestRateWindow)).Post("/otp", handlers.GenerateForgottenPasswordOTP)
					r.With(httprate.LimitByIP(otpInputAttemptLimit, otpInputAttemptWindow)).Post("/otp/validate", handlers.ValidateForgottenPasswordOTP)
					r.Put("/", handlers.ResetForgottenPassword)
				})
			})
		})

		// Protected routes
		r.Group(buildProtectedRoutes(config))
	})

	// Return router
	return r
}

func buildProtectedRoutes(config server.Config) func(r chi.Router) {
	return func(r chi.Router) {

		// Apply auth middleware to protected routes
		r.Use(middleware.AuthMiddleware(config))

		// Users
		r.Route("/user", func(r chi.Router) {

			// User's self profile : retrieving userID through context
			r.Route("/me", func(r chi.Router) {
				r.Get("/", handlers.GetUserSelf)
				r.Put("/", handlers.UpdateUserSelf)
				r.Delete("/", handlers.DeleteUserSelf)

				// Logged-in user changing password (security verification)
				r.Route("/password", func(r chi.Router) {

					// Retrieve config values for OTP request : rate limiting
					otpRequestRateLimit := viper.GetInt("OTP_REQUEST_RATE_LIMIT")
					otpRequestRateWindow := viper.GetDuration("OTP_REQUEST_RATE_WINDOW")
					if otpRequestRateWindow == 0 {
						otpRequestRateWindow = 24 * time.Hour
					}

					// Retrieve config values for OTP input : attempts limiting
					otpInputAttemptLimit := viper.GetInt("OTP_INPUT_ATTEMPT_LIMIT")
					otpInputAttemptWindow := viper.GetDuration("OTP_INPUT_ATTEMPT_WINDOW")
					if otpInputAttemptWindow == 0 {
						otpInputAttemptWindow = 1 * time.Hour
					}

					r.With(httprate.LimitByIP(otpRequestRateLimit, otpRequestRateWindow)).Post("/otp", handlers.GenerateChangePasswordOTP)
					r.With(httprate.LimitByIP(otpInputAttemptLimit, otpInputAttemptWindow)).Post("/otp/validate", handlers.ValidateChangePasswordOTP)
					r.Put("/", handlers.SubmitChangePassword)
				})
			})

			// User specific
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetUser)
			})
		})

		// Security
		r.Route("/security", func(r chi.Router) {

			// Permission
			r.Route("/permission", func(r chi.Router) {
				r.Post("/", handlers.CreatePermission)
				r.Get("/", handlers.ListPermissions)

				// Permission specific
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", handlers.GetPermission)
					r.Put("/", handlers.UpdatePermission)
					r.Delete("/", handlers.DeletePermission)
				})
			})

			// Role
			r.Route("/role", func(r chi.Router) {
				r.Post("/", handlers.CreateRole)
				r.Get("/", handlers.ListRoles)

				// Role specific
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", handlers.GetRole)
					r.Put("/", handlers.UpdateRole)
					r.Delete("/", handlers.DeleteRole)

					// Permissions
					r.Route("/permission", func(r chi.Router) {
						r.Get("/", handlers.GetRolePermissions)
						r.Put("/", handlers.SetRolePermissions)
					})

					// Users
					r.Route("/user", func(r chi.Router) {
						r.Get("/", handlers.ListUsersForRole)
						r.Put("/", handlers.AddUsersToRole)
						r.Delete("/", handlers.RemoveUsersFromRole)
					})
				})

				// User roles specific
				r.Route("/user", func(r chi.Router) {
					r.Get("/", handlers.ListUsersWithRoles)

					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", handlers.ListRolesWithPermissionsForUser)
						r.Put("/", handlers.SetRolesForUser)
					})
				})
			})
		})

		// Broker
		r.Route("/broker", func(r chi.Router) {
			r.Post("/", handlers.CreateBroker)
			r.Get("/", handlers.ListBrokers)

			// Broker specific
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

				// User specific : retrieving userID through context
				r.Delete("/user", handlers.DeleteUserBroker)
			})

			// Broker user specific : retrieving userID through context
			r.Route("/user", func(r chi.Router) {
				r.Post("/", handlers.CreateUserBroker)
				r.Get("/", handlers.ListUserBrokers)
			})
		})

		// Transaction : retrieving userID through context
		r.Route("/transaction", func(r chi.Router) {
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
