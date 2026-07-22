package login_session

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginSession struct {
	gorm.Model

	UserID    uuid.UUID `gorm:"type:uuid;index"`
	Device    string
	IP        string
	TokenHash string `sql:"index"`
	ExpiredAt time.Time

	Time   string `gorm:"-"`
	Status string `gorm:"-"`
}
