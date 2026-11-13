package userdocs

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/sunfmin/reflectutils"
)

// The field of the permissions of a role, of an access key (vuetifyx's
// vx-permission-input): its value one line (perm.JoinLines of the lines
// stored), shown as the tree of the permissions of the admin kept to what
// the value touches (PermValueTree) — the views, the groups of the menu, the
// search and the kinds as the tab Árvore —, each item edited: its prefix,
// what may follow it (the + of a dialog: PermChoicesOf, PermRecordsOf), out.
// Each edit is an operation of the event of the field
// (PermissionInputEventFunc), answering the value, canonical, and its tree.
// Stored, the value is the items of its root, one a scope (perm.SplitRoot),
// what is given to a resource by its unique name, no placeholder
// (perm.StoreLines).

// PermissionInputOptions are where the tree of a PermissionInput comes from
// and what it may write.
type PermissionInputOptions struct {
	// URL, Event: the event of the field (PermissionInputEventFunc), at the
	// address of the model that registered it; Query, what it is asked with
	// besides
	URL, Event string
	Query      map[string]string
	// Prefixes are the prefixes a line may have ("grant", "deny", "allow",
	// "reset"); none, all — of an access key: deny only
	Prefixes []string
	// Value: how the tree of the value is made — what it is kept to (Base),
	// what its lines allow (Effect)
	Value PermValueOptions
	// Readonly: shown, not changed
	Readonly bool
	// StringView: the permission as text, indented, in a card beside the
	// tree (under it when the field is narrow)
	StringView bool
	// NoTextInput: no field of text to add lines from (on: one)
	NoTextInput bool
	// HelpURL, HelpTitle: a ? at the end of the bar of the field, opening the
	// help of the permissions (a new tab); "": none
	HelpURL, HelpTitle string
}

// PermissionInputEventFunc is the event of a PermissionInput (r.Data), of
// the options opts gives of the request:
//
//   - op "tree": {value, items, formatted}: the tree of the form's value
//     (PermValueTree) of the form's view, groups ("1") and only ("1": only
//     what is defined), and the value to
//     read (perm.FormatLine, indented by PermIndent);
//   - op "set", "remove", "add": the value edited (EditValue) of the form's
//     resources (a JSON list) and prefix — of a group taken out, what is in
//     it too —, and its tree;
//   - op "merge": the lines of the form's text — any form: indented, by
//     commas — merged in the value (EditValue), and its tree; what is wrong
//     with the text, {error};
//   - op "children": the parts that may follow the form's path
//     (PermChoicesOf), of the groups or not;
//   - op "records": the records of the model of the path whose title or id
//     has q (PermRecordsOf).
func PermissionInputEventFunc(p *presets.Builder, opts func(ctx *web.EventContext) (PermissionInputOptions, error)) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		o, err := opts(ctx)
		if err != nil {
			return r, err
		}
		groups := ctx.R.FormValue("groups") == "1"
		switch op := ctx.R.FormValue("op"); op {
		case "children":
			ch, err := PermChoicesOf(ctx.Context(), p, ctx.R.FormValue("path"), groups, o.Value.Base)
			if err != nil {
				return r, err
			}
			if ch.Choices == nil {
				ch.Choices = []PermChoice{}
			}
			r.Data = ch
		case PermEditMerge:
			// a text of lines, any form (indented, by commas): merged, or what
			// is wrong with it (data.error)
			value, idle, err := mergeText(ctx.R.FormValue("value"), ctx.R.FormValue("text"), o.Prefixes)
			if err != nil {
				r.Data = map[string]any{"error": err.Error()}
				return r, nil
			}
			if r, err = permValueAnswer(ctx, p, o, value); err != nil || len(idle) == 0 {
				return r, err
			}
			// the unsets that took nothing out: said, under the text
			r.Data.(map[string]any)["warning"] = fmt.Sprintf(GetMessages(ctx.Context()).PermInputUnsetIdle,
				strings.Join(idle, ", "))
			return r, nil
		case "records":
			recs, err := PermRecordsOf(ctx, p, ctx.R.FormValue("path"), ctx.R.FormValue("q"))
			if err != nil {
				return r, err
			}
			if recs == nil {
				recs = []presets.PermRecord{}
			}
			r.Data = recs
		default:
			value := ctx.R.FormValue("value")
			if op != "tree" {
				var resources []string
				if s := ctx.R.FormValue("resources"); s != "" {
					if err := json.Unmarshal([]byte(s), &resources); err != nil {
						return r, err
					}
				}
				if op == PermEditRemove {
					for _, res := range append([]string{}, resources...) {
						if _, _, ok := groupOfLine(perm.StripLabels(res)); ok {
							resources = append(resources, groupMembers(ctx, p, perm.StripLabels(res))...)
						}
					}
				}
				prefix := ctx.R.FormValue("prefix")
				if !prefixAllowed(o.Prefixes, prefix) {
					return r, fmt.Errorf("userdocs: the prefix %q is not of this field", prefix)
				}
				if value, err = EditValue(value, op, resources, prefix); err != nil {
					return r, err
				}
			}
			return permValueAnswer(ctx, p, o, value)
		}
		return
	}
}

