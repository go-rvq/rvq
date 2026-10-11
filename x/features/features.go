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
package features

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

// Feature is a functionality an installation may turn on.
type Feature struct {
	// Name is its key in `features:`.
	Name string
	// Label and Hint say what it is, in the documentation of the server.
	Label, Hint string
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

// Enabled says the configuration turned it on.
func (f *Feature) Enabled() bool {
	if f == nil {
		return false
	}
	return Enabled(f.Name)
}

// Enabled says the configuration turned the feature name on.
func Enabled(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	return enabled[name]
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

// State is a feature declared and whether it is on.
type State struct {
	*Feature
	On bool
}

// All are the features declared, by name.
func All() (out []State) {
	mu.RLock()
	defer mu.RUnlock()
	for _, f := range declared {
		out = append(out, State{Feature: f, On: enabled[f.Name]})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return
}

// Unknown are the keys the configuration has that no feature declared.
func Unknown() (out []string) {
	mu.RLock()
	defer mu.RUnlock()
	for k := range enabled {
		if _, ok := declared[k]; !ok {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return
}
