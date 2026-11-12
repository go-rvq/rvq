package userdocs

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/go-rvq/rvq/admin/presets"
)

// The tree of the permissions of the admin (PermissionsTree), down to its
// leaves: a leaf is one permission — a resource with the verb asked of it,
// "admin:posts:<*>:#Body:@edit" —, a branch what holds them. Two views of
// the same permissions:
//
//   - by verb (PermTreeByVerb): a model by what is asked of it — listing,
//     creating, seeing a record, editing it… —, each verb with what it asks
//     in turn: the fields of the listing's columns under @list, the record's
//     fields, sections (each with its fields), actions, pages and nested
//     models under @get;
//   - by kind (PermTreeByKind): a model by what it holds — its permissions,
//     actions, fields, sections, pages, nested models —, each part with the
//     verbs asked of it.
//
// The groups of the menu are in the tree only when asked (Groups): without
// them, the resources are by the unique names alone — what decides first
// (presets Verifier.Prefer) —, the models and pages of every group together.

// PermTreeView is how the tree of the permissions is drawn.
type PermTreeView string

const (
	// PermTreeByVerb: a model by the verbs asked of it, each with what it
	// asks.
	PermTreeByVerb PermTreeView = "verb"
	// PermTreeByKind: a model by what it holds, each part with its verbs.
	PermTreeByKind PermTreeView = "kind"
)

// PermTreeOptions are how the tree of the permissions is made.
type PermTreeOptions struct {
	View PermTreeView
	// Groups: the groups of the menu as branches, and the resources by them
	// too (beside the unique names)
	Groups bool
}

// PermTreeItem is a node of the tree of the permissions: its title, what it
// is (Description) and its resources — a leaf's, each with its verb.
// Subtitle is both, for a VTreeview (the documentation's).
type PermTreeItem struct {
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle,omitempty"`
	Description string   `json:"description,omitempty"`
	Resources   []string `json:"resources,omitempty"`
	Value       string   `json:"value"`
	// Kind is what it is (PermKindModel, PermKindField…); "" a branch of
	// the tree only ("Fields (3)")
	Kind string `json:"kind,omitempty"`
	// Path is where a part of the admin (a group, a model, a page, a
	// check) is: the groups of the menu and the models it is in, the
	// outermost first, each [id, label] — breadcrumbs ("Site › SEO"), two of
	// the same title told apart
	Path []PermPathPart `json:"path,omitempty"`
	// Active: of a user's stack (SetActive), whether it allows what the item
	// asks and what is under it — PermActiveAll, PermActiveSome,
	// PermActiveNone
	Active   string          `json:"active,omitempty"`
	Children []*PermTreeItem `json:"children,omitempty"`
	// Checks are the permissions the item asks, each as a Verifier asks it
	Checks []PermCheck `json:"-"`

	// Of the tree of a value (PermValueTree), what edits it:
	//
	// Edit is the resource the item's line writes — a resource by its unique
	// name, a group by the groups —; "" not editable (a branch of the tree
	// only)
	Edit string `json:"edit,omitempty"`
	// Set is the resource the prefix of the item writes: of a branch, what is
	// under it ("admin:posts:*"); of a leaf, Edit
	Set string `json:"set,omitempty"`
	// Prefix is of an item of a line of its own (Explicit — of a branch, its
	// "*"): "grant", "deny", "allow", "reset"
	Prefix   string `json:"prefix,omitempty"`
	Explicit bool   `json:"explicit,omitempty"`
	// CanAdd: parts may follow it (the + of the field)
	CanAdd bool `json:"canAdd,omitempty"`
	// Incomplete: a line of its own ending by ":" — what follows it missing,
	// no "*" of it
	Incomplete bool `json:"incomplete,omitempty"`
	// Error: a record not found (or the user may not read)
	Error bool `json:"error,omitempty"`
	// Expand: on the way to a line of its own — opened when shown
	Expand bool `json:"expand,omitempty"`
}

