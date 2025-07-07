package user

import (
	"time"

	"github.com/qor5/admin/v3/role"
	"github.com/qor5/x/v3/login"
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
