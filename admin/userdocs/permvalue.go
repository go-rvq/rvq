package userdocs

import (
	"fmt"
	"strings"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
)

// The tree of a value — the permissions of a role, of an access key, as one
// line (perm.JoinLines) — as the field that edits it shows it
// (PermValueTree): the tree of the permissions of the admin
// (PermissionsTree), of the view and the groups asked, kept to what the value
// touches — what it allows, the items of a line of their own (Explicit, of
// their Prefix), and what they are under —; the records a line names, each a
// branch of its model ("Registro: título"); "*" the prefix of the branch
// before it (Set); the lines that are of no item, under "Outras". An edit (EditValue) changes
// the lines of the value — one put, one taken out, a prefix changed — and
// makes the value of them again, canonical: what is given to a resource by
// its unique name, out of the groups of the menu (perm.UngroupPattern).

// The prefixes of the items (PermTreeItem.Prefix).
const (
	PermPrefixGrant = "grant"
	PermPrefixDeny  = "deny"
	PermPrefixAllow = "allow"
	PermPrefixReset = "reset"
)

var prefixNames = map[perm.LineKind]string{perm.LineGrant: PermPrefixGrant, perm.LineDeny: PermPrefixDeny,
	perm.LineAllow: PermPrefixAllow, perm.LineReset: PermPrefixReset}

// lineKindOf is the kind of the prefix name; ok false when it is none ("",
// "none": no line).
func lineKindOf(name string) (perm.LineKind, bool) {
	for k, n := range prefixNames {
		if n == name {
			return k, true
		}
	}
	return perm.LineGrant, false
}

// valueLine is a line of a value: its kind, its pattern as written (its
// labels with it) and without its labels (key).
type valueLine struct {
	kind         perm.LineKind
	pattern, key string
	used         bool
}

func valueLines(value string) ([]*valueLine, error) {
	ls, err := perm.LinesOf(value)
	if err != nil {
		return nil, err
	}
	out := make([]*valueLine, 0, len(ls))
	for _, l := range ls {
		k, p := perm.ParseLine(l)
		out = append(out, &valueLine{kind: k, pattern: p, key: perm.StripLabels(p)})
	}
	return out, nil
}

// PermValueOptions are how the tree of a value is made.
type PermValueOptions struct {
	View PermTreeView
	// Groups: the groups of the menu as branches, each resource in its group
	Groups bool
	// Base, when not nil, is what the value is kept to — of an access key,
	// what its user may —: the parts of which it allows nothing are not
	// shown, and an item is active when both allow it
	Base func(PermCheck) bool
	// Effect, when not nil, is what the lines of the value allow (an access
	// key's: all but its denies and the protected); nil, the lines alone (a
	// stack of them)
	Effect func(lines []string) func(PermCheck) bool
	// Only: the items of a line of their own and what holds them — not all
	// the value allows
	Only bool
}

// PermValueTree is the tree of the value as the field shows it, in the
// language of ctx, as o says.
func PermValueTree(ctx *web.EventContext, p *presets.Builder, value string, o PermValueOptions) ([]*PermTreeItem, error) {
	vls, err := valueLines(value)
	if err != nil {
		return nil, err
	}
	var stripped []string
	for _, l := range vls {
		stripped = append(stripped, l.kind.Prefix()+l.key)
	}
	allows := o.Effect
	if allows == nil {
		allows = func(lines []string) func(PermCheck) bool {
			st := (&perm.Stack{}).Add("value", lines...).Resolve()
			return func(c PermCheck) bool { return st.Allows(c.Resource, c.Preferred, c.Action) }
		}
	}
	effect := allows(stripped)
	t := &permTree{ctx: ctx.Context(), msgs: GetMessages(ctx.Context())}
	items := PermissionsTree(ctx.Context(), p, PermTreeOptions{View: o.View, Groups: o.Groups})
	if o.Base != nil {
		SetActive(items, o.Base)
		items = prune(items, func(it *PermTreeItem) bool {
			return it.Active == PermActiveAll || it.Active == PermActiveSome
		})
	}
	setEdits(items)
	v := &valueTree{ctx: ctx, p: p, t: t, ix: newPermIndex(p), lines: vls}
	v.records(items)
	v.explicit(items)
	if !o.Groups {
		v.groupLines(items)
	}
	v.sets(items)
	SetActive(items, func(c PermCheck) bool { return (o.Base == nil || o.Base(c)) && effect(c) })
	items = prune(items, func(it *PermTreeItem) bool {
		return it.Explicit || !o.Only && (it.Active == PermActiveAll || it.Active == PermActiveSome)
	})
	items = v.others(items)
	setExpand(items)
	return items, nil
}

