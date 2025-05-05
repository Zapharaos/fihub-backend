package handlers

import (
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers/render"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/internal/mappers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"go.uber.org/zap"
	"net/http"
)

// CreatePermission godoc
//
//	@Id				CreatePermission
//
//	@Summary		Create a new permission
//	@Description	Create a new permission. (Permission: <b>admin.permissions.create</b>)
//	@Tags			Security, Permission
//	@Accept			json
//	@Produce		json
//	@Param			permission	body	models.Permission	true	"permission (json)"
//	@Security		Bearer
//	@Success		200	{object}	models.Permission	"permission"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/permission [post]
func CreatePermission(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var permission models.Permission
	err := json.NewDecoder(r.Body).Decode(&permission)
	if err != nil {
		zap.L().Warn("Permission json decode", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	// Create gRPC gen.CreatePermissionRequest
	permissionRequest := &securitypb.CreatePermissionRequest{
		Value:       permission.Value,
		Scope:       permission.Scope,
		Description: permission.Description,
	}

	// Create the Permission
	response, err := clients.C().Security().CreatePermission(r.Context(), permissionRequest)
	if err != nil {
		zap.L().Error("Create Permission", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the models.Permission struct
	render.JSON(w, r, mappers.PermissionFromProto(response.Permission))
}

// GetPermission godoc
//
//	@Id				GetPermission
//
//	@Summary		Get a permission
//	@Description	Get a permission by id. (Permission: <b>permissions.read</b>)
//	@Tags			Security, Permission
//	@Produce		json
//	@Param			id	path	string	true	"permission id"
//	@Security		Bearer
//	@Success		200	{object}	models.Permission	"permission"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		404	{string}	string					"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/permission/{id} [get]
func GetPermission(w http.ResponseWriter, r *http.Request) {
	permissionID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Create gRPC gen.CreatePermissionRequest
	permissionRequest := &securitypb.GetPermissionRequest{
		Id: permissionID.String(),
	}

	// Get the Permission
	response, err := clients.C().Security().GetPermission(r.Context(), permissionRequest)
	if err != nil {
		zap.L().Error("Get Permission", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the models.Permission struct
	render.JSON(w, r, mappers.PermissionFromProto(response.Permission))
}

// UpdatePermission godoc
//
//	@Id				UpdatePermission
//
//	@Summary		Update permission
//	@Description	Updates the permission. (Permission: <b>admin.permissions.update</b>)
//	@Tags			Security, Permission
//	@Accept			json
//	@Produce		json
//	@Param			id			path	string					true	"permission ID"
//	@Param			permission	body	models.Permission	true	"permission (json)"
//	@Security		Bearer
//	@Success		200	{object}	models.Permission	"permission"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/permission/{id} [put]
func UpdatePermission(w http.ResponseWriter, r *http.Request) {
	permissionID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	var permission models.Permission
	err := json.NewDecoder(r.Body).Decode(&permission)
	if err != nil {
		zap.L().Warn("Permission json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}

	// Create gRPC gen.UpdatePermissionRequest
	permissionRequest := &securitypb.UpdatePermissionRequest{
		Id:          permissionID.String(),
		Value:       permission.Value,
		Scope:       permission.Scope,
		Description: permission.Description,
	}

	// Update the Permission
	response, err := clients.C().Security().UpdatePermission(r.Context(), permissionRequest)
	if err != nil {
		zap.L().Error("Update Permission", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	// Map the response to the models.Permission struct
	render.JSON(w, r, mappers.PermissionFromProto(response.Permission))
}

// DeletePermission godoc
//
//	@Id				DeletePermission
//
//	@Summary		Delete permission
//	@Description	Deletes a permission. (Permission: <b>admin.permissions.delete</b>)
//	@Tags			Security, Permission
//	@Produce		json
//	@Param			id	path	string	true	"permission ID"
//	@Security		Bearer
//	@Success		200	{string}	string					"status OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad PasswordRequest"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/permission/{id} [delete]
func DeletePermission(w http.ResponseWriter, r *http.Request) {
	permissionID, ok := U().ParseParamUUID(w, r, "id")
	if !ok {
		return
	}

	// Create gRPC gen.DeletePermissionRequest
	permissionRequest := &securitypb.DeletePermissionRequest{
		Id: permissionID.String(),
	}

	// Delete the Permission
	_, err := clients.C().Security().DeletePermission(r.Context(), permissionRequest)
	if err != nil {
		zap.L().Error("Delete Permission", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.OK(w, r)
}

// ListPermissions godoc
//
//	@Id				ListPermissions
//
//	@Summary		Get all permissions
//	@Description	Gets a list of all permissions. (Permission: <b>admin.permissions.list</b>)
//	@Tags			Security, Permission
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{array}		models.Permission	"list of permissions"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/security/permission [get]
func ListPermissions(w http.ResponseWriter, r *http.Request) {
	// List the Broker
	response, err := clients.C().Security().ListPermissions(r.Context(), &securitypb.ListPermissionsRequest{})
	if err != nil {
		zap.L().Error("List Permissions", zap.Error(err))
		render.ErrorCodesCodeToHttpCode(w, r, err)
		return
	}

	render.JSON(w, r, mappers.PermissionsFromProto(response.Permissions))
}
