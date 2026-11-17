// Package appconfig is the app.yml as read (admin/app/env): the sections of
// the packages that have their own (`fstools:`), each read by its package
// when it needs it (Section).
//
// A package names its section by itself (its module, short: `fstools`),
// keeps its defaults in the struct it reads into — the file says only what
// changes — and reads it once the configuration is loaded (in Install, not
// in an init: env loads the app.yml in its own init, which may run after):
//
//	var cfg struct {
//		Terminal TerminalConfig `yaml:"terminal"`
//	}
//	cfg.Terminal = DefaultTerminalConfig()
//	if _, err := appconfig.Section("fstools", &cfg); err != nil {
//		// a wrong value: say it, keep what is safe
//	}
//
// A section is the package's to document: its keys, its defaults, what a
// wrong value does. The keys of the app.yml env does not know (`env:`,
// `features:`, `security:` aside) are the packages'; a section nobody reads
// is silently ignored.
package appconfig

import (
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	mu   sync.RWMutex
	docs [][]byte
)

// Load adds content, an app.yml read: a later one's sections over an
// earlier one's.
func Load(content []byte) {
	mu.Lock()
	defer mu.Unlock()
	docs = append(docs, append([]byte(nil), content...))
}

// Reset forgets what was loaded (tests).
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	docs = nil
}

// Section reads the section name of the app.yml into out (a pointer to a
// struct of yaml tags): the defaults out has stay where the file says
// nothing; false, no app.yml has it.
func Section(name string, out any) (bool, error) {
	mu.RLock()
	defer mu.RUnlock()
	found := false
	for _, d := range docs {
		var doc map[string]yaml.Node
		if err := yaml.Unmarshal(d, &doc); err != nil {
			return found, err
		}
		n, ok := doc[name]
		if !ok {
			continue
		}
		if err := n.Decode(out); err != nil {
			return found, err
		}
		found = true
	}
	return found, nil
}