// valueTree makes the tree of a value.
type valueTree struct {
	ctx   *web.EventContext
	p     *presets.Builder
	t     *permTree
	ix    *permIndex
	lines []*valueLine
	// the tree with the groups of the menu (of the groups of lines shown
	// without them), made once asked
	grouped []*PermTreeItem
}

// walkItems calls f with each item, the children after their parent.
func walkItems(items []*PermTreeItem, f func(it *PermTreeItem)) {
	for _, it := range items {
		f(it)
		walkItems(it.Children, f)
	}
}

// setEdits sets the resource each item's line writes: a group's by the
// groups, a resource's by its unique name; a branch of the tree only
// ("Campos (3)"), none.
func setEdits(items []*PermTreeItem) {
	walkItems(items, func(it *PermTreeItem) {
		switch it.Kind {
		case "":
			it.Edit = ""
		case PermKindGroup:
			it.Edit = it.Value
		default:
			it.Edit = perm.UngroupPattern(it.Value)
		}
		it.CanAdd = strings.HasSuffix(it.Edit, ":")
		it.Set = setOf(it.Edit)
	})
}

// setOf is the resource the prefix of the item of edit writes: of a branch
// (":" at its end), what is under it ("*").
func setOf(edit string) string {
	if strings.HasSuffix(edit, ":") {
		return edit + "*"
	}
	return edit
}

// prune are the items keep keeps, and those under which some is kept.
func prune(items []*PermTreeItem, keep func(it *PermTreeItem) bool) (out []*PermTreeItem) {
	for _, it := range items {
		it.Children = prune(it.Children, keep)
		if keep(it) || len(it.Children) > 0 {
			out = append(out, it)
		}
	}
	return
}

// setExpand opens the items on the way to a line of their own; whether
// some item of items is one or is on the way.
func setExpand(items []*PermTreeItem) (any bool) {
	for _, it := range items {
		it.Expand = setExpand(it.Children)
		if it.Expand || it.Explicit {
			any = true
		}
	}
	return
}

// recordOf is the first record of the line that is one (not "*", no glob):
// the beginning of the line before it (its model's listing, as written),
// its id and its label; ok false when none.
func recordOf(pattern string) (model, id, label string, ok bool) {
	e, err := perm.ParseExpr(pattern)
	if err != nil {
		return
	}
	for i, n := range e.Nodes {
		r, isRecord := n.(*perm.ExprRecord)
		if !isRecord {
			continue
		}
		if r.ID == "*" || strings.ContainsAny(r.ID, "*?[") {
			continue
		}
		return perm.NodesString(e.Nodes[:i]), r.ID, r.Label, true
	}
	return
}

// records puts, under the item of each model whose records the lines name,
// a branch of each record ("Registro: título"): what the model's items ask
// of any record (<*>), asked of it.
func (v *valueTree) records(items []*PermTreeItem) {
	type rec struct{ model, id, label string }
	var recs []rec
	seen := map[string]bool{}
	for _, l := range v.lines {
		model, id, label, ok := recordOf(l.pattern)
		if !ok {
			continue
		}
		key := perm.StripLabels(model) + "<" + id + ">"
		if seen[key] {
			continue
		}
		seen[key] = true
		recs = append(recs, rec{model, id, label})
	}
	for _, r := range recs {
		modelKey := perm.StripLabels(r.model)
		var owner *PermTreeItem
		walkItems(items, func(it *PermTreeItem) {
			if owner == nil && (it.Kind == PermKindModel) && it.Edit == modelKey {
				owner = it
			}
		})
		if owner == nil {
			continue
		}
		label, found := recordLabel(v.ctx, v.ix, modelKey, r.id)
		if !found {
			label = r.label
			if label == "" {
				label = r.id
			}
		}
		fromName, toName := owner.Value+"<*>:", owner.Value+"<"+r.id+">:"
		fromUnique, toUnique := modelKey+"<*>:", modelKey+"<"+r.id+">:"
		labeled := modelKey + "<" + r.id + "<" + perm.EscapeLabel(label) + ">>:"
		children := cloneUnder(owner.Children, fromName, toName, fromUnique, toUnique)
		setEdits(children)
		walkItems(children, func(it *PermTreeItem) {
			it.Edit = strings.Replace(it.Edit, toUnique, labeled, 1)
		})
		owner.Children = append(owner.Children, &PermTreeItem{
			Title: fmt.Sprintf(v.t.msgs.PermInputRecordOf, label), Kind: PermKindRecord,
			Description: fmt.Sprintf(v.t.msgs.PermDescRecord, owner.Title),
			Value:       toName, Resources: []string{toUnique}, Edit: labeled, Set: setOf(labeled), CanAdd: true, Error: !found,
			Path:     append(append([]PermPathPart{}, owner.Path...), PermPathPart{owner.Value, owner.Title}),
			Children: children,
		})
	}
}

