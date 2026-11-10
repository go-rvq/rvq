package userdocs

import "testing"

// A document's `{{…}}` reaches Vue as entities: shown as written, never an
// expression.
func TestVueSafe(t *testing.T) {
	in := `<p>keep <code>{{…}}</code> and {{ x }}</p>`
	want := `<p>keep <code>&#123;&#123;…&#125;&#125;</code> and &#123;&#123; x &#125;&#125;</p>`
	if got := VueSafe(in); got != want {
		t.Errorf("%s", got)
	}
}
