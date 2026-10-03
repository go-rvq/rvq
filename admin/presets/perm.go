package presets

import (
	"context"
	"sort"
	"strings"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/mpvl/unique"
)

// FieldPerm is the resource of a record's field: "#Title".
func FieldPerm(name string) string {
	return "#" + name
}

// SectionPerm is the resource of a section of a record's detail page:
// "$Main". A section is seen (get) and edited in place (create/update) by its
// own permission; the fields it writes are still asked theirs (FieldPerm).
func SectionPerm(name string) string {
	return "$" + name
}

// The permissions of the admin, as a tree (Builder.Permissions): what may be
// given to a role, each part by its resource — by the groups of the menu
// and, for a model or a page, by its unique name, which decides first — and
// the permissions and actions asked of it.
//
//	group           :site/:
//	  group         :site/:seo/:
//	    model       :site/:seo/:seo_config:      @list @create !bulk…
//	      record    :site/:seo/:seo_config:<*>:  @get @edit @delete !action…
//	        field   …:<*>:#Title:                       @get @edit @create
//	          field …:<*>:#Config:#Title:               (nested)
//	        section …:<*>:$Main:                        @get @edit
//	        page    …:<*>:/report:
//	        model   …:<*>:revisions:                    (nested, the same way)
//	      page      :site/:seo/:seo_config:/import:
//	  page          :site/:/report:
//
// A singleton has no listing: its node is its record.

// PermNodeKind is what a node of the permissions is: a group of the menu, a
// model's listing, …
type PermNodeKind string

const (
	PermNodeGroup     PermNodeKind = "group"     // a group of the menu
	PermNodeModel     PermNodeKind = "model"     // a model: its listing
	PermNodeSingleton PermNodeKind = "singleton" // a singleton: its record
	PermNodeRecord    PermNodeKind = "record"    // a record of a model, <*>
	PermNodeField     PermNodeKind = "field"     // a field of a record
	PermNodeSection   PermNodeKind = "section"   // a section of a detail
	PermNodePage      PermNodeKind = "page"      // a page
	PermNodeCheck     PermNodeKind = "check"     // a verifier of its own
)

// PermMenu is a group of the menu in the permissions, or a page of the admin
// (no Resources).
type PermMenu struct {
	Parent *PermMenu `yaml:"-" json:"-"`
	// Name is the resource by the groups, Unique the one by the unique name
	// (a page)
	Name   string                           `yaml:",omitempty" json:",omitempty"`
	Unique string                           `yaml:",omitempty" json:",omitempty"`
	Kind   PermNodeKind                     `yaml:",omitempty" json:",omitempty"`
	Title  func(ctx context.Context) string `yaml:"-" json:"-"`
	// Description is what it is, in the language of the request; nil or
	// "": none of its own
	Description func(ctx context.Context) string `yaml:"-" json:"-"`
	Resources   []*ModelPerm                     `yaml:",omitempty" json:",omitempty"`
	Children    []*PermMenu                      `yaml:",omitempty" json:",omitempty"`
}

func (m *PermMenu) AddChildren(children ...*PermMenu) {
	m.Children = append(m.Children, children...)
	for _, child := range children {
		child.Parent = m
	}
}

// Tree is the menu as a tree of nodes.
func (m *PermMenu) Tree() (n *PermNode) {
	n = &PermNode{Name: m.Name, Unique: m.Unique, Kind: m.Kind, Title: m.Title, Description: m.Description}
	for _, res := range m.Resources {
		n.AddChildren(res.Tree())
	}
	for _, c := range m.Children {
		n.AddChildren(c.Tree())
	}
	sortNodes(n.Children)
	return
}

// ModelPerm is a model in the permissions: its node — its listing's, its
// record's for a singleton — with everything under it.
type ModelPerm struct {
	Model *ModelBuilder                    `yaml:"-" json:"-"`
	Name  string                           `yaml:",omitempty" json:",omitempty"`
	Title func(ctx context.Context) string `yaml:"-" json:"-"`
	node  *PermNode
}

