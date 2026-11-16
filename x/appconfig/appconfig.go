// Package appconfig is the app.yml as read (admin/app/env): the sections of
// the packages that have their own (`fstools:`), each read by its package
// when it needs it (Section).
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
