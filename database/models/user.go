package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	BirthDate    string    `json:"birth_date"`
	Email        string    `json:"email"`
	PasswordHash string
	gorm.Model
}
