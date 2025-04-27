package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	"github.com/Zapharaos/fihub-backend/cmd/security/app/repositories"
	"github.com/Zapharaos/fihub-backend/cmd/user/app/service"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"net/http"

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
//	@Param			role	body	models.RoleWithPermissions	true	"role (json)"
//	@Security		Bearer
//	@Success		200	{object}	models.RoleWithPermissions				"role"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles [post]
func CreateRole(w http.ResponseWriter, r *http.Request) {
	var role models.RoleWithPermissions
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		zap.L().Warn("Role json decode", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Create gRPC protogen.CreateRoleRequest
	roleRequest := &protogen.CreateRoleRequest{
		Name: role.Name,
	}

	// Create the role
	response, err := clients.C().Security().CreateRole(r.Context(), roleRequest)
	if err != nil {
		zap.L().Error("Create Role", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the RoleWithPermissions model
	newRole, err := models.FromProtogenRoleWithPermissions(response.GetRole())
	if err != nil {
		zap.L().Error("Bad protogen role", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
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
//	@Success		200	{object}	models.RoleWithPermissions "role"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		404	{string}	string					"Role not found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id} [get]
func GetRole(w http.ResponseWriter, r *http.Request) {
	roleID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Get the role
	response, err := clients.C().Security().GetRole(r.Context(), &protogen.GetRoleRequest{
		Id: roleID.String(),
	})
	if err != nil {
		zap.L().Error("Get Role", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the RoleWithPermissions model
	role, err := models.FromProtogenRoleWithPermissions(response.GetRole())
	if err != nil {
		zap.L().Error("Bad protogen roles", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
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
//	@Success		200	{array}		models.RoleWithPermissions				"list of roles"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles [get]
func GetRoles(w http.ResponseWriter, r *http.Request) {
	// List the roles
	response, err := clients.C().Security().ListRoles(r.Context(), nil)
	if err != nil {
		zap.L().Error("List Roles", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the RolesWithPermissions model
	roles, err := models.FromProtogenRolesWithPermissions(response.GetRoles())
	if err != nil {
		zap.L().Error("Bad protogen roles", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, roles)
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
//	@Param			role	body	models.RoleWithPermissions	true	"role (json)"
//	@Security		Bearer
//	@Success		200	{object}	models.RoleWithPermissions				"role"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id} [put]
func UpdateRole(w http.ResponseWriter, r *http.Request) {
	roleID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	var role models.RoleWithPermissions
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		zap.L().Warn("Role json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}

	// Create gRPC protogen.UpdateRoleRequest
	roleRequest := &protogen.UpdateRoleRequest{
		Id:   roleID.String(),
		Name: role.Name,
	}

	// Update the role
	response, err := clients.C().Security().UpdateRole(r.Context(), roleRequest)
	if err != nil {
		zap.L().Error("Update Role", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the RoleWithPermissions model
	newRole, err := models.FromProtogenRoleWithPermissions(response.GetRole())
	if err != nil {
		zap.L().Error("Bad protogen role", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, newRole)
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
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id} [delete]
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	roleID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Delete the role
	_, err := clients.C().Security().DeleteRole(r.Context(), &protogen.DeleteRoleRequest{
		Id: roleID.String(),
	})
	if err != nil {
		zap.L().Error("Delete Role", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
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
//	@Success		200	{array}		models.Permission	"list of permissions"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id}/permissions [get]
func GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	roleId, ok := U().ParseParamUUID(w, r, "id")
	if !ok || !U().CheckPermission(w, r, "admin.roles.permissions.list") {
		return
	}

	perms, err := repositories.R().P().GetAllByRoleId(roleId)
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
//	@Param			role	body	models.Permissions	true	"List of permissions (json)"
//	@Security		Bearer
//	@Success		200	{object}	models.RoleWithPermissions	"role"
//	@Failure		400	{object}	render.ErrorResponse		"Bad PasswordRequest"
//	@Failure		401	{string}	string						"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse		"Internal Server Error"
//	@Router			/api/v1/roles/{id}/permissions [put]
func SetRolePermissions(w http.ResponseWriter, r *http.Request) {
	roleId, ok := U().ParseParamUUID(w, r, "id")
	if !ok || !U().CheckPermission(w, r, "admin.roles.permissions.update") {
		return
	}

	var perms models.Permissions
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

	err = repositories.R().R().SetRolePermissions(roleId, perms.GetUUIDs())
	if err != nil {
		zap.L().Error("Failed to set permissions", zap.Error(err))
		render.Error(w, r, err, "Set role permissions")
		return
	}

	render.OK(w, r)
}

// TODO : move to user private

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
//	@Success		200	{array}		models.User				"list of users"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id}/users [get]
/*func GetRoleUsers(w http.ResponseWriter, r *http.Request) {
	roleId, ok := U().ParseParamUUID(w, r, "id")
	if !ok || !U().CheckPermission(w, r, "admin.roles.users.list") {
		return
	}

	roleUsers, err := repositories.R().U().GetUsersByRoleID(roleId)
	if err != nil {
		render.Error(w, r, err, "Get role users")
		return
	}

	render.JSON(w, r, roleUsers)
}*/

// TODO : move to user private

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
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id}/users [put]
/*func PutUsersRole(w http.ResponseWriter, r *http.Request) {
	roleId, ok := U().ParseParamUUID(w, r, "id")
	if !ok || !U().CheckPermission(w, r, "admin.roles.users.update") {
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

	err = repositories.R().U().AddUsersRole(userUUIDs, roleId)
	if err != nil {
		render.Error(w, r, err, "Set user roles")
		return
	}

	render.OK(w, r)
}*/

// TODO : move to user private

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
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles/{id}/users [delete]
/*func DeleteUsersRole(w http.ResponseWriter, r *http.Request) {
	roleId, ok := U().ParseParamUUID(w, r, "id")
	if !ok || !U().CheckPermission(w, r, "admin.roles.users.delete") {
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

	err = repositories.R().U().RemoveUsersRole(userUUIDs, roleId)
	if err != nil {
		render.Error(w, r, err, "Delete user roles")
		return
	}

	render.OK(w, r)
}*/

// TODO : keep but return user with roleUUIDs instead of models.UserWithRoles

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
//	@Success		200	{array}		models.RoleWithPermissions	"list of roles"
//	@Failure		400	{object}	render.ErrorResponse		"Bad PasswordRequest"
//	@Failure		401	{string}	string						"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse		"Internal Server Error"
//	@Router			/api/v1/users/{id}/roles [get]
func GetUserRoles(w http.ResponseWriter, r *http.Request) {
	userId, ok := U().ParseParamUUID(w, r, "id")
	if !ok || !U().CheckPermission(w, r, "admin.users.roles.list") {
		return
	}

	userRolesWithPermissions, err := service.LoadUserRoles(userId)
	if err != nil {
		render.Error(w, r, err, "Cannot load roles")
		return
	}

	render.JSON(w, r, userRolesWithPermissions)
}

// TODO : keep but return user with roleUUIDs instead of models.UserWithRoles

/*// GetAllUsersWithRoles godoc
//
//	@Id				GetAllUsersWithRoles
//
//	@Summary		Get all users with their roles
//	@Description	Gets a list of all users with their roles. (Permission: <b>admin.users.list</b>)
//	@Tags			Users, UserRoles
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{array}		models.UserWithRoles	"list of users"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users [get]
func GetAllUsersWithRoles(w http.ResponseWriter, r *http.Request) {
	if !U().CheckPermission(w, r, "admin.users.list") {
		return
	}

	usersWithRoles, err := repositories.R().GetAllWithRoles()
	if err != nil {
		zap.L().Error("GetAllWithRoles", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, usersWithRoles)
}*/

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
//	@Param			roles	body	[]string			true	"array of role UUIDs"
//	@Security		Bearer
//	@Success		200	{object}	models.UserWithRoles		"user"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/users/{id}/roles [put]
func SetUserRoles(w http.ResponseWriter, r *http.Request) {
	userId, ok := U().ParseParamUUID(w, r, "id")
	if !ok || !U().CheckPermission(w, r, "admin.users.roles.update") {
		return
	}

	var stringRoles []string
	err := json.NewDecoder(r.Body).Decode(&stringRoles)
	if err != nil {
		zap.L().Warn("Role UUIDs json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}

	// Parse role UUIDs
	uuidRoles := make([]uuid.UUID, 0, len(stringRoles))
	for _, stringRole := range stringRoles {
		uuidRole, err := uuid.Parse(stringRole)
		if err != nil {
			zap.L().Warn("Invalid role UUID", zap.String("uuid", stringRole), zap.Error(err))
			render.BadRequest(w, r, fmt.Errorf("invalid role UUID: %s", stringRole))
			return
		}
		uuidRoles = append(uuidRoles, uuidRole)
	}

	// Set roles on user
	err = repositories.R().SetUserRoles(userId, uuidRoles)
	if err != nil {
		zap.L().Error("PutUser.Update", zap.Error(err))
		render.Error(w, r, err, "Set roles on user")
		return
	}

	render.OK(w, r)
}
