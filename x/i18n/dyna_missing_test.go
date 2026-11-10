package i18n

import (
	"strings"
	"testing"
)

// A key allowed to be empty (a field's hint, its description) missing is not
// reported; one that is not, is — and only it.
func TestMissingKeysAllowEmpty(t *testing.T) {
	d := DynaNew().Language("pt-BR")
	d.putMissingKey("m", "Field_Hint", "", true)
	if d.HaveMissingKeys() || d.PrettyMissingKeys() != "" {
		t.Errorf("a hint reported: %q", d.PrettyMissingKeys())
	}
	d.putMissingKey("m", "Title", "Title")
	out := d.PrettyMissingKeys()
	if !d.HaveMissingKeys() || !strings.Contains(out, "Title string") || strings.Contains(out, "Field_Hint") {
		t.Errorf("reported: %q", out)
	}
}
