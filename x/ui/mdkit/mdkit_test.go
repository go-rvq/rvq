package mdkit

import (
	"strings"
	"testing"
)

// TestJSComponentsPack: the bundle is in the binary, its window.RvqMdKit.
func TestJSComponentsPack(t *testing.T) {
	if js := string(JSComponentsPack()); !strings.Contains(js, "RvqMdKit") || len(js) < 100000 {
		t.Errorf("the pack of mdkit: %d bytes, RvqMdKit %v", len(js), strings.Contains(js, "RvqMdKit"))
	}
}