// cloneUnder are copies of the items of any record (their resources from
// fromName, by the groups, or fromUnique) — and of what holds them —, of one
// record (toName, toUnique).
func cloneUnder(items []*PermTreeItem, fromName, toName, fromUnique, toUnique string) (out []*PermTreeItem) {
	swap := func(s string) string {
		if strings.HasPrefix(s, fromName) {
			return toName + s[len(fromName):]
		}
		if strings.HasPrefix(s, fromUnique) {
			return toUnique + s[len(fromUnique):]
		}
		return s
	}
	for _, it := range items {
		children := cloneUnder(it.Children, fromName, toName, fromUnique, toUnique)
		if !strings.HasPrefix(it.Value, fromName) && len(children) == 0 {
			continue
		}
		c := *it
		c.Value = swap(it.Value)
		c.Resources = make([]string, len(it.Resources))
		for i, r := range it.Resources {
			c.Resources[i] = swap(r)
		}
		c.Checks = make([]PermCheck, len(it.Checks))
		for i, ch := range it.Checks {
			c.Checks[i] = PermCheck{Resource: swap(ch.Resource), Preferred: swap(ch.Preferred), Action: ch.Action}
		}
		c.Children = children
		out = append(out, &c)
	}
	return
}

// explicit marks the items of a line of their own, of its prefix: of a
// branch, the line of what is under it ("…:*"); the line of the branch
// alone ("…:", what follows it missing), incomplete.
func (v *valueTree) explicit(items []*PermTreeItem) {
	walkItems(items, func(it *PermTreeItem) {
		if it.Edit != "" {
			v.mark(it)
		}
	})
}

// mark marks the item of the lines of its own.
func (v *valueTree) mark(it *PermTreeItem) {
	key, set := perm.StripLabels(it.Edit), perm.StripLabels(it.Set)
	for _, l := range v.lines {
		switch {
		case l.key == set:
			it.Explicit, it.Prefix, it.Incomplete = true, prefixNames[l.kind], false
			l.used = true
		case l.key == key && key != set:
			if !it.Explicit {
				it.Explicit, it.Prefix, it.Incomplete = true, prefixNames[l.kind], true
			}
			l.used = true
		}
	}
}

// groupLines shows, without the groups of the menu, the lines given to
// groups: each group a row of its scope, not what is in it.
func (v *valueTree) groupLines(items []*PermTreeItem) {
	shown := map[string]*PermTreeItem{}
	for _, l := range v.lines {
		if l.used {
			continue
		}
		group, _, ok := groupOfLine(l.key)
		if !ok {
			continue
		}
		var scope *PermTreeItem
		for _, it := range items {
			if it.Kind == PermKindScope && strings.HasPrefix(group, it.Edit) {
				scope = it
			}
		}
		if scope == nil {
			continue
		}
		g := shown[group]
		if g == nil {
			g = &PermTreeItem{Title: group, Kind: PermKindGroup, Value: group, Resources: []string{group}, Edit: group,
				Set: setOf(group), CanAdd: true}
			if gi := v.groupItem(group); gi != nil {
				g.Title, g.Description, g.Path = gi.Title, gi.Description, gi.Path
			}
			shown[group] = g
			scope.Children = append(scope.Children, g)
		}
		v.mark(g)
	}
}

