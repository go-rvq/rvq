package user

import (
	"time"

	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/x/login"
)

type User interface {
	login.UserPasser
	GetID() uint
	GetName() string
	SetName(v string)
	SetEmail(v string)
	SetRegistrationDate(v time.Time)
	GetStatus() string
	GetAccountName() string
	GetRoles() role.Roles
	SetRoles(roles role.Roles)
}
