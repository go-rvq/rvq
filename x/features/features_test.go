package features

import (
	"reflect"
	"testing"
)

// A feature is off until the configuration turns it on; nested keys are
// dotted; what is not a bool is refused; keys nobody declared are told.
func TestFeatures(t *testing.T) {
	f := Register("test.a", "A", "the a")
	Register("test.b", "B", "")
	if f.Enabled() {
		t.Fatal("on before Load")
	}
	if err := Load(map[string]any{
		"test.a": true,
		"test":   map[string]any{"b": "false"},
		"tset.a": 1,
	}); err != nil {
		t.Fatal(err)
	}
	if !f.Enabled() || Enabled("test.b") {
		t.Errorf("a=%v b=%v", f.Enabled(), Enabled("test.b"))
	}
	if got := Unknown(); !reflect.DeepEqual(got, []string{"tset.a"}) {
		t.Errorf("Unknown: %v", got)
	}
	var names []string
	for _, s := range All() {
		if s.Name == "test.a" || s.Name == "test.b" {
			names = append(names, s.Name)
		}
	}
	if !reflect.DeepEqual(names, []string{"test.a", "test.b"}) {
		t.Errorf("All: %v", names)
	}
	if err := Load(map[string]any{"test.a": "maybe"}); err == nil {
		t.Error("a string not a bool taken")
	}
	var none *Feature
	if none.Enabled() {
		t.Error("nil enabled")
	}
	defer func() {
		if recover() == nil {
			t.Error("declared twice, no panic")
		}
	}()
	Register("test.a", "", "")
}