// permValueAnswer is the answer of the value: {value, items, formatted} —
// its tree of the form's view, groups and only, and the value to read.
func permValueAnswer(ctx *web.EventContext, p *presets.Builder, o PermissionInputOptions, value string) (r web.EventResponse, err error) {
	vo := o.Value
	vo.Groups = ctx.R.FormValue("groups") == "1"
	vo.Only = ctx.R.FormValue("only") == "1"
	vo.View = PermTreeByVerb
	if PermTreeView(ctx.R.FormValue("view")) == PermTreeByKind {
		vo.View = PermTreeByKind
	}
	items, err := PermValueTree(ctx, p, value, vo)
	if err != nil {
		return r, err
	}
	if items == nil {
		items = []*PermTreeItem{}
	}
	// the value to read: indented, each alternative a line
	formatted, _ := perm.FormatLine(value, perm.WithIndent(PermIndent))
	r.Data = map[string]any{"value": value, "items": items, "formatted": formatted}
	return
}

// mergeText is the value with the lines of the text merged
// (perm.ExprAlt.MergeReport: an "unset" takes out the lines before it it
// covers), as stored (perm.UngroupPattern, the placeholders others are under
// out), and the unsets that took nothing out (idle); what is wrong with the
// text — it cannot be read, a prefix the field does not take (prefixes) —,
// its error.
func mergeText(value, text string, prefixes []string) (merged string, idle []string, err error) {
	v, err := perm.ParseExpr(value)
	if err != nil {
		return "", nil, err
	}
	t, err := perm.ParseExpr(text)
	if err != nil {
		return "", nil, err
	}
	if err := linesAllowed(prefixes, t.Normalize().Lines()); err != nil {
		return "", nil, err
	}
	m, idle := v.MergeReport(t)
	merged, err = EditValue(m.String(), PermEditMerge, nil, "")
	return merged, idle, err
}

// linesAllowed says what is wrong with the lines: a prefix the field does
// not take (prefixes; none: any) — of an access key, a line that does not
// deny.
func linesAllowed(prefixes []string, lines []string) error {
	for _, l := range lines {
		k, _ := perm.ParseLine(l)
		if k == perm.LineUnset {
			// an edit, not a line: any field takes it
			continue
		}
		if name := prefixNames[k]; !prefixAllowed(prefixes, name) {
			return fmt.Errorf("%q: the prefix %q is not of this field", l, name)
		}
	}
	return nil
}

// PermIndent indents the value of a PermissionInput to read: each
// alternative of a group a line.
const PermIndent = "  "

// prefixAllowed says the prefix may be written: none ("", "none": no line),
// or one of prefixes (none: any).
func prefixAllowed(prefixes []string, prefix string) bool {
	if prefix == "" || prefix == "none" || len(prefixes) == 0 {
		return true
	}
	for _, p := range prefixes {
		if p == prefix {
			return true
		}
	}
	return false
}

// PermKnown is what a part of a line is, by its path (labels none of it):
// for the field to show it before asking.
type PermKnown struct {
	Title       string `json:"title,omitempty"`
	Kind        string `json:"kind,omitempty"`
	Description string `json:"description,omitempty"`
	// Error: a record not found (or the user may not read)
	Error bool `json:"error,omitempty"`
}

// DescribeLine is the line canonical, the labels of its records as they are
// now, and what each of its parts is, by its path (PermKnown). A line that
// cannot be read is as it is, nothing known.
func DescribeLine(ctx *web.EventContext, p *presets.Builder, line string) (string, map[string]PermKnown) {
	known := map[string]PermKnown{}
	e, err := perm.ParseExpr(line)
	if err != nil {
		return line, known
	}
	e = e.Normalize()
	ix := newPermIndex(p)
	t := &permTree{ctx: ctx.Context(), msgs: GetMessages(ctx.Context())}
	choices := map[string]map[string]PermChoice{}
	choicesOf := func(path string) map[string]PermChoice {
		if m, ok := choices[path]; ok {
			return m
		}
		m := map[string]PermChoice{}
		// by the menu and by the unique names: a group, a resource of a group
		for _, groups := range []bool{true, false} {
			if r, err := PermChoicesOf(ctx.Context(), p, path, groups, nil); err == nil {
				for _, c := range r.Choices {
					m[c.Segment] = c
				}
			}
		}
		choices[path] = m
		return m
	}
	var walk func(before string, a *perm.ExprAlt)
	walk = func(before string, a *perm.ExprAlt) {
		segs, group := exprSegments(a.Nodes)
		for _, seg := range segs {
			key := perm.NodesText(seg)
			path := before + key
			c, ok := choicesOf(before)[key]
			if r, isRecord := seg[0].(*perm.ExprRecord); isRecord && r.ID != "*" {
				k := PermKnown{Kind: PermKindRecord, Title: r.Label}
				if label, found := recordLabel(ctx, ix, before, r.ID); found {
					r.Label, k.Title = label, label
				} else {
					k.Error, k.Title = true, r.ID
				}
				known[path] = k
			} else if ok {
				known[path] = PermKnown{Title: c.Title, Kind: c.Kind, Description: c.Description}
			} else if key == "*" {
				known[path] = PermKnown{Title: t.msgs.PermInputEverything, Kind: PermKindAll}
			}
			before = path
		}
		if group != nil {
			for _, c := range group.Alts {
				walk(before, c)
			}
		}
	}
	walk("", e)
	return e.String(), known
}

