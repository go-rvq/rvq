package user

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/packages/user/userref"
)

// Ref is a user as another record names it — who made it, who changed it —:
// its id and its name, read by a LEFT JOIN from the table of the users,
// whatever the model of the users is (userref). The way a record of rvq
// names a user.
type Ref = userref.Ref

// RefIDOf is the id of the user of r, to keep in a record beside its Ref;
// nobody, zero.
func RefIDOf(r *http.Request) model.UUID { return userref.IDOf(r) }
