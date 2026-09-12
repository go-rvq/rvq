package admin

import (
	"bytes"
	"encoding/json"
	"html"
	"reflect"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gorm.io/gorm"
)

// FieldDiffInput carries everything a field-diff handler needs: the field, both
// records (rebuilt with each revision applied) and their raw snapshots, plus the
// request context. Handlers are reusable across fields and models.
type FieldDiffInput struct {
	MH        *ModelHistory
	Field     string
	OldObj    any
	NewObj    any
	OldSnap   map[string]json.RawMessage
	NewSnap   map[string]json.RawMessage
	RecordKey string
	Ctx       *web.EventContext
	Msgr      *Messages
}

// FieldDiffFunc renders one field's diff: the OLD side, the NEW side and an
// optional merged view. handled=false lets the dispatcher fall through to the
// next handler. Handlers are reusable — register one for a field with
// ModelHistory.FieldDiffHandler, or let the dispatcher pick a built-in.
type FieldDiffFunc func(in *FieldDiffInput) (oldC, newC, mergedC h.HTMLComponent, handled bool)

// FieldDiffHandler registers a reusable diff handler for one field (by name),
// overriding the built-in chosen by type. e.g. FieldDiffHandler("Author",
// history.ReferenceDiff).
func (mh *ModelHistory) FieldDiffHandler(field string, fn FieldDiffFunc) *ModelHistory {
	mh.fieldDiffers[field] = fn
	return mh
}

// diffField computes a field's OLD/NEW/merged diff components: a per-field
// override if registered, otherwise the built-in handler its value selects. The
// model diff is the summary of its fields' diffs (one panel per changed field).
func (mh *ModelHistory) diffField(f string, oldObj, newObj any, am, bm map[string]json.RawMessage, recordKey string, ctx *web.EventContext, msgr *Messages) (oldC, newC, mergedC h.HTMLComponent) {
	in := &FieldDiffInput{
		MH: mh, Field: f, OldObj: oldObj, NewObj: newObj,
		OldSnap: am, NewSnap: bm, RecordKey: recordKey, Ctx: ctx, Msgr: msgr,
	}
	if fn := mh.fieldDiffers[f]; fn != nil {
		if o, n, m, ok := fn(in); ok {
			return o, n, m
		}
	}
	o, n, m, _ := mh.builtinHandler(f)(in)
	return o, n, m
}

// builtinHandler picks the reusable handler for a field by its nature: declared
// HTML, a model reference (ModelSelector / relation), boolean, structured
// (JSON), or the generic rendered-value diff.
func (mh *ModelHistory) builtinHandler(f string) FieldDiffFunc {
	if mh.IsHTML(f) {
		return HTMLValueDiff
	}
	if mh.isReferenceField(f) {
		return ReferenceDiff
	}
	if ft, ok := structFieldType(mh.mb.NewModel(), f); ok {
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		switch {
		case ft.Kind() == reflect.Bool:
			return BoolDiff
		case ft.Kind() == reflect.Struct, ft.Kind() == reflect.Slice,
			ft.Kind() == reflect.Map, ft.Kind() == reflect.Array:
			return PrismDiff("json")
		}
	}
	// Plain text / numbers / time / unknown: Prism plaintext with line numbers.
	return PrismDiff("")
}

// HTMLValueDiff diffs an HTML field by its own value (rendering can hide the real
// value). Top-level fields render through the DetailingBuilder (so the TipTap
// field emits diff-friendly HTML); nested ones use the raw value.
func HTMLValueDiff(in *FieldDiffInput) (oldC, newC, mergedC h.HTMLComponent, handled bool) {
	f := in.Field
	if !strings.ContainsAny(f, ".[") {
		oldHTML, _ := in.MH.detailHTML(in.OldObj, in.RecordKey, []string{f}, in.Ctx)
		newHTML, _ := in.MH.detailHTML(in.NewObj, in.RecordKey, []string{f}, in.Ctx)
		oldC, newC, mergedC = htmlSideDiff(stripScaffold(oldHTML), stripScaffold(newHTML))
		return oldC, newC, mergedC, true
	}
	oldC, newC, mergedC = htmlSideDiff(fieldValue(in.OldSnap, f), fieldValue(in.NewSnap, f))
	return oldC, newC, mergedC, true
}

