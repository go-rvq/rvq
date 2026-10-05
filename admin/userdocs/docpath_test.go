package userdocs

import "testing"

// The address of a document is its path under the documentation
// (presets.TreePage).
func TestDocPath(t *testing.T) {
	for node, want := range map[string]string{
		"guides/policies/06-roles": "/admin/docs/guides/policies/06-roles",
		"posts/actions/Localize":   "/admin/docs/posts/actions/Localize",
		"a b/c?d":                  "/admin/docs/a%20b/c%3Fd",
	} {
		if got := DocPath("/admin/docs", node); got != want {
			t.Errorf("%q: %q, want %q", node, got, want)
		}
	}
}