// The kinds of the items of the tree (PermTreeItem.Kind).
const (
	PermKindScope     = "scope"     // the scope of the admin, the root
	PermKindGroup     = "group"     // a group of the menu
	PermKindModel     = "model"     // a model: its listing, its records
	PermKindSingleton = "singleton" // a model of a single record
	PermKindVerb      = "verb"      // a permission asked (@list, @edit, a page's @get…)
	PermKindAction    = "action"    // an action (!publish)
	PermKindField     = "field"     // a field
	PermKindInline    = "inline"    // a field holding a model edited in place
	PermKindSection   = "section"   // a section of a detail
	PermKindPage      = "page"      // a page
	PermKindCheck     = "check"     // a permission of its own (a verifier)
)

// PermKinds are the kinds of the items, in the order a filter offers them.
var PermKinds = []string{PermKindGroup, PermKindModel, PermKindSingleton, PermKindVerb, PermKindAction,
	PermKindField, PermKindInline, PermKindSection, PermKindPage, PermKindCheck}

// kindOf is the kind of the item of the node n.
func kindOf(n *presets.PermNode) string {
	switch n.Kind {
	case presets.PermNodeGroup:
		return PermKindGroup
	case presets.PermNodeModel, presets.PermNodeRecord:
		return PermKindModel
	case presets.PermNodeSingleton:
		return PermKindSingleton
	case presets.PermNodeField:
		return PermKindField
	case presets.PermNodeInline:
		return PermKindInline
	case presets.PermNodeSection:
		return PermKindSection
	case presets.PermNodePage:
		return PermKindPage
	}
	return PermKindCheck
}

// verbKind is the kind of the permission asked by verb: a verb (@…), an
// action (!…).
func verbKind(verb string) string {
	if strings.HasPrefix(verb, "!") {
		return PermKindAction
	}
	return PermKindVerb
}

// PermCheck is a permission asked: the resource with its verb at the end, by
// the groups and by the unique name (Preferred, "" none), and the verb.
type PermCheck struct {
	Resource, Preferred, Action string
}

// The states of PermTreeItem.Active.
const (
	PermActiveAll  = "all"
	PermActiveSome = "some"
	PermActiveNone = "none"
)

// SetActive sets Active of the items and of their children: allows says
// whether a permission is allowed (perm.Stack.Allows).
func SetActive(items []*PermTreeItem, allows func(c PermCheck) bool) {
	for _, it := range items {
		SetActive(it.Children, allows)
		var yes, no int
		for _, c := range it.Checks {
			if allows(c) {
				yes++
			} else {
				no++
			}
		}
		for _, c := range it.Children {
			switch c.Active {
			case PermActiveAll:
				yes++
			case PermActiveNone:
				no++
			case PermActiveSome:
				yes, no = yes+1, no+1
			}
		}
		switch {
		case yes > 0 && no == 0:
			it.Active = PermActiveAll
		case yes > 0:
			it.Active = PermActiveSome
		case no > 0:
			it.Active = PermActiveNone
		}
	}
}

// SetSubtitles sets Subtitle of the items and of their children: the
// description and the resources, by " · ".
func SetSubtitles(items []*PermTreeItem) {
	for _, it := range items {
		var parts []string
		if it.Description != "" {
			parts = append(parts, it.Description)
		}
		it.Subtitle = strings.Join(append(parts, it.Resources...), " · ")
		SetSubtitles(it.Children)
	}
}

// PermissionsTree is the tree of the permissions of the admin p
// (presets.Builder.BuildPermissions, made now: from the menu as it is), in
// the language of ctx, as opts says: its root the scope of the admin.
func PermissionsTree(ctx context.Context, p *presets.Builder, opts PermTreeOptions) []*PermTreeItem {
	t := &permTree{ctx: ctx, msgs: GetMessages(ctx), opts: opts}
	scope := presets.PermModule
	root := &PermTreeItem{Title: t.msgs.PermScopeAdmin, Description: t.msgs.PermScopeAdminDesc,
		Resources: []string{scope + ":*"}, Value: scope + ":", Kind: PermKindScope}
	root.Children = t.nodes(p.BuildPermissions().Tree().Children)
	if !opts.Groups {
		// the models and pages of every group together: by their titles,
		// then where they are
		sort.SliceStable(root.Children, func(i, j int) bool {
			a, b := root.Children[i], root.Children[j]
			if ta, tb := strings.ToLower(a.Title), strings.ToLower(b.Title); ta != tb {
				return ta < tb
			}
			return strings.ToLower(pathText(a.Path)) < strings.ToLower(pathText(b.Path))
		})
	}
	return []*PermTreeItem{root}
}

