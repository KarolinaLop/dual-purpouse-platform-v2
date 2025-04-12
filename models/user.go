package models

// User represents a user in the system.
type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	CreatedAt    string `json:"created_at"`
	Password     string `json:"password"`
	PasswordHash string `json:"-"`
}
