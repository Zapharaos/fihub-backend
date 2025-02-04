package users

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/google/uuid"
	"time"
)

var (
	ErrEmailRequired        = errors.New("email-required")
	ErrEmailInvalid         = errors.New("email-invalid")
	ErrCheckboxInvalid      = errors.New("checkbox-invalid")
	ErrPasswordRequired     = errors.New("password-required")
	ErrPasswordInvalid      = errors.New("password-invalid")
	ErrConfirmationRequired = errors.New("confirmation-required")
	ErrConfirmationInvalid  = errors.New("confirmation-invalid")
)

const (
	PasswordMinLength = 8
	PasswordMaxLength = 64
)

// UserInputCreate extends UserInputPassword with a checkbox
type UserInputCreate struct {
	UserInputPassword
	Checkbox bool `json:"checkbox"`
}

// UserInputPassword extends UserWithPassword with a password-confirmation
type UserInputPassword struct {
	UserWithPassword
	Confirmation string `json:"confirmation"`
}

// UserWithPassword extends User with a password field for authentication purposes
type UserWithPassword struct {
	User
	Password string `json:"password"`
}

// UserWithRoles extends User with a roles field for authorization purposes
type UserWithRoles struct {
	User
	Roles roles.RolesWithPermissions `json:"roles"`
}

// User represents a User entity in the system
type User struct {
	ID        uuid.UUID `json:"ID"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsValid checks if a User is valid and has no missing mandatory PGFields
// * Email must not be empty
// * Email must not be valid
func (u User) IsValid() (bool, error) {
	if u.Email == "" {
		return false, ErrEmailRequired
	}
	if !email.IsValid(u.Email) {
		return false, ErrEmailInvalid
	}
	return true, nil
}

// IsValid checks if a User with password is valid and has no missing mandatory PGFields
// * User must be valid (see User struct)
// * Password must not be valid (see isValidPassword function)
func (u UserWithPassword) IsValid() (bool, error) {
	if ok, err := u.User.IsValid(); !ok {
		return false, err
	}
	return isValidPassword(u.Password)
}

// IsValid checks if a User input is valid and has no missing mandatory PGFields
// * UserWithPassword must be valid (see UserWithPassword struct)
// * Confirmation must not be valid (see isValidConfirmation function)
func (u UserInputPassword) IsValid() (bool, error) {
	if ok, err := u.UserWithPassword.IsValid(); !ok {
		return false, err
	}
	return isValidConfirmation(u.Password, u.Confirmation)
}

// IsValid checks if a User input is valid and has no missing mandatory PGFields
// * UserInputPassword must be valid (see UserInputPassword struct)
// * Checkbox must be true
func (u UserInputCreate) IsValid() (bool, error) {
	if ok, err := u.UserInputPassword.IsValid(); !ok {
		return false, err
	}
	if !u.Checkbox {
		return false, ErrCheckboxInvalid
	}
	return true, nil
}

// isValidPassword checks if a password is valid
// * Password must not be empty
// * Password must not be shorter than 8 characters
// * Password must not be longer than 64 characters
func isValidPassword(password string) (bool, error) {
	if password == "" {
		return false, ErrPasswordRequired
	}
	if len(password) < PasswordMinLength {
		return false, ErrPasswordInvalid
	}
	if len(password) > PasswordMaxLength {
		return false, ErrPasswordInvalid
	}
	return true, nil
}

// isValidConfirmation checks if a password and a confirmation are valid
// * Confirmation must not be empty
// * Confirmation must be equal to password
func isValidConfirmation(password, confirmation string) (bool, error) {
	if confirmation == "" {
		return false, ErrConfirmationRequired
	}
	if confirmation != password {
		return false, ErrConfirmationInvalid
	}
	return true, nil
}

// IsValidPassword checks if a User input password is valid
// * Password must not be valid (see isValidPassword function)
// * Confirmation must not be valid (see isValidConfirmation function)
func (u UserInputPassword) IsValidPassword() (bool, error) {
	if ok, err := isValidPassword(u.Password); !ok {
		return false, err
	}
	if ok, err := isValidConfirmation(u.Password, u.Confirmation); !ok {
		return false, err
	}
	return true, nil
}

// HasPermission returns true if the User has the given permission.
// Wildcards (*) in permissions are supported.
func (u *UserWithRoles) HasPermission(permission string) bool {
	for _, r := range u.Roles {
		for _, p := range r.Permissions {
			if p.Match(permission) {
				return true
			}
		}
	}
	return false
}