// permTree makes the tree of the permissions.
type permTree struct {
	ctx  context.Context
	msgs *Messages
	opts PermTreeOptions
}

// title is the label of n in the language of the request: its own, else by
// its name.
func (t *permTree) title(n *presets.PermNode) string {
	if n == nil {
		return ""
	}
	if n.Title != nil {
		if s := n.Title(t.ctx); s != "" {
			return s
		}
	}
	parts := strings.Split(strings.TrimSuffix(n.Name, ":"), ":")
	return fmt.Sprintf(t.msgs.PermUntitled, presets.HumanizeString(strings.TrimLeft(parts[len(parts)-1], "#$&/<")))
}

// describe is what n is, in the language of the request: its own
// description, else by its kind.
func (t *permTree) describe(n *presets.PermNode, label string) string {
	if n.Description != nil {
		if d := n.Description(t.ctx); d != "" {
			return d
		}
	}
	msgs, parent := t.msgs, t.title(n.Parent)
	switch n.Kind {
	case presets.PermNodeGroup:
		return fmt.Sprintf(msgs.PermDescGroup, label)
	case presets.PermNodeModel:
		if n.Parent != nil && n.Parent.Kind == presets.PermNodeRecord {
			return fmt.Sprintf(msgs.PermDescNested, label, parent)
		}
		return fmt.Sprintf(msgs.PermDescModel, label)
	case presets.PermNodeSingleton:
		return fmt.Sprintf(msgs.PermDescSingleton, label)
	case presets.PermNodeRecord:
		return fmt.Sprintf(msgs.PermDescRecord, label)
	case presets.PermNodeField:
		return fmt.Sprintf(msgs.PermDescField, label, parent)
	case presets.PermNodeInline:
		return fmt.Sprintf(msgs.PermDescInline, label, parent)
	case presets.PermNodeSection:
		return fmt.Sprintf(msgs.PermDescSection, label, parent)
	case presets.PermNodePage:
		if n.Parent != nil && n.Parent.Kind != presets.PermNodeGroup && n.Parent.Kind != "" {
			return fmt.Sprintf(msgs.PermDescPageOf, label, parent)
		}
		return fmt.Sprintf(msgs.PermDescPage, label)
	}
	return fmt.Sprintf(msgs.PermDescCheck, label)
}

// resources are the resources of n with verb at the end: by the unique
// name; with the groups, by the groups too (first). A resource with no
// unique name, by the groups.
func (t *permTree) resources(n *presets.PermNode, verb string) []string {
	if n.Unique == "" || n.Unique == n.Name {
		return []string{n.Name + verb}
	}
	if !t.opts.Groups {
		return []string{n.Unique + verb}
	}
	return []string{n.Name + verb, n.Unique + verb}
}

// check is the permission of n asked with verb, as a Verifier asks it.
func (t *permTree) check(n *presets.PermNode, verb string) PermCheck {
	c := PermCheck{Resource: n.Name + verb, Action: verb}
	if n.Unique != "" && n.Unique != n.Name {
		c.Preferred = n.Unique + verb
	}
	return c
}

// branch is the item of n: its title, description and resources, no verb;
// of a part of the admin, where it is (Path).
func (t *permTree) branch(n *presets.PermNode) *PermTreeItem {
	label := t.title(n)
	it := &PermTreeItem{Title: label, Description: t.describe(n, label), Resources: t.resources(n, ""), Value: n.Name,
		Kind: kindOf(n)}
	switch it.Kind {
	case PermKindGroup, PermKindModel, PermKindSingleton, PermKindPage, PermKindCheck:
		it.Path = t.path(n)
	}
	return it
}

// PermPathPart is a part of PermTreeItem.Path: [its id — the resource of
// the group or model, by the groups —, its label].
type PermPathPart [2]string

// PathSep separates the labels of a path shown.
const PathSep = " › "

// pathText is the labels of path, by PathSep.
func pathText(path []PermPathPart) string {
	labels := make([]string, len(path))
	for i, p := range path {
		labels[i] = p[1]
	}
	return strings.Join(labels, PathSep)
}

