package role

import (
	"time"

	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Role struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string `admin:"required" gorm:"unique"`
	// SystemKey, when set, marks a role the application needs (SystemRole):
	// not deleted nor renamed, its permissions reset to its originals.
	SystemKey string `gorm:"index"`
	// SystemPolicies are the originals of a system role it was last given —
	// its lines, one a line —: what the next version of them changes — the
	// originals it gained added, those it lost taken — leaving what was
	// changed by hand.
	SystemPolicies string `gorm:"type:text" admin:"-"`
	// Permissions are its lines, in order: a resource each, maybe after a
	// prefix — "deny ", "allow ", "reset " — (perm.Stack). In the stack of a
	// user they come after the code's and after those of the user's roles
	// before it.
	Permissions pq.StringArray `gorm:"type:text[]"`
	// Policies are the permissions of the roles of before the lines: the
	// policies whose ReferID is its key — made lines once (MigrateLines).
	// Read by hand, not a relation: gorm would give ReferID — text, it also
	// holds free-form references — the type of the key (uuid).
	Policies []*perm.DefaultDBPolicy `gorm:"-"`
}

func (r *Role) String() string {
	return r.Name
}

type Roles []*Role

func (s Roles) Contains(name string) bool {
	for _, r := range s {
		if r.Name == name {
			return true
		}
	}
	return false
}

func (s Roles) Names() []string {
	names := make([]string, len(s))
	for i, r := range s {
		names[i] = r.Name
	}
	return names
}
func (s Roles) FirstName() string {
	if len(s) == 0 {
		return ""
	}
	return s[0].Name
}
