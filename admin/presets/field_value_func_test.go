package presets

import (
	"context"
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/timezone"
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

// A field DefaultNow of a new record starts with the admin's clock, in the
// admin's time zone (a day: the admin's day); not with an error, not editing.
func TestDefaultNow(t *testing.T) {
	type rec struct {
		ID  uint
		Em  time.Time
		Dia time.Time
	}
	now := time.Date(2026, 6, 16, 1, 30, 0, 0, time.UTC) // 15/06 22:30 in São Paulo
	p := New(i18n.New()).NowFunc(func() time.Time { return now })
	mb := p.Model(&rec{})
	sp, _ := time.LoadLocation("America/Sao_Paulo")
	r := httptest.NewRequest("POST", "/", nil)
	ctx := &web.EventContext{R: r.WithContext(timezone.WithLocation(r.Context(), sp))}
	render2 := func(name string, comp FieldComponentFunc, modes FieldModeStack, errs ...string) string {
		f := &FieldBuilder{mode: ALL}
		f.name = name
		f.DefaultNow(true)
		fc := &FieldContext{Field: f, Obj: &rec{}, Name: name, FormKey: name, EventContext: ctx,
			ModelInfo: mb.Info(), Mode: modes, Errors: errs}
		return h.MustString(comp(fc, ctx), context.Background())
	}
	render := func(name string, comp FieldComponentFunc, mode FieldMode, errs ...string) string {
		return render2(name, comp, FieldModeStack{mode}, errs...)
	}
	if got := render("Em", DateTimeComponentFunc, NEW); !strings.Contains(got, "2026-06-15 22:30") {
		t.Errorf("an instant: now in the admin's time zone: %s", got)
	}
	if got := render("Dia", DateComponentFunc, NEW); !strings.Contains(got, "2026-06-15") || strings.Contains(got, "2026-06-16") {
		t.Errorf("a day: the admin's today: %s", got)
	}
	if got := render("Em", DateTimeComponentFunc, NEW, "vazio"); strings.Contains(got, "2026-06-15 22:30") {
		t.Errorf("refused with an error: not now again: %s", got)
	}
	// a new record's field drawn editable (the form of an action)
	if got := render2("Em", DateTimeComponentFunc, FieldModeStack{NEW, EDIT}); !strings.Contains(got, "2026-06-15 22:30") {
		t.Errorf("a new record, its field editable: now: %s", got)
	}
	if got := render("Em", DateTimeComponentFunc, EDIT); strings.Contains(got, "2026-06-15 22:30") {
		t.Errorf("editing: not now: %s", got)
	}
}