// BoolDiff shows a boolean field's before/after as check/cross icons.
func BoolDiff(in *FieldDiffInput) (oldC, newC, mergedC h.HTMLComponent, handled bool) {
	return boolSide(fieldValue(in.OldSnap, in.Field)), boolSide(fieldValue(in.NewSnap, in.Field)), nil, true
}

// JSONDiff compares a structured field (struct / slice / map) as JSON — the real
// value, because the rendered form may not reflect it (e.g. Page.LayoutConfig
// renders as JSON and must be compared as JSON). Rendered with Prism (json).
func JSONDiff(in *FieldDiffInput) (oldC, newC, mergedC h.HTMLComponent, handled bool) {
	return PrismDiff("json")(in)
}

// PrismDiff returns a reusable handler that renders both sides with Prism: syntax
// highlighting for the given language (a Prism grammar name such as "json",
// "gad", "markup"; empty renders escaped plain text), numbered lines, and the
// changed lines marked (removed on OLD, added on NEW), computed line by line. It
// is the handler for JSON and plain-text fields, and can be registered for any
// field whose value is code — FieldDiffHandler("Config", history.PrismDiff("yaml")).
func PrismDiff(language string) FieldDiffFunc {
	return func(in *FieldDiffInput) (oldC, newC, mergedC h.HTMLComponent, handled bool) {
		oldTxt := prismText(in.OldSnap, in.Field, language)
		newTxt := prismText(in.NewSnap, in.Field, language)
		removed, added := lineDiffSides(oldTxt, newTxt)
		oldC = vx.VXCode(oldTxt).Language(language).MarkKind("del").MarkLines(removed)
		newC = vx.VXCode(newTxt).Language(language).MarkKind("add").MarkLines(added)
		return oldC, newC, nil, true
	}
}

// prismText extracts the field value for Prism: pretty JSON for the json
// language, the plain value otherwise.
func prismText(m map[string]json.RawMessage, field, language string) string {
	if language == "json" {
		return prettyJSONPath(m, field)
	}
	return fieldValue(m, field)
}

// lineDiffSides reports which 1-based lines were removed (present in old, not
// new) and added (present in new, not old), from a line-level diff.
func lineDiffSides(oldStr, newStr string) (removed, added []int) {
	// A trailing final line without "\n" is a different token from the same line
	// with one; normalize so an appended last line is an insert, not a change.
	if oldStr != "" && !strings.HasSuffix(oldStr, "\n") {
		oldStr += "\n"
	}
	if newStr != "" && !strings.HasSuffix(newStr, "\n") {
		newStr += "\n"
	}
	dmp := diffmatchpatch.New()
	a, b, lineArray := dmp.DiffLinesToChars(oldStr, newStr)
	diffs := dmp.DiffCharsToLines(dmp.DiffMain(a, b, false), lineArray)
	oldLine, newLine := 0, 0
	for _, d := range diffs {
		n := countLines(d.Text)
		switch d.Type {
		case diffmatchpatch.DiffEqual:
			oldLine += n
			newLine += n
		case diffmatchpatch.DiffDelete:
			for i := 0; i < n; i++ {
				oldLine++
				removed = append(removed, oldLine)
			}
		case diffmatchpatch.DiffInsert:
			for i := 0; i < n; i++ {
				newLine++
				added = append(added, newLine)
			}
		}
	}
	return
}

// countLines counts the lines in a line-diff text chunk (each token is a line,
// the last possibly without a trailing newline).
func countLines(s string) int {
	if s == "" {
		return 0
	}
	n := strings.Count(s, "\n")
	if !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}

// ReferenceDiff compares a field that selects another model (ModelSelector /
// foreign key) by EXACT value — never a partial text diff — and shows the
// referenced record rendered: its title with the id highlighted beside it. When
// the value changed the whole side is highlighted (OLD red, NEW green). A
// to-many relation (a list value, e.g. Post.Tags) is compared as a set: each
// side lists its items, with removed items highlighted red on OLD and added
// items green on NEW.
func ReferenceDiff(in *FieldDiffInput) (oldC, newC, mergedC h.HTMLComponent, handled bool) {
	oldRaw, _ := navigateRaw(in.OldSnap, in.Field)
	newRaw, _ := navigateRaw(in.NewSnap, in.Field)
	if isJSONArray(oldRaw) || isJSONArray(newRaw) {
		return manyReferenceDiff(oldRaw, newRaw)
	}
	changed := fieldValue(in.OldSnap, in.Field) != fieldValue(in.NewSnap, in.Field)
	oldC = referenceSide(in, in.OldObj, in.OldSnap, changed, false)
	newC = referenceSide(in, in.NewObj, in.NewSnap, changed, true)
	return oldC, newC, nil, true
}

