package dashboard

import (
	"testing"
)

// A dynamic widget of a group is known as a registered one, in its group's
// section on the panel — none of another group's.
func TestDynamic(t *testing.T) {
	b, pn := newTest()
	b.Group("g").Widget("static", nil).Dynamic(func() []*WidgetBuilder {
		return []*WidgetBuilder{NewWidget("g", "dyn"), NewWidget("other", "x")}
	})
	if pn.widget("g/dyn") == nil || pn.widget("other/x") != nil {
		t.Fatal("the dynamic widgets of the group")
	}
	var p Prefs
	if s := sectionOf(pn.layout(&p), "g/dyn"); s == nil || s.ID != "g" {
		t.Fatalf("the dynamic widget's section: %v", s)
	}
	pn.remove(&p, "g/dyn")
	if sectionOf(pn.layout(&p), "g/dyn") != nil {
		t.Fatal("removed, still there")
	}
}
