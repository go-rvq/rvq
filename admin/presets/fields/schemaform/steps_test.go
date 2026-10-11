package schemaform

import (
	"strings"
	"testing"
)

const stepsForm = `
class Start {
    name str
    count int
    methods {
        next(r) { return this.count < 10 ? "$1" : "done" }
    }
}
[steps, autosave]
class Form {
    [label="Início"]
    a Start
    [label="Antes de 10", next=(r) => r.a.count > 5 ? Form.done : nil]
    {
        valInB int
        other? str
    }
    [label="Conclusão"]
    done class {
        valInDone bool
    }
}
return Form`

// A form of steps: each field of its class a step — a group `{ … }` one, its
// fields the record's own —, labeled; the next by the step's [next=…], else
// its class's method next, else the order of the class; "" after the last.
func TestFlow(t *testing.T) {
	st, err := readClass(stepsForm)
	if err != nil {
		t.Fatal(err)
	}
	fl := st.Flow()
	if fl == nil || !fl.Autosave || len(fl.Steps) != 3 {
		t.Fatalf("flow: %+v", fl)
	}
	want := []struct {
		key, label string
		group      bool
		fields     int
	}{{"a", "Início", false, 1}, {"$1", "Antes de 10", true, 2}, {"done", "Conclusão", false, 1}}
	for i, w := range want {
		s := fl.Steps[i]
		if s.Key != w.key || s.Label != w.label || s.Group != w.group || len(s.Fields) != w.fields {
			t.Errorf("step %d: %s %q group=%v fields=%d", i, s.Key, s.Label, s.Group, len(s.Fields))
		}
	}
	for _, c := range []struct {
		key  string
		data map[string]any
		want string
	}{
		// the method of its class
		{"a", map[string]any{"a": map[string]any{"name": "x", "count": int64(3)}}, "$1"},
		{"a", map[string]any{"a": map[string]any{"name": "x", "count": int64(12)}}, "done"},
		// its [next=…]: a class field, or nil — the order
		{"$1", map[string]any{"a": map[string]any{"count": int64(7)}}, "done"},
		{"$1", map[string]any{"a": map[string]any{"count": int64(3)}}, "done"},
		{"done", nil, ""},
	} {
		got, err := fl.Next(c.key, c.data)
		if err != nil || got != c.want {
			t.Errorf("next of %s (%v): %q %v, want %q", c.key, c.data, got, err, c.want)
		}
	}
	if _, err := fl.Next("nope", nil); err == nil {
		t.Error("the next of no step")
	}

	// a step decoded alone: its fields, checked; the others not asked
	d := &Decoder{Schema: st}
	rec, errs := d.DecodeStep(fl, "$1", map[string]any{"valInB": "x"})
	if len(errs) != 1 || errs[0].Path != "valInB" || len(rec) != 2 {
		t.Errorf("the step $1: %v %v", rec, errs)
	}
	if _, errs := d.DecodeStep(fl, "a", map[string]any{"a": map[string]any{"name": "x", "count": "3"}}); len(errs) != 0 {
		t.Errorf("the step a: %v", errs)
	}

	// the way of an answer: the steps it goes through; what it skips not
	// asked
	path, err := fl.Path(d.Typed(map[string]any{"a": map[string]any{"name": "x", "count": "12"}}).Map())
	if err != nil || strings.Join(path, ",") != "a,done" {
		t.Errorf("the path: %v %v", path, err)
	}
	if _, errs := d.DecodeFlow(fl, map[string]any{"a": map[string]any{"name": "x", "count": "12"}, "done": map[string]any{"valInDone": true}}); len(errs) != 0 {
		t.Errorf("a step skipped asked: %v", errs)
	}
	if _, errs := d.DecodeFlow(fl, map[string]any{"a": map[string]any{"name": "x", "count": "3"}, "done": map[string]any{"valInDone": true}}); len(errs) == 0 {
		t.Error("a step of the way not asked")
	}

	// a form not of steps: none
	plain, _ := readClass("class Form { a str }; return Form")
	if plain.Flow() != nil {
		t.Error("a flow of a form not of steps")
	}
}

// A next that answers what is no step is the form's error.
func TestFlowNextNoStep(t *testing.T) {
	st, err := readClass(`[steps]
class Form {
    [next=(r) => "zz"]
    a str
    b str
}
return Form`)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.Flow().Next("a", nil); err == nil {
		t.Error("no error for a next of no step")
	}
	if fl := st.Flow(); fl.Autosave {
		t.Error("autosave without [autosave]")
	}
}
