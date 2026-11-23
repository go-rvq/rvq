package presets

import (
	"context"
	h "github.com/go-rvq/htmlgo"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/web"
)

// A virtual field — no struct field of its name — read by its ValueFunc: its
// Required reports it empty, and not when the record has it.
func TestFieldValueFuncRequired(t *testing.T) {
	type rec struct{ OrigemID string }
	f := &FieldBuilder{mode: ALL}
	f.name = "Origem"
	f.Required(true).ValueFunc(func(obj any) any { return obj.(*rec).OrigemID })
	validate := func(r *rec) web.ValidationErrors {
		ctx := &web.EventContext{R: httptest.NewRequest("POST", "/", nil)}
		return f.Validate(&FieldContext{Field: f, Obj: r, Name: "Origem", FormKey: "Origem", EventContext: ctx})
	}
	if verr := validate(&rec{}); !verr.HaveErrors() || len(verr.GetFieldErrors("Origem")) == 0 {
		t.Errorf("an empty virtual field required: %v", verr)
	}
	if verr := validate(&rec{OrigemID: "a"}); verr.HaveErrors() {
		t.Errorf("a virtual field told: %v", verr)
	}
}

// The input of an instant shows the errors of its field.
func TestDateTimeComponentFuncErrors(t *testing.T) {
	type rec struct{ Data time.Time }
	f := &FieldBuilder{mode: ALL}
	f.name = "Data"
	ctx := &web.EventContext{R: httptest.NewRequest("POST", "/", nil)}
	fc := &FieldContext{Field: f, Obj: &rec{}, Name: "Data", FormKey: "Data", EventContext: ctx,
		Errors: []string{"obrigatório"}}
	html := h.MustString(DateTimeComponentFunc(fc, ctx), context.Background())
	if !strings.Contains(html, "obrigatório") {
		t.Errorf("the errors of the field not shown: %s", html)
	}
}
