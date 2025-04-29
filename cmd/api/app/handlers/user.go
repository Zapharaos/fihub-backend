package handlers

import (
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"go.uber.org/zap"
	"net/http"
)

// CreateUser godoc
//
//	@Id				CreateUser
//
//	@Summary		Create a new user
//	@Description	Create a new user.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body	models.UserInputCreate	true	"user (json)"
//	@Security		Bearer
//	@Success		200	{object}	models.User				"user"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/user [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var userInputCreate models.UserInputCreate
	err := json.NewDecoder(r.Body).Decode(&userInputCreate)
	if err != nil {
		zap.L().Warn("User json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Map UserInputCreate to gRPC CreateUserRequest
	createUserRequest := &protogen.CreateUserRequest{
		Email:        userInputCreate.Email,
		Password:     userInputCreate.Password,
		Confirmation: userInputCreate.Confirmation,
		Checkbox:     userInputCreate.Checkbox,
	}

	// Create user
	createUserResponse, err := clients.C().User().CreateUser(r.Context(), createUserRequest)
	if err != nil {
		zap.L().Error("Create user", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the models.User struct
	user, err := models.FromProtogenUser(createUserResponse.User)
	if err != nil {
		zap.L().Error("Bad protogen user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, user)
}

// GetUser godoc
//
//	@Id				GetUser
//
//	@Summary		Get a user by ID
//	@Description	Get a user by ID. (Permission: <b>admin.users.read</b>)
//	@Tags			User
//	@Produce		json
//	@Param			id	path	string	true	"user ID"
//	@Security		Bearer
//	@Success		200	{object}	models.User				"user"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		404	{string}	string					"User not found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/user/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Get user
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

	render.JSON(w, r, user)
}

// GetUserSelf godoc
//
//	@Id				GetUserSelf
//
//	@Summary		Get the currently authenticated user
//	@Description	Retrieves the currently authenticated user.
//	@Tags			User
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{object}	models.User				"user"
//	@Failure		400	{string}	string					"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/user/me [get]
func GetUserSelf(w http.ResponseWriter, r *http.Request) {
	userCtx, found := U().GetUserFromContext(r)
	if !found {
		zap.L().Debug("No context user provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, userCtx)
}

// UpdateUserSelf godoc
//
//	@Id				UpdateUserSelf
//
//	@Summary		Update the currently authenticated user
//	@Description	Updates the currently authenticated user.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body	models.User	true	"user (json)"
//	@Security		Bearer
//	@Success		200	{object}	models.User				"user"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/user/me [put]
func UpdateUserSelf(w http.ResponseWriter, r *http.Request) {
	userCtx, found := U().GetUserFromContext(r)
	if !found {
		zap.L().Debug("No context user provided")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Parse request body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		zap.L().Warn("User json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Map User to gRPC UpdateUserRequest
	updateUserRequest := &protogen.UpdateUserRequest{
		Id:    userCtx.ID.String(),
		Email: user.Email,
	}

	// Update user
	updateUserResponse, err := clients.C().User().UpdateUser(r.Context(), updateUserRequest)
	if err != nil {
		zap.L().Error("Update user", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the models.User struct
	user, err = models.FromProtogenUser(updateUserResponse.User)
	if err != nil {
		zap.L().Error("Bad protogen user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, user)
}

// UpdateUserPassword godoc
//
//	@Id				UpdateUserPassword
//
//	@Summary		Update the password of the currently authenticated user
//	@Description	Update the password of the currently authenticated user.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			password	body	models.UserInputPassword	true	"password (json)"
//	@Security		Bearer
//	@Success		200	{string}	string					"status OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/user/me/password [put]
func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	userCtx, found := U().GetUserFromContext(r)
	if !found {
		zap.L().Debug("No context user provided")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Parse request body
	var userPassword models.UserInputPassword
	err := json.NewDecoder(r.Body).Decode(&userPassword)
	if err != nil {
		zap.L().Warn("User json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Map UserInputPassword to gRPC UpdateUserRequest
	updateUserRequest := &protogen.UpdateUserPasswordRequest{
		Id:           userCtx.ID.String(),
		Password:     userPassword.Password,
		Confirmation: userPassword.Confirmation,
	}

	// Update user password
	_, err = clients.C().User().UpdateUserPassword(r.Context(), updateUserRequest)
	if err != nil {
		zap.L().Error("Update user password", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.OK(w, r)
}

// DeleteUserSelf godoc
//
//	@Id				DeleteUserSelf
//
//	@Summary		Delete the currently authenticated user
//	@Description	Deletes the currently authenticated user.
//	@Tags			User
//	@Security		Bearer
//	@Success		200	{string}	string					"status OK"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/user/me [delete]
func DeleteUserSelf(w http.ResponseWriter, r *http.Request) {
	userCtx, found := U().GetUserFromContext(r)
	if !found {
		zap.L().Debug("No context user provided")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Delete user
	_, err := clients.C().User().DeleteUser(r.Context(), &protogen.DeleteUserRequest{
		Id: userCtx.ID.String(),
	})
	if err != nil {
		zap.L().Error("Delete user", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.OK(w, r)
}
