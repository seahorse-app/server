package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey;"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	BirthDate    string    `json:"birth_date"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Sessions     []Session `json:"sessions"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `gorm:"index"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}
