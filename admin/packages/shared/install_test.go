package shared

import (
	"bytes"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/google/uuid"
)

func TestShareBodyTemplatesVsCheckboxes(t *testing.T) {
	m := Messages_en_US

	// with templates -> a template selector (no verb checkboxes)
	templates := []ShareTemplate{{ID: uuid.New(), Name: "Viewer", Actions: VerbView}}
	var withT bytes.Buffer
	if err := h.Fprint(&withT, shareBody(m, nil, templates), nil); err != nil {
		t.Fatal(err)
	}
	ht := withT.String()
	if !strings.Contains(ht, "Viewer") || !strings.Contains(ht, fieldTemplate) {
		t.Errorf("with templates the dialog should offer a template selector; got:\n%s", ht)
	}
	if strings.Contains(ht, fieldView) {
		t.Errorf("with templates the verb checkboxes should be hidden")
	}

	// without templates -> raw verb checkboxes
	var noT bytes.Buffer
	if err := h.Fprint(&noT, shareBody(m, nil, nil), nil); err != nil {
		t.Fatal(err)
	}
	if hn := noT.String(); !strings.Contains(hn, fieldView) {
		t.Errorf("without templates the dialog should offer verb checkboxes; got:\n%s", hn)
	}
}
