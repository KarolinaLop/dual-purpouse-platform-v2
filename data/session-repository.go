package data

import (
	"database/sql"

	"github.com/KarolinaLop/dp/models"
)

// CreateSession creates a new session in the database.
func CreateSession(db *sql.DB, sessionID string, userID int) error {
	_, err := db.Exec(`
		REPLACE INTO sessions (id, user_id, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
	`, sessionID, userID)
	return err
}

// GetSessionUser retrieves the session's user.
func GetSessionUser(db *sql.DB, sessionID string) (models.User, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM sessions WHERE id = ?", sessionID).Scan(&userID)
	if err == sql.ErrNoRows {
		return models.User{}, err
	}

	return GetUserByID(db, userID)
}

// DeleteSessions deletes all sessions for a user.
func DeleteSessions(db *sql.DB, userID int) error {
	query := "DELETE FROM sessions WHERE user_id = ?"
	_, err := db.Exec(query, userID)
	return err
}
