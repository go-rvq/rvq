package dashboard

import (
	"slices"
	"testing"
)

// A dynamic widget is known as a registered one; an opt-in one is on the
// dashboard only once added, and leaves when removed.
func TestDynamicOptIn(t *testing.T) {
	b := New(nil, nil)
	b.Group(&Group{ID: "g"})
	b.Widget(&Widget{ID: "static", Group: "g"})
	b.Dynamic(func() []*Widget {
		return []*Widget{{ID: "dyn", Group: "g", OptIn: true}}
	})
	if b.widget("dyn") == nil {
		t.Fatal("the dynamic widget is not known")
	}
	in := func(p *Prefs, id string) bool { return sectionOf(b.layout(p), id) != nil }

	var p Prefs
	if !in(&p, "static") || in(&p, "dyn") {
		t.Fatalf("before adding: static %v, dyn %v", in(&p, "static"), in(&p, "dyn"))
	}
	b.add(&p, "dyn")
	if !in(&p, "dyn") || !slices.Contains(p.Added, "dyn") {
		t.Fatalf("added: %v %v", in(&p, "dyn"), p.Added)
	}
	b.remove(&p, "dyn")
	if in(&p, "dyn") || slices.Contains(p.Added, "dyn") || slices.Contains(p.Removed, "dyn") {
		t.Fatalf("removed: %v %v %v", in(&p, "dyn"), p.Added, p.Removed)
	}
	if !b.moveTo(&p, "g", "dyn") || !in(&p, "dyn") || !slices.Contains(p.Added, "dyn") {
		t.Fatalf("moved to a section: %v %v", in(&p, "dyn"), p.Added)
	}
}
