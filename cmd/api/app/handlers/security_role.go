package handlers

import (
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	apimodels "github.com/Zapharaos/fihub-backend/cmd/api/app/models"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/mappers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"net/http"

	"go.uber.org/zap"
)

// CreateRole godoc
//
//	@Id				CreateRole
//
//	@Summary		Create a new role
//	@Description	Create a new role. (Permission: <b>admin.roles.create</b>)
//	@Tags			Security, Role
//	@Accept			json
//	@Produce		json
//	@Param			role	body	models.RoleWithPermissions	true	"role (json)"
//	@Security		Bearer
//	@Success		200	{object}	models.RoleWithPermissions				"role"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/role [post]
func CreateRole(w http.ResponseWriter, r *http.Request) {
	var role models.RoleWithPermissions
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		zap.L().Warn("Role json decode", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Create gRPC gen.CreateRoleRequest
	roleRequest := &securitypb.CreateRoleRequest{
		Name:        role.Name,
		Permissions: role.Permissions.GetUUIDsAsStrings(),
	}

	// Create the role
	response, err := clients.C().Security().CreateRole(r.Context(), roleRequest)
	if err != nil {
		zap.L().Error("Create Role", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the RoleWithPermissions model
	render.JSON(w, r, mappers.RoleWithPermissionsFromProto(response.GetRole()))
}

// GetRole godoc
//
//	@Id				GetRole
//
//	@Summary		Get a role
//	@Description	Get a role by id. (Permission: <b>admin.roles.read</b>)
//	@Tags			Security, Role
//	@Produce		json
//	@Param			id	path	string	true	"role id"
//	@Security		Bearer
//	@Success		200	{object}	models.RoleWithPermissions "role"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		404	{string}	string					"Role not found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/role/{id} [get]
func GetRole(w http.ResponseWriter, r *http.Request) {
	roleID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Get the role
	response, err := clients.C().Security().GetRole(r.Context(), &securitypb.GetRoleRequest{
		Id: roleID.String(),
	})
	if err != nil {
		zap.L().Error("Get Role", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the RoleWithPermissions model
	render.JSON(w, r, mappers.RoleWithPermissionsFromProto(response.GetRole()))
}

// ListRoles godoc
//
//	@Id				ListRoles
//
//	@Summary		List all roles
//	@Description	List all roles. (Permission: <b>admin.roles.list</b>)
//	@Tags			Security, Role
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{array}		models.RoleWithPermissions				"list of roles"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/role [get]
func ListRoles(w http.ResponseWriter, r *http.Request) {
	// List the roles
	response, err := clients.C().Security().ListRoles(r.Context(), nil)
	if err != nil {
		zap.L().Error("List Roles", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the RolesWithPermissions model
	render.JSON(w, r, mappers.RolesWithPermissionsFromProto(response.GetRoles()))
}

// UpdateRole godoc
//
//	@Id				UpdateRole
//
//	@Summary		Update role
//	@Description	Updates the role. (Permission: <b>admin.roles.update</b>)
//	@Tags			Security, Role
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string		true	"role ID"
//	@Param			role	body	models.RoleWithPermissions	true	"role (json)"
//	@Security		Bearer
//	@Success		200	{object}	models.RoleWithPermissions				"role"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/role/{id} [put]
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

	// Create gRPC gen.UpdateRoleRequest
	roleRequest := &securitypb.UpdateRoleRequest{
		Id:          roleID.String(),
		Name:        role.Name,
		Permissions: role.Permissions.GetUUIDsAsStrings(),
	}

	// Update the role
	response, err := clients.C().Security().UpdateRole(r.Context(), roleRequest)
	if err != nil {
		zap.L().Error("Update Role", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the RoleWithPermissions model
	render.JSON(w, r, mappers.RoleWithPermissionsFromProto(response.GetRole()))
}

// DeleteRole godoc
//
//	@Id				DeleteRole
//
//	@Summary		Delete role
//	@Description	Deletes a role. (Permission: <b>admin.roles.delete</b>)
//	@Tags			Security, Role
//	@Produce		json
//	@Param			id	path	string	true	"role ID"
//	@Security		Bearer
//	@Success		200	{string}	string					"status OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/role/{id} [delete]
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	roleID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Delete the role
	_, err := clients.C().Security().DeleteRole(r.Context(), &securitypb.DeleteRoleRequest{
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
//	@Tags			Security, Role, Permission
//	@Produce		json
//	@Param			id	path	string	true	"role ID"
//	@Security		Bearer
//	@Success		200	{array}		models.Permission	"list of permissions"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/role/{id}/permission [get]
func GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	roleId, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// List the role permissions
	response, err := clients.C().Security().ListRolePermissions(r.Context(), &securitypb.ListRolePermissionsRequest{
		Id: roleId.String(),
	})
	if err != nil {
		zap.L().Error("List Role Permissions", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the Permissions model
	render.JSON(w, r, mappers.PermissionsFromProto(response.GetPermissions()))
}

// SetRolePermissions godoc
//
//	@Id				SetRolePermissions
//
//	@Summary		Set permissions for a given role
//	@Description	Updates the role. (Permission: <b>admin.roles.permissions.update</b>)
//	@Tags			Security, Role, Permission
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string						true	"role ID"
//	@Param			role	body	models.Permissions	true	"List of permissions (json)"
//	@Security		Bearer
//	@Success		200	{array}		models.Permission	"list of permissions"
//	@Failure		400	{object}	render.ErrorResponse		"Bad PasswordRequest"
//	@Failure		401	{string}	string						"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse		"Internal Server Error"
//	@Router			/api/v1/security/role/{id}/permission [put]
func SetRolePermissions(w http.ResponseWriter, r *http.Request) {
	roleId, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	var perms models.Permissions
	err := json.NewDecoder(r.Body).Decode(&perms)
	if err != nil {
		zap.L().Warn("Permissions json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}

	// Create gRPC gen.UpdateRoleRequest
	roleRequest := &securitypb.SetRolePermissionsRequest{
		Id:          roleId.String(),
		Permissions: perms.GetUUIDsAsStrings(),
	}

	// Set the role permissions
	response, err := clients.C().Security().SetRolePermissions(r.Context(), roleRequest)
	if err != nil {
		zap.L().Error("Set Role Permissions", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the Permissions model
	render.JSON(w, r, mappers.PermissionsFromProto(response.GetPermissions()))
}

// AddUsersToRole godoc
//
//	@Id				AddUsersToRole
//
//	@Summary		Add users to a given role
//	@Description	Updates the role. (Permission: <b>admin.roles.users.update</b>)
//	@Tags			Security, Role, User
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string		true	"role ID"
//	@Param			role	body	[]string	true	"List of user UUIDs (json)"
//	@Security		Bearer
//	@Success		200	{array}		string				"list of users ids"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/role/{id}/user [put]
func AddUsersToRole(w http.ResponseWriter, r *http.Request) {
	roleId, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	var userUUIDs []string
	err := json.NewDecoder(r.Body).Decode(&userUUIDs)
	if err != nil {
		zap.L().Warn("User UUIDs json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}

	// Add users to the role
	response, err := clients.C().Security().AddUsersToRole(r.Context(), &securitypb.AddUsersToRoleRequest{
		RoleId:  roleId.String(),
		UserIds: userUUIDs,
	})
	if err != nil {
		zap.L().Error("Add Users To Role", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.JSON(w, r, response.GetUserIds())
}

// RemoveUsersFromRole godoc
//
//	@Id				RemoveUsersFromRole
//
//	@Summary		Remove users from a given role
//	@Description	Updates the role. (Permission: <b>admin.roles.users.delete</b>)
//	@Tags			Security, Role, User
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string		true	"role ID"
//	@Param			role	body	[]string	true	"List of user UUIDs (json)"
//	@Security		Bearer
//	@Success		200	{array}		string				"list of users ids"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/role/{id}/user [delete]
func RemoveUsersFromRole(w http.ResponseWriter, r *http.Request) {
	roleId, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	var userUUIDs []string
	err := json.NewDecoder(r.Body).Decode(&userUUIDs)
	if err != nil {
		zap.L().Warn("User UUIDs json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}

	// Remove users from the role
	response, err := clients.C().Security().RemoveUsersFromRole(r.Context(), &securitypb.RemoveUsersFromRoleRequest{
		RoleId:  roleId.String(),
		UserIds: userUUIDs,
	})
	if err != nil {
		zap.L().Error("Remove Users from Role", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.JSON(w, r, response.GetUserIds())
}

// ListUsersForRole godoc
//
//	@Id				ListUsersForRole
//
//	@Summary		List all users for a specified role id
//	@Description	List all role users. (Permission: <b>admin.roles.users.list</b>)
//	@Tags			Security, Role, User
//	@Produce		json
//	@Param			id	path	string	true	"role ID"
//	@Security		Bearer
//	@Success		200	{array}		string				"list of users ids"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/role/{id}/user [get]
func ListUsersForRole(w http.ResponseWriter, r *http.Request) {
	roleId, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// List users for the role
	response, err := clients.C().Security().ListUsersForRole(r.Context(), &securitypb.ListUsersForRoleRequest{
		RoleId: roleId.String(),
	})
	if err != nil {
		zap.L().Error("List Users for Role", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.JSON(w, r, response.GetUserIds())
}

// SetRolesForUser godoc
//
//	@Id				SetRolesForUser
//
//	@Summary		Set roles on a user
//	@Description	Set roles on a user. (Permission: <b>admin.users.roles.update</b>)
//	@Tags			Security, Role, User
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string				true	"user ID"
//	@Param			roles	body	[]string			true	"array of role UUIDs"
//	@Security		Bearer
//	@Success		200	{array}		models.RoleWithPermissions	"list of roles with permissions"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/role/user/{id} [put]
func SetRolesForUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	var roleUUIDs []string
	err := json.NewDecoder(r.Body).Decode(&roleUUIDs)
	if err != nil {
		zap.L().Warn("Role UUIDs json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}

	// Set roles for the user
	_, err = clients.C().Security().SetRolesForUser(r.Context(), &securitypb.SetRolesForUserRequest{
		UserId:  userId.String(),
		RoleIds: roleUUIDs,
	})
	if err != nil {
		zap.L().Error("Set Roles for User", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.OK(w, r)
}

// ListRolesWithPermissionsForUser godoc
//
//	@Id				ListRolesWithPermissionsForUser
//
//	@Summary		List all roles with permissions for a specified user id
//	@Description	List of all roles with permissions. (Permission: <b>admin.users.roles.list</b>)
//	@Tags			Security, Role, User
//	@Produce		json
//	@Param			id	path	string	true	"user ID"
//	@Security		Bearer
//	@Success		200	{array}		models.RoleWithPermissions	"list of roles with permissions"
//	@Failure		400	{object}	render.ErrorResponse		"Bad PasswordRequest"
//	@Failure		401	{string}	string						"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse		"Internal Server Error"
//	@Router			/api/v1/security/role/user/{id} [get]
func ListRolesWithPermissionsForUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// List roles for the user
	response, err := clients.C().Security().ListRolesWithPermissionsForUser(r.Context(), &securitypb.ListRolesWithPermissionsForUserRequest{
		UserId: userId.String(),
	})
	if err != nil {
		zap.L().Error("List Roles with Permissions for User", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the RolesWithPermissions model
	render.JSON(w, r, mappers.RolesWithPermissionsFromProto(response.GetRoles()))
}

// ListUsersWithRoles godoc
//
//	@Id				ListUsersWithRoles
//
//	@Summary		List users with their roles
//	@Description	List of all users with their roles. (Permission: <b>admin.users.list</b>)
//	@Tags			Security, Role, User
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{array}		apimodels.UserWithRoles "list of users with roles"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/role/user [get]
func ListUsersWithRoles(w http.ResponseWriter, r *http.Request) {
	// List roles for the users
	respSecurityUsers, err := clients.C().Security().ListUsersFull(r.Context(), &securitypb.ListUsersFullRequest{})
	if err != nil {
		zap.L().Error("List Security Users", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Create a map of user IDs to roles (for faster lookup)
	userRolesMap := make(map[string]models.Roles)
	for _, user := range respSecurityUsers.GetUsers() {
		userID := user.GetUserId()
		userRoles := mappers.RolesFromProto(user.GetRoles())
		userRolesMap[userID] = userRoles
	}

	// Retrieve all app users
	respUserUsers, err := clients.C().User().ListUsers(r.Context(), &userpb.ListUsersRequest{})
	if err != nil {
		zap.L().Error("List User Users", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Link the users with their roles
	users := make([]apimodels.UserWithRoles, len(respUserUsers.GetUsers()))
	for i, u := range respUserUsers.GetUsers() {
		// Map the user to the UserWithRoles model
		user := mappers.UserFromProto(u)

		// Retrieve user roles from userRolesMap
		userRoles, ok := userRolesMap[user.ID.String()]
		if !ok {
			userRoles = models.Roles{}
		}

		// Link user to his roles
		users[i] = apimodels.UserWithRoles{
			User:  user,
			Roles: userRoles,
		}
	}

	render.JSON(w, r, users)
}
