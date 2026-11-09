package presets

import (
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
)

// The headers of the listing's layout and its columns: only the fields the
// request may read ("#Field" of the model) — a header of a field not read
// would stand over the cells of the next one.
func TestListingHeadersByPermission(t *testing.T) {
	_, posts, _ := uniqueApp([]string{"reader"},
		perm.PolicyFor("reader").WhoAre(perm.Allowed).ToDo(PermList).On("admin:posts:@list"),
		perm.PolicyFor("reader").WhoAre(perm.Allowed).ToDo(PermList).On("admin:posts:#Title:*"),
	)
	lb := posts.Listing("ID", "Title")
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/admin/posts", nil)}
	fields, headers := lb.layoutFields(ctx)
	if len(fields) != 1 || fields[0].name != "Title" || len(headers) != 1 {
		t.Fatalf("the fields %d, the headers %d", len(fields), len(headers))
	}
	if cols := lb.columns(ctx, fields); len(cols) != 1 || cols[0].Name != "Title" {
		t.Errorf("the columns: %+v", cols)
	}
}
