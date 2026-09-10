package user

import (
	"time"

	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/x/login"
	"github.com/google/uuid"
)

type User interface {
	login.UserPasser
	GetID() uuid.UUID
	SetID(v uuid.UUID)
	// Anonymous reports whether this is the static anonymous user — its ID is
	// the immutable AnonymousID. There is always a user on a request (a save or
	// an activity log never has a nil author): when the session carries none,
	// GetCurrentUser returns the anonymous user instead of nil.
	Anonymous() bool
	GetName() string
	SetName(v string)
	SetEmail(v string)
	SetRegistrationDate(v time.Time)
	GetStatus() string
	GetAccountName() string
	GetRoles() role.Roles
	SetRoles(roles role.Roles)
}
