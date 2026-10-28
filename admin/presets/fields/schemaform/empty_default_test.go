package schemaform

import (
	"strings"
	"testing"
)

// A record chosen in a choice, or added to a list, begins with the defaults
// its class declares — an int's, an enum member's —, as a record loaded is
// filled; a field with none, empty.
func TestEmptyItemDefaults(t *testing.T) {
	s, err := Parse(`export enum state { open; closed }

export class default {
	n int = 20
	s state = state.closed
	m int
	label str
}`)
	if err != nil {
		t.Fatal(err)
	}
	got := emptyItemJS(s)
	for _, want := range []string{`"n": 20`, `"s": "closed"`, `"m": 0`, `"label": ""`} {
		if !strings.Contains(got, want) {
			t.Errorf("no %s in %s", want, got)
		}
	}
}
