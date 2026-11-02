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
	// Anonymous reports whether this is the anonymous user — no id
	// (AnonymousID: uuid.Nil), never a row of the users. There is always a
	// user on a request: when the session carries none, GetCurrentUser
	// returns the anonymous user (in memory) instead of nil; what it does is
	// recorded with no author (NULL).
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
