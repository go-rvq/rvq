package integration

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets/actions"
	. "github.com/go-rvq/rvq/web/multipartestutils"
)

// Saved, a schema form is checked by the server — the limits of a field, a
// `[validation=…]` function — and the form comes back with what is wrong
// under each field (its key's error-messages), nothing saved; right, it is.
func TestSchemaformSaveChecks(t *testing.T) {
	handler, err := NewSchemaformSeededHandler()
	if err != nil {
		t.Fatal(err)
	}
	post := func(rooms, word string) string {
		t.Helper()
		req := NewMultipartBuilder().
			PageURL("/admin/"+SchemaformQuoteURI).
			EventFunc(actions.Update).
			Query("id", "1").
			AddField("Options.contact.name", "Ana").
			AddField("Options.contact.phone", "(781) 555-0100").
			AddField("Options.project.rooms", rooms).
			AddField("Options.project.word", word).
			BuildEventFuncRequest()
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		b, _ := io.ReadAll(w.Result().Body)
		// (the form is HTML in JSON: its quotes escaped, once or twice)
		return strings.ReplaceAll(string(b), `\`, "")
	}
	body := post("11", "bad")
	for _, want := range []string{`"Options.project.rooms":["the largest value is 10"]`, `"Options.project.word":["not that word"]`} {
		if !strings.Contains(body, want) {
			t.Errorf("no %q in the answer:\n%.3000s", want, body)
		}
	}
	if body := post("3", "good"); strings.Contains(body, `"Options.project.`) {
		t.Errorf("a form right refused:\n%.2000s", body)
	}
}
