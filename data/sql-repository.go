package data

import (
	"database/sql"

	"github.com/KarolinaLop/dp/models"
)

// GetUser retrieves a user by email.
func GetUser(db *sql.DB, email string) (models.User, error) {
	query := "SELECT id, name, email, created_at FROM users WHERE email = ?;"
	row := db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil // No user found
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