// path is where n is: the groups of the menu and the models above it, the
// outermost first.
func (t *permTree) path(n *presets.PermNode) (parts []PermPathPart) {
	for p := n.Parent; p != nil; p = p.Parent {
		switch p.Kind {
		case presets.PermNodeGroup, presets.PermNodeModel, presets.PermNodeSingleton:
			parts = append([]PermPathPart{{p.Name, t.title(p)}}, parts...)
		}
	}
	return
}

// leaf is the permission of n asked with verb: its title, its
// description, its kind.
func (t *permTree) leaf(n *presets.PermNode, title, desc, verb, kind string) *PermTreeItem {
	return &PermTreeItem{Title: title, Description: desc, Resources: t.resources(n, verb), Value: n.Name + verb, Kind: kind,
		Checks: []PermCheck{t.check(n, verb)}}
}

// verbLeaf is the leaf of the verb or action a of n.
func (t *permTree) verbLeaf(n *presets.PermNode, a *presets.PermNodeAction, label string) *PermTreeItem {
	title := verbLabel(t.msgs, a, t.ctx)
	if n.Kind == presets.PermNodePage {
		title = pageVerbLabel(t.msgs, a, t.ctx)
	}
	desc := ""
	if a.Description != nil {
		desc = a.Description(t.ctx)
	}
	if desc == "" && !strings.HasPrefix(a.Name, "@") {
		desc = fmt.Sprintf(t.msgs.PermDescAction, title, label)
	}
	return t.leaf(n, title, desc, a.Name, verbKind(a.Name))
}

// group is the branch of children under parent (its value), key telling it
// apart, its title with how many they are; nil for none.
func (t *permTree) group(parent, key, title, desc, label string, children []*PermTreeItem) *PermTreeItem {
	if len(children) == 0 {
		return nil
	}
	return &PermTreeItem{Title: fmt.Sprintf("%s (%d)", title, len(children)), Description: fmt.Sprintf(desc, label),
		Value: parent + "(" + key + ")", Children: children}
}

// add appends the items not nil.
func add(out []*PermTreeItem, items ...*PermTreeItem) []*PermTreeItem {
	for _, it := range items {
		if it != nil {
			out = append(out, it)
		}
	}
	return out
}

// nodes are the items of nodes: without the groups of the menu, the groups'
// children in their place.
func (t *permTree) nodes(nodes []*presets.PermNode) (out []*PermTreeItem) {
	for _, n := range nodes {
		if n.Kind == presets.PermNodeGroup && !t.opts.Groups {
			out = append(out, t.nodes(n.Children)...)
			continue
		}
		out = append(out, t.node(n))
	}
	return
}

// node is the item of n: a group with its children; a model; a page or a
// check with its verbs.
func (t *permTree) node(n *presets.PermNode) *PermTreeItem {
	switch n.Kind {
	case presets.PermNodeModel, presets.PermNodeSingleton:
		return t.model(n)
	case presets.PermNodeGroup:
		it := t.branch(n)
		it.Children = t.nodes(n.Children)
		return it
	}
	it := t.branch(n)
	if len(n.Actions) == 0 {
		it.Checks = []PermCheck{t.check(n, "")}
	}
	for _, a := range n.Actions {
		it.Children = append(it.Children, t.verbLeaf(n, a, it.Title))
	}
	it.Children = append(it.Children, t.nodes(n.Children)...)
	return it
}

// parts are the children of a model's or record's node, by their kinds.
type parts struct {
	record                                           *presets.PermNode
	fields, inlines, sections, pages, models, checks []*presets.PermNode
	verbs, actions                                   []*presets.PermNodeAction
}

func partsOf(n *presets.PermNode) (p parts) {
	for _, a := range n.Actions {
		if strings.HasPrefix(a.Name, "@") {
			p.verbs = append(p.verbs, a)
		} else {
			p.actions = append(p.actions, a)
		}
	}
	for _, c := range n.Children {
		switch c.Kind {
		case presets.PermNodeRecord:
			p.record = c
		case presets.PermNodeField:
			p.fields = append(p.fields, c)
		case presets.PermNodeInline:
			p.inlines = append(p.inlines, c)
		case presets.PermNodeSection:
			p.sections = append(p.sections, c)
		case presets.PermNodePage:
			p.pages = append(p.pages, c)
		case presets.PermNodeModel, presets.PermNodeSingleton:
			p.models = append(p.models, c)
		default:
			p.checks = append(p.checks, c)
		}
	}
	return
}

