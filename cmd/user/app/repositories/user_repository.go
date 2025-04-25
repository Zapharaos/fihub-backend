package repositories

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
)

// UserRepository is a storage interface which can be implemented by multiple backend
// (in-memory map, sql database, in-memory cache, file system, ...)
// It allows standard CRUD operation on User
type UserRepository interface {
	Create(user models.UserWithPassword) (uuid.UUID, error)
	Get(userID uuid.UUID) (models.User, bool, error)
	GetByEmail(email string) (models.User, bool, error)
	Exists(email string) (bool, error)
	Authenticate(email string, password string) (models.User, bool, error)
	Update(user models.User) error
	UpdateWithPassword(user models.UserWithPassword) error
	Delete(userID uuid.UUID) error

	GetWithRoles(userID uuid.UUID) (models.UserWithRoles, error)
	GetAllWithRoles() ([]models.UserWithRoles, error)
	GetUsersByRoleID(roleUUID uuid.UUID) ([]models.User, error)
	UpdateWithRoles(user models.UserWithRoles, roleUUIDs []uuid.UUID) error
	SetUserRoles(userUUID uuid.UUID, roleUUIDs []uuid.UUID) error
	AddUsersRole(userUUIDs []uuid.UUID, id uuid.UUID) error
	RemoveUsersRole(userUUIDs []uuid.UUID, roleUUID uuid.UUID) error
}
