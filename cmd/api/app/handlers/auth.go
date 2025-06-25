package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
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
	response, err := clients.C().Auth().GenerateToken(r.Context(), &authpb.GenerateTokenRequest{
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

func generateOTP(w http.ResponseWriter, r *http.Request, purpose authpb.OtpPurpose) {
	// Parse request body
	var requestUserOtp models.RequestUserOtp
	err := json.NewDecoder(r.Body).Decode(&requestUserOtp)
	if err != nil {
		zap.L().Warn("RequestUserOtp json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve user language from query parameters
	userLanguage := U().ParseParamLanguage(w, r)

	// Generate OTP
	response, err := clients.C().Auth().GenerateOTP(r.Context(), &authpb.GenerateOTPRequest{
		Email:    requestUserOtp.Email,
		Language: userLanguage.String(),
		Purpose:  purpose,
	})
	if err != nil {
		// TODO : differentiate between otp already exist and other errors
		zap.L().Error("Generate OTP", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Return user ID and expires_at in JSON response
	render.JSON(w, r, models.ResponseUserOtp{
		ExpiresAt:  response.GetExpiresAt().AsTime(),
		Identifier: response.GetIdentifier(),
	})
}

func validateOTP(w http.ResponseWriter, r *http.Request, purpose authpb.OtpPurpose) {
	// Parse request body
	var validateUserOtp models.ValidateUserOtp
	err := json.NewDecoder(r.Body).Decode(&validateUserOtp)
	if err != nil {
		zap.L().Warn("ValidateUserOtp json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate OTP
	response, err := clients.C().Auth().ValidateOTP(r.Context(), &authpb.ValidateOTPRequest{
		Identifier: validateUserOtp.UserID.String(),
		Otp:        validateUserOtp.Otp,
		Purpose:    purpose,
	})
	if err != nil {
		zap.L().Error("Validate OTP", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.JSON(w, r, response.GetRequestId())
}

// GenerateSignupOTP godoc
//
//	@Id				GenerateSignupOTP
//
//	@Summary		Starts verification process for new user
//	@Description	Handles OTP for the user with the provided email.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			lang	query	string					false	"Language code"
//	@Param			request	body	models.RequestUserOtp	true	"request (json)"
//	@Success		200	{object}	models.ResponseUserOtp	"ResponseUserOtp"
//	@Failure		400	{object}	render.ErrorResponse	"Bad RequestUserOtp"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/verify/otp [post]
func GenerateSignupOTP(w http.ResponseWriter, r *http.Request) {
	generateOTP(w, r, authpb.OtpPurpose_USER_SIGNUP)
}

// GenerateForgottenPasswordOTP godoc
//
//	@Id				GenerateForgottenPasswordOTP
//
//	@Summary		Starts password reset process for user
//	@Description	Handles OTP for the user with the provided email.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			lang	query	string					false	"Language code"
//	@Param			request	body	models.RequestUserOtp	true	"request (json)"
//	@Success		200	{object}	models.ResponseUserOtp	"ResponseUserOtp"
//	@Failure		400	{object}	render.ErrorResponse	"Bad RequestUserOtp"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/password/reset/otp [post]
func GenerateForgottenPasswordOTP(w http.ResponseWriter, r *http.Request) {
	generateOTP(w, r, authpb.OtpPurpose_PASSWORD_RESET)
}

// GenerateChangePasswordOTP godoc
//
//	@Id				GenerateChangePasswordOTP
//
//	@Summary		Starts password change process for user
//	@Description	Handles OTP for the user with the provided email.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			lang	query	string					false	"Language code"
//	@Param			request	body	models.RequestUserOtp	true	"request (json)"
//	@Success		200	{object}	models.ResponseUserOtp	"ResponseUserOtp"
//	@Failure		400	{object}	render.ErrorResponse	"Bad RequestUserOtp"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/password/change/otp [post]
func GenerateChangePasswordOTP(w http.ResponseWriter, r *http.Request) {
	generateOTP(w, r, authpb.OtpPurpose_PASSWORD_CHANGE)
}

// ValidateSignupOTP godoc
//
//	@Id				ValidateSignupOTP
//
//	@Summary		Validates OTP verification process for new user and activates account
//	@Description	Handles OTP for the user with the provided email.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.ValidateUserOtp	true	"request (json)"
//	@Success		200	{object}	string					"request ID"
//	@Failure		400	{object}	render.ErrorResponse	"Bad ValidateUserOtp"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/verify/otp/validate [post]
func ValidateSignupOTP(w http.ResponseWriter, r *http.Request) {
	validateOTP(w, r, authpb.OtpPurpose_USER_SIGNUP)
}

// ValidateForgottenPasswordOTP godoc
//
//	@Id				ValidateForgottenPasswordOTP
//
//	@Summary		Validates password reset OTP process for user
//	@Description	Handles OTP for the user with the provided email.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.ValidateUserOtp	true	"request (json)"
//	@Success		200	{object}	string					"request ID"
//	@Failure		400	{object}	render.ErrorResponse	"Bad ValidateUserOtp"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/password/reset/otp/validate [post]
func ValidateForgottenPasswordOTP(w http.ResponseWriter, r *http.Request) {
	validateOTP(w, r, authpb.OtpPurpose_PASSWORD_RESET)
}

// ValidateChangePasswordOTP godoc
//
//	@Id				ValidateChangePasswordOTP
//
//	@Summary		Validates password change OTP process for user
//	@Description	Handles OTP for the user with the provided email.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.ValidateUserOtp	true	"request (json)"
//	@Success		200	{object}	string					"request ID"
//	@Failure		400	{object}	render.ErrorResponse	"Bad ValidateUserOtp"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/password/change/otp/validate [post]
func ValidateChangePasswordOTP(w http.ResponseWriter, r *http.Request) {
	validateOTP(w, r, authpb.OtpPurpose_PASSWORD_CHANGE)
}

// ResetForgottenPassword
//
// @Id ResetForgottenPassword
//
// @Summary Resets the user password
// @Description Executes it as part of forgotten password OTP auth process
//
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.UserInputResetPassword	true	"request (json)"
//	@Success		200	{object}	string					"OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad UserInputResetPassword"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/password/reset [put]
func ResetForgottenPassword(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var inputResetPassword models.UserInputResetPassword
	err := json.NewDecoder(r.Body).Decode(&inputResetPassword)
	if err != nil {
		zap.L().Warn("UserInputResetPassword json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Completes request and updates the user password
	_, err = clients.C().Auth().ResetForgottenPassword(r.Context(), &authpb.ResetForgottenPasswordRequest{
		RequestId:    inputResetPassword.OtpRequestID.String(),
		UserId:       inputResetPassword.UserID.String(),
		Password:     inputResetPassword.Password,
		Confirmation: inputResetPassword.Confirmation,
	})
	if err != nil {
		zap.L().Error("ResetForgottenPassword", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.OK(w, r)
}

// SubmitChangePassword
//
// @Id SubmitChangePassword
//
// @Summary Resets the user password
// @Description Executes it as part of user settings change OTP auth process
//
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.UserInputChangePassword	true	"request (json)"
//	@Success		200	{object}	string					"OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad UserInputChangePassword"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/password/change [put]
func SubmitChangePassword(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var inputChangePassword models.UserInputChangePassword
	err := json.NewDecoder(r.Body).Decode(&inputChangePassword)
	if err != nil {
		zap.L().Warn("UserInputChangePassword json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the user password
	_, err = clients.C().Auth().UpdatePassword(r.Context(), &authpb.UpdatePasswordRequest{
		RequestId:    inputChangePassword.OtpRequestID.String(),
		Password:     inputChangePassword.Password,
		Confirmation: inputChangePassword.Confirmation,
	})
	if err != nil {
		zap.L().Error("UpdatePassword", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.OK(w, r)
}

// RegisterUser
//
// @Id RegisterUser
//
// @Summary Registers the user
// @Description Executes it as part of user creation OTP auth process
//
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.UserInputCreate	true	"request (json)"
//	@Success		200	{object}	string					"OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad UserInputCreate"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/auth/register [post]
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var userInputCreate models.UserInputCreate
	err := json.NewDecoder(r.Body).Decode(&userInputCreate)
	if err != nil {
		zap.L().Warn("User json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the user password
	_, err = clients.C().Auth().CreateUser(r.Context(), &authpb.CreateUserRequest{
		Email:        userInputCreate.Email,
		Password:     userInputCreate.Password,
		Confirmation: userInputCreate.Confirmation,
		Checkbox:     userInputCreate.Checkbox,
	})
	if err != nil {
		zap.L().Error("CreateUser", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.OK(w, r)
}
