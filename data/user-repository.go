package data

import (
	"database/sql"
	"errors"

	"github.com/KarolinaLop/dp/models"
)

// ErrUserNotFound is returned when a user is not found in the database.
var ErrUserNotFound = errors.New("user not found")

// GetUserByEmail retrieves a user by email.
func GetUserByEmail(db *sql.DB, email string) (models.User, error) {
	query := "SELECT id, name, email, created_at, password FROM users WHERE email = ?;"
	row := db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err // Some other error occurred
	}
	return user, nil
}

// GetUserByID retrieves a user by ID.
func GetUserByID(db *sql.DB, ID int) (models.User, error) {
	query := "SELECT id, name, email, created_at, password FROM users WHERE id = ?;"
	row := db.QueryRow(query, ID)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err // Some other error occurred
	}
	return user, nil
}

// UserExists checks if a user with a given email already exists in the database.
func UserExists(db *sql.DB, email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	var count int
	err := db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateUser inserts a new user into the database.
func CreateUser(db *sql.DB, user models.User) (models.User, error) {
	query := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
	result, err := db.Exec(query, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return models.User{}, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return models.User{}, err
	}
	user.ID = int(userID)
	return user, nil
}
