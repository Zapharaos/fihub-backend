package users

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// PostgresRepository is a repository containing the Issue definition based on a PSQL database and
// implementing the repository interface
type PostgresRepository struct {
	conn *sqlx.DB
}

// NewPostgresRepository returns a new instance of PostgresRepository
func NewPostgresRepository(dbClient *sqlx.DB) Repository {
	r := PostgresRepository{
		conn: dbClient,
	}
	var repo Repository = &r
	return repo
}

// Create method used to create a user
func (r *PostgresRepository) Create(user *UserWithPassword) (uuid.UUID, error) {

	// UUID
	userID := uuid.New()

	// Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.UUID{}, err
	}

	// Get timestamps
	creationTS := time.Now().Truncate(1 * time.Millisecond).UTC()
	updateTS := creationTS

	// Prepare query
	query := `INSERT INTO users (id, email, first_name, last_name, password, created_at, updated_at)
				VALUES (:id, :email, :first_name, :last_name, :password, :created_at, :updated_at)`
	params := map[string]interface{}{
		"id":         userID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"password":   hashedPassword,
		"created_at": creationTS,
		"updated_at": updateTS,
	}

	// Execute query
	_, err = r.conn.NamedQuery(query, params)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

// Get use to retrieve a user by id
func (r *PostgresRepository) Get(userID uuid.UUID) (*User, error) {

	// Prepare query
	query := `SELECT *
			  FROM users as u
			  WHERE u.id = :id`
	params := map[string]interface{}{
		"id": userID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Retrieve user
	var user *User
	if rows.Next() {
		user, err = scanUser(rows)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return user, nil
}

func scanUser(rows *sqlx.Rows) (*User, error) {
	var userWithPassword UserWithPassword
	err := rows.Scan(
		&userWithPassword.ID,
		&userWithPassword.Email,
		&userWithPassword.Password,
		&userWithPassword.FirstName,
		&userWithPassword.LastName,
		&userWithPassword.CreatedAt,
		&userWithPassword.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Create a new User struct without the password hash
	user := User{
		ID:        userWithPassword.ID,
		Email:     userWithPassword.Email,
		FirstName: userWithPassword.FirstName,
		LastName:  userWithPassword.LastName,
		CreatedAt: userWithPassword.CreatedAt,
		UpdatedAt: userWithPassword.UpdatedAt,
	}

	return &user, nil
}