// model is the item of the model (or singleton) n, in the view asked.
func (t *permTree) model(n *presets.PermNode) *PermTreeItem {
	it := t.branch(n)
	mp := partsOf(n)
	rec, rp := n, mp
	if mp.record != nil {
		rec, rp = mp.record, partsOf(mp.record)
	}
	if t.opts.View == PermTreeByKind {
		it.Children = t.byKind(it, n, rec, mp, rp)
	} else {
		it.Children = t.byVerb(it, n, rec, mp, rp)
	}
	return it
}

// verbOrder is the order of the verbs of a model: listing, creating, seeing,
// editing, deleting; the others after, by name.
var verbOrder = map[string]int{presets.PermList: 1, presets.PermCreate: 2, presets.PermGet: 3, presets.PermUpdate: 4,
	presets.PermDelete: 5, presets.PermDeleteWithRelated: 6}

// byVerb are the children of the model n by its verbs: each with what it
// asks — the fields of its form; under @list, what the listing has (its
// actions, pages, checks); under @get, what the record has (its sections,
// actions, pages, nested models, checks); under @edit, the sections too.
func (t *permTree) byVerb(it *PermTreeItem, n, rec *presets.PermNode, mp, rp parts) (out []*PermTreeItem) {
	label, msgs := it.Title, t.msgs
	type owned struct {
		node *presets.PermNode
		p    parts
		a    *presets.PermNodeAction
	}
	var verbs []owned
	for _, a := range mp.verbs {
		verbs = append(verbs, owned{n, mp, a})
	}
	if rec != n {
		for _, a := range rp.verbs {
			verbs = append(verbs, owned{rec, rp, a})
		}
	}
	rank := func(name string) int {
		if r := verbOrder[name]; r > 0 {
			return r
		}
		return len(verbOrder) + 1
	}
	sort.SliceStable(verbs, func(i, j int) bool {
		ri, rj := rank(verbs[i].a.Name), rank(verbs[j].a.Name)
		if ri != rj {
			return ri < rj
		}
		return verbs[i].a.Name < verbs[j].a.Name
	})
	// what the listing has goes under @list, what the record has under @get
	// (a singleton's: its node is its record); with no such verb, in the
	// model itself
	extras := func(owner *presets.PermNode, p parts, under string) []*PermTreeItem {
		var actions, pages, models, checks []*PermTreeItem
		for _, a := range p.actions {
			actions = append(actions, t.verbLeaf(owner, a, label))
		}
		for _, c := range p.pages {
			pages = append(pages, t.node(c))
		}
		for _, c := range p.models {
			models = append(models, t.model(c))
		}
		for _, c := range p.checks {
			checks = append(checks, t.node(c))
		}
		return add(nil,
			t.group(under, "actions", msgs.PermGroupActions, msgs.PermGroupActionsDesc, label, actions),
			t.group(under, "pages", msgs.PermGroupPages, msgs.PermGroupPagesDesc, label, pages),
			t.group(under, "models", msgs.PermGroupModels, msgs.PermGroupModelsDesc, label, models),
			t.group(under, "checks", msgs.PermGroupChecks, msgs.PermGroupChecksDesc, label, checks))
	}
	extrasOf := map[*presets.PermNode]string{n: presets.PermList, rec: presets.PermGet}
	placed := map[*presets.PermNode]bool{}
	for _, o := range verbs {
		v := t.verbLeaf(o.node, o.a, label)
		v.Children = add(v.Children,
			t.group(v.Value, "fields", msgs.PermGroupFields, msgs.PermGroupFieldsDesc, label, t.verbFields(o.p.fields, o.a.Name)),
			t.group(v.Value, "inlines", msgs.PermGroupInlines, msgs.PermGroupInlinesDesc, label, t.verbFields(o.p.inlines, o.a.Name)))
		var sections []*PermTreeItem
		for _, s := range o.p.sections {
			if a := action(s, o.a.Name); a != nil {
				sections = append(sections, t.verbSection(s, a, o.p))
			}
		}
		v.Children = add(v.Children, t.group(v.Value, "sections", msgs.PermGroupSections, msgs.PermGroupSectionsDesc, label, sections))
		if extrasOf[o.node] == o.a.Name && !placed[o.node] {
			placed[o.node] = true
			v.Children = append(v.Children, extras(o.node, o.p, v.Value)...)
		}
		out = append(out, v)
	}
	for _, x := range []struct {
		node *presets.PermNode
		p    parts
	}{{n, mp}, {rec, rp}} {
		if !placed[x.node] {
			placed[x.node] = true
			out = append(out, extras(x.node, x.p, x.node.Name)...)
		}
	}
	return
}