// Tree is the node of the model.
func (m *ModelPerm) Tree() *PermNode { return m.node }

// PermNodeAction is a permission ("@edit") or an action ("!publish") asked of
// a resource; Fields, the fields a form of it holds.
type PermNodeAction struct {
	Name        string                           `yaml:",omitempty" json:",omitempty"`
	Fields      []string                         `yaml:",omitempty" json:",omitempty"`
	Title       func(ctx context.Context) string `yaml:"-" json:"-"`
	Description func(ctx context.Context) string `yaml:"-" json:"-"`
}

// PermNode is a resource of the permissions: by the groups (Name) and by the
// unique name (Unique, when it has one), what is asked of it, and the
// resources under it.
type PermNode struct {
	Parent *PermNode                        `yaml:"-" json:"-"`
	Name   string                           `yaml:",omitempty" json:",omitempty"`
	Unique string                           `yaml:",omitempty" json:",omitempty"`
	Kind   PermNodeKind                     `yaml:",omitempty" json:",omitempty"`
	Title  func(ctx context.Context) string `yaml:"-" json:"-"`
	// Description is what it is, in the language of the request; nil or
	// "": none of its own
	Description func(ctx context.Context) string `yaml:"-" json:"-"`
	Actions     []*PermNodeAction                `yaml:",omitempty" json:",omitempty"`
	Children    []*PermNode                      `yaml:",omitempty" json:",omitempty"`
}

func (n *PermNode) AddChildren(children ...*PermNode) {
	n.Children = append(n.Children, children...)
	for _, child := range children {
		child.Parent = n
	}
}

func (n *PermNode) Walk(f func(parents []*PermNode, node *PermNode)) {
	n.walk(nil, f)
}

func (n *PermNode) walk(parents []*PermNode, f func(parents []*PermNode, node *PermNode)) {
	for _, child := range n.Children {
		f(parents, child)
		child.walk(append(parents, n), f)
	}
}

// PermZipEntry is a node of the tree in a list: its resources and what is
// asked of them.
type PermZipEntry struct {
	Resource string            `yaml:",omitempty" json:",omitempty"`
	Unique   string            `yaml:",omitempty" json:",omitempty"`
	Actions  []*PermNodeAction `yaml:",omitempty" json:",omitempty"`
}

// Zip is the tree as a list, every node under n, in order.
func (n *PermNode) Zip() (enties []*PermZipEntry) {
	n.Walk(func(parents []*PermNode, node *PermNode) {
		enties = append(enties, &PermZipEntry{Resource: node.Name, Unique: node.Unique, Actions: node.Actions})
	})
	return
}

func sortNodes(nodes []*PermNode) {
	sort.SliceStable(nodes, func(i, j int) bool { return nodes[i].Name < nodes[j].Name })
}

func sortActions(actions []*PermNodeAction) {
	sort.SliceStable(actions, func(i, j int) bool { return actions[i].Name < actions[j].Name })
}

