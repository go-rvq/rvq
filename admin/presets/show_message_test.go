package presets

import (
	"strings"
	"testing"

	"github.com/go-rvq/rvq/web"
)

// The snackbar in presets_layout.go binds vars.presetsMessage.message and
// .htmlMessage. A message written under any other key leaves it open and empty,
// which is what "htmlText" — the field's JSON name — used to do.
func TestShowMessageKeys(t *testing.T) {
	for _, c := range []struct {
		name    string
		msg     any
		wantKey string
		wantTxt string
	}{
		{"plain text", "moved", "message", "moved"},
		{"text field", &web.FlashMessage{Text: "moved"}, "message", "moved"},
		{"html body", &web.FlashMessage{HtmlText: "moved"}, "htmlMessage", "moved"},
	} {
		t.Run(c.name, func(t *testing.T) {
			r := &web.EventResponse{}
			ShowMessage(r, c.msg)

			if !strings.Contains(r.RunScript, c.wantKey+": ") {
				t.Errorf("script does not set %q:\n%s", c.wantKey, r.RunScript)
			}
			if !strings.Contains(r.RunScript, c.wantTxt) {
				t.Errorf("script does not carry %q:\n%s", c.wantTxt, r.RunScript)
			}
		})
	}
}

// The body is embedded in a script, so its angle brackets are escaped rather
// than closing the tag around it.
func TestShowMessageEscapesHTML(t *testing.T) {
	r := &web.EventResponse{}
	ShowMessage(r, &web.FlashMessage{HtmlText: "<b>moved</b>"})

	if strings.Contains(r.RunScript, "<b>") {
		t.Errorf("angle brackets left unescaped:\n%s", r.RunScript)
	}
	if !strings.Contains(r.RunScript, `\u003cb\u003e`) {
		t.Errorf("escaped body not found:\n%s", r.RunScript)
	}
}

// A message with nothing to say writes no script at all, rather than opening an
// empty snackbar.
func TestShowMessageEmpty(t *testing.T) {
	r := &web.EventResponse{}
	ShowMessage(r, &web.FlashMessage{})

	if r.RunScript != "" {
		t.Errorf("expected no script, got:\n%s", r.RunScript)
	}
}
