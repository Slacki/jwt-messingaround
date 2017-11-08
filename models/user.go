package models

// User model representing user table
type User struct {
	ID           uint `gorm:"primary_key"`
	Identifier   string
	PasswordHash string
	Email        string
}