// manyReferenceDiff renders a to-many relation (e.g. Post.Tags) as two lists,
// diffed as whole lists: OLD lists its items with the removed ones marked (red,
// struck), NEW lists its items with the added ones marked (green). Membership is
// by id, so the diff says exactly which item was added or removed. Each item is a
// VListItem with the id as a VChip.
func manyReferenceDiff(oldRaw, newRaw json.RawMessage) (oldC, newC, mergedC h.HTMLComponent, handled bool) {
	oldItems, newItems, removed, added := manyRefChanges(oldRaw, newRaw)
	oldC = refList(oldItems, func(it refItem) refState {
		if removed[it.id] {
			return refRemovedState
		}
		return refKept
	})
	newC = refList(newItems, func(it refItem) refState {
		if added[it.id] {
			return refAddedState
		}
		return refKept
	})
	return oldC, newC, nil, true
}

// refItem is one referenced record: its id and a display label.
type refItem struct {
	id, label string
}

// refItems decodes a to-many snapshot value (a JSON array) into display items.
func refItems(raw json.RawMessage) []refItem {
	if !isJSONArray(raw) {
		return nil
	}
	var arr []json.RawMessage
	if json.Unmarshal(raw, &arr) != nil {
		return nil
	}
	out := make([]refItem, 0, len(arr))
	for _, el := range arr {
		id, label := refElemLabel(el)
		out = append(out, refItem{id: id, label: label})
	}
	return out
}

// manyRefChanges compares two to-many snapshot values by id, returning both item
// lists and the sets of ids removed (in OLD, not NEW) and added (in NEW, not
// OLD) — the whole-list diff that says which item changed.
func manyRefChanges(oldRaw, newRaw json.RawMessage) (oldItems, newItems []refItem, removed, added map[string]bool) {
	oldItems, newItems = refItems(oldRaw), refItems(newRaw)
	oldIDs := map[string]bool{}
	for _, it := range oldItems {
		oldIDs[it.id] = true
	}
	newIDs := map[string]bool{}
	for _, it := range newItems {
		newIDs[it.id] = true
	}
	removed, added = map[string]bool{}, map[string]bool{}
	for _, it := range oldItems {
		if !newIDs[it.id] {
			removed[it.id] = true
		}
	}
	for _, it := range newItems {
		if !oldIDs[it.id] {
			added[it.id] = true
		}
	}
	return
}

// refState is how a to-many item changed on its side of the diff.
type refState int

const (
	refKept refState = iota
	refRemovedState
	refAddedState
)

// refList renders a to-many item list as a VList, each item a VListItem with the
// id as a VChip and colored by its state.
func refList(items []refItem, state func(refItem) refState) h.HTMLComponent {
	if len(items) == 0 {
		return h.Span("—").Class("text-medium-emphasis")
	}
	lis := make(h.HTMLComponents, 0, len(items))
	for _, it := range items {
		lis = append(lis, refListItem(it, state(it)))
	}
	return v.VList(lis...).Density(v.DensityCompact).Class("py-0 bg-transparent")
}

// refListItem is one item of a to-many list: its label, the id as a VChip, and a
// theme color per state (removed → error struck, added → success, kept → none).
func refListItem(it refItem, state refState) h.HTMLComponent {
	row := h.HTMLComponents{h.Span(it.label)}
	if it.id != "" {
		row = append(row, v.VChip(h.Text("#"+it.id)).Size(v.SizeXSmall).Variant(v.VariantTonal).Class("ms-2"))
	}
	cls := "px-0"
	switch state {
	case refRemovedState:
		cls += " text-error text-decoration-line-through"
	case refAddedState:
		cls += " text-success"
	}
	return v.VListItem(h.Div(row...).Class("d-flex align-center")).Density(v.DensityCompact).Class(cls)
}

