package integration_test

import (
	"context"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	. "github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

type switchRecord struct {
	Kind     string
	AURL     string
	AKey     string
	BURL     string
	Shared   string
	Note     string
	Advanced string
}

func switchBuilder(layout ...any) *FieldsBuilder {
	pb := New(i18n.New())
	pb.Model(&switchRecord{})
	return NewFieldsBuilder(pb).Model(&switchRecord{}).Defaults(NewFieldDefaults(WRITE)).Only(layout...)
}

func kindSwitch() *FieldsSwitch {
	return NewFieldsSwitch("engine", "form.Kind", func(obj any, _ *web.EventContext) string { return obj.(*switchRecord).Kind }).
		Case("a", "AURL", "AKey", "Shared").
		Case("b", "BURL", "Shared").
		Else("Note")
}

func renderFields(fb *FieldsBuilder, obj *switchRecord, mode FieldMode) string {
	ctx := &web.EventContext{R: httptest.NewRequest("POST", "/", nil)}
	return h.MustString(fb.ToComponent(nil, nil, obj, FieldModeStack{mode}, ctx), context.TODO())
}

// A switch in a form: every case drawn, a v-if / v-else-if / v-else chain on
// the expression — a field of two cases in both —; in the detail, only the
// case of the record.
func TestFieldsSwitch(t *testing.T) {
	fb := switchBuilder("Kind", kindSwitch())
	got := renderFields(fb, &switchRecord{Kind: "b"}, WRITE)
	for _, want := range []string{
		`v-if='(form.Kind) === "a"'`, `v-else-if='(form.Kind) === "b"'`, "v-else", "data-fields-switch='engine'",
		`form["AURL"]`, `form["AKey"]`, `form["BURL"]`, `form["Note"]`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("no %q:\n%s", want, got)
		}
	}
	if n := strings.Count(got, `v-model='form["Shared"]'`); n != 2 {
		t.Errorf("Shared, a field of two cases, drawn %d times", n)
	}
	// the chain in order: a, b, else
	ia, ib, ie := strings.Index(got, `=== "a"`), strings.Index(got, `=== "b"`), strings.Index(got, "<template v-else>")
	if !(ia < ib && ib < ie) {
		t.Errorf("the chain out of order: %d %d %d\n%s", ia, ib, ie, got)
	}

	// v-show: all mounted, the default on none of the values
	sw := kindSwitch()
	sw.Show = true
	got = renderFields(switchBuilder("Kind", sw), &switchRecord{}, WRITE)
	if strings.Contains(got, "v-if") || !strings.Contains(got, `v-show='(form.Kind) === "a"'`) ||
		!strings.Contains(got, `v-show='![&#34;a&#34;, &#34;b&#34;].includes(form.Kind)'`) && !strings.Contains(got, `v-show='!["a", "b"].includes(form.Kind)'`) {
		t.Errorf("v-show:\n%s", got)
	}

	// the detail: the case of the record, no condition of the form
	for kind, want := range map[string][]string{"a": {"AURL", "AKey", "Shared"}, "b": {"BURL", "Shared"}, "z": {"Note"}} {
		got = renderFields(fb, &switchRecord{Kind: kind}, DETAIL)
		if strings.Contains(got, "v-if") || strings.Contains(got, "v-else") {
			t.Errorf("%s: a condition of the form in the detail", kind)
		}
		for _, f := range []string{"AURL", "AKey", "BURL", "Note"} {
			if strings.Contains(got, f) != slices.Contains(want, f) {
				t.Errorf("detail of %s: %s drawn=%v", kind, f, strings.Contains(got, f))
			}
		}
	}

	// the active names: those of the case of the record
	l := FieldsLayout{"Kind", kindSwitch()}
	if got := l.ActiveNames(&switchRecord{Kind: "a"}, nil); !slices.Equal(got, []string{"Kind", "AURL", "AKey", "Shared"}) {
		t.Errorf("active of a: %v", got)
	}
	if got := l.ActiveNames(&switchRecord{Kind: "?"}, nil); !slices.Equal(got, []string{"Kind", "Note"}) {
		t.Errorf("active of none: %v", got)
	}
	// every name, for the listing and the setters
	if got := l.Names(); !slices.Equal(got, []string{"Kind", "AURL", "AKey", "Shared", "BURL", "Shared", "Note"}) {
		t.Errorf("names: %v", got)
	}
	// filtered: the switch kept, its cases filtered
	f, err := l.Filter(func(typ FieldLayoutEntryType, name string) bool { return name != "AKey" })
	if err != nil {
		t.Fatal(err)
	}
	if fs := f[1].(*FieldsSwitch); fs.Expr != "form.Kind" || slices.Contains(f.Names(), "AKey") || len(fs.Cases) != 2 {
		t.Errorf("filtered: %v", f.Names())
	}
	// the tree of the listing: a group of every case's fields
	if g, ok := l.ToGroup()[1].(*FieldsGroup); !ok || len(g.Items) != 6 {
		t.Errorf("to group: %#v", l.ToGroup()[1])
	}
}

