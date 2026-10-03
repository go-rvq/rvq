package perm_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/gorm"
)

// fakeDBPolicies are the policies "of the database": what LoadDBPolicies
// gives.
type fakeDBPolicies struct{ policies *[]*perm.PolicyBuilder }

func (f fakeDBPolicies) LoadDBPolicies(*gorm.DB, *time.Time) ([]*perm.PolicyBuilder, []*perm.PolicyBuilder) {
	return *f.policies, nil
}

// With LoadFrequency(0) the policies of the database are loaded by
// ReloadDBPolicies alone — a LISTEN announcing their changes —: each reload
// brings what is there and drops from memory what is not any more; the
// policies of the code stay.
func TestReloadDBPolicies(t *testing.T) {
	var db []*perm.PolicyBuilder
	allow := func(id, res string) *perm.PolicyBuilder {
		return perm.PolicyFor("editor").WhoAre(perm.Allowed).ToDo(perm.Anything).On(res).ID(id)
	}
	p := perm.New().Policies(allow("code", "app:code:*")).SubjectsFunc(sf("editor")).
		DBPolicy(perm.NewDBPolicy(nil).Model(fakeDBPolicies{&db}).LoadFrequency(0))
	allowed := func(res string) bool {
		return perm.NewVerifier("app", p).Spawn().On(res).Do("@get").WithReq(httptest.NewRequest("GET", "/", nil)).Allowed()
	}
	if allowed("posts") {
		t.Fatal("a policy of the database before any reload")
	}
	db = []*perm.PolicyBuilder{allow("a", "app:posts:*"), allow("b", "app:pages:*")}
	p.ReloadDBPolicies(nil)
	if !allowed("posts") || !allowed("pages") || !allowed("code") {
		t.Fatal("not loaded")
	}
	// one deleted from the database: gone from memory; the code's stays
	db = db[:1]
	p.ReloadDBPolicies(nil)
	if !allowed("posts") || allowed("pages") || !allowed("code") {
		t.Errorf("after a deletion: posts %v pages %v code %v", allowed("posts"), allowed("pages"), allowed("code"))
	}
}
