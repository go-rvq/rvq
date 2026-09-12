package admin

import (
	"encoding/json"
	"html"
	"reflect"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// FieldDiffFunc renders the diff of one field from its OLD and NEW values,
// returning the OLD side, the NEW side and an optional merged view. handled=false
// falls through to the built-in handler for the field's type. Registered per
// field with ModelHistory.FieldDiff — the escape hatch for a field whose real
// value must drive the comparison, not its rendered form.
type FieldDiffFunc func(old, new any, ctx *web.EventContext) (oldC, newC, mergedC h.HTMLComponent, handled bool)

// FieldDiffHandler registers a custom diff renderer for one field (by name),
// overriding the type-based default. See FieldDiffFunc.
func (mh *ModelHistory) FieldDiffHandler(field string, fn FieldDiffFunc) *ModelHistory {
	mh.fieldDiffers[field] = fn
	return mh
}

// diffKind is how a field is compared: the model diff is the summary of its
// fields' diffs, and each field is diffed by the handler its kind selects.
type diffKind int

const (
	kindText diffKind = iota // plain text / numbers / time — escaped text diff
	kindHTML                 // declared HTML (e.g. Body) — HTML diff of the value
	kindBool                 // boolean — before/after icons
	kindJSON                 // struct / slice / map — JSON diff of the real value
)

// fieldDiffKind classifies a (possibly nested) field. HTMLFields win; then the
// Go type decides. time.Time is shown as text, not expanded as a struct.
func (mh *ModelHistory) fieldDiffKind(field string) diffKind {
	if mh.IsHTML(field) {
		return kindHTML
	}
	ft, ok := structFieldType(mh.mb.NewModel(), field)
	if !ok {
		return kindText
	}
	for ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
	}
	switch {
	case ft.Kind() == reflect.Bool:
		return kindBool
	case ft.String() == "time.Time":
		return kindText
	case ft.Kind() == reflect.Struct, ft.Kind() == reflect.Slice,
		ft.Kind() == reflect.Map, ft.Kind() == reflect.Array:
		return kindJSON
	default:
		return kindText
	}
}

// diffField computes a field's OLD/NEW/merged diff components, dispatching to the
// handler its kind selects (or a per-field override). A field with no declared
// kind falls back to the generic "diff the rendered value" path.
func (mh *ModelHistory) diffField(f string, oldObj, newObj any, am, bm map[string]json.RawMessage, recordKey string, ctx *web.EventContext, msgr *Messages) (oldC, newC, mergedC h.HTMLComponent) {
	if fn := mh.fieldDiffers[f]; fn != nil {
		if o, n, m, ok := fn(navigateAny(am, f), navigateAny(bm, f), ctx); ok {
			return o, n, m
		}
	}
	switch mh.fieldDiffKind(f) {
	case kindHTML:
		// The value IS HTML: diff the value itself (rendering can hide the real
		// value). Top-level fields render through the DetailingBuilder (so the
		// TipTap field emits diff-friendly HTML); nested ones use the raw value.
		if !strings.ContainsAny(f, ".[") {
			oldHTML, _ := mh.detailHTML(oldObj, recordKey, []string{f}, ctx)
			newHTML, _ := mh.detailHTML(newObj, recordKey, []string{f}, ctx)
			return htmlSideDiff(stripScaffold(oldHTML), stripScaffold(newHTML))
		}
		return htmlSideDiff(fieldValue(am, f), fieldValue(bm, f))
	case kindBool:
		return boolSide(fieldValue(am, f)), boolSide(fieldValue(bm, f)), nil
	case kindJSON:
		// Compare the real value as JSON (the rendered form may not reflect it —
		// e.g. Page.LayoutConfig renders as JSON and must be compared as JSON).
		oj := "<pre>" + html.EscapeString(prettyJSONPath(am, f)) + "</pre>"
		nj := "<pre>" + html.EscapeString(prettyJSONPath(bm, f)) + "</pre>"
		return htmlSideDiff(oj, nj)
	default: // kindText
		// A top-level non-HTML field: the generic path — diff the rendered value.
		if !strings.ContainsAny(f, ".[") {
			oldHTML, _ := mh.detailHTML(oldObj, recordKey, []string{f}, ctx)
			newHTML, _ := mh.detailHTML(newObj, recordKey, []string{f}, ctx)
			if strings.TrimSpace(stripScaffold(oldHTML)) != "" || strings.TrimSpace(stripScaffold(newHTML)) != "" {
				return htmlSideDiff(stripScaffold(oldHTML), stripScaffold(newHTML))
			}
		}
		// Nested/plain value: escaped text diff.
		ov := "<p>" + html.EscapeString(fieldValue(am, f)) + "</p>"
		nv := "<p>" + html.EscapeString(fieldValue(bm, f)) + "</p>"
		return htmlSideDiff(ov, nv)
	}
}

// htmlSideDiff runs the visual HTML diff of two HTML fragments and projects the
// merged result onto an OLD side (removals highlighted), a NEW side (additions),
// and the merged view. With no usable diff it returns the two sides verbatim.
func htmlSideDiff(oldHTML, newHTML string) (oldC, newC, mergedC h.HTMLComponent) {
	oldC, newC = h.RawHTML(oldHTML), h.RawHTML(newHTML)
	if diffs, err := htmlDiffConfig.HTMLdiff([]string{oldHTML, newHTML}); err == nil && len(diffs) > 0 {
		oldC = h.RawHTML(oldSide(diffs[0]))
		newC = h.RawHTML(newSide(diffs[0]))
		mergedC = h.RawHTML(diffs[0])
	}
	return
}

// boolSide renders a boolean snapshot value as a check/cross icon.
func boolSide(raw string) h.HTMLComponent {
	if raw == "true" {
		return v.VIcon("mdi-check-circle").Color("success")
	}
	return v.VIcon("mdi-close-circle").Color("error")
}

// navigateAny returns the value at a field path in a snapshot as a decoded Go
// value (for a per-field FieldDiffFunc override).
func navigateAny(m map[string]json.RawMessage, field string) any {
	raw, ok := navigateRaw(m, field)
	if !ok {
		return nil
	}
	var v any
	if json.Unmarshal(raw, &v) != nil {
		return nil
	}
	return v
}
