package schemaform

import (
	"context"
	"regexp"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

var portalAttr = regexp.MustCompile(`portal-name='(schemaform_portal_\d+)'`)

// The root of a form and of a detail embeds a portal of its own, and a
// component of the application — registered with Type or Display — gets its
// name (Context.PortalName). A listing cell embeds
// none, and its fields get "".
func TestPortal(t *testing.T) {
	b := New()
	b.Type("opener", func(c *Context) h.HTMLComponent { return h.Span("edit:" + c.PortalName()) })
	b.Display("opener", func(c *Context) h.HTMLComponent { return h.Span("show:" + c.PortalName()) })
	schema, err := b.Parse(`{x? opener}`)
	if err != nil {
		t.Fatal(err)
	}
	marshal := func(comp h.HTMLComponent) string {
		out, err := h.Marshal(comp, context.Background())
		if err != nil {
			t.Fatal(err)
		}
		return string(out)
	}
	value := map[string]any{"x": "1"}

	for name, got := range map[string]string{
		"edit": marshal(b.ComponentFunc(schema)(&presets.FieldContext{
			ToComponentOptions: &presets.ToComponentOptions{},
			Name:               "Value",
			FormKey:            "Value",
			Mode:               presets.FieldModeStack{presets.EDIT},
		}, &web.EventContext{})),
		"show": marshal(b.DetailComponent(&web.EventContext{}, schema, value)),
	} {
		m := portalAttr.FindStringSubmatch(got)
		if m == nil {
			t.Fatalf("%s: no portal embedded: %s", name, got)
		}
		mustContain(t, got, name+":"+m[1])
	}

	cell := marshal(b.ListComponent(&web.EventContext{}, schema, value))
	if portalAttr.MatchString(cell) {
		t.Errorf("a listing cell embeds a portal: %s", cell)
	}
	mustContain(t, cell, "show:<")

	if NewPortalName() == NewPortalName() {
		t.Error("two portals share a name")
	}
}
