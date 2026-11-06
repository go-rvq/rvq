package integration_test

import (
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	. "github.com/go-rvq/rvq/admin/presets/integration"
)

// has: the HTML has the pattern re, whatever its escapes (HTML in JSON in
// HTML): a quote may be ', &#39;, \u0026#39;; a < may be \u003c.
func has(body, re string) bool {
	q := `(?:'|&#39;|\\u0026#39;|"|\\")`
	lt := `(?:<|\\u003c|&lt;)`
	gt := `(?:>|\\u003e|&gt;)`
	re = strings.NewReplacer("'", q, "<", lt, ">", gt).Replace(re)
	return regexp.MustCompile(re).MatchString(body)
}

// A list editor in a table is drawn as a table (a wide screen's), with the
// watch of the browser's width; the form's layout key "mobile" (a phone's,
// asked by the browser: ListEditorLayoutEvent) draws its items as expansion
// panels of their forms instead — the table not drawn.
func TestListEditorTableMobilePanels(t *testing.T) {
	h, err := NewSeededHandler()
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/admin/"+ProductURI+"/1/edit", nil))
	body := w.Body.String()
	if !has(body, `<v-table`) || has(body, `data-list-editor-panels`) {
		t.Errorf("a wide screen's: the table only (%d)", w.Code)
	}
	if !has(body, `<vx-layout-watch :mobile='false'`) || !strings.Contains(body, actions.ListEditorLayoutEvent) {
		t.Errorf("no watch of the width asking the layout")
	}
	// the page of the record (its id in the path): its events carry the id —
	// drawn again, the record, not a new one
	if !regexp.MustCompile(`listEditor_addRowEvent[^)]*\)\.queries\(\{[^}]*\\*"id\\*":\[\\*"1\\*"\]`).MatchString(body) {
		t.Errorf("the add of a row without the record's id")
	}
	if !regexp.MustCompile(`listEditor_layoutEvent[^)]*\)\.queries\(\{\}\)\.query\(\\*"id\\*", \\*"1\\*"\)`).MatchString(body) {
		t.Errorf("the layout's event without the record's id")
	}

	// the browser asks a phone's layout: the form posted with the key
	form := url.Values{presets.ListEditorLayoutKey: {presets.ListEditorLayoutMobile}}
	r := httptest.NewRequest("POST", "/admin/"+ProductURI+"?__execute_event__="+actions.ListEditorLayoutEvent+"&id=1", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	body = w.Body.String()
	for _, want := range []string{`data-list-editor-panels`, `<v-expansion-panel-title>`, `<v-expansion-panel-text :eager='true'`, `text-truncate'>Item 1<`, `text-truncate'>Item 2<`, `<vx-layout-watch :mobile='true'`} {
		if !has(body, want) {
			t.Errorf("a phone's: no %s (%d)", want, w.Code)
		}
	}
	if has(body, `<v-table`) {
		t.Errorf("a phone's: the table drawn too")
	}
}
