package presets

import (
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

func textComp(s string) FieldComponentFunc {
	return func(*FieldContext, *web.EventContext) h.HTMLComponent { return h.Text(s) }
}

// A wrap does not give a field without a component one: the field stays
// skipped, and the wrap waits for the component it gets.
func TestWrapComponentFuncWithoutComponent(t *testing.T) {
	f := &FieldBuilder{}
	f.WrapComponentFunc(func(old FieldComponentFunc) FieldComponentFunc {
		if old == nil {
			t.Fatal("old is nil")
		}
		return func(fc *FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(old(fc, ctx))
		}
	})
	if f.compFunc != nil {
		t.Fatal("the wrap gave the field a component: it would not be skipped")
	}

	c := f.Clone()
	c.ComponentFunc(textComp("x"))
	if f.compFunc != nil || len(f.compWraps) != 1 {
		t.Fatal("setting the clone's component changed the original")
	}
	got := strings.TrimSpace(h.MustString(c.compFunc(nil, nil), t.Context()))
	if got != "<div>x</div>" {
		t.Fatalf("the wrap was not applied: %q", got)
	}
	if len(c.compWraps) != 0 {
		t.Fatal("the wraps are applied once")
	}
}
