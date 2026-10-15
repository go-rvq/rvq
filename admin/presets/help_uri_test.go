package presets

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

// A drawer and a dialog of an overlay say the address that filled them (the
// help of what is in them).
func TestOverlayHelpURI(t *testing.T) {
	uri := "/admin/posts?__execute_event__=presets_Edit&id=7"
	ctx := &web.EventContext{R: httptest.NewRequest("POST", uri, nil)}
	attr := HelpURIAttr + `='` + uri + `'`
	render := func(r *web.EventResponse) string {
		if len(r.UpdatePortals) == 0 {
			t.Fatal("no portal updated")
		}
		return h.MustString(r.UpdatePortals[0].Body, context.Background())
	}

	p := New(i18n.New())
	var r web.EventResponse
	p.Drawer(actions.RightDrawer).SetHelpURI(HelpURI(ctx)).Respond(&r, h.Div(h.Text("form")))
	if s := render(&r); !strings.Contains(s, attr) {
		t.Errorf("the drawer: %s", s)
	}

	r = web.EventResponse{}
	p.Dialog().Respond(ctx, &r, h.Div(h.Text("form")))
	if s := render(&r); !strings.Contains(s, attr) {
		t.Errorf("the dialog: %s", s)
	}

	r = web.EventResponse{}
	p.Overlay(ctx, &r, h.Div(h.Text("x")), "")
	if s := render(&r); !strings.Contains(s, attr) {
		t.Errorf("the right drawer: %s", s)
	}
}