// BuildPermissions is the tree of the permissions of the admin: the menu,
// group by group, recursively — each model with its listing, its record, the
// fields (nested ones too), the sections, the actions, the pages and the
// models nested in it —, the pages of the admin, and, at the root, what has
// permissions out of the menu.
func (b *Builder) BuildPermissions() (rootMenu *PermMenu) {
	rootMenu = &PermMenu{}
	seen := map[string]bool{}

	var group func(g *MenuGroupBuilder, into *PermMenu)
	group = func(g *MenuGroupBuilder, into *PermMenu) {
		for _, it := range g.Items() {
			switch v := it.Value.(type) {
			case *MenuGroupBuilder:
				parts := []string{PermModule}
				for _, name := range v.PathNames() {
					parts = append(parts, GroupPermPart(name))
				}
				sub := &PermMenu{Name: strings.Join(parts, ":") + ":", Kind: PermNodeGroup, Title: v.TTitle, Description: v.TDescription}
				group(v, sub)
				into.AddChildren(sub)
			case *ModelBuilder:
				if !v.IsInMenu() {
					continue
				}
				mp := b.modelPerm(v)
				seen[mp.Name] = true
				into.Resources = append(into.Resources, mp)
			case *HttpPageBuilder:
				if p := b.pagePerm(v); p != nil {
					seen[p.Name] = true
					into.AddChildren(p)
				}
			}
		}
		sortMenu(into)
	}
	group(b.MenuTree(), rootMenu)

	// what has permissions out of the menu: the models (not nested), the
	// pages and the verifiers of the builder
	for _, mb := range b.models {
		if mb.parent != nil || mb.IsInMenu() {
			continue
		}
		if mp := b.modelPerm(mb); !seen[mp.Name] {
			seen[mp.Name] = true
			rootMenu.Resources = append(rootMenu.Resources, mp)
		}
	}
	// the builder's verifiers and its pages' — not added to the builder:
	// the tree is made again on every call (the menu may have changed)
	verifiers := append(perm.PermVerifiers{}, b.verifiers...)
	for _, page := range b.pagesRegistrator.httpPages {
		if page.verififer != nil {
			verifiers = append(verifiers, page.GetVerifier())
		}
	}
	for _, verifier := range verifiers {
		v := verifier.Build(b.verifier.Spawn())
		if !seen[v.Resource()] {
			seen[v.Resource()] = true
			rootMenu.AddChildren(&PermMenu{Name: v.Resource(), Unique: v.PreferredResource(), Kind: PermNodeCheck, Title: verifier.GetTitle(),
				Description: verifier.GetDescription()})
		}
	}
	sortMenu(rootMenu)
	return
}

func sortMenu(m *PermMenu) {
	sort.SliceStable(m.Resources, func(i, j int) bool { return m.Resources[i].Name < m.Resources[j].Name })
	sort.SliceStable(m.Children, func(i, j int) bool { return m.Children[i].Name < m.Children[j].Name })
}

// AppendListPermActions adds, to the permissions of the model's listing,
// actions of its own no builder of it knows — the kinds of jobs of a worker,
// "!upload_posts" —: f is asked when the permissions are made.
func (mb *ModelBuilder) AppendListPermActions(f func() []*PermNodeAction) *ModelBuilder {
	mb.listPermActions = append(mb.listPermActions, f)
	return mb
}

// pagePerm is a page of the admin in the permissions, nil when it has no
// permission of its own.
func (b *Builder) pagePerm(page *HttpPageBuilder) *PermMenu {
	if page.verififer == nil {
		return nil
	}
	v := page.GetVerifier().Build(b.verifier.Spawn())
	return &PermMenu{Name: v.Resource(), Unique: v.PreferredResource(), Kind: PermNodePage, Title: page.TTitle, Description: page.TDescription}
}

// anyID is the id of any record: "<*>".
func anyID() ID {
	return ID{Fields: []model.Field{model.SingleField("ID")}, Values: []any{"*"}}
}

// modelPerm is mb in the permissions.
func (b *Builder) modelPerm(mb *ModelBuilder) *ModelPerm {
	n := b.modelPermNode(mb)
	return &ModelPerm{Model: mb, Name: n.Name, Title: n.Title, node: n}
}

