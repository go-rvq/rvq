package features

import (
	"reflect"
	"strings"
	"testing"
)

// A feature is off until the configuration turns it on; nested keys are
// dotted; what is not a bool is refused; keys nobody declared are told.
func TestFeatures(t *testing.T) {
	f := Register("test.a", "A", "the a").CompiledIf(true)
	Register("test.b", "B", "").CompiledIf(true)
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

// A feature on needs the ones it depends on on: off with one off — whatever
// the configuration says of it —, in a chain; said Unmet; a cycle is off.
func TestDependsOn(t *testing.T) {
	site := Register("zz.dep.site", "", "").CompiledIf(true)
	forms := Register("zz.dep.forms", "", "").CompiledIf(true).DependsOn(site)
	drafts := Register("zz.dep.forms.drafts", "", "").CompiledIf(true).DependsOn(forms)
	t.Cleanup(func() {
		mu.Lock()
		for _, n := range []string{site.Name, forms.Name, drafts.Name, "zz.dep.a", "zz.dep.b"} {
			delete(declared, n)
			delete(enabled, n)
		}
		mu.Unlock()
	})
	Set(forms.Name, true)
	Set(drafts.Name, true)
	if forms.Enabled() || drafts.Enabled() {
		t.Error("on without the site")
	}
	if got := Unmet(); !reflect.DeepEqual(got, []string{"zz.dep.forms (needs zz.dep.site)", "zz.dep.forms.drafts (needs zz.dep.forms)"}) {
		t.Errorf("unmet: %v", got)
	}
	Set(site.Name, true)
	if !forms.Enabled() || !drafts.Enabled() || !Enabled(drafts.Name) {
		t.Error("off with the site on")
	}
	if len(Unmet()) != 0 {
		t.Errorf("unmet: %v", Unmet())
	}
	for _, s := range All() {
		if s.Name == drafts.Name && (!s.On || !s.Asked || len(s.Missing) != 0) {
			t.Errorf("state: %+v", s)
		}
	}
	// a cycle: off
	a := Register("zz.dep.a", "", "").CompiledIf(true)
	b := Register("zz.dep.b", "", "").CompiledIf(true).DependsOn(a)
	a.DependsOn(b)
	Set(a.Name, true)
	Set(b.Name, true)
	if a.Enabled() || b.Enabled() {
		t.Error("a cycle on")
	}
}

// A feature the binary has not (CompiledIf false) — or one it depends on
// has not — is none: off, whatever the configuration says; not listed; the
// configuration's key Unknown, not Unmet; its tag is its name.
func TestCompiled(t *testing.T) {
	off := Register("zz.cmp.off", "", "")
	dep := Register("zz.cmp.dep", "", "").CompiledIf(true).DependsOn(off)
	t.Cleanup(func() {
		mu.Lock()
		for _, n := range []string{off.Name, dep.Name} {
			delete(declared, n)
			delete(enabled, n)
		}
		mu.Unlock()
	})
	Set(off.Name, true)
	Set(dep.Name, true)
	if off.Compiled() || dep.Compiled() || off.Enabled() || dep.Enabled() {
		t.Errorf("not compiled: compiled=%v/%v on=%v/%v", off.Compiled(), dep.Compiled(), off.Enabled(), dep.Enabled())
	}
	for _, st := range All() {
		if st.Name == off.Name || st.Name == dep.Name {
			t.Errorf("listed: %s", st.Name)
		}
	}
	var unknown []string
	for _, k := range Unknown() {
		if strings.HasPrefix(k, "zz.cmp.") {
			unknown = append(unknown, k)
		}
	}
	if got := unknown; !reflect.DeepEqual(got, []string{"zz.cmp.dep", "zz.cmp.off"}) {
		t.Errorf("Unknown: %v", got)
	}
	if got := Unmet(); len(got) != 0 {
		t.Errorf("Unmet: %v", got)
	}
	if off.Tag() != "zz.cmp.off" {
		t.Errorf("tag: %s", off.Tag())
	}
	off.CompiledIf(true)
	if !dep.Compiled() || !dep.Enabled() {
		t.Error("compiled, still off")
	}
	if Enabled("zz.cmp.none") {
		t.Error("an undeclared feature on")
	}
}
