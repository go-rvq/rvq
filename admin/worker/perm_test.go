package worker

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
)

// Doing a job of a kind is an action of the model of the jobs, scope
// presets: "!" and the job's name — allowed or denied kind by kind.
func TestJobPermissionIsAnActionOfTheModel(t *testing.T) {
	pb := presets.New(i18n.New()).URIPrefix("/admin")
	pb.Permission(perm.New().Policies(
		perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything),
		perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(perm.Anything).On(":jobs:!upload_posts"),
	).SubjectsFunc(func(*http.Request) []string { return []string{"editor"} }))
	saved := jobsModel
	defer func() { jobsModel = saved }()
	jobsModel = pb.Model(&Job{})
	r := httptest.NewRequest("GET", "/admin/workers", nil)
	if editIsAllowed(r, "UploadPosts") == nil {
		t.Error("a kind denied by the unique name of the model was allowed")
	}
	if err := editIsAllowed(r, "export_books"); err != nil {
		t.Errorf("another kind: %v", err)
	}
	if jobsModel.UniquePermName() != "jobs" {
		t.Errorf("the model %q", jobsModel.UniquePermName())
	}
}
