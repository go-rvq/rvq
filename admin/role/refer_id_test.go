package role

import (
	"sync"
	"testing"

	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/gorm/schema"
)

// A role's policies are no relation of gorm: it would give their ReferID —
// text, free-form references too ("sharedaccess:…") — the type of the key
// of a role (uuid), and the column made of it would refuse them.
func TestPolicyReferIDStaysText(t *testing.T) {
	cache := &sync.Map{}
	ns := schema.NamingStrategy{}
	if _, err := schema.Parse(&Role{}, cache, ns); err != nil {
		t.Fatal(err)
	}
	s, err := schema.Parse(&perm.DefaultDBPolicy{}, cache, ns)
	if err != nil {
		t.Fatal(err)
	}
	if dt := s.LookUpField("ReferID").DataType; dt != schema.String {
		t.Errorf("ReferID: %q, want %q", dt, schema.String)
	}
}
