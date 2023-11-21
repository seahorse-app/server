package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
}
