package note

import (
	"errors"
	"strings"

	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	uuidkey.Model

	UserID       uuid.UUID `gorm:"type:uuid;index"`
	Creator      string
	ResourceType string `gorm:"index"`
	ResourceID   string `gorm:"index"`
	Content      string `sql:"size:5000"`
}

func (this *Note) BeforeCreate(tx *gorm.DB) (err error) {
	if strings.TrimSpace(this.Content) == "" {
		err = errors.New("Note cannot be empty")
	}

	return
}

type UserNote struct {
	uuidkey.Model

	UserID       uuid.UUID `gorm:"type:uuid;index"`
	ResourceType string    `gorm:"index"`
	ResourceID   string    `gorm:"index"`
	Number       int64
}

func GetUnreadNotesCount(db *gorm.DB, userID uuid.UUID, resourceType, resourceID string) int64 {
	var total int64
	db.Model(&Note{}).Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).Count(&total)

	if total == 0 {
		return 0
	}

	userNote := UserNote{}
	db.Where("user_id = ? AND resource_type = ? AND resource_id = ?", userID, resourceType, resourceID).First(&userNote)
	return total - userNote.Number
}
