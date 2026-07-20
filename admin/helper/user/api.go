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
	GetName() string
	SetName(v string)
	SetEmail(v string)
	SetRegistrationDate(v time.Time)
	GetStatus() string
	GetAccountName() string
	GetRoles() role.Roles
	SetRoles(roles role.Roles)
}
