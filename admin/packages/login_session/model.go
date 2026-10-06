package login_session

import (
	"time"

	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/go-rvq/rvq/x/place"
	"github.com/google/uuid"
)

// The kinds of access a LoginSession records.
const (
	// KindLogin is a login of the admin (a session, its cookie).
	KindLogin = "login"
	// KindWebDAV is the WebDAV of the files.
	KindWebDAV = "webdav"
	// KindGit is the git of the admin (the site's files).
	KindGit = "git"
)

// LoginSession is a place the user is in the admin from: a login (its
// session), or the accesses of the WebDAV or the git — of one device and IP,
// by one way in (Auth), in a window of time (Manager.AccessWindow) — one row.
type LoginSession struct {
	uuidkey.Model

	UserID uuid.UUID `gorm:"type:uuid;index"`
	// Device is the family of the browser and its system ("Chrome -
	// Linux"), or what the access was by (an access key's name).
	Device string
	IP     string
	// UserAgent is the browser's, whole.
	UserAgent string
	// Place is where IP is (Manager.PlaceFunc).
	Place     place.Place `gorm:"embedded"`
	TokenHash string      `sql:"index"`
	ExpiredAt time.Time

	// Kind is the way of the access: KindLogin, KindWebDAV, KindGit ("" a
	// login recorded before it).
	Kind string `gorm:"size:16;index"`
	// Auth is how the user was told: login.AuthPassword,
	// login.AuthAccessKey, login.AuthSecureKey, login.AuthSession; or
	// AuthLocked, the account locked by this access.
	Auth string `gorm:"size:16"`
	// LastAccessAt is the last request of the access (a window's), nil for a
	// login.
	LastAccessAt *time.Time
	// Requests is how many requests the access made.
	Requests int

	Time   string `gorm:"-"`
	Status string `gorm:"-"`
	// Access is Kind and Auth as the list shows them; PlaceText, Place.
	Access    string `gorm:"-"`
	PlaceText string `gorm:"-"`
}
