package models

// User represents a user in the system.
type User struct {
	ID           int
	Name         string
	Email        string
	CreatedAt    string
	Password     string
	PasswordHash string
}
