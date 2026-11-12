package accesskeys

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func TestMigrateUngroup(t *testing.T) {
	db := testDB(t)
	k, _ := newKey(t, db, uuid.New(), func(k *AccessKey) {
		k.Permissions = pq.StringArray{"deny admin:content/:posts:@delete", "deny admin:content/:*"}
	})
	if err := MigrateUngroup(db); err != nil {
		t.Fatal(err)
	}
	var got AccessKey
	db.First(&got, "id = ?", k.ID)
	if s, want := strings.Join(got.Permissions, ";"), "deny admin:posts:@delete;deny admin:content/:*"; s != want {
		t.Errorf("%q, want %q", s, want)
	}
}
