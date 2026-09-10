package user

import (
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AnonymousID is the immutable primary key of the static anonymous user — the
// author recorded when a change is made without a logged-in user. It is a fixed
// UUID (never generated), so User.Anonymous() is a plain comparison against it
// and foreign keys to users.id resolve to a single well-known row.
var AnonymousID = uuid.MustParse("00000000-0000-4000-a000-000000000001")

// IsAnonymous reports whether id is the anonymous user's id.
func IsAnonymous(id uuid.UUID) bool { return id == AnonymousID }

// anonymousFactory builds the request's anonymous user: the app's concrete User
// model, carrying AnonymousID and a name translated for the request. It is set
// by the user Builder; GetCurrentUser falls back to it when the session has no
// user, so a request always has a user (possibly Anonymous()).
var anonymousFactory func(r *http.Request) User

// NewAnonymous builds an anonymous user of the concrete model type, with
// AnonymousID and the given (already translated) name.
func NewAnonymous(newModel func() User, name string) User {
	u := newModel()
	u.SetID(AnonymousID)
	u.SetName(name)
	return u
}

// SeedAnonymous inserts the anonymous user row when missing, so foreign keys to
// users.id (activity logs, revisions, …) always resolve. The stored row is only
// a placeholder for referential integrity: at runtime the name shown is the
// per-request translated one (never this value), and the admin keeps the row
// non-editable and non-deletable.
func SeedAnonymous(db *gorm.DB, newModel func() User, name string) error {
	var count int64
	if err := db.Unscoped().Model(newModel()).Where("id = ?", AnonymousID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return db.Create(NewAnonymous(newModel, name)).Error
}
