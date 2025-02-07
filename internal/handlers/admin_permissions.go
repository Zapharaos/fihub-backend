package handlers

import (
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/internal/handlers/render"
	"go.uber.org/zap"
	"net/http"
)

// CreatePermission godoc
//
//	@Id				CreatePermission
//
//	@Summary		Create a new permission
//	@Description	Create a new permission. (Permission: <b>admin.permissions.create</b>)
//	@Tags			Permissions
//	@Accept			json
//	@Produce		json
//	@Param			permission	body	permissions.Permission	true	"permission (json)"
//	@Security		Bearer
//	@Success		200	{object}	permissions.Permission	"permission"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/permissions [post]
func CreatePermission(w http.ResponseWriter, r *http.Request) {
	if !U().CheckPermission(w, r, "admin.permissions.create") {
		return
	}

	var permission permissions.Permission
	err := json.NewDecoder(r.Body).Decode(&permission)
	if err != nil {
		zap.L().Warn("Permission json decode", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	if ok, err := permission.IsValid(); !ok {
		zap.L().Warn("Permission is not valid", zap.Error(err))
		render.BadRequest(w, r, err)
		return
	}

	permissionID, err := permissions.R().Create(permission)
	if err != nil {
		render.Error(w, r, err, "Create permission")
		return
	}

	newPermission, found, err := permissions.R().Get(permissionID)
	if err != nil {
		zap.L().Error("Cannot get permission", zap.String("uuid", permissionID.String()), zap.Error(err))
		render.Error(w, r, nil, "")
		return
	}
	if !found {
		zap.L().Error("Permission not found after creation", zap.String("uuid", permissionID.String()))
		render.Error(w, r, nil, "")
		return
	}

	render.JSON(w, r, newPermission)
}

// GetPermission godoc
//
//	@Id				GetPermission
//
//	@Summary		Get a permission
//	@Description	Get a permission by id. (Permission: <b>permissions.read</b>)
//	@Tags			Permissions
//	@Produce		json
//	@Param			id	path	string	true	"permission id"
//	@Security		Bearer
//	@Success		200	{object}	permissions.Permission	"permission"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		404	{string}	string					"Not Found"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/permissions/{id} [get]
func GetPermission(w http.ResponseWriter, r *http.Request) {
	permissionID, ok := U().ParseParamUUID(w, r, "id")
	if !ok || !U().CheckPermission(w, r, "admin.permissions.read") {
		return
	}

	permission, found, err := permissions.R().Get(permissionID)
	if err != nil {
		zap.L().Error("Cannot load permission", zap.String("uuid", permissionID.String()), zap.Error(err))
		render.Error(w, r, nil, "")
		return
	}

	if !found {
		zap.L().Debug("Permission not found", zap.String("uuid", permissionID.String()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	render.JSON(w, r, permission)
}

// GetPermissions godoc
//
//	@Id				GetPermissions
//
//	@Summary		Get all permissions
//	@Description	Gets a list of all permissions. (Permission: <b>admin.permissions.list</b>)
//	@Tags			Permissions
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{array}		permissions.Permission	"list of permissions"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/permissions [get]
func GetPermissions(w http.ResponseWriter, r *http.Request) {
	if !U().CheckPermission(w, r, "admin.permissions.list") {
		return
	}

	result, err := permissions.R().GetAll()
	if err != nil {
		render.Error(w, r, err, "Get permissions")
		return
	}

	render.JSON(w, r, result)
}

// UpdatePermission godoc
//
//	@Id				UpdatePermission
//
//	@Summary		Update permission
//	@Description	Updates the permission. (Permission: <b>admin.permissions.update</b>)
//	@Tags			Permissions
//	@Accept			json
//	@Produce		json
//	@Param			id			path	string					true	"permission ID"
//	@Param			permission	body	permissions.Permission	true	"permission (json)"
//	@Security		Bearer
//	@Success		200	{object}	permissions.Permission	"permission"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/permissions/{id} [put]
func UpdatePermission(w http.ResponseWriter, r *http.Request) {
	permissionID, ok := U().ParseParamUUID(w, r, "id")
	if !ok || !U().CheckPermission(w, r, "admin.permissions.update") {
		return
	}

	var permission permissions.Permission
	err := json.NewDecoder(r.Body).Decode(&permission)
	if err != nil {
		zap.L().Warn("Permission json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}
	permission.Id = permissionID

	if ok, err := permission.IsValid(); !ok {
		render.BadRequest(w, r, err)
		return
	}

	err = permissions.R().Update(permission)
	if err != nil {
		render.Error(w, r, err, "Update permission")
		return
	}

	permission, found, err := permissions.R().Get(permissionID)
	if err != nil {
		zap.L().Error("Cannot get permission", zap.String("uuid", permissionID.String()), zap.Error(err))
		render.Error(w, r, nil, "")
		return
	}
	if !found {
		zap.L().Error("Permission not found after update", zap.String("uuid", permissionID.String()))
		render.Error(w, r, nil, "")
		return
	}

	render.JSON(w, r, permission)
}

// DeletePermission godoc
//
//	@Id				DeletePermission
//
//	@Summary		Delete permission
//	@Description	Deletes a permission. (Permission: <b>admin.permissions.delete</b>)
//	@Tags			Permissions
//	@Produce		json
//	@Param			id	path	string	true	"permission ID"
//	@Security		Bearer
//	@Success		200	{string}	string					"status OK"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/permissions/{id} [delete]
func DeletePermission(w http.ResponseWriter, r *http.Request) {
	permissionID, ok := U().ParseParamUUID(w, r, "id")
	if !ok || !U().CheckPermission(w, r, "admin.permissions.delete") {
		return
	}

	err := permissions.R().Delete(permissionID)
	if err != nil {
		render.Error(w, r, err, "Cannot delete permission")
		return
	}

	render.OK(w, r)
}