// recordLabel is the title of the record id of the model of the beginning
// path (its listing).
func recordLabel(ctx *web.EventContext, ix *permIndex, path, id string) (string, bool) {
	pp, err := ix.readPath(path)
	if err != nil {
		return "", false
	}
	n := ix.byName[pp.key]
	if n == nil || n.Model == nil {
		return "", false
	}
	return n.Model.PermRecordLabel(ctx, pp.parents(), id)
}

// exprSegments are the segments of nodes, each up to its ":" (with it) — a
// record with the text that follows it —, and the group that ends them.
func exprSegments(nodes []perm.ExprNode) (segs [][]perm.ExprNode, group *perm.ExprGroup) {
	var cur []perm.ExprNode
	for _, n := range nodes {
		switch n := n.(type) {
		case *perm.ExprGroup:
			if len(cur) > 0 {
				segs = append(segs, cur)
			}
			return segs, n
		case *perm.ExprRecord:
			cur = append(cur, n)
		case *perm.ExprText:
			t := n.Text
			for {
				j := colonAt(t)
				if j < 0 {
					if t != "" {
						cur = append(cur, &perm.ExprText{Text: t})
					}
					break
				}
				segs = append(segs, append(cur, &perm.ExprText{Text: t[:j+1]}))
				cur = nil
				t = t[j+1:]
			}
		}
	}
	if len(cur) > 0 {
		segs = append(segs, cur)
	}
	return
}

// PermInputKindLabels are the labels of the kinds of the parts, in the
// language of ctx: the tree's and the field's own.
func PermInputKindLabels(ctx context.Context) map[string]string {
	m := GetMessages(ctx)
	l := PermKindLabels(ctx)
	l[PermKindScope] = m.PermKindScope
	l[PermKindRecord] = m.PermKindRecord
	l[PermKindAll] = m.PermKindAll
	l[PermKindOther] = m.PermKindOther
	l[PermKindSet] = m.PermKindSet
	return l
}

// permInputMessages are the words of the field, in the language of ctx.
func permInputMessages(ctx context.Context) map[string]string {
	m := GetMessages(ctx)
	return map[string]string{
		"add": m.PermInputAdd, "addTo": m.PermInputAddTo, "addScope": m.PermInputAddScope,
		"remove": m.PermInputRemove, "children": m.PermInputChildren, "record": m.PermInputRecord,
		"recordSearch": m.PermInputRecordSearch, "recordNotFound": m.PermInputRecordNotFound,
		"noChoices": m.PermInputNoChoices, "empty": m.PermInputEmpty, "ok": m.PermInputOk,
		"cancel": m.PermInputCancel, "copy": m.PermInputCopy, "prefix": m.PermInputPrefix,
		"prefixNone": m.PermInputPrefixNone, "loadError": m.PermInputLoadError, "invalid": m.PermInputInvalid,
		"incomplete": m.PermInputIncomplete, "search": m.PermTreeSearch, "kinds": m.PermTreeKinds,
		"name": m.PermTreeName, "kind": m.PermTreeKind, "permission": m.PermTreePermission,
		"everything": m.PermInputEverything, "allRecords": m.PermInputAllRecords,
		"grant": m.PermInputGrant, "deny": m.PermInputDeny, "allow": m.PermInputAllow, "reset": m.PermInputReset, "asText": m.PermInputAsText, "onlyDefined": m.PermInputOnlyDefined,
		"fromText": m.PermInputFromText, "fromTextHint": m.PermInputFromTextHint, "paste": m.PermInputPaste,
		"apply": m.PermInputApply, "definedUnder": m.PermInputDefinedUnder, "showText": m.PermInputShowText, "hideText": m.PermInputHideText, "groups": m.PermTreeGroups, "byVerb": m.PermTreeByVerb, "byKind": m.PermTreeByKind,
		"expandAll": m.PermTreeExpandAll, "collapseAll": m.PermTreeCollapseAll, "path": m.PermTreePath,
		"description": m.PermTreeDescription, "active": m.PermTreeActive, "activeAll": m.PermTreeActiveAll,
		"activeSome": m.PermTreeActiveSome, "activeNone": m.PermTreeActiveNone,
		"sets": m.PermInputSets, "setGives": m.PermInputSetGives, "via": m.PermInputVia,
	}
}

