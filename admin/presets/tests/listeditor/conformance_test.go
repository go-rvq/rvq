package listeditor

import (
	"net/http/httptest"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	. "github.com/go-rvq/rvq/admin/presets/integration"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/web/multipartestutils"
)

// formKeyRe finds every LITERAL key the rendered form binds — `v-model='form["X"]'`
// and the `v-assign='[form, {"X": …}]'` that seeds it. The response travels as
// JSON, so the quotes arrive escaped; a key the page computes in JS
// (`form[k + ".__pos"]`) is not a literal and does not match.
var formKeyRe = regexp.MustCompile(`form\[\\?"([^"\\]+)\\?"\]|form, \{\\?"([^"\\]+)\\?":`)

// controlKey reports whether a key is a control field of the form machinery
// rather than a value of the record: the signed stamp itself and the list
// editor's per-row flags.
func controlKey(key string) bool {
	if strings.HasPrefix(key, "__") {
		return true
	}
	for _, suffix := range []string{".__index", ".__pos", ".__deleted", ".__new", ".__present", ".__purged"} {
		if strings.HasSuffix(key, suffix) {
			return true
		}
	}
	return false
}

// The record stamp is only as true as the convention the components follow: it
// hashes the fields of the editing builder, and the form is what the components
// rendered. This walks the RENDERED form and requires every key it binds to be
// a key the stamp covers — a component that binds something else would leave
// the guard comparing what the form does not carry.
//
// See docs/fields.md ("The convention a component follows").
func TestStampCoversEveryKeyTheFormBinds(t *testing.T) {
	app, err := twoPersistedItems()
	if err != nil {
		t.Fatal(err)
	}

	// the edit form of product 1, as a browser asks for it
	req := NewMultipartBuilder().
		PageURL("/admin/products").
		EventFunc(actions.Edit).
		Query(presets.ParamID, "1").
		BuildEventFuncRequest()

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	body := w.Body.String()

	var bound []string
	for _, m := range formKeyRe.FindAllStringSubmatch(body, -1) {
		key := m[1]
		if key == "" {
			key = m[2]
		}
		if key != "" && !controlKey(key) {
			bound = append(bound, key)
		}
	}
	if len(bound) == 0 {
		t.Fatalf("o form renderizado não ligou chave nenhuma:\n%s", body)
	}

	mb := app.GetModel(&Product{})
	obj := mb.NewModel()
	mid, err := mb.ParseRecordID("1")
	if err != nil {
		t.Fatal(err)
	}
	if err = mb.Editing().Fetcher(obj, mid, &web.EventContext{R: req}); err != nil {
		t.Fatal(err)
	}

	covered := map[string]bool{}
	for _, k := range mb.Editing().RecordStampKeys(obj) {
		covered[strings.TrimSuffix(k, "#")] = true
	}

	var missing []string
	for _, key := range bound {
		if !covered[key] {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		var have []string
		for k := range covered {
			have = append(have, k)
		}
		sort.Strings(have)
		t.Errorf("o form liga chaves que o stamp não cobre: %v\nstamp cobre: %v",
			unique(missing), have)
	}
}

func unique(in []string) (out []string) {
	seen := map[string]bool{}
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return
}