// groupOfLine is the group the line is given to — "admin:content/:sub/:" of
// "admin:content/:sub/:" (all false) or of "admin:content/:sub/:*" (all) —;
// ok false when it is not one's.
func groupOfLine(key string) (group string, all bool, ok bool) {
	rest := key
	if strings.HasSuffix(rest, ":*") {
		rest, all = strings.TrimSuffix(rest, "*"), true
	}
	i := strings.Index(rest, ":")
	if i < 0 || i == len(rest)-1 {
		return "", false, false
	}
	for _, s := range strings.SplitAfter(rest[i+1:], ":") {
		if s == "" {
			continue
		}
		if !strings.HasSuffix(s, "/:") || strings.ContainsAny(s, "*?[{<") {
			return "", false, false
		}
	}
	return rest, all, true
}

// groupItem is the item of the group of the menu (by the groups), nil when
// none.
func (v *valueTree) groupItem(group string) (found *PermTreeItem) {
	if v.grouped == nil {
		v.grouped = PermissionsTree(v.ctx.Context(), v.p, PermTreeOptions{Groups: true})
	}
	walkItems(v.grouped, func(it *PermTreeItem) {
		if found == nil && it.Kind == PermKindGroup && it.Value == group {
			found = it
		}
	})
	return
}

// sets puts each line of a set under the item of its resource: a row of
// the set ("Leitura", %readers), of the line's prefix, with the lines it
// gives there — none, an error: it gives nothing there, ignored. A line of
// a set of no item (a glob, a group) is of "Outras" (others).
func (v *valueTree) sets(items []*PermTreeItem) {
	for _, l := range v.lines {
		if l.used {
			continue
		}
		res, name, ok := perm.SetRef(l.key)
		if !ok {
			continue
		}
		var owner *PermTreeItem
		walkItems(items, func(it *PermTreeItem) {
			if owner == nil && it.Edit == res && it.Kind != PermKindRecord {
				owner = it
			}
		})
		if owner == nil {
			continue
		}
		l.used = true
		it := v.setItem(l, res, name)
		it.Path = append(append([]PermPathPart{}, owner.Path...), PermPathPart{owner.Value, owner.Title})
		owner.Children = append(owner.Children, it)
	}
}

// setItem is the row of the line l of the set name given to res.
func (v *valueTree) setItem(l *valueLine, res, name string) *PermTreeItem {
	ctx := v.ctx.Context()
	title := ""
	if n := v.ix.byName[res]; n != nil {
		for _, o := range v.p.PermSetsOf(ctx, n) {
			if o.Name == name {
				title = o.Title
			}
		}
	}
	if s := v.p.PermSet(name); title == "" && s != nil {
		title = s.TitleOf(ctx)
	}
	if title == "" {
		title = name
	}
	lines, ok := perm.ExpandSet(res, name)
	it := &PermTreeItem{Title: title, Subtitle: perm.SetPrefix + name, Kind: PermKindSet, Value: l.key,
		Resources: []string{l.pattern}, Edit: l.pattern, Set: l.pattern, Explicit: true, Prefix: prefixNames[l.kind],
		Lines: relativeLines(res, lines), Error: !ok}
	if s := v.p.PermSet(name); s != nil {
		it.Description = s.DescriptionOf(ctx)
	}
	if !ok {
		it.Description = fmt.Sprintf(v.t.msgs.PermInputSetIgnored, name)
	}
	return it
}

// others are the items, and a branch "Outras" of the lines that are of no
// item (their own globs, resources there are no more).
func (v *valueTree) others(items []*PermTreeItem) []*PermTreeItem {
	var rows []*PermTreeItem
	for _, l := range v.lines {
		if l.used {
			continue
		}
		if res, name, ok := perm.SetRef(l.key); ok {
			it := v.setItem(l, res, name)
			it.Title, it.Value = it.Title+" · "+l.pattern, "(others)"+l.key
			rows = append(rows, it)
			continue
		}
		rows = append(rows, &PermTreeItem{Title: l.pattern, Value: "(others)" + l.key, Resources: []string{l.pattern},
			Edit: l.pattern, Set: l.pattern, Explicit: true, Prefix: prefixNames[l.kind], Kind: PermKindOther,
			Incomplete: strings.HasSuffix(l.key, ":"), CanAdd: strings.HasSuffix(l.key, ":")})
	}
	if len(rows) == 0 {
		return items
	}
	return append(items, &PermTreeItem{Title: v.t.msgs.PermInputOthers, Description: v.t.msgs.PermInputOthersDesc,
		Value: "(others)", Children: rows})
}