// action is the action of n of the name, nil when it has none.
func action(n *presets.PermNode, name string) *presets.PermNodeAction {
	for _, a := range n.Actions {
		if a.Name == name {
			return a
		}
	}
	return nil
}

// verbFields are the fields among nodes asked verb, each its permission
// with the fields nested in it asked the same.
func (t *permTree) verbFields(nodes []*presets.PermNode, verb string) (out []*PermTreeItem) {
	for _, f := range nodes {
		if action(f, verb) != nil {
			out = append(out, t.verbField(f, verb))
		}
	}
	return
}

// verbField is the permission of the field f asked verb, the fields nested
// in it (a model edited in place) under it.
func (t *permTree) verbField(f *presets.PermNode, verb string) *PermTreeItem {
	label := t.title(f)
	it := t.leaf(f, label, t.describe(f, label), verb, kindOf(f))
	p := partsOf(f)
	it.Children = add(nil,
		t.group(it.Value, "fields", t.msgs.PermGroupFields, t.msgs.PermGroupFieldsDesc, label, t.verbFields(p.fields, verb)),
		t.group(it.Value, "inlines", t.msgs.PermGroupInlines, t.msgs.PermGroupInlinesDesc, label, t.verbFields(p.inlines, verb)))
	return it
}

// sectionFields are the fields of the record (p) the section holds — by
// the names of a — asked verb.
func (t *permTree) sectionFields(p parts, a *presets.PermNodeAction, verb string) (out []*presets.PermNode) {
	for _, name := range a.Fields {
		for _, f := range append(append([]*presets.PermNode{}, p.fields...), p.inlines...) {
			key := f.Name[strings.LastIndex(strings.TrimSuffix(f.Name, ":"), ":")+1:]
			if key == presets.FieldPerm(name)+":" || key == presets.InlinePerm(name)+":" {
				if verb == "" || action(f, verb) != nil {
					out = append(out, f)
				}
			}
		}
	}
	return
}

// verbSection is the permission of the section s asked a, the fields of the
// record it holds under it, asked the same.
func (t *permTree) verbSection(s *presets.PermNode, a *presets.PermNodeAction, p parts) *PermTreeItem {
	label := t.title(s)
	it := t.leaf(s, label, t.describe(s, label), a.Name, PermKindSection)
	var fields []*PermTreeItem
	for _, f := range t.sectionFields(p, a, a.Name) {
		fields = append(fields, t.verbField(f, a.Name))
	}
	it.Children = add(nil, t.group(it.Value, "fields", t.msgs.PermGroupFields, t.msgs.PermGroupFieldsDesc, label, fields))
	return it
}

