package userdocs

import "testing"

func TestServicesText(t *testing.T) {
	if got := mdText("a <name>.git, `<b>` <c>"); got != "a &lt;name>.git, `<b>` &lt;c>" {
		t.Errorf("%q", got)
	}
}
