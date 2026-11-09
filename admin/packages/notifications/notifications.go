// Package notifications is the notification center of the users of the
// admin: what an application tells a user — a transfer left a wallet below
// zero, an invite… — kept until they read it: in the bell (its unread ones,
// a source of the notification center of presets), and in their listing (My
// notifications).
package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rvq/rvq/x/i18n"
	"github.com/google/uuid"
	"github.com/sunfmin/reflectutils"
	"gorm.io/gorm"
)

// The levels of importance of a notification, the least first: each one is
// the name of its color in the palette of the theme (LevelColor).
const (
	LevelInfo    = "info"    // information
	LevelSuccess = "success" // something went well
	LevelWarning = "warning" // an alert: to look at
	LevelError   = "error"   // critical: to act on
)

// Levels are the levels of importance, the least first.
var Levels = []string{LevelInfo, LevelSuccess, LevelWarning, LevelError}

// SourceSystem is the source of what the system tells (System): said by
// messages, in the language of who reads it.
const SourceSystem = "system"

// SourceUser is the source of what the user id sent: their own text, not
// translated.
func SourceUser(id uuid.UUID) string { return sourceUserPrefix + id.String() }

const sourceUserPrefix = "user:"

// Notification is what is told to a user.
type Notification struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index:idx_user_notifications_user"`
	// Title and Body are the text as it was written: what is shown when it
	// has no message (Module) or the language of the reader lacks it, and
	// what the search finds.
	Title string `gorm:"size:300;not null"`
	Body  string `gorm:"size:2000"`
	// Module, TitleKey and BodyKey are the messages it is said by, in the
	// language of who reads it (Text): the fields TitleKey and BodyKey of the
	// messages of the i18n module Module, formatted (fmt; the arguments are text: %s) with
	// TitleArgs, BodyArgs — the names, the amounts.
	Module    i18n.ModuleKey `gorm:"size:120"`
	TitleKey  string         `gorm:"size:120"`
	TitleArgs []string       `gorm:"serializer:json"`
	BodyKey   string         `gorm:"size:120"`
	BodyArgs  []string       `gorm:"serializer:json"`
	// Link is where it takes (an address of the admin); "" none
	Link  string `gorm:"size:1000"`
	Level string `gorm:"size:12;not null;default:info"`
	// Source is who made it: the system (SourceSystem, the default), a user
	// (SourceUser) or another one of the app's
	Source string `gorm:"size:120;index"`
	// System: told by the system (Source is SourceSystem) — its title and
	// body said by messages (Say), translated; otherwise the text as written
	System    bool `gorm:"not null;default:false;index"`
	CreatedAt time.Time
	// ReadAt is when the user read it (nil: unread)
	ReadAt *time.Time `gorm:"index:idx_user_notifications_user"`
	// DeletedAt: deleted softly — in the trash of the listing, restored
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// MessageID is the sending it is of: the copies of a message sent to
	// many share it, and its attachments (File)
	MessageID uuid.UUID `gorm:"type:uuid;index"`
	// ThreadID is the conversation it is of: the MessageID of its first
	// message; its answers (Reply) share it
	ThreadID uuid.UUID `gorm:"type:uuid;index"`
	// ParentID is the message it answers (MessageID; none: the first)
	ParentID uuid.UUID `gorm:"type:uuid;index"`
	// HTML: its Body is rich text (HTML, a user's: shown sanitized)
	HTML bool `gorm:"not null;default:false"`
	// ExternalID is where it came from out of the admin (the Message-ID of
	// an answer by mail): what keeps it from being kept twice
	ExternalID string `gorm:"size:255;index"`
}

func (Notification) TableName() string { return "user_notifications" }

// String is its title (its name in the admin), as it was written.
func (n *Notification) String() string { return n.Title }

// Text are its title and body in the language of ctx: when it is the
// system's (System), its messages (Module, TitleKey, BodyKey) with their args;
// each one its Title, Body when it is another's, has no message or the
// language lacks it.
func (n *Notification) Text(ctx context.Context) (title, body string) {
	return n.message(ctx, n.TitleKey, n.TitleArgs, n.Title), n.message(ctx, n.BodyKey, n.BodyArgs, n.Body)
}

func (n *Notification) message(ctx context.Context, key string, a []string, written string) string {
	if !n.System || n.Module == "" || key == "" {
		return written
	}
	msgs := i18n.MustGetModuleMessages(ctx, n.Module, nil)
	if msgs == nil {
		return written
	}
	v, err := reflectutils.Get(msgs, key)
	format, _ := v.(string)
	if err != nil || format == "" {
		return written
	}
	if len(a) == 0 {
		return format
	}
	args := make([]any, len(a))
	for i, x := range a {
		args[i] = x
	}
	return fmt.Sprintf(format, args...)
}

// Message is a message of an i18n module: the name of its field, the
// arguments it is formatted with (fmt).
type Message struct {
	Key  string
	Args []string
}

// Msg is the message key with args.
func Msg(key string, args ...string) Message { return Message{Key: key, Args: args} }

// Say makes n the system's (SourceSystem), said by the messages of module in
// the language of who reads it (Text): its title the message title, its body
// body (a zero Message: none). Its Title and Body are written in the language
// of ctx — the one of who made it.
func (n Notification) Say(ctx context.Context, module i18n.ModuleKey, title, body Message) Notification {
	n.Source, n.System = SourceSystem, true
	n.Module = module
	n.TitleKey, n.TitleArgs = title.Key, title.Args
	n.BodyKey, n.BodyArgs = body.Key, body.Args
	n.Title, n.Body = n.Text(ctx)
	return n
}

func (n *Notification) BeforeCreate(*gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	if n.Level == "" {
		n.Level = LevelInfo
	}
	if n.Source == "" {
		n.Source = SourceSystem
	}
	n.System = n.Source == SourceSystem
	return nil
}

// AutoMigrate makes the tables of the notifications and their attachments;
// one of before a message (MessageID) and a conversation (ThreadID) of its
// own.
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&Notification{}, &File{}); err != nil {
		return err
	}
	return db.Model(&Notification{}).Unscoped().
		Where("message_id IS NULL OR message_id = ?", uuid.Nil).
		UpdateColumns(map[string]any{"message_id": gorm.Expr("id"), "thread_id": gorm.Expr("id")}).Error
}

// Notify tells n to each of users (one notification each; a user twice,
// once): one message (MessageID, a new one when n has none), of its
// conversation (ThreadID; a new message, its own).
func Notify(db *gorm.DB, n Notification, users ...uuid.UUID) error {
	if n.MessageID == uuid.Nil {
		n.MessageID = uuid.New()
	}
	if n.ThreadID == uuid.Nil {
		n.ThreadID = n.MessageID
	}
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

// MarkUnread marks the notifications ids of the user unread again (another's
// left).
func MarkUnread(db *gorm.DB, user uuid.UUID, ids ...uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	return db.Model(&Notification{}).Where("user_id = ? AND id IN ? AND read_at IS NOT NULL", user, ids).
		Update("read_at", nil).Error
}

// MarkAllRead marks every notification of the user read.
func MarkAllRead(db *gorm.DB, user uuid.UUID) error {
	return db.Model(&Notification{}).Where("user_id = ? AND read_at IS NULL", user).Update("read_at", time.Now()).Error
}
