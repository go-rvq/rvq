package role

import (
	"time"

	"github.com/qor5/x/v3/perm"
)

type Role struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string                  `admin:"required" gorm:"unique"`
	Permissions []*perm.DefaultDBPolicy `gorm:"foreignKey:ReferID"`
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
