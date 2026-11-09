// Package notifications is the notification center of the users of the
// admin: what an application tells a user — a transfer left a wallet below
// zero, an invite… — kept until they read it: in the bell (its unread ones,
// a source of the notification center of presets), and in their listing (My
// notifications).
package notifications

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// The levels of a notification.
const (
	LevelInfo    = "info"
	LevelWarning = "warning"
	LevelError   = "error"
)

// Notification is what is told to a user.
type Notification struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index:idx_user_notifications_user"`
	Title  string    `gorm:"size:300;not null"`
	Body   string    `gorm:"size:2000"`
	// Link is where it takes (an address of the admin); "" none
	Link  string `gorm:"size:1000"`
	Level string `gorm:"size:12;not null;default:info"`
	// Source is what made it ("finance.transferencia"): to find or clear
	// the ones of a kind
	Source    string `gorm:"size:120;index"`
	CreatedAt time.Time
	// ReadAt is when the user read it (nil: unread)
	ReadAt *time.Time `gorm:"index:idx_user_notifications_user"`
}

func (Notification) TableName() string { return "user_notifications" }

func (n *Notification) BeforeCreate(*gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	if n.Level == "" {
		n.Level = LevelInfo
	}
	return nil
}

// AutoMigrate makes the table of the notifications.
func AutoMigrate(db *gorm.DB) error { return db.AutoMigrate(&Notification{}) }

// Notify tells n to each of users (one notification each; a user twice,
// once).
func Notify(db *gorm.DB, n Notification, users ...uuid.UUID) error {
	seen := map[uuid.UUID]bool{}
	for _, u := range users {
		if u == uuid.Nil || seen[u] {
			continue
		}
		seen[u] = true
		x := n
		x.ID, x.UserID, x.ReadAt = uuid.Nil, u, nil
		if err := db.Create(&x).Error; err != nil {
			return err
		}
	}
	return nil
}

// Unread is how many notifications of the user are unread.
func Unread(db *gorm.DB, user uuid.UUID) (n int64, err error) {
	err = db.Model(&Notification{}).Where("user_id = ? AND read_at IS NULL", user).Count(&n).Error
	return
}

// Latest are the newest unread notifications of the user, at most limit.
func Latest(db *gorm.DB, user uuid.UUID, limit int) (out []Notification, err error) {
	err = db.Where("user_id = ? AND read_at IS NULL", user).Order("created_at DESC").Limit(limit).Find(&out).Error
	return
}

// MarkRead marks the notifications ids of the user read (another's left).
func MarkRead(db *gorm.DB, user uuid.UUID, ids ...uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	return db.Model(&Notification{}).Where("user_id = ? AND id IN ? AND read_at IS NULL", user, ids).
		Update("read_at", time.Now()).Error
}

// MarkAllRead marks every notification of the user read.
func MarkAllRead(db *gorm.DB, user uuid.UUID) error {
	return db.Model(&Notification{}).Where("user_id = ? AND read_at IS NULL", user).Update("read_at", time.Now()).Error
}
