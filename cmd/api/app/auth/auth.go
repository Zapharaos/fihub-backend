package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	"github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/app"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"
	"time"
)

var (
	ErrLoginInvalid = errors.New("login-invalid")
)

const (
	CheckHeader = 1 << iota
	CheckQuery
)

type JwtToken struct {
	Token string `json:"token"`
}

type Config struct {
	CORS        bool
	Security    bool
	GatewayMode bool
	AllowOrigin string
}

type Auth struct {
	SigningKey []byte
	Checks     int8
	Config     Config
}

// New initialize a new instance of Auth and returns a pointer of it
// The signing key is generated randomly and is used to sign the JWT tokens
// The checks parameter is a bitfield to enable or disable checks
// how to use it: auth.NewAuth(auth.CheckHeader | auth.CheckQuery)
func New(checks int8, config Config) *Auth {
	if checks == 0 {
		zap.L().Fatal("no checks are enabled")
	}
	return &Auth{
		SigningKey: []byte(utils.RandString(128)),
		Checks:     checks,
		Config:     config,
	}
}

// GetToken godoc
//
//	@Id				GetToken
//
//	@Summary		Get a JWT token (authenticate)
//	@Description	Login and get a JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body	models.UserWithPassword	true	"login & user (json)"
//	@Security		Bearer
//	@Success		200	{object}	auth.JwtToken			"event"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/token [post]
func (a *Auth) GetToken(w http.ResponseWriter, r *http.Request) {
	var userCredentials models.UserWithPassword

	body, err := io.ReadAll(r.Body)
	if err != nil {
		zap.L().Error("GetToken.ReadAll:", zap.Error(err))
		render.BadRequest(w, r, ErrLoginInvalid)
		return
	}

	err = json.Unmarshal(body, &userCredentials)
	if err != nil {
		zap.L().Error("GetToken.Decode:", zap.Error(err))
		render.BadRequest(w, r, ErrLoginInvalid)
		return
	}

	user, found, err := repositories.R().Authenticate(userCredentials.Email, userCredentials.Password)
	if err != nil {
		zap.L().Warn("GetToken.Authenticate", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		zap.L().Warn("GetToken.Authenticate", zap.Error(err))
		render.BadRequest(w, r, ErrLoginInvalid)
		return
	}

	token, err := a.GenerateToken(user)
	if err != nil {
		zap.L().Warn("GetToken.GenerateToken", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, token)
}

// GenerateToken generate a JWT token for a specific user
func (a *Auth) GenerateToken(user models.User) (JwtToken, error) {
	claims := &jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		"iat": jwt.NewNumericDate(time.Now()),
		"nbf": jwt.NewNumericDate(time.Now()),
		"id":  user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our signing key
	tokenString, err := token.SignedString(a.SigningKey)
	if err != nil {
		return JwtToken{}, err
	}

	return JwtToken{Token: tokenString}, nil
}

// ValidateToken validate a JWT token
func (a *Auth) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

// Middleware is a middleware to authenticate and validate JWT tokens
func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiBasePath := viper.GetString("API_BASE_PATH")
		if r.URL.Path == apiBasePath+"/users" && r.Method == "POST" {
			// No need for middleware for user creation
			next.ServeHTTP(w, r)
			return
		}

		if a.Config.GatewayMode {
			// No need for middleware in gateway mode
			next.ServeHTTP(w, r)
			return
		}

		// first check header if enabled
		var tokenString string
		if a.Checks&CheckHeader != 0 {
			tokenString = r.Header.Get("Authorization")
		} else if a.Checks&CheckQuery != 0 {
			tokenString = r.URL.Query().Get("token")
		}

		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := a.ValidateToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// get the id from the claims
		rawUserID, ok := claims["id"]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// parse uuid
		userId, err := uuid.Parse(rawUserID.(string))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Setup metadata for gRPC clients as context
		md := metadata.Pairs("x-user-id", userId.String())
		ctx := metadata.NewOutgoingContext(r.Context(), md)
		r = r.WithContext(ctx)

		// Get user for app context from the database
		userResponse, err := clients.C().User().GetUser(r.Context(), &protogen.GetUserRequest{
			Id: userId.String(),
		})
		if err != nil {
			zap.L().Error("Get user", zap.Error(err))
			render.ErrorCodesCodeToHttpCode(w, r, err)
			return
		}

		// Map the response to the models.User struct
		user, err := models.FromProtogenUser(userResponse.User)
		if err != nil {
			zap.L().Error("Bad protogen user", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(r.Context(), app.ContextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
