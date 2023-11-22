package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (session *Session) BeforeCreate(tx *gorm.DB) (err error) {
	session.ID = uuid.New()
	return
}
