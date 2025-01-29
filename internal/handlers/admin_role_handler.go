package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
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
//	@Param			role	body	roles.Role	true	"role (json)"
//	@Security		Bearer
//	@Success		200	{object}	roles.Role				"role"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles [post]
func CreateRole(w http.ResponseWriter, r *http.Request) {
	if !checkPermission(w, r, "admin.roles.create") {
		return
	}

	var role roles.Role
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

	roleID, err := roles.R().Create(role)
	if err != nil {
		render.Error(w, r, err, "Create role")
		return
	}

	newRole, found, err := roles.R().Get(roleID)
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
//	@Success		200	{object}	roles.Role				"role"
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

	role, found, err := roles.R().Get(roleID)
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
//	@Success		200	{array}		roles.Role				"list of roles"
//	@Failure		400	{object}	render.ErrorResponse	"Bad Request"
//	@Failure		401	{string}	string					"Permission denied"
//	@Failure		500	{object}	render.ErrorResponse	"Internal Server Error"
//	@Router			/api/v1/roles [get]
func GetRoles(w http.ResponseWriter, r *http.Request) {
	if !checkPermission(w, r, "admin.roles.list") {
		return
	}

	result, err := roles.R().GetAll()
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
//	@Param			role	body	roles.Role	true	"role (json)"
//	@Security		Bearer
//	@Success		200	{object}	roles.Role				"role"
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

	var role roles.Role
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		zap.L().Warn("Role json decode", zap.Error(err))
		render.BadRequest(w, r, nil)
		return
	}
	role.Id = roleID

	if ok, err := role.IsValid(); !ok {
		render.BadRequest(w, r, err)
		return
	}

	err = roles.R().Update(role)
	if err != nil {
		render.Error(w, r, err, "Update role")
		return
	}

	role, found, err := roles.R().Get(roleID)
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