// refElemLabel pulls a display label and id from one referenced record's JSON:
// a common title field (Name/Title/Label/…) for the label, an id field for the
// id, falling back to the id when there is no title.
func refElemLabel(raw json.RawMessage) (id, label string) {
	var obj map[string]json.RawMessage
	if json.Unmarshal(raw, &obj) != nil {
		return "", strings.Trim(string(raw), `"`)
	}
	for _, k := range []string{"ID", "Id", "id"} {
		if r, ok := obj[k]; ok {
			id = strings.Trim(string(r), `"`)
			break
		}
	}
	for _, k := range []string{"Name", "Title", "Label", "Nome", "Titulo", "Slug", "Code"} {
		if r, ok := obj[k]; ok {
			var s string
			if json.Unmarshal(r, &s) == nil && s != "" {
				label = s
				break
			}
		}
	}
	if label == "" {
		label = "#" + id
	}
	return
}

// isJSONArray reports whether raw is a JSON array value.
func isJSONArray(raw json.RawMessage) bool {
	t := bytes.TrimSpace(raw)
	return len(t) > 0 && t[0] == '['
}

// RenderedDiff is the generic default: diff the value as rendered by the
// DetailingBuilder (top-level), falling back to an escaped-text diff.
func RenderedDiff(in *FieldDiffInput) (oldC, newC, mergedC h.HTMLComponent, handled bool) {
	f := in.Field
	if !strings.ContainsAny(f, ".[") {
		oldHTML, _ := in.MH.detailHTML(in.OldObj, in.RecordKey, []string{f}, in.Ctx)
		newHTML, _ := in.MH.detailHTML(in.NewObj, in.RecordKey, []string{f}, in.Ctx)
		oldHTML, newHTML = stripScaffold(oldHTML), stripScaffold(newHTML)
		if strings.TrimSpace(oldHTML) != "" || strings.TrimSpace(newHTML) != "" {
			oldC, newC, mergedC = htmlSideDiff(oldHTML, newHTML)
			return oldC, newC, mergedC, true
		}
	}
	ov := "<p>" + html.EscapeString(fieldValue(in.OldSnap, f)) + "</p>"
	nv := "<p>" + html.EscapeString(fieldValue(in.NewSnap, f)) + "</p>"
	oldC, newC, mergedC = htmlSideDiff(ov, nv)
	return oldC, newC, mergedC, true
}

// referenceSide renders one side of a reference field: the record as the model
// renders it (its title chip), the id highlighted beside it, and — when the
// value changed — the whole side highlighted (red for OLD, green for NEW).
func referenceSide(in *FieldDiffInput, obj any, snap map[string]json.RawMessage, changed, isNew bool) h.HTMLComponent {
	rendered, _ := in.MH.detailHTML(obj, in.RecordKey, []string{in.Field}, in.Ctx)
	body := h.HTMLComponents{h.RawHTML(stripScaffold(rendered))}
	if id := referenceID(snap, in.Field); id != "" {
		body = append(body, v.VChip(h.Text("#"+id)).Size(v.SizeXSmall).Variant(v.VariantTonal).Class("ms-2"))
	}
	cls := "d-flex align-center"
	if changed {
		if isNew {
			cls += " text-success"
		} else {
			cls += " text-error"
		}
	}
	return h.Div(body...).Class(cls)
}

// referenceID extracts the referenced record's id from a snapshot value — the
// "ID" field of the embedded object, or the raw scalar when the value is the id.
func referenceID(m map[string]json.RawMessage, field string) string {
	raw, ok := navigateRaw(m, field)
	if !ok {
		return ""
	}
	var obj map[string]json.RawMessage
	if json.Unmarshal(raw, &obj) == nil {
		for _, k := range []string{"ID", "Id", "id"} {
			if idRaw, ok := obj[k]; ok {
				return strings.Trim(string(idRaw), `"`)
			}
		}
		return ""
	}
	return strings.Trim(string(raw), `"`)
}

// isReferenceField reports whether a field selects another model — a manual
// ReferenceFields override, or a relation detected from the gorm schema.
func (mh *ModelHistory) isReferenceField(field string) bool {
	top := topField(field)
	if mh.refFields[top] {
		return true
	}
	return mh.relationFields()[top]
}

// relationFields is the set of the model's relation field names (belongs-to,
// has-one, …) from the gorm schema, computed once.
func (mh *ModelHistory) relationFields() map[string]bool {
	if mh.schemaRefs != nil {
		return mh.schemaRefs
	}
	set := map[string]bool{}
	stmt := &gorm.Statement{DB: mh.db}
	if stmt.Parse(mh.mb.NewModel()) == nil && stmt.Schema.Relationships.Relations != nil {
		for name := range stmt.Schema.Relationships.Relations {
			set[name] = true
		}
	}
	mh.schemaRefs = set
	return set
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