// modelPermNode is the node of mb: its listing — a singleton's, its record —
// with everything under it; a nested model under any record of its parents.
func (b *Builder) modelPermNode(mb *ModelBuilder) *PermNode {
	var parents []ID
	for p := mb.parent; p != nil; p = p.parent {
		if !p.singleton {
			parents = append(parents, anyID())
		}
	}
	lv := mb.Permissioner().ListVerifier(parents...)
	n := &PermNode{Name: lv.Resource(), Unique: lv.PreferredResource(), Kind: PermNodeModel,
		Title: func(ctx context.Context) string { return mb.TTitleAuto(ctx) }, Description: mb.TDescription}
	if mb.singleton {
		n.Kind = PermNodeSingleton
	}

	// the verifiers of its own (pages, custom checks)
	for node := range perm.WalkPermVerififierBuilders(mb.AllVerifiers()) {
		v := node.Elem.Build(mb.Permissioner().ListVerifier(parents...))
		n.AddChildren(&PermNode{Name: v.Resource(), Unique: v.PreferredResource(), Kind: PermNodeCheck, Title: node.Elem.GetTitle(),
			Description: node.Elem.GetDescription()})
	}

	record := n
	if !mb.singleton {
		n.Actions = append(n.Actions, &PermNodeAction{Name: PermList, Fields: permFields(&mb.listing.FieldsBuilder)})
		if !mb.creatingDisabled {
			n.Actions = append(n.Actions, &PermNodeAction{Name: PermCreate,
				Fields: permFields(&mb.editing.CreatingBuilder().FieldsBuilder)})
		}
		for _, a := range mb.listing.bulkActions {
			n.Actions = append(n.Actions, &PermNodeAction{Name: ActionPerm(a.name), Title: a.RequestTitle, Description: a.RequestDescription})
		}
		for _, a := range mb.listing.actions {
			action := a
			n.Actions = append(n.Actions, &PermNodeAction{Name: action.PermName(),
				Title:       func(ctx context.Context) string { return action.RequestTitle(mb, ctx) },
				Description: func(ctx context.Context) string { return action.RequestDescription(mb, ctx) }})
		}
		for _, f := range mb.listPermActions {
			n.Actions = append(n.Actions, f()...)
		}
		// the pages of its listing
		if mb.listing.pagesRegistrator != nil {
			for _, p := range mb.listing.pagesRegistrator.HttpPages() {
				n.AddChildren(&PermNode{Name: n.Name + p.path + ":", Unique: suffixed(n.Unique, p.path+":"), Kind: PermNodePage, Title: p.TTitle, Description: p.TDescription})
			}
		}

		rv := mb.Permissioner().Verifier(anyID(), parents...)
		record = &PermNode{Name: rv.Resource(), Unique: rv.PreferredResource(), Kind: PermNodeRecord,
			Title: func(ctx context.Context) string { return mb.TTitle(ctx) }}
		n.AddChildren(record)
	}

	record.Actions = append(record.Actions, &PermNodeAction{Name: PermGet, Fields: permFields(&mb.detailing.FieldsBuilder)})
	if !mb.editingDisabled {
		record.Actions = append(record.Actions, &PermNodeAction{Name: PermUpdate, Fields: permFields(&mb.editing.FieldsBuilder)})
	}
	if !mb.singleton && !mb.deletingDisabled {
		record.Actions = append(record.Actions,
			&PermNodeAction{Name: PermDelete},
			&PermNodeAction{Name: PermDeleteWithRelated})
	}
	actions := append(append([]*ActionBuilder{}, mb.detailing.actions...), mb.listing.itemActions...)
	for _, a := range actions {
		action := a
		record.Actions = append(record.Actions, &PermNodeAction{Name: action.PermName(),
			Title:       func(ctx context.Context) string { return action.RequestTitle(mb, ctx) },
			Description: func(ctx context.Context) string { return action.RequestDescription(mb, ctx) }})
	}
	for _, verifier := range mb.detailing.verifiers {
		v := verifier.Build(mb.Permissioner().Verifier(anyID(), parents...))
		record.AddChildren(&PermNode{Name: v.Resource(), Unique: v.PreferredResource(), Kind: PermNodeCheck, Title: verifier.GetTitle(),
			Description: verifier.GetDescription()})
	}

	// its fields — of the detail, of the edit, of the new record —, nested
	// ones under theirs; its sections; its pages; the models nested in it
	record.AddChildren(fieldPermNodes(mb, record)...)
	for _, s := range mb.detailing.GetSections() {
		section := s
		part := SectionPerm(section.name) + ":"
		record.AddChildren(&PermNode{Name: record.Name + part, Unique: suffixed(record.Unique, part), Kind: PermNodeSection,
			// its label, in the language of the request
			Title: func(ctx context.Context) string {
				label := section.label
				if label == "" {
					label = section.name
				}
				if t := i18n.Translate(mb.FieldTranslator(), ctx, label); t != "" {
					return t
				}
				return HumanizeString(section.name)
			},
			Description: func(ctx context.Context) string { return section.TDescription(mb, ctx) },
			Actions:     []*PermNodeAction{{Name: PermUpdate}, {Name: PermGet}}})
	}
	if mb.detailing.pagesRegistrator != nil {
		for _, p := range mb.detailing.pagesRegistrator.HttpPages() {
			record.AddChildren(&PermNode{Name: record.Name + p.path + ":", Unique: suffixed(record.Unique, p.path+":"), Kind: PermNodePage, Title: p.TTitle, Description: p.TDescription})
		}
	}
	for _, child := range mb.children {
		record.AddChildren(b.modelPermNode(child))
	}

	sortActions(n.Actions)
	if record != n {
		sortActions(record.Actions)
	}
	sortNodes(n.Children)
	sortNodes(record.Children)
	return n
}

