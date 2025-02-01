package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/handlers/render"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateRole godoc
//
//	@Id				CreateRole
//
//	@Summary		Create a new role
//	@Description	Create a new role. (Permission: <b>admin.roles.create</b>)
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Param			role	body	roles.RoleWithPermissions	true	"role (json)"
//	@Security		Bearer
//	@Success		200	{object}	roles.RoleWithPermissions				"role"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles [post]
func CreateRole(w http.ResponseWriter, r *http.Request) {
	if !checkPermission(w, r, "admin.roles.create") {
		return
	}

	var role roles.RoleWithPermissions
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		zap.L().Warn("Role json decode", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	if ok, err := role.IsValid(); !ok {
		zap.L().Warn("Role is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	var perms []uuid.UUID

	// Enable permission update
	if checkPermission(w, r, "admin.roles.permissions.update") {

		if ok, err := role.Permissions.IsValid(); !ok {
			zap.L().Warn("Permissions are not valid", zap.Error(err))
			render.BadRequest(w, r, err)
			return
		}

		perms = role.Permissions.GetUUIDs()
	}

	roleID, err := roles.R().Create(role.Role, perms)
	if err != nil {
		zap.L().Error("Cannot create role", zap.Error(err))
		render.Error(w, r, err, "Create role")
		return
	}

	newRole, found, err := roles.R().GetWithPermissions(roleID)
	if err != nil {
		zap.L().Error("Cannot get role", zap.String("uuid", roleID.String()), zap.Error(err))
		render.Error(w, r, nil, "")
		return
	}
	if !found {
		zap.L().Error("Role not found after creation", zap.String("uuid", roleID.String()))
		render.Error(w, r, nil, "")
		return
	}

	render.JSON(w, r, newRole)
}

// GetRole godoc
//
//	@Id				GetRole
//
//	@Summary		Get a role
//	@Description	Get a role by id. (Permission: <b>admin.roles.read</b>)
//	@Tags			Roles
//	@Produce		json
//	@Param			id	path	string	true	"role id"
//	@Security		Bearer
//	@Success		200	{object}	roles.RoleWithPermissions "role"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		404	{string}	string					"Role not found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id} [get]
func GetRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	roleID, err := uuid.Parse(id)
	if err != nil {
		zap.L().Warn("Parse role id", zap.Error(err))
		render.BadRequest(w, r, fmt.Errorf("invalid role id"))
		return
	}

	if !checkPermission(w, r, "admin.roles.read") {
		return
	}

	role, found, err := roles.R().GetWithPermissions(roleID)
	if err != nil {
		zap.L().Error("Cannot load role", zap.String("uuid", roleID.String()), zap.Error(err))
		render.Error(w, r, nil, "")
		return
	}

	if !found {
		zap.L().Debug("Role not found", zap.String("uuid", roleID.String()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	render.JSON(w, r, role)
}

// GetRoles godoc
//
//	@Id				GetRoles
//
//	@Summary		Get all roles
//	@Description	Gets a list of all roles. (Permission: <b>admin.roles.list</b>)
//	@Tags			Roles
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{array}		roles.RoleWithPermissions				"list of roles"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles [get]
func GetRoles(w http.ResponseWriter, r *http.Request) {
	if !checkPermission(w, r, "admin.roles.list") {
		return
	}

	result, err := roles.R().GetAllWithPermissions()
	if err != nil {
		render.Error(w, r, err, "Get roles")
		return
	}

	render.JSON(w, r, result)
}

// UpdateRole godoc
//
//	@Id				UpdateRole
//
//	@Summary		Update role
//	@Description	Updates the role. (Permission: <b>admin.roles.update</b>)
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string		true	"role ID"
//	@Param			role	body	roles.RoleWithPermissions	true	"role (json)"
//	@Security		Bearer
//	@Success		200	{object}	roles.RoleWithPermissions				"role"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id} [put]
func UpdateRole(w http.ResponseWriter, r *http.Request) {
	roleID, ok := parseParamUUID(w, r, "id")
	if !ok {
		return
	}

	if !checkPermission(w, r, "admin.roles.update") {
		return
	}

	var role roles.RoleWithPermissions
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		zap.L().Warn("Role json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}

	if ok, err := role.IsValid(); !ok {
		render.BadRequest(w, r, err)
		return
	}

	role.Id = roleID
	var perms []uuid.UUID

	// Enable permission update
	if checkPermission(w, r, "admin.roles.permissions.update") {

		if ok, err := role.Permissions.IsValid(); !ok {
			zap.L().Warn("Permissions are not valid", zap.Error(err))
			render.BadRequest(w, r, err)
			return
		}

		perms = role.Permissions.GetUUIDs()
	}

	err = roles.R().Update(role.Role, perms)
	if err != nil {
		zap.L().Error("Cannot update role", zap.Error(err))
		render.Error(w, r, err, "Update role")
		return
	}

	role, found, err := roles.R().GetWithPermissions(roleID)
	if err != nil {
		zap.L().Error("Cannot get role", zap.String("uuid", roleID.String()), zap.Error(err))
		render.Error(w, r, nil, "")
		return
	}
	if !found {
		zap.L().Error("Role not found after update", zap.String("uuid", roleID.String()))
		render.Error(w, r, nil, "")
		return
	}

	render.JSON(w, r, role)
}

// DeleteRole godoc
//
//	@Id				DeleteRole
//
//	@Summary		Delete role
//	@Description	Deletes a role. (Permission: <b>admin.roles.delete</b>)
//	@Tags			Roles
//	@Produce		json
//	@Param			id	path	string	true	"role ID"
//	@Security		Bearer
//	@Success		200	{string}	string					"status OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id} [delete]
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	roleID, ok := parseParamUUID(w, r, "id")
	if !ok || !checkPermission(w, r, "admin.roles.delete") {
		return
	}

	err := roles.R().Delete(roleID)
	if err != nil {
		render.Error(w, r, err, "Cannot delete role")
		return
	}

	render.OK(w, r)
}

// GetRolePermissions godoc
//
//	@Id				GetRolePermissions
//
//	@Summary		Get all permissions for a specified role id
//	@Description	Gets a list of all role permissions. (Permission: <b>admin.roles.permissions.list</b>)
//	@Tags			Roles, RolePermissions
//	@Produce		json
//	@Param			id	path	string	true	"role ID"
//	@Security		Bearer
//	@Success		200	{array}		permissions.Permission	"list of permissions"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id}/permissions [get]
func GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	roleId, ok := parseParamUUID(w, r, "id")
	if !ok || !checkPermission(w, r, "admin.roles.permissions.list") {
		return
	}

	perms, err := permissions.R().GetAllByRoleId(roleId)
	if err != nil {
		render.Error(w, r, err, "Get role permissions")
		return
	}

	render.JSON(w, r, perms)
}

// SetRolePermissions godoc
//
//	@Id				SetRolePermissions
//
//	@Summary		Set permissions for a given role
//	@Description	Updates the role. (Permission: <b>admin.roles.permissions.update</b>)
//	@Tags			Roles, RolePermissions
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string						true	"role ID"
//	@Param			role	body	permissions.Permissions	true	"List of permissions (json)"
//	@Security		Bearer
//	@Success		200	{object}	roles.RoleWithPermissions	"role"
//	@Failure		400	{object}	render.ErrorResponse		"Bad Request"
//	@Failure		401	{string}	string						"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse		"Internal Server Error"
//	@Router			/api/v1/roles/{id}/permissions [put]
func SetRolePermissions(w http.ResponseWriter, r *http.Request) {
	roleId, ok := parseParamUUID(w, r, "id")
	if !ok || !checkPermission(w, r, "admin.roles.permissions.update") {
		return
	}

	var perms permissions.Permissions
	err := json.NewDecoder(r.Body).Decode(&perms)
	if err != nil {
		zap.L().Warn("Permissions json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}

	if ok, err := perms.IsValid(); !ok {
		zap.L().Warn("Permissions are not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	err = roles.R().SetRolePermissions(roleId, perms.GetUUIDs())
	if err != nil {
		zap.L().Error("Failed to set permissions", zap.Error(err))
		render.Error(w, r, err, "Set role permissions")
		return
	}

	render.OK(w, r)
}

// GetRoleUsers godoc
//
//	@Id				GetRoleUsers
//
//	@Summary		Get all users for a specified role id
//	@Description	Gets a list of all role users. (Permission: <b>admin.roles.users.list</b>)
//	@Tags			Roles, RoleUsers
//	@Produce		json
//	@Param			id	path	string	true	"role ID"
//	@Security		Bearer
//	@Success		200	{array}		users.User				"list of users"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id}/users [get]
func GetRoleUsers(w http.ResponseWriter, r *http.Request) {
	roleId, ok := parseParamUUID(w, r, "id")
	if !ok || !checkPermission(w, r, "admin.roles.users.list") {
		return
	}

	roleUsers, err := users.R().GetUsersByRoleID(roleId)
	if err != nil {
		render.Error(w, r, err, "Get role users")
		return
	}

	render.JSON(w, r, roleUsers)
}

// PutUsersRole godoc
//
//	@Id				PutUsersRole
//
//	@Summary		Add users to a given role
//	@Description	Updates the role. (Permission: <b>admin.roles.users.update</b>)
//	@Tags			Roles, RoleUsers
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string		true	"role ID"
//	@Param			role	body	[]string	true	"List of user UUIDs (json)"
//	@Security		Bearer
//	@Success		200	{string}	string					"status OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id}/users [put]
func PutUsersRole(w http.ResponseWriter, r *http.Request) {
	roleId, ok := parseParamUUID(w, r, "id")
	if !ok || !checkPermission(w, r, "admin.roles.users.update") {
		return
	}

	var userUUIDs []uuid.UUID
	err := json.NewDecoder(r.Body).Decode(&userUUIDs)
	if err != nil {
		zap.L().Warn("User UUIDs json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}

	if len(userUUIDs) == 0 {
		render.BadRequest(w, r, fmt.Errorf("empty user list"))
		return
	}

	err = users.R().AddUsersRole(userUUIDs, roleId)
	if err != nil {
		render.Error(w, r, err, "Set user roles")
		return
	}

	render.OK(w, r)
}

// DeleteUsersRole godoc
//
//	@Id				DeleteUsersRole
//
//	@Summary		Remove users from a given role
//	@Description	Updates the role. (Permission: <b>admin.roles.users.delete</b>)
//	@Tags			Roles, RoleUsers
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string		true	"role ID"
//	@Param			role	body	[]string	true	"List of user UUIDs (json)"
//	@Security		Bearer
//	@Success		200	{string}	string					"status OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id}/users [delete]
func DeleteUsersRole(w http.ResponseWriter, r *http.Request) {
	roleId, ok := parseParamUUID(w, r, "id")
	if !ok || !checkPermission(w, r, "admin.roles.users.delete") {
		return
	}

	var userUUIDs []uuid.UUID
	err := json.NewDecoder(r.Body).Decode(&userUUIDs)
	if err != nil {
		zap.L().Warn("User UUIDs json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}

	if len(userUUIDs) == 0 {
		render.BadRequest(w, r, fmt.Errorf("empty user list"))
		return
	}

	err = users.R().RemoveUsersRole(userUUIDs, roleId)
	if err != nil {
		render.Error(w, r, err, "Delete user roles")
		return
	}

	render.OK(w, r)
}
