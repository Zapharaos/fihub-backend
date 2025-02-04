package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/auth"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/handlers/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

// CreateUser godoc
//
//	@Id				CreateUser
//
//	@Summary		Create a new user
//	@Description	Create a new user.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body	users.UserInputCreate	true	"user (json)"
//	@Security		Bearer
//	@Success		200	{object}	users.User				"user"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var userInputCreate users.UserInputCreate
	err := json.NewDecoder(r.Body).Decode(&userInputCreate)
	if err != nil {
		zap.L().Warn("User json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate user
	if ok, err := userInputCreate.IsValid(); !ok {
		zap.L().Warn("User is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Convert to UserWithPassword
	userWithPassword := userInputCreate.UserWithPassword

	// Verify user existence
	exists, err := users.R().Exists(userWithPassword.Email)
	if err != nil {
		zap.L().Error("Check user exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists {
		zap.L().Warn("User already exists", zap.String("email", userWithPassword.Email))
		render.BadRequest(w, r, fmt.Errorf("email-used"))
		return
	}

	// Create user
	userID, err := users.R().Create(userWithPassword)
	if err != nil {
		zap.L().Error("PostUser.Create", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get user back from database
	user, found, err := users.R().Get(userID)
	if err != nil {
		zap.L().Error("Cannot get user", zap.String("uuid", userID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		zap.L().Error("User not found after creation", zap.String("uuid", userID.String()))
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
//	@Tags			Users
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{object}	users.UserWithRoles		"user"
//	@Failure		400	{string}	string					"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users/me [get]
func GetUserSelf(w http.ResponseWriter, r *http.Request) {
	userCtx, found := getUserFromContext(r)
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
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body	users.User	true	"user (json)"
//	@Security		Bearer
//	@Success		200	{object}	users.User				"user"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users/me [put]
func UpdateUserSelf(w http.ResponseWriter, r *http.Request) {
	userCtx, found := getUserFromContext(r)
	if !found {
		zap.L().Debug("No context user provided")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Parse request body
	var user users.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		zap.L().Warn("User json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate user
	if ok, err := user.IsValid(); !ok {
		zap.L().Warn("User is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Replace ID with the one from the context
	user.ID = userCtx.ID

	// Update user
	err = users.R().Update(user)
	if err != nil {
		zap.L().Error("PutUser.Update", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get user back from database
	user, found, err = users.R().Get(userCtx.ID)
	if err != nil {
		zap.L().Error("Cannot get user", zap.String("uuid", userCtx.ID.String()), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		zap.L().Error("User not found after update", zap.String("uuid", userCtx.ID.String()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, user)
}

// ChangeUserPassword godoc
//
//	@Id				ChangeUserPassword
//
//	@Summary		Change the password of the currently authenticated user
//	@Description	Changes the password of the currently authenticated user.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			password	body	users.UserInputPassword	true	"password (json)"
//	@Security		Bearer
//	@Success		200	{string}	string					"status OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users/me/password [put]
func ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	userCtx, found := getUserFromContext(r)
	if !found {
		zap.L().Debug("No context user provided")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Parse request body
	var userPassword users.UserInputPassword
	err := json.NewDecoder(r.Body).Decode(&userPassword)
	if err != nil {
		zap.L().Warn("User json decode", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate password
	if ok, err := userPassword.IsValidPassword(); !ok {
		zap.L().Warn("User is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Convert to UserWithPassword
	userWithPassword := userPassword.UserWithPassword
	userWithPassword.ID = userCtx.ID

	// Update password
	err = users.R().UpdateWithPassword(userWithPassword)
	if err != nil {
		zap.L().Error("PutUser.UpdatePassword", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
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
//	@Tags			Users
//	@Security		Bearer
//	@Success		200	{string}	string					"status OK"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users/me [delete]
func DeleteUserSelf(w http.ResponseWriter, r *http.Request) {
	userCtx, found := getUserFromContext(r)
	if !found {
		zap.L().Debug("No context user provided")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := users.R().Delete(userCtx.ID)
	if err != nil {
		zap.L().Error("DeleteUser.Delete", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.OK(w, r)
}

// GetUser godoc
//
//	@Id				GetUser
//
//	@Summary		Get a user by ID
//	@Description	Get a user by ID. (Permission: <b>admin.users.read</b>)
//	@Tags			Users
//	@Produce		json
//	@Param			id	path	string	true	"user ID"
//	@Security		Bearer
//	@Success		200	{object}	users.UserWithRoles		"user"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		404	{string}	string					"User not found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := parseParamUUID(w, r, "id")
	if !ok || !checkPermission(w, r, "admin.users.read") {
		return
	}

	var user users.UserWithRoles
	var found bool
	var err error

	// the user default accessible data
	user.User, found, err = users.R().Get(userId)
	if err != nil {
		zap.L().Error("GetUser.Get", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		zap.L().Warn("User not found", zap.String("uuid", userId.String()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if the user can retrieve the user roles as well
	if checkPermission(w, r, "admin.users.roles.list") {
		roles, err := auth.LoadUserRoles(userId)
		if err != nil {
			zap.L().Error("GetUser.LoadUserRoles", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user.Roles = roles
	}

	render.JSON(w, r, user)
}

// SetUser godoc
//
//	@Id				SetUser
//
//	@Summary		Update a user by ID
//	@Description	Update a user by ID. (Permission: <b>admin.users.update</b>)
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string				true	"user ID"
//	@Param			user	body	users.UserWithRoles	true	"user (json)"
//	@Security		Bearer
//	@Success		200	{object}	users.UserWithRoles		"user"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users/{id} [put]
func SetUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := parseParamUUID(w, r, "id")
	if !ok || !checkPermission(w, r, "admin.users.update") {
		return
	}

	var user users.UserWithRoles
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		zap.L().Warn("User json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}

	user.ID = userId

	err = users.R().UpdateWithRoles(user, user.Roles.GetUUIDs())
	if err != nil {
		zap.L().Error("PutUser.Update", zap.Error(err))
		render.Error(w, r, err, "Update user")
		return
	}

	render.OK(w, r)
}

// SetUserRoles godoc
//
//	@Id				SetUserRoles
//
//	@Summary		Set roles on a user
//	@Description	Set roles on a user. (Permission: <b>admin.users.roles.update</b>)
//	@Tags			Users, UserRoles
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string				true	"user ID"
//	@Param			user	body	users.UserWithRoles	true	"user (json)"
//	@Security		Bearer
//	@Success		200	{object}	users.UserWithRoles		"user"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users/{id}/roles [put]
func SetUserRoles(w http.ResponseWriter, r *http.Request) {
	userId, ok := parseParamUUID(w, r, "id")
	if !ok || !checkPermission(w, r, "admin.users.roles.update") {
		return
	}

	var user users.UserWithRoles
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		zap.L().Warn("User json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}

	roleUUIDs := make([]uuid.UUID, 0)

	for _, role := range user.Roles {
		roleUUIDs = append(roleUUIDs, role.Id)
	}

	err = users.R().SetUserRoles(userId, roleUUIDs)
	if err != nil {
		zap.L().Error("PutUser.Update", zap.Error(err))
		render.Error(w, r, err, "Set roles on user")
		return
	}

	render.OK(w, r)
}

// GetUserRoles godoc
//
//	@Id				GetUserRoles
//
//	@Summary		Get all roles for a specified user id
//	@Description	Gets a list of all roles. (Permission: <b>admin.users.roles.list</b>)
//	@Tags			Users, UserRoles
//	@Produce		json
//	@Param			id	path	string	true	"user ID"
//	@Security		Bearer
//	@Success		200	{array}		roles.RoleWithPermissions	"list of roles"
//	@Failure		400	{object}	render.ErrorResponse		"Bad Request"
//	@Failure		401	{string}	string						"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse		"Internal Server Error"
//	@Router			/api/v1/users/{id}/roles [get]
func GetUserRoles(w http.ResponseWriter, r *http.Request) {
	userId, ok := parseParamUUID(w, r, "id")
	if !ok || !checkPermission(w, r, "admin.users.roles.list") {
		return
	}

	userRolesWithPermissions, err := auth.LoadUserRoles(userId)
	if err != nil {
		render.Error(w, r, err, "Cannot load roles")
		return
	}

	render.JSON(w, r, userRolesWithPermissions)
}

// GetAllUsersWithRoles godoc
//
//	@Id				GetAllUsersWithRoles
//
//	@Summary		Get all users with their roles
//	@Description	Gets a list of all users with their roles. (Permission: <b>admin.users.list</b>)
//	@Tags			Users, UserRoles
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{array}		users.UserWithRoles	"list of users"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users [get]
func GetAllUsersWithRoles(w http.ResponseWriter, r *http.Request) {
	if !checkPermission(w, r, "admin.users.list") {
		return
	}

	usersWithRoles, err := users.R().GetAllWithRoles()
	if err != nil {
		zap.L().Error("GetAllWithRoles", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, usersWithRoles)
}