// suffixed is s with part after it, "" for no s.
func suffixed(s, part string) string {
	if s == "" {
		return ""
	}
	return s + part
}

// permFields are the names of the fields of fb, sorted.
func permFields(fb *FieldsBuilder) (names []string) {
	names = append(names, fb.CurrentLayout().Names()...)
	unique.Strings(&names)
	return
}

// fieldPermNodes are the nodes of the fields of mb under record: each with
// the permissions its forms ask — of the detail (@get), of the edit (@edit),
// of the new record (@create) —, its nested fields under it.
func fieldPermNodes(mb *ModelBuilder, record *PermNode) []*PermNode {
	type form struct {
		fb   *FieldsBuilder
		verb string
	}
	forms := []form{{&mb.detailing.FieldsBuilder, PermGet}}
	if !mb.editingDisabled {
		forms = append(forms, form{&mb.editing.FieldsBuilder, PermUpdate})
	}
	if !mb.singleton && !mb.creatingDisabled {
		forms = append(forms, form{&mb.editing.CreatingBuilder().FieldsBuilder, PermCreate})
	}

	root := &PermNode{Name: record.Name, Unique: record.Unique}
	var add func(fb *FieldsBuilder, info *ModelInfo, verb string, under *PermNode, depth int)
	add = func(fb *FieldsBuilder, info *ModelInfo, verb string, under *PermNode, depth int) {
		if fb == nil || depth > 8 {
			return
		}
		for _, name := range fb.CurrentLayout().Names() {
			f := fb.GetField(name)
			if f == nil {
				continue
			}
			part := FieldPerm(name) + ":"
			var node *PermNode
			for _, c := range under.Children {
				if c.Name == under.Name+part {
					node = c
				}
			}
			if node == nil {
				field, fieldInfo := f, info
				node = &PermNode{Name: under.Name + part, Unique: suffixed(under.Unique, part), Kind: PermNodeField,
					// its label, in the language of the request
					Title:       func(ctx context.Context) string { return field.ContextLabel(fieldInfo, ctx) },
					Description: func(ctx context.Context) string { return field.ContextDescription(fieldInfo, ctx) }}
				under.AddChildren(node)
			}
			if !containsAction(node.Actions, verb) {
				node.Actions = append(node.Actions, &PermNodeAction{Name: verb})
			}
			if nested := f.GetNested(); nested != nil && nested.FieldsBuilder() != nil {
				sub := info
				if m := nested.Model(); m != nil {
					sub = m.Info()
				}
				add(nested.FieldsBuilder(), sub, verb, node, depth+1)
			}
		}
	}
	for _, fm := range forms {
		add(fm.fb, mb.Info(), fm.verb, root, 0)
	}
	var tidy func(nodes []*PermNode)
	tidy = func(nodes []*PermNode) {
		sortNodes(nodes)
		for _, n := range nodes {
			sortActions(n.Actions)
			tidy(n.Children)
		}
	}
	tidy(root.Children)
	return root.Children
}

func containsAction(actions []*PermNodeAction, name string) bool {
	for _, a := range actions {
		if a.Name == name {
			return true
		}
	}
	return false
}
