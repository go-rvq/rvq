package schemaform

import (
	"context"
	"net/url"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

// An optional bool (`x? bool`) is a select of yes and no that may be left
// empty — nil, not false —; a required one, a switch.
func TestOptionalBool(t *testing.T) {
	s, err := Parse(`interface Form { on? bool; must bool; [required] chosen? bool }`)
	if err != nil {
		t.Fatal(err)
	}
	b := New()
	html := func(name string) string {
		f := s.FieldAt(name)
		c := &Context{Field: f, Value: "form." + name, Path: name, Builder: b, Event: &web.EventContext{}}
		return h.MustString(BoolComponentFunc(c), context.Background())
	}
	if got := html("on"); !strings.Contains(got, "v-select") || !strings.Contains(got, ":clearable='true'") {
		t.Errorf("optional: %s", got)
	}
	if got := html("must"); !strings.Contains(got, "v-switch") || !strings.Contains(got, "primary") {
		t.Errorf("required: %s", got)
	}
	if got := html("chosen"); !strings.Contains(got, "v-select") || !strings.Contains(got, ":clearable='false'") || !strings.Contains(got, " required ") {
		t.Errorf("[required] optional: %s", got)
	}
	for raw, want := range map[string]any{"": nil, "true": true, "false": false} {
		v, err := b.DecodeForm(&web.EventContext{}, s, url.Values{"k.on": {raw}, "k.must": {"true"}, "k.chosen": {"false"}}, "k")
		if err != nil {
			t.Fatal(err)
		}
		if got := recordValue(v, "on"); got != want {
			t.Errorf("%q: %#v, want %#v", raw, got, want)
		}
	}
	if _, err := b.DecodeForm(&web.EventContext{}, s, url.Values{"k.chosen": {""}}, "k"); err == nil {
		t.Error("[required] optional left empty: taken")
	}
	v, _ := b.DecodeForm(&web.EventContext{}, s, url.Values{"k.chosen": {"true"}}, "k")
	if got := recordValue(v, "on"); got != nil {
		t.Errorf("not posted: %#v", got)
	}
	if got := recordValue(v, "must"); got != false {
		t.Errorf("a switch not posted: %#v", got)
	}
}
