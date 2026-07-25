package models

import (
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

type LoginSession struct {
	gorm.Model

	UserID    uuid.UUID `gorm:"type:uuid" sql:"index"`
	Device    string
	IP        string
	TokenHash string `sql:"index"`
	ExpiredAt time.Time

	Time   string `gorm:"-"`
	Status string `gorm:"-"`
}
