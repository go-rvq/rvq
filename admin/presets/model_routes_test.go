package presets

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

type routesProduct struct {
	ID   uint
	Name string
}

// Creating and editing are governed by DIFFERENT flags: a model may allow one
// and forbid the other. The pages they mount have to follow the matching flag —
// `/{listing}/new` renders the CREATE form, `/{listing}/{id}/edit` the edit one.
func TestModelRoutesCreatingEditingFlags(t *testing.T) {
	cases := []struct {
		name             string
		creatingDisabled bool
		editingDisabled  bool
		wantNew          bool
		wantEdit         bool
	}{
		{name: "both enabled", wantNew: true, wantEdit: true},
		{name: "creating disabled", creatingDisabled: true, wantNew: false, wantEdit: true},
		{name: "editing disabled", editingDisabled: true, wantNew: true, wantEdit: false},
		{name: "both disabled", creatingDisabled: true, editingDisabled: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := New(i18n.New())
			mb := b.Model(&routesProduct{}).URIName("products")
			mb.SetCreatingDisabled(c.creatingDisabled)
			mb.SetEditingDisabled(c.editingDisabled)

			mux := http.NewServeMux()
			mb.SetupRoutes(mux)

			mounted := func(path string) bool {
				_, pattern := mux.Handler(httptest.NewRequest("GET", path, nil))
				return pattern != ""
			}

			if got := mounted("/products/new"); got != c.wantNew {
				t.Errorf("/products/new mounted = %v, want %v", got, c.wantNew)
			}
			if got := mounted("/products/1/edit"); got != c.wantEdit {
				t.Errorf("/products/1/edit mounted = %v, want %v", got, c.wantEdit)
			}
		})
	}
}

// The page func carries the same guard as the route: `GetCreatingPageFunc` is
// public, so an application may mount it anywhere.
func TestCreatingPageFuncFollowsCreatingFlag(t *testing.T) {
	cases := []struct {
		name             string
		creating         bool
		creatingDisabled bool
		editingDisabled  bool
		wantErr          error
	}{
		{name: "create allowed", creating: true},
		{name: "create forbidden", creating: true, creatingDisabled: true, wantErr: ErrCreateRecordNotAllowed},
		{name: "create is not edit", creating: true, editingDisabled: true},
		{name: "edit forbidden", editingDisabled: true, wantErr: ErrUpdateRecordNotAllowed},
		{name: "edit is not create", creatingDisabled: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := New(i18n.New())
			mb := b.Model(&routesProduct{}).URIName("products")
			mb.SetCreatingDisabled(c.creatingDisabled)
			mb.SetEditingDisabled(c.editingDisabled)

			ctx := &web.EventContext{R: httptest.NewRequest("GET", "/products/new", nil)}
			_, err := mb.Editing().DefaultPageFuncMode(c.creating, ctx)

			// only the guard is under test: with it open the page func goes on to
			// fetch/render, which needs a data operator this test has no reason to
			// provide — any OTHER error is fine.
			if c.wantErr != nil {
				if err != c.wantErr {
					t.Errorf("err = %v, want %v", err, c.wantErr)
				}
			} else if err == ErrCreateRecordNotAllowed || err == ErrUpdateRecordNotAllowed {
				t.Errorf("err = %v, want the guard to let it through", err)
			}
		})
	}
}
