package presets

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/web"
)

// A page under another (SubPage): its URL under the other's wherever the menu
// puts it, a pattern of the mux ("{user}"); out of the menu tree; the other's
// title; no permission of its own (Parent's asked).
func TestSubPage(t *testing.T) {
	b := menuBuilder().URIPrefix("/admin")
	pb := b.PagesRegistrator().New(HttpPage("/files").
		TitleFunc(func(context.Context) string { return "Files" })).
		Raw(func(*web.EventContext) (web.PageResponse, error) { return web.PageResponse{}, nil })
	var got string
	sub := pb.SubPage("/u/{user}").Raw(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		got = ctx.R.PathValue("user")
		return
	})

	a := b.MenuGroup("a")
	deep := b.MenuGroup("deep")
	a.Add(deep)
	for _, c := range []struct {
		dst  *MenuGroupBuilder
		want string
	}{
		{nil, "/admin/files/u/{user}"},
		{deep, "/admin/a/deep/files/u/{user}"},
	} {
		if err := b.MoveMenuItem(PageItem("/files"), c.dst); err != nil {
			t.Fatal(err)
		}
		sub.Page().Build("admin")
		if p := sub.Page().FullPath(); p != c.want {
			t.Errorf("FullPath() = %q, want %q", p, c.want)
		}
	}
	if sub.Page().Parent() != pb.Page() || sub.Page().IsInMenu() || sub.Page().GetVerifier() != nil {
		t.Error("not the page's: its parent, out of the menu, no permission of its own")
	}
	if sub.Page().TTitle(context.Background()) != "Files" {
		t.Errorf("the title: %q", sub.Page().TTitle(context.Background()))
	}
	// a key of its own: not in the tree
	if b.menuGroupOf(PageItem("/files/u/{user}")) != nil {
		t.Error("in the menu tree")
	}

	// served by its pattern
	pb.Build()
	sub.Build()
	mux := http.NewServeMux()
	b.Build(mux)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("GET", "/admin/a/deep/files/u/42", nil))
	if w.Code != http.StatusOK || got != "42" {
		t.Errorf("served: %d, user %q", w.Code, got)
	}
}