// byKind are the children of the model n by what it holds: its permissions,
// actions, fields, models edited in place, sections, pages, nested models,
// checks — each with the verbs asked of it.
func (t *permTree) byKind(it *PermTreeItem, n, rec *presets.PermNode, mp, rp parts) []*PermTreeItem {
	label, msgs := it.Title, t.msgs
	var verbs, actions []*PermTreeItem
	owners := []struct {
		node *presets.PermNode
		p    parts
	}{{n, mp}}
	if rec != n {
		owners = append(owners, struct {
			node *presets.PermNode
			p    parts
		}{rec, rp})
	}
	var fields, inlines, pages, models, checks []*presets.PermNode
	for _, o := range owners {
		for _, a := range o.p.verbs {
			verbs = append(verbs, t.verbLeaf(o.node, a, label))
		}
		for _, a := range o.p.actions {
			actions = append(actions, t.verbLeaf(o.node, a, label))
		}
		fields = append(fields, o.p.fields...)
		inlines = append(inlines, o.p.inlines...)
		pages = append(pages, o.p.pages...)
		models = append(models, o.p.models...)
		checks = append(checks, o.p.checks...)
	}
	var sections, pageItems, modelItems, checkItems []*PermTreeItem
	for _, s := range rp.sections {
		sections = append(sections, t.kindSection(s, rp))
	}
	for _, c := range pages {
		pageItems = append(pageItems, t.node(c))
	}
	for _, c := range models {
		modelItems = append(modelItems, t.model(c))
	}
	for _, c := range checks {
		checkItems = append(checkItems, t.node(c))
	}
	return add(nil,
		t.group(n.Name, "verbs", msgs.PermGroupVerbs, msgs.PermGroupVerbsDesc, label, verbs),
		t.group(n.Name, "actions", msgs.PermGroupActions, msgs.PermGroupActionsDesc, label, actions),
		t.group(n.Name, "fields", msgs.PermGroupFields, msgs.PermGroupFieldsDesc, label, t.kindFields(fields)),
		t.group(n.Name, "inlines", msgs.PermGroupInlines, msgs.PermGroupInlinesDesc, label, t.kindFields(inlines)),
		t.group(n.Name, "sections", msgs.PermGroupSections, msgs.PermGroupSectionsDesc, label, sections),
		t.group(n.Name, "pages", msgs.PermGroupPages, msgs.PermGroupPagesDesc, label, pageItems),
		t.group(n.Name, "models", msgs.PermGroupModels, msgs.PermGroupModelsDesc, label, modelItems),
		t.group(n.Name, "checks", msgs.PermGroupChecks, msgs.PermGroupChecksDesc, label, checkItems))
}

// kindFields are the fields among nodes, one item a field — the model's and
// the record's of the same name together —, the verbs asked of each under
// it, then the fields nested in it the same way.
func (t *permTree) kindFields(nodes []*presets.PermNode) (out []*PermTreeItem) {
	var order []string
	same := map[string][]*presets.PermNode{}
	for _, n := range nodes {
		parts := strings.Split(strings.TrimSuffix(n.Name, ":"), ":")
		key := parts[len(parts)-1]
		if same[key] == nil {
			order = append(order, key)
		}
		same[key] = append(same[key], n)
	}
	for _, key := range order {
		ns := same[key]
		label := t.title(ns[0])
		it := &PermTreeItem{Title: label, Description: t.describe(ns[0], label), Value: ns[0].Name, Kind: kindOf(ns[0])}
		var fields, inlines []*presets.PermNode
		for _, n := range ns {
			it.Resources = append(it.Resources, t.resources(n, "")...)
			for _, a := range n.Actions {
				it.Children = append(it.Children, t.verbLeaf(n, a, label))
			}
			p := partsOf(n)
			fields, inlines = append(fields, p.fields...), append(inlines, p.inlines...)
		}
		it.Children = add(it.Children,
			t.group(it.Value, "fields", t.msgs.PermGroupFields, t.msgs.PermGroupFieldsDesc, label, t.kindFields(fields)),
			t.group(it.Value, "inlines", t.msgs.PermGroupInlines, t.msgs.PermGroupInlinesDesc, label, t.kindFields(inlines)))
		out = append(out, it)
	}
	return
}

// kindSection is the section s: its verbs, then the fields of the record
// (p) it holds, each with its verbs.
func (t *permTree) kindSection(s *presets.PermNode, p parts) *PermTreeItem {
	it := t.branch(s)
	var fields []*presets.PermNode
	seen := map[*presets.PermNode]bool{}
	for _, a := range s.Actions {
		it.Children = append(it.Children, t.verbLeaf(s, a, it.Title))
		for _, f := range t.sectionFields(p, a, "") {
			if !seen[f] {
				seen[f] = true
				fields = append(fields, f)
			}
		}
	}
	it.Children = add(it.Children, t.group(it.Value, "fields", t.msgs.PermGroupFields, t.msgs.PermGroupFieldsDesc, it.Title, t.kindFields(fields)))
	return it
}
