package presets

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

type bitsRecord struct{ Reuse uint8 }

// A field of bits: a chip an option, with a `?` that opens its description;
// to edit, the chips toggle and the value is the sum of the marked bits; in
// the detail, the set ones are marked.
func TestBitsComponentFunc(t *testing.T) {
	opts := []BitOption{
		{Bit: 1, Label: func(context.Context) string { return "User" }, Description: func(context.Context) string { return "the user of it" }},
		{Bit: 2, Label: func(context.Context) string { return "Password" }, Description: func(context.Context) string { return "its password" }},
		{Bit: 4, Label: func(context.Context) string { return "Host" }},
	}
	f := BitsComponentFunc(opts...)
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/", nil)}
	render := func(mode FieldMode, ro bool) string {
		fc := &FieldContext{Name: "Reuse", FormKey: "Inbound.Reuse", Label: "Reuse", Obj: &bitsRecord{Reuse: 5},
			Mode: FieldModeStack{mode}, ReadOnly: ro}
		return h.MustString(f(fc, ctx), context.TODO())
	}

	edit := render(EDIT, false)
	for _, want := range []string{"<v-chip-group", "User", "Password", "Host", "the user of it", "its password",
		"mdi-help-circle-outline", "<v-menu", `fieldValue.value = (sel || []).reduce`, "[1, 2, 4].filter",
		"Inbound.Reuse"} {
		if !strings.Contains(edit, want) {
			t.Errorf("edit: no %q:\n%s", want, edit)
		}
	}
	if strings.Contains(edit, "<v-tooltip") {
		t.Error("edit: a tooltip")
	}
	if strings.Count(edit, "data-bit-help") != 2 {
		t.Errorf("a ? for each option with a description: %d", strings.Count(edit, "data-bit-help"))
	}

	detail := render(DETAIL, true)
	if strings.Contains(detail, "<v-chip-group") || strings.Count(detail, "data-bit-set") != 2 {
		t.Errorf("detail: the set bits (1, 4) marked, nothing to toggle:\n%s", detail)
	}
}