// colonAt is where the first ":" of a text is, out of brackets and escapes;
// -1 when none.
func colonAt(s string) int {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case '[':
			depth++
		case ']':
			if depth > 0 {
				depth--
			}
		case ':':
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// PermissionInput is a vx-permission-input of o, its v-model (and the rest
// of its binding) the caller's: its tree asked of the event of o.
func PermissionInput(ctx *web.EventContext, o PermissionInputOptions) *h.HTMLTagBuilder {
	keys := make([]string, 0, len(o.Query))
	for k := range o.Query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	query := ""
	for _, k := range keys {
		query += ".query(" + strconv.Quote(k) + ", " + strconv.Quote(o.Query[k]) + ")"
	}
	// the params in the body (a value may be long): the op in the address
	loader := fmt.Sprintf(`(op, p) => plaid().url(%s).eventFunc(%s)%s.query("op", op).form(p || {}).go().then((r) => r && r.data)`,
		strconv.Quote(o.URL), strconv.Quote(o.Event), query)
	tag := h.Tag("vx-permission-input").
		Attr(":loader", loader).
		Attr(":kind-labels", h.JSONString(PermInputKindLabels(ctx.Context()))).
		Attr(":messages", h.JSONString(permInputMessages(ctx.Context())))
	if len(o.Prefixes) > 0 {
		tag.Attr(":prefixes", h.JSONString(o.Prefixes))
	}
	if o.Readonly {
		tag.Attr(":readonly", "true")
	}
	if o.StringView {
		tag.Attr(":string-view", "true")
	}
	if o.NoTextInput {
		tag.Attr(":text-input", "false")
	}
	if o.HelpURL != "" {
		tag.AppendChildren(h.Template(
			v.VBtn("").Icon("mdi-help-circle-outline").Variant(v.VariantText).Size(v.SizeSmall).
				Attr("href", o.HelpURL).Attr("target", "_blank").Attr("rel", "noopener").
				Attr("title", o.HelpTitle).Attr("data-permissions-help", true),
		).Attr("#bar-append", true))
	}
	return tag
}

// PermissionValue is the value of lines stored: one line (perm.JoinLines),
// canonical, the labels of its records as they are now (DescribeLine); lines
// it cannot join, one a line of text.
func PermissionValue(ctx *web.EventContext, p *presets.Builder, lines []string) string {
	v, err := perm.JoinLines(lines)
	if err != nil {
		return strings.Join(lines, "\n")
	}
	v, _ = DescribeLine(ctx, p, v)
	return v
}

// PermissionFieldComponentFunc is the field of the permissions of a stack
// (a role's, an access key's): a PermissionInput of its value
// (PermissionValue of its lines: a []string, a pq.StringArray), as opts says
// of the request. Its setter: PermissionFieldSetter.
func PermissionFieldComponentFunc(p *presets.Builder, opts func(ctx *web.EventContext) (PermissionInputOptions, error)) presets.FieldComponentFunc {
	return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		o, err := opts(ctx)
		if err != nil {
			return v.VAlert(h.Text(err.Error())).Type("error").Density(v.DensityCompact)
		}
		if field.ReadOnly {
			o.Readonly = true
		}
		return PermissionInput(ctx, o).
			Attr(web.VField(field.FormKey, PermissionValue(ctx, p, stringsOf(field.Value())))...).
			Attr("label", field.InputLabel()).
			Attr(":error-messages", h.JSONString(field.Errors)).
			Attr("data-permission-input", true)
	}
}

// PermissionFieldSetter sets the lines of the field of the form's value
// (perm.StoreLines): the items of its root, one a scope, what is given to a
// resource by its unique name, no placeholder others are under; a value it
// cannot read, as one line — the field's check says what is wrong.
func PermissionFieldSetter(obj any, field *presets.FieldContext, ctx *web.EventContext) error {
	value := ctx.R.FormValue(field.FormKey)
	lines, err := perm.StoreLines(value)
	if err != nil {
		lines = []string{value}
	}
	if lines == nil {
		lines = []string{}
	}
	return reflectutils.Set(obj, field.Name, lines)
}

// stringsOf are the texts of v, a list of them.
func stringsOf(v any) (out []string) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		return nil
	}
	for i := 0; i < rv.Len(); i++ {
		if e := rv.Index(i); e.Kind() == reflect.String {
			out = append(out, e.String())
		}
	}
	return
}
