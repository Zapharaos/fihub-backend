package models

import "github.com/Zapharaos/fihub-backend/internal/models"

type UserWithRoles struct {
	models.User
	Roles models.Roles `json:"roles"`
}