// PermKindOther is the kind of a line of no item (PermValueTree's "Outras").
const PermKindOther = "other"

// The edits of a value (EditValue).
const (
	// PermEditSet: the line of the resource of the prefix ("" or "none":
	// no line) — in the place of its line, else after the others; of what
	// is under a branch ("…:*"), the line of the branch alone ("…:") out
	PermEditSet = "set"
	// PermEditRemove: no line of the resources, nor of what is under them
	PermEditRemove = "remove"
	// PermEditAdd: a line of each resource, of the prefix ("": grant), those
	// there are already kept
	PermEditAdd = "add"
	// PermEditMerge: the lines given (resources: lines, each of its
	// prefix) merged in the value (perm.ExprAlt.Merge) — none: the value
	// only as stored
	PermEditMerge = "merge"
)

// EditValue is the value with the edit op of the resources (the prefix, of
// set and add), canonical: each line out of the groups of the menu
// (perm.UngroupPattern), the placeholders others are under out
// (perm.DropPlaceholders), the lines joined (perm.JoinLines).
func EditValue(value, op string, resources []string, prefix string) (string, error) {
	lines, err := perm.LinesOf(value)
	if err != nil {
		return "", err
	}
	keyOf := func(l string) string {
		_, p := perm.ParseLine(l)
		return perm.StripLabels(p)
	}
	var out []string
	switch op {
	case PermEditSet:
		if len(resources) == 0 {
			return value, nil
		}
		r := resources[0]
		key := perm.StripLabels(r)
		alone := ""
		if strings.HasSuffix(key, ":*") {
			alone = strings.TrimSuffix(key, "*")
		}
		k, ok := lineKindOf(prefix)
		placed := false
		for _, l := range lines {
			if lk := keyOf(l); lk == key || alone != "" && lk == alone {
				if ok && !placed {
					out, placed = append(out, k.Prefix()+r), true
				}
				continue
			}
			out = append(out, l)
		}
		if ok && !placed {
			out = append(out, k.Prefix()+r)
		}
	case PermEditRemove:
		keys := make([]string, len(resources))
		for i, r := range resources {
			keys[i] = perm.StripLabels(r)
		}
	lines:
		for _, l := range lines {
			lk := keyOf(l)
			for _, k := range keys {
				if lk == k || strings.HasSuffix(k, ":") && strings.HasPrefix(lk, k) {
					continue lines
				}
			}
			out = append(out, l)
		}
	case PermEditMerge:
		v, err := perm.ParseExpr(value)
		if err != nil {
			return "", err
		}
		var others []*perm.ExprAlt
		for _, r := range resources {
			e, err := perm.ParseExpr(r)
			if err != nil {
				return "", err
			}
			others = append(others, e)
		}
		out = v.Merge(others...).Lines()
	case PermEditAdd:
		k, ok := lineKindOf(prefix)
		if !ok {
			k = perm.LineGrant
		}
		out = lines
		have := map[string]bool{}
		for _, l := range lines {
			have[keyOf(l)] = true
		}
		for _, r := range resources {
			if key := perm.StripLabels(r); r != "" && !have[key] {
				have[key] = true
				out = append(out, k.Prefix()+r)
			}
		}
	default:
		return "", fmt.Errorf("userdocs: no edit %q of a value", op)
	}
	for i, l := range out {
		k, p := perm.ParseLine(l)
		out[i] = k.Prefix() + perm.UngroupPattern(p)
	}
	return perm.JoinLines(perm.DropPlaceholders(out))
}

// groupMembers are, of the group of the menu (by the groups), the resources
// of what is in it: what taking the group out takes too.
func groupMembers(ctx *web.EventContext, p *presets.Builder, group string) (out []string) {
	v := &valueTree{ctx: ctx, p: p}
	g := v.groupItem(group)
	if g == nil {
		return nil
	}
	items := []*PermTreeItem{g}
	setEdits(items)
	walkItems(g.Children, func(it *PermTreeItem) {
		if it.Edit != "" && strings.HasSuffix(it.Edit, ":") {
			out = append(out, it.Edit)
		}
	})
	return
}
