package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/auth/password"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/handlers/render"
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

// CreatePasswordResetRequest godoc
//
//	@Id				CreatePasswordResetRequest
//
//	@Summary		Request a password reset
//	@Description	Requests a password reset for the user with the provided email.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	password.InputRequest	true	"request (json)"
//	@Success		200	{string}	string	"status OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/password [post]
func CreatePasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var inputRequest password.InputRequest
	err := json.NewDecoder(r.Body).Decode(&inputRequest)
	if err != nil {
		zap.L().Warn("Request json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve user and validate email
	user, exists, err := users.R().GetByEmail(inputRequest.Email)
	if err != nil {
		zap.L().Error("Check resetPassword exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		zap.L().Warn("ResetPassword not found", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if user already has a token and that it is valid
	exists, err = password.R().ValidForUser(user.ID)
	if err != nil {
		zap.L().Error("Check resetPassword exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists {
		zap.L().Warn("ResetPassword already exists", zap.Error(err))
		render.BadRequest(w, r, errors.New("request-active"))
		return
	}

	// TODO : check for rate limiting

	// Create request
	request := password.InitRequest(user.ID)

	// Store request
	err = password.R().Create(request)
	if err != nil {
		zap.L().Error("RequestPasswordReset", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Prepare email
	subject := "Your OTP Code"
	plainTextContent := "Your OTP code is: " + request.Token
	htmlContent := "<p>Your OTP code is: <strong>" + request.Token + "</strong></p>"

	// Send email
	err = email.S().Send(user.Email, subject, plainTextContent, htmlContent)
	if err != nil {
		zap.L().Error("Failed to send OTP email", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.OK(w, r)
}

// GetPasswordResetRequestID godoc
//
//	@Id				GetPasswordResetRequestID
//
//	@Summary		Get password reset request ID
//	@Description	Returns the request ID for a given user ID and token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string	true	"User ID"
//	@Param			token	body	password.InputToken	true	"token (json)"
//	@Success		200	{string}	string	"request_id"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/{id}/password/{token} [get]
func GetPasswordResetRequestID(w http.ResponseWriter, r *http.Request) {
	userID, ok := ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Parse request body
	var inputToken password.InputToken
	err := json.NewDecoder(r.Body).Decode(&inputToken)
	if err != nil {
		zap.L().Warn("Request json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if request exists and is valid
	requestID, err := password.R().GetRequestID(userID, inputToken.Token)
	if err != nil {
		zap.L().Error("GetPasswordResetRequestID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if requestID == uuid.Nil {
		zap.L().Warn("ResetPassword request not found", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, requestID)
}

// ResetPassword godoc
//
//	@Id				ResetPassword
//
//	@Summary		Reset the user's password
//	@Description	Resets the user's password using the provided token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			id				path	string					true	"User ID"
//	@Param			request_id		path	string					true	"Reset token"
//	@Param			password		body	users.UserInputPassword	true	"password (json)"
//	@Security		Bearer
//	@Success		200	{string}	string					"status OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/{id}/password/{request_id} [put]
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	userID, requestID, ok := parseUUIDPair(w, r, "request_id")
	if !ok {
		return
	}

	// Check if request exists and is valid
	exists, err := password.R().Valid(userID, requestID)
	if err != nil {
		zap.L().Error("Check resetPassword exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		zap.L().Warn("ResetPassword not found", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse request body
	var userPassword users.UserInputPassword
	err = json.NewDecoder(r.Body).Decode(&userPassword)
	if err != nil {
		zap.L().Warn("Reset json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate password
	if ok, err := userPassword.IsValidPassword(); !ok {
		zap.L().Warn("ResetPassword", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Convert to UserWithPassword
	userWithPassword := userPassword.ToUserWithPassword()
	userWithPassword.ID = userID

	// Reset password
	err = users.R().UpdateWithPassword(userWithPassword)
	if err != nil {
		zap.L().Error("PutUser.UpdatePassword", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.OK(w, r)
}
