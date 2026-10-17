// Package features is what an installation turns on: a functionality a
// module declares (Register) — off unless the configuration says so —, the
// `features:` of the app.yml (admin/app/env, Load).
//
//	features:
//	  hermon.forms.drafts: true
//
// A key is its module and what it is for, "." apart; nested maps are the same
// keys (`hermon: {forms: {drafts: true}}`). A key nobody declared is logged
// once the application asks about the features (Unknown), a mistyped one
// being otherwise silent.
//
// A feature may depend on others (DependsOn): it is on only when they are —
// whatever the configuration says of it —, so a module asks Enabled alone.
// One the configuration turns on and a dependency keeps off is Unmet.
//
// # Compiled
//
// A feature is on only in a binary it was compiled into: its build tag is
// its name (Tag) — `go build -tags hermon.site,hermon.forms` —, and the
// package that declares it says so (CompiledIf), by a constant of a file of
// that tag:
//
//	// compiled_on.go
//	//go:build hermon.forms
//	const Compiled = true
//
//	// compiled_off.go
//	//go:build !hermon.forms
//	const Compiled = false
//
//	var Feature = features.Register("hermon.forms", …).CompiledIf(Compiled)
//
//	func Enabled() bool { return Compiled && Feature.Enabled() }
//
// The constant false, what is only reached by Enabled is no code of the
// binary: the compiler drops the branch, the linker what only it called.
//
// A feature the binary has not — or one it depends on has not — is as if it
// were never declared: never on, never listed (All), never named; the
// configuration that turns it on has a key of no feature (Unknown).
package features

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Feature is a functionality an installation may turn on.
type Feature struct {
	// Name is its key in `features:`.
	Name string
	// Label and Hint say what it is, in the documentation of the server.
	Label, Hint string
	// deps are the features it needs (DependsOn).
	deps []*Feature
	// compiled: the binary has it (CompiledIf).
	compiled bool
}

// CompiledIf says whether the binary has the feature — its package's
// constant of the build tag (Tag) —: one not compiled is never on.
func (f *Feature) CompiledIf(compiled bool) *Feature {
	mu.Lock()
	defer mu.Unlock()
	f.compiled = compiled
	return f
}

// Compiled says the binary has the feature (CompiledIf) and the features it
// depends on.
func (f *Feature) Compiled() bool {
	if f == nil {
		return false
	}
	mu.RLock()
	defer mu.RUnlock()
	return built(f, map[*Feature]bool{})
}

// built is Compiled, mu held; seen: the features asked on the way (a cycle
// is not).
func built(f *Feature, seen map[*Feature]bool) bool {
	if f == nil || !f.compiled || seen[f] {
		return false
	}
	seen[f] = true
	defer delete(seen, f)
	for _, d := range f.deps {
		if !built(d, seen) {
			return false
		}
	}
	return true
}

// Tag is the build tag that compiles the feature: its name.
func (f *Feature) Tag() string { return f.Name }

// DependsOn says the feature needs deps: it is on only when they are too.
func (f *Feature) DependsOn(deps ...*Feature) *Feature {
	mu.Lock()
	defer mu.Unlock()
	f.deps = append(f.deps, deps...)
	return f
}

// Dependencies are the features it needs.
func (f *Feature) Dependencies() []*Feature {
	mu.RLock()
	defer mu.RUnlock()
	return append([]*Feature(nil), f.deps...)
}

var (
	mu       sync.RWMutex
	declared = map[string]*Feature{}
	enabled  = map[string]bool{}
)

// Register declares a feature: off unless the configuration turns it on. A
// name declared twice panics.
func Register(name, label, hint string) *Feature {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := declared[name]; ok {
		panic(fmt.Sprintf("features: %q declared twice", name))
	}
	f := &Feature{Name: name, Label: label, Hint: hint}
	declared[name] = f
	return f
}

// Enabled says the configuration turned it on, and the features it depends
// on are on.
func (f *Feature) Enabled() bool {
	if f == nil {
		return false
	}
	return Enabled(f.Name)
}

// Enabled says the feature name is declared and compiled, the configuration
// turned it on, and the features it depends on are on.
func Enabled(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	return on(name, map[string]bool{})
}

// on is Enabled, mu held; seen: the features asked on the way (a cycle is
// off).
func on(name string, seen map[string]bool) bool {
	f := declared[name]
	if seen[name] || !enabled[name] || f == nil || !built(f, map[*Feature]bool{}) {
		return false
	}
	seen[name] = true
	defer delete(seen, name)
	for _, d := range f.deps {
		if !on(d.Name, seen) {
			return false
		}
	}
	return true
}

// missing are the dependencies of f that keep it off — mu held —, by name.
func missing(f *Feature) (out []string) {
	for _, d := range f.deps {
		if !on(d.Name, map[string]bool{}) {
			out = append(out, d.Name)
		}
	}
	return
}

// Set turns the feature name on or off: what Load does, and a test.
func Set(name string, on bool) {
	mu.Lock()
	defer mu.Unlock()
	enabled[name] = on
}

// Load reads the `features:` of a configuration: its keys, flattened ("."
// apart), each a bool (true, "true", 1…). What is not one is an error.
func Load(m map[string]any) error {
	flat := map[string]bool{}
	if err := flatten("", m, flat); err != nil {
		return err
	}
	mu.Lock()
	defer mu.Unlock()
	for k, v := range flat {
		enabled[k] = v
	}
	return nil
}

func flatten(prefix string, m map[string]any, out map[string]bool) error {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch x := v.(type) {
		case map[string]any:
			if err := flatten(key, x, out); err != nil {
				return err
			}
		case bool:
			out[key] = x
		case nil:
			out[key] = false
		case int:
			out[key] = x != 0
		case string:
			b, err := strconv.ParseBool(x)
			if err != nil {
				return fmt.Errorf("features: %s: %q is not a bool", key, x)
			}
			out[key] = b
		default:
			return fmt.Errorf("features: %s: %v is not a bool", key, v)
		}
	}
	return nil
}

// State is a feature declared, compiled, and whether it is on: Asked, the
// configuration's word; On, Asked and its dependencies on; Missing, the
// dependencies off.
type State struct {
	*Feature
	On      bool
	Asked   bool
	Missing []string
}

// All are the features declared the binary has (Compiled), by name: one it
// has not is none.
func All() (out []State) {
	mu.RLock()
	defer mu.RUnlock()
	for _, f := range declared {
		if !built(f, map[*Feature]bool{}) {
			continue
		}
		out = append(out, State{Feature: f, On: on(f.Name, map[string]bool{}), Asked: enabled[f.Name], Missing: missing(f)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return
}

// Unmet are the features the configuration turns on, compiled, that a
// dependency keeps off: "name (needs dep, …)".
func Unmet() (out []string) {
	mu.RLock()
	defer mu.RUnlock()
	for name, f := range declared {
		if enabled[name] && built(f, map[*Feature]bool{}) {
			if m := missing(f); len(m) > 0 {
				out = append(out, fmt.Sprintf("%s (needs %s)", name, strings.Join(m, ", ")))
			}
		}
	}
	sort.Strings(out)
	return
}

// Unknown are the keys the configuration has that no feature declared — or
// the binary has not (Compiled).
func Unknown() (out []string) {
	mu.RLock()
	defer mu.RUnlock()
	for k := range enabled {
		if f, ok := declared[k]; !ok || !built(f, map[*Feature]bool{}) {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return
}
