package userdocs

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

// The application's functions in the templates: admin.<name>(args…), its
// markdown; its error, the template's.
func TestFunc(t *testing.T) {
	b := &Builder{p: presets.New(i18n.New())}
	b.Func("roles", func(ctx *web.EventContext, args ...string) (string, error) {
		return "roles of " + strings.Join(args, "+") + " for " + ctx.R.URL.Path, nil
	})
	b.Func("broken", func(*web.EventContext, ...string) (string, error) { return "", errors.New("no roles") })
	r := &renderer{b: b, ctx: &web.EventContext{R: httptest.NewRequest("GET", "/admin/docs", nil)}}

	md, err := RenderTemplate(`{%= admin.roles("a", "b") %}`, r.globals())
	if err != nil || md != "roles of a+b for /admin/docs" {
		t.Errorf("%q %v", md, err)
	}
	if _, err := RenderTemplate(`{%= admin.broken() %}`, r.globals()); err == nil || !strings.Contains(err.Error(), "admin.broken: no roles") {
		t.Errorf("the error: %v", err)
	}
}
