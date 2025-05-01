package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"go.uber.org/zap"
	"io"
	"net/http"
)

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
//	@Success		200	{object}	string			"jwt token"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/token [post]
func GetToken(w http.ResponseWriter, r *http.Request) {
	var creds models.UserWithPassword

	var (
		ErrLoginInvalid = errors.New("login-invalid")
	)

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil || json.Unmarshal(body, &creds) != nil {
		zap.L().Error("GetToken: invalid input", zap.Error(err))
		render.BadRequest(w, r, ErrLoginInvalid)
		return
	}

	// Validate the credentials and return the token
	response, err := clients.C().Auth().GenerateToken(r.Context(), &protogen.GenerateTokenRequest{
		Email:    creds.Email,
		Password: creds.Password,
	})
	if err != nil {
		zap.L().Warn("GetToken: failed", zap.Error(err))
		render.BadRequest(w, r, ErrLoginInvalid)
		return
	}

	render.JSON(w, r, response.Token)
}