// A group on a condition: v-if / v-show in a form; Condition on the server —
// in the detail, and in the form too.
func TestFieldsGroupCondition(t *testing.T) {
	advanced := &FieldsGroup{Name: "adv", Items: FieldsLayout{"Advanced"}, VIf: `form.Kind === "x"`}
	shown := &FieldsSection{Name: "notes", Rows: [][]string{{"Note"}}, VShow: "form.Kind"}
	fb := switchBuilder("Kind", advanced, shown)

	got := renderFields(fb, &switchRecord{}, WRITE)
	if !strings.Contains(got, `<template v-if='form.Kind === "x"'>`) || !strings.Contains(got, "v-show='form.Kind'") ||
		!strings.Contains(got, `form["Advanced"]`) {
		t.Errorf("form:\n%s", got)
	}
	got = renderFields(fb, &switchRecord{}, DETAIL)
	if strings.Contains(got, "v-if") || strings.Contains(got, "v-show") || !strings.Contains(got, "Advanced") {
		t.Errorf("detail:\n%s", got)
	}

	// Condition: not drawn, not active
	advanced.Condition = func(obj any, _ *web.EventContext) bool { return obj.(*switchRecord).Kind == "x" }
	for _, mode := range []FieldMode{WRITE, DETAIL} {
		if got := renderFields(fb, &switchRecord{Kind: "y"}, mode); strings.Contains(got, "Advanced") {
			t.Errorf("mode %v: a group its Condition refuses drawn", mode)
		}
	}
	if got := renderFields(fb, &switchRecord{Kind: "x"}, DETAIL); !strings.Contains(got, "Advanced") {
		t.Error("a group its Condition accepts not drawn")
	}
	if got := (FieldsLayout{advanced}).ActiveNames(&switchRecord{Kind: "y"}, nil); len(got) != 0 {
		t.Errorf("active: %v", got)
	}
	// kept through Filter and ToGroup
	f, _ := FieldsLayout{advanced, shown}.Filter(func(FieldLayoutEntryType, string) bool { return true })
	if g := f[0].(*FieldsGroup); g.VIf == "" || g.Condition == nil {
		t.Error("Filter lost the condition of a group")
	}
	if s := f[1].(*FieldsSection); s.VShow == "" {
		t.Error("Filter lost the condition of a section")
	}
	if g := (FieldsLayout{shown}).ToGroup()[0].(*FieldsGroup); g.VShow != "form.Kind" {
		t.Error("ToGroup lost the condition of a section")
	}

	// a tab on a condition: its head and its window
	tabs := NewFieldsTabs("t",
		&FieldsGroup{Name: "one", Title: func(context.Context) string { return "One" }, Items: FieldsLayout{"AURL"}},
		&FieldsGroup{Name: "two", Title: func(context.Context) string { return "Two" }, Items: FieldsLayout{"BURL"}, VShow: `form.Kind === "b"`},
	)
	got = renderFields(switchBuilder("Kind", tabs), &switchRecord{}, WRITE)
	if strings.Count(got, `v-show='form.Kind === "b"'`) != 2 {
		t.Errorf("tab on a condition:\n%s", got)
	}
}

// A case with nothing to draw is not drawn, and its value is not the
// default's: in a form, the default on none of the values (not a v-else); in
// the detail, nothing. Nothing in any case, no switch.
func TestFieldsSwitchEmptyCase(t *testing.T) {
	sw := kindSwitch().Case("c")
	fb := switchBuilder("Kind", sw)
	got := renderFields(fb, &switchRecord{Kind: "c"}, WRITE)
	if strings.Contains(got, `=== "c"`) || strings.Contains(got, "v-else>") ||
		!strings.Contains(got, `v-if='![&#34;a&#34;, &#34;b&#34;, &#34;c&#34;].includes(form.Kind)'`) && !strings.Contains(got, `v-if='!["a", "b", "c"].includes(form.Kind)'`) {
		t.Errorf("form, an empty case:\n%s", got)
	}
	if got := renderFields(fb, &switchRecord{Kind: "c"}, DETAIL); strings.Contains(got, "Note") || strings.Contains(got, "data-fields-switch") {
		t.Errorf("detail of the empty case:\n%s", got)
	}

	empty := NewFieldsSwitch("none", "form.Kind", func(obj any, _ *web.EventContext) string { return obj.(*switchRecord).Kind }).Case("a")
	if got := renderFields(switchBuilder("Kind", empty), &switchRecord{Kind: "a"}, WRITE); strings.Contains(got, "data-fields-switch") {
		t.Errorf("an empty switch drawn:\n%s", got)
	}
}
