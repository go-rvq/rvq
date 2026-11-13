// Package userref is a user as another record names it — who made it, who
// changed it —: its id and its name, read from the table of the users (Table)
// whatever the model of the users of the application is. A record keeps the
// id (a model.UUID column) and reads the user by a LEFT JOIN of a belongs-to
// association:
//
//	CreatedByID model.UUID    `gorm:"type:uuid"`
//	CreatedBy   *userref.Ref `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
//
//	db.Joins("CreatedBy").Find(&records)
//
// A Ref is never written through the record: save it with
// Omit(clause.Associations). The package user gives it as user.Ref; it is
// here, of no other package of the admin, so the packages user depends on
// (role…) may use it too.
package userref

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/x/login"
	"github.com/google/uuid"
	"github.com/sunfmin/reflectutils"
)

// Table is the table of the users (its columns id and name).
var Table = "users"

// Ref is a user as a record names it: its id and its name, read only.
type Ref struct {
	ID   model.UUID `gorm:"type:uuid;primaryKey"`
	Name string     `gorm:"->;-:migration"`
}

// TableName is Table.
func (Ref) TableName() string { return Table }

// String is the name; "" of none.
func (r *Ref) String() string {
	if r == nil {
		return ""
	}
	return r.Name
}

// IDOf is the id of the user of r (login.CurrentUserFromRequest: its field
// ID, a uuid.UUID, a model.UUID or a string of one); nobody, zero.
func IDOf(r *http.Request) (id model.UUID) {
	if r == nil {
		return
	}
	u := login.CurrentUserFromRequest(r)
	if u == nil {
		return
	}
	v, err := reflectutils.Get(u, "ID")
	if err != nil {
		return
	}
	switch t := v.(type) {
	case uuid.UUID:
		id = model.UUID(t)
	case model.UUID:
		id = t
	case string:
		if p, err := uuid.Parse(t); err == nil {
			id = model.UUID(p)
		}
	}
	return
}
