package login_session

import (
	"time"

	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/google/uuid"
)

type LoginSession struct {
	uuidkey.Model

	UserID    uuid.UUID `gorm:"type:uuid;index"`
	Device    string
	IP        string
	TokenHash string `sql:"index"`
	ExpiredAt time.Time

	Time   string `gorm:"-"`
	Status string `gorm:"-"`
}
