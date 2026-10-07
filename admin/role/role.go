package role

import (
	"time"

	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
)

type Role struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string `admin:"required" gorm:"unique"`
	// SystemKey, when set, marks a role the application needs (SystemRole):
	// not deleted nor renamed, its permissions reset to its originals.
	SystemKey string `gorm:"index"`
	// SystemPolicies are the originals of a system role it was last given
	// (policyKey, one a line): what the next version of them changes —
	// the originals it gained added, those it lost taken — leaving what was
	// changed by hand.
	SystemPolicies string `gorm:"type:text" admin:"-"`
	// Permissions are the policies whose ReferID is this role's key. The column
	// is text — it also holds free-form references — so the relation is read
	// but creates no foreign key (-:migration).
	Permissions []*perm.DefaultDBPolicy `gorm:"-:migration;foreignKey:ReferID"`
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
