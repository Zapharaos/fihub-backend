package auth

import (
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/handlers/render"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

const (
	CheckHeader = 1 << iota
	CheckQuery
)

type JwtToken struct {
	Token string `json:"token"`
}

type Auth struct {
	signingKey []byte
	checks     int8
}

// New initialize a new instance of Auth and returns a pointer of it
// The signing key is generated randomly and is used to sign the JWT tokens
// The checks parameter is a bitfield to enable or disable checks
// how to use it: auth.NewAuth(auth.CheckHeader | auth.CheckQuery)
func New(checks int8) *Auth {
	if checks == 0 {
		zap.L().Fatal("no checks are enabled")
	}
	return &Auth{
		signingKey: []byte(utils.RandString(128)),
		checks:     checks,
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
//	@Param			user	body	users.UserWithPassword	true	"login & user (json)"
//	@Security		Bearer
//	@Success		200	{object}	auth.JwtToken			"event"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/token [post]
func (a *Auth) GetToken(w http.ResponseWriter, r *http.Request) {
	var userCredentials users.UserWithPassword

	body, err := io.ReadAll(r.Body)
	if err != nil {
		zap.L().Error("GetToken.ReadAll:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &userCredentials)
	if err != nil {
		zap.L().Error("GetToken.Decode:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := users.R().Authenticate(userCredentials.Email, userCredentials.Password)
	if err != nil {
		zap.L().Warn("GetToken.Authenticate", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
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
func (a *Auth) GenerateToken(user *users.User) (JwtToken, error) {
	claims := &jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		"iat": jwt.NewNumericDate(time.Now()),
		"nbf": jwt.NewNumericDate(time.Now()),
		"id":  user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our signing key
	tokenString, err := token.SignedString(a.signingKey)
	if err != nil {
		return JwtToken{}, err
	}

	return JwtToken{Token: tokenString}, nil
}

// TODO : Middleware
