package admin

import (
	"encoding/json"
	"reflect"
	"testing"
)

type opts struct {
	Layout string
	Order  int
}

type docWithStruct struct {
	ID      uint
	Title   string
	Options opts
	OptPtr  *opts
}

func TestStructFieldType(t *testing.T) {
	obj := &docWithStruct{}
	cases := map[string]reflect.Type{
		"Title":         reflect.TypeOf(""),
		"Options":       reflect.TypeOf(opts{}),
		"Options.Order": reflect.TypeOf(0),
		"OptPtr":        reflect.TypeOf((*opts)(nil)),
	}
	for field, want := range cases {
		got, ok := structFieldType(obj, field)
		if !ok || got != want {
			t.Errorf("structFieldType(%q) = %v, %v; want %v", field, got, ok, want)
		}
	}
	if _, ok := structFieldType(obj, "Nope"); ok {
		t.Error("unknown field should not resolve")
	}
}

// TestSetFieldFromJSONStruct is the regression for the reported panic: a struct
// field must be rebuilt as its own type, not left as a map[string]any.
func TestSetFieldFromJSONStruct(t *testing.T) {
	obj := &docWithStruct{}
	raw := json.RawMessage(`{"Layout":"grid","Order":3}`)
	if err := setFieldFromJSON(obj, "Options", raw); err != nil {
		t.Fatalf("setFieldFromJSON struct: %v", err)
	}
	if obj.Options.Layout != "grid" || obj.Options.Order != 3 {
		t.Fatalf("Options = %+v, want {grid 3}", obj.Options)
	}

	// Pointer-to-struct field.
	if err := setFieldFromJSON(obj, "OptPtr", json.RawMessage(`{"Layout":"list"}`)); err != nil {
		t.Fatalf("setFieldFromJSON ptr: %v", err)
	}
	if obj.OptPtr == nil || obj.OptPtr.Layout != "list" {
		t.Fatalf("OptPtr = %+v, want &{list 0}", obj.OptPtr)
	}

	// Plain scalar still works.
	if err := setFieldFromJSON(obj, "Title", json.RawMessage(`"hi"`)); err != nil {
		t.Fatalf("setFieldFromJSON scalar: %v", err)
	}
	if obj.Title != "hi" {
		t.Fatalf("Title = %q, want hi", obj.Title)
	}
}
