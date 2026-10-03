package media

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// What the forms do with the files is asked of the model media_libraries,
// scope presets: uploading is @create, deleting @delete, describing @edit —
// by its unique name here, which decides before its groups.
func TestPermissionsOfTheModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	gone := uuid.New()
	allow := perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything)
	deny := func(res string) *perm.PolicyBuilder {
		return perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(perm.Anything).On(res)
	}
	pb := presets.New(i18n.New()).URIPrefix("/admin")
	pb.DataOperator(gorm2op.DataOperator(db)) // the key of a record
	pb.Permission(perm.New().Policies(allow,
		deny(":media_libraries:@create"),
		deny(":media_libraries:<"+gone.String()+">:@delete"),
		deny(":media_libraries:<"+gone.String()+">:@edit"),
	).SubjectsFunc(func(*http.Request) []string { return []string{"editor"} }))
	b := New(db)
	if err := b.Install(pb); err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("GET", "/admin/media-library", nil)
	kept, deleted := &media_library.MediaLibrary{}, &media_library.MediaLibrary{}
	kept.ID, deleted.ID = uuid.New(), gone

	if b.uploadIsAllowed(r) == nil {
		t.Error("uploading: allowed, denied @create")
	}
	if b.deleteIsAllowed(r, kept) != nil || b.deleteIsAllowed(r, deleted) == nil {
		t.Error("deleting: by the record's @delete")
	}
	if b.updateDescIsAllowed(r, kept) != nil || b.updateDescIsAllowed(r, deleted) == nil {
		t.Error("describing: by the record's @edit")
	}
}
