package userdocs

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"strings"

	"github.com/go-rvq/rvq/admin/presets"
)

// The permissions of a part of the admin, documented: a node of each model
// (MODEL/permissions), group (groups/GROUP/permissions) and page of the
// admin (pages/PAGE/permissions). The resources are read from the menu as it
// is — the chain of the groups changes with it —, never written in a
// document: admin.permissions() makes them.

// permissionsNode is the node of the permissions of the part of node id.
func permissionsNode(id string, msgs *Messages) *Node {
	return &Node{ID: id + "/permissions", Title: msgs.Permissions, Icon: "mdi-shield-key-outline"}
}

// permissionsKind is the kind of the part of the permissions node node — a
// group, a page of the admin, a model —, and the part's node.
func permissionsKind(node string) (kind, part string) {
	part = strings.TrimSuffix(node, "/permissions")
	switch {
	case strings.HasPrefix(part, "groups/"):
		return "group", part
	case strings.HasPrefix(part, "pages/") && !strings.Contains(strings.TrimPrefix(part, "pages/"), "/"):
		return "page", part
	}
	return "model", part
}

// code is s as markdown code, a cell's bars escaped.
func code(s string) string { return "`" + strings.ReplaceAll(s, "|", `\|`) + "`" }

// table is a markdown table of the head and rows.
func table(head []string, rows [][]string) string {
	var b strings.Builder
	b.WriteString("| " + strings.Join(head, " | ") + " |\n|")
	for range head {
		b.WriteString(" --- |")
	}
	b.WriteString("\n")
	for _, r := range rows {
		for i := range r {
			r[i] = cell(r[i])
		}
		b.WriteString("| " + strings.Join(r, " | ") + " |\n")
	}
	return b.String()
}

// permissions are the permissions of the part of the node of the document
// rendered, in markdown.
func (r *renderer) permissions() string {
	kind, part := permissionsKind(r.node)
	switch kind {
	case "group":
		return r.groupPermissions(strings.TrimPrefix(part, "groups/"))
	case "page":
		return r.pagePermissions("/" + strings.TrimPrefix(part, "pages/"))
	}
	if mb := r.b.modelOf(part); mb != nil {
		return r.modelPermissions(mb)
	}
	return GetMessages(r.ctx.Context()).PermNone
}

// resources is the table of the resources of what verifier v is of: by its
// unique name (first), by its groups.
func (r *renderer) resources(v interface {
	Resource() string
	PreferredResource() string
}) string {
	msgs := GetMessages(r.ctx.Context())
	var rows [][]string
	if u := v.PreferredResource(); u != "" {
		rows = append(rows, []string{msgs.PermUnique, code(u + "*")})
	}
	rows = append(rows, []string{msgs.PermByGroups, code(v.Resource() + "*")})
	return table([]string{msgs.PermBy, msgs.PermResource}, rows)
}

// modelPermissions are the permissions of mb: its resources, the verbs of
// its records, its fields and sections, its actions, its pages.
func (r *renderer) modelPermissions(mb *presets.ModelBuilder) string {
	ctx := r.ctx
	msgs := GetMessages(ctx.Context())
	lv := mb.Permissioner().ListVerifier()
	base := lv.PreferredResource()
	if base == "" {
		base = lv.Resource()
	}
	record := base + "<id>:"
	if mb.GetSingleton() {
		record = base
	}
	var out []string
	out = append(out, "### "+msgs.PermResources+"\n\n"+r.resources(lv))

	verbs := [][]string{
		{msgs.PermList, code(presets.PermList), code(base)},
		{msgs.PermGet, code(presets.PermGet), code(record)},
	}
	if !mb.GetSingleton() && !mb.CreatingDisabled() {
		verbs = append(verbs, []string{msgs.PermCreate, code(presets.PermCreate), code(base)})
	}
	if !mb.EditingDisabled() {
		verbs = append(verbs, []string{msgs.PermUpdate, code(presets.PermUpdate), code(record)})
	}
	if !mb.DeletingDisabled() && !mb.GetSingleton() {
		verbs = append(verbs, []string{msgs.PermDelete, code(presets.PermDelete), code(record)})
	}
	out = append(out, "### "+msgs.PermRecords+"\n\n"+table([]string{msgs.PermWhat, msgs.PermVerb, msgs.PermResource}, verbs))

	// the fields: of the edit (written) and of the detail (read)
	var fields [][]string
	seen := map[string]bool{}
	add := func(fb *presets.FieldsBuilder, verb string) {
		if fb == nil {
			return
		}
		for _, name := range fb.CurrentLayout().Names() {
			f := fb.GetField(name)
			if f == nil || seen[name] {
				continue
			}
			seen[name] = true
			label := f.ContextLabel(mb.Info(), ctx.Context())
			if label == "" {
				label = name
			}
			fields = append(fields, []string{label, code(record + presets.FieldPerm(name) + ":"), verb})
		}
	}
	if mb.HasDetailing() {
		add(&mb.Detailing().FieldsBuilder, code(presets.PermGet))
	}
	if !mb.EditingDisabled() {
		add(&mb.Editing().FieldsBuilder, code(presets.PermGet)+", "+code(presets.PermUpdate)+", "+code(presets.PermCreate))
	}
	if len(fields) > 0 {
		out = append(out, "### "+msgs.PermFields+"\n\n"+msgs.PermFieldsHint+"\n\n"+
			table([]string{msgs.Field, msgs.PermResource, msgs.PermVerb}, fields))
	}

	if mb.HasDetailing() {
		var sections [][]string
		for _, s := range mb.Detailing().GetSections() {
			label := s.NameLabel.Label()
			if label == "" {
				label = presets.HumanizeString(s.Name())
			}
			sections = append(sections, []string{label, code(record + presets.SectionPerm(s.Name()) + ":"),
				code(presets.PermGet) + ", " + code(presets.PermUpdate)})
		}
		if len(sections) > 0 {
			out = append(out, "### "+msgs.PermSections+"\n\n"+msgs.PermSectionsHint+"\n\n"+
				table([]string{msgs.PermSection, msgs.PermResource, msgs.PermVerb}, sections))
		}
	}

	// the actions: of a record, of the listing, in bulk
	var acts [][]string
	for _, a := range detailingActions(mb) {
		acts = append(acts, []string{a.RequestTitle(mb, ctx.Context()), code(record), code(a.PermName())})
	}
	for _, a := range mb.Listing().GetBulkActions() {
		acts = append(acts, []string{a.RequestTitle(ctx.Context()), code(base), code("bulk:" + a.Name())})
	}
	if len(acts) > 0 {
		out = append(out, "### "+msgs.Actions+"\n\n"+table([]string{msgs.PermWhat, msgs.PermResource, msgs.PermVerb}, acts))
	}

	var pages [][]string
	for _, p := range modelPages(mb) {
		pages = append(pages, []string{p.TTitle(ctx.Context()), code(record + p.Path() + ":")})
	}
	if len(pages) > 0 {
		out = append(out, "### "+msgs.Pages+"\n\n"+table([]string{msgs.PermWhat, msgs.PermResource}, pages))
	}
	return strings.Join(out, "\n\n") + "\n"
}

// groupPermissions are the permissions of the group name: its resource,
// which holds every part under it, and those parts by their unique names.
func (r *renderer) groupPermissions(name string) string {
	ctx := r.ctx
	msgs := GetMessages(ctx.Context())
	g := findGroup(r.b.p.MenuTree(), name)
	if g == nil {
		return msgs.PermNone
	}
	res := presets.PermModule + ":" + strings.Join(g.PathNames(), ":") + ":*"
	out := []string{"### " + msgs.PermResources + "\n\n" + table([]string{msgs.PermBy, msgs.PermResource},
		[][]string{{msgs.PermByGroups, code(res)}})}

	var parts [][]string
	var walk func(g *presets.MenuGroupBuilder)
	walk = func(g *presets.MenuGroupBuilder) {
		for _, it := range g.Items() {
			switch v := it.Value.(type) {
			case *presets.ModelBuilder:
				if v.IsInMenu() {
					parts = append(parts, []string{v.TTitleAuto(ctx.Context()), code(presets.PermModule + ":" + v.UniquePermName() + ":*")})
				}
			case *presets.HttpPageBuilder:
				parts = append(parts, []string{v.TTitle(ctx.Context()), code(presets.PermModule + ":" + v.UniquePermName() + ":*")})
			case *presets.MenuGroupBuilder:
				walk(v)
			}
		}
	}
	walk(g)
	if len(parts) > 0 {
		out = append(out, "### "+msgs.PermParts+"\n\n"+msgs.PermPartsHint+"\n\n"+
			table([]string{msgs.PermWhat, msgs.PermUnique}, parts))
	}
	return strings.Join(out, "\n\n") + "\n"
}

// pagePermissions are the permissions of the page of the admin of path.
func (r *renderer) pagePermissions(path string) string {
	msgs := GetMessages(r.ctx.Context())
	page := r.b.p.PagesRegistrator().GetHttpPage(path)
	if page == nil || page.GetVerifier() == nil {
		return msgs.PermNone
	}
	return "### " + msgs.PermResources + "\n\n" + r.resources(page.Verifier(r.ctx.R))
}

// findGroup is the group name under g, nil when there is none.
func findGroup(g *presets.MenuGroupBuilder, name string) *presets.MenuGroupBuilder {
	for _, it := range g.Items() {
		if sub, ok := it.Value.(*presets.MenuGroupBuilder); ok {
			if sub.Name() == name {
				return sub
			}
			if f := findGroup(sub, name); f != nil {
				return f
			}
		}
	}
	return nil
}

// permTreeItem is a node of the tree of the permissions as the VTreeview
// shows it: its title, its resources and what is asked of them.
type permTreeItem struct {
	Title    string          `json:"title"`
	Subtitle string          `json:"subtitle,omitempty"`
	Value    string          `json:"value"`
	Children []*permTreeItem `json:"children,omitempty"`
}

// permissionsTree is the tree of the permissions of the admin
// (presets.Builder.Permissions) as a VTreeview: each part by its title, its
// resource by the groups and by the unique name, and the permissions and
// actions asked of it — made from the menu as it is.
func (r *renderer) permissionsTree() string {
	ctx := r.ctx.Context()
	msgs := GetMessages(ctx)
	title := func(n *presets.PermNode) string {
		if n == nil {
			return ""
		}
		if n.Title != nil {
			if t := n.Title(ctx); t != "" {
				return t
			}
		}
		parts := strings.Split(strings.TrimSuffix(n.Name, ":"), ":")
		return fmt.Sprintf(msgs.PermUntitled, presets.HumanizeString(strings.TrimLeft(parts[len(parts)-1], "#$/<")))
	}
	// what the node is, in the language of the request
	describe := func(n *presets.PermNode, label string) string {
		// its own, when it has one
		if n.Description != nil {
			if d := n.Description(ctx); d != "" {
				return d
			}
		}
		parent := title(n.Parent)
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
	// resources is the subtitle part of the resources of n: by the groups,
	// by the unique name
	resources := func(n *presets.PermNode, suffix string) []string {
		r := []string{n.Name + suffix}
		if n.Unique != "" && n.Unique != n.Name {
			r = append(r, n.Unique+suffix)
		}
		return r
	}
	var items func(nodes []*presets.PermNode) []*permTreeItem
	var item func(n *presets.PermNode) *permTreeItem
	// grouped is what is inside a model — its record's too —, in groups:
	// its permissions, actions, fields, sections, pages, nested models and
	// permissions of its own
	grouped := func(n *presets.PermNode, label string) (out []*permTreeItem) {
		type owned struct {
			node *presets.PermNode
			a    *presets.PermNodeAction
		}
		var verbs, actions []owned
		var fields, inlines, sections, pages, models, checks []*presets.PermNode
		var collect func(n *presets.PermNode)
		collect = func(n *presets.PermNode) {
			for _, a := range n.Actions {
				if strings.HasPrefix(a.Name, "@") {
					verbs = append(verbs, owned{n, a})
				} else {
					actions = append(actions, owned{n, a})
				}
			}
			for _, c := range n.Children {
				switch c.Kind {
				case presets.PermNodeRecord:
					collect(c) // the record is inside the model
				case presets.PermNodeField:
					fields = append(fields, c)
				case presets.PermNodeInline:
					inlines = append(inlines, c)
				case presets.PermNodeSection:
					sections = append(sections, c)
				case presets.PermNodePage:
					pages = append(pages, c)
				case presets.PermNodeModel, presets.PermNodeSingleton:
					models = append(models, c)
				default:
					checks = append(checks, c)
				}
			}
		}
		collect(n)
		group := func(key, title, desc string, children []*permTreeItem) {
			if len(children) == 0 {
				return
			}
			out = append(out, &permTreeItem{Title: fmt.Sprintf("%s (%d)", title, len(children)),
				Subtitle: fmt.Sprintf(desc, label), Value: n.Name + "(" + key + ")", Children: children})
		}
		var vi, ai []*permTreeItem
		for _, o := range verbs {
			vi = append(vi, &permTreeItem{Title: verbLabel(msgs, o.a, ctx),
				Subtitle: strings.Join(resources(o.node, o.a.Name), " · "), Value: o.node.Name + o.a.Name})
		}
		for _, o := range actions {
			t := verbLabel(msgs, o.a, ctx)
			d := ""
			if o.a.Description != nil {
				d = o.a.Description(ctx)
			}
			if d == "" {
				d = fmt.Sprintf(msgs.PermDescAction, t, label)
			}
			ai = append(ai, &permTreeItem{Title: t,
				Subtitle: strings.Join(append([]string{d}, resources(o.node, o.a.Name)...), " · "), Value: o.node.Name + o.a.Name})
		}
		group("verbs", msgs.PermGroupVerbs, msgs.PermGroupVerbsDesc, vi)
		group("actions", msgs.PermGroupActions, msgs.PermGroupActionsDesc, ai)
		group("fields", msgs.PermGroupFields, msgs.PermGroupFieldsDesc, items(fields))
		group("inlines", msgs.PermGroupInlines, msgs.PermGroupInlinesDesc, items(inlines))
		group("sections", msgs.PermGroupSections, msgs.PermGroupSectionsDesc, items(sections))
		group("pages", msgs.PermGroupPages, msgs.PermGroupPagesDesc, items(pages))
		group("models", msgs.PermGroupModels, msgs.PermGroupModelsDesc, items(models))
		group("checks", msgs.PermGroupChecks, msgs.PermGroupChecksDesc, items(checks))
		return
	}
	item = func(n *presets.PermNode) *permTreeItem {
		label := title(n)
		sub := append([]string{describe(n, label)}, resources(n, "")...)
		it := &permTreeItem{Title: label, Value: n.Name}
		switch n.Kind {
		case presets.PermNodeModel, presets.PermNodeSingleton, presets.PermNodeInline:
			// a model edited in place too: its structure inside the field
			it.Children = grouped(n, label)
		default:
			// a field, a section: what is asked of it, beside it
			var verbs []string
			for _, a := range n.Actions {
				verbs = append(verbs, verbLabel(msgs, a, ctx)+" ("+a.Name+")")
			}
			if len(verbs) > 0 {
				sub = append(sub, strings.Join(verbs, " "))
			}
			it.Children = items(n.Children)
		}
		it.Subtitle = strings.Join(sub, " · ")
		return it
	}
	items = func(nodes []*presets.PermNode) (out []*permTreeItem) {
		for _, n := range nodes {
			out = append(out, item(n))
		}
		return
	}
	// the scopes first, by their labels: the admin's — its resources begin
	// with "admin:"
	scope := presets.PermModule
	tree := []*permTreeItem{{Title: msgs.PermScopeAdmin, Subtitle: msgs.PermScopeAdminDesc + " · " + scope + ":*", Value: scope + ":",
		// made now, as the dump: from the menu as it is
		Children: items(r.b.p.BuildPermissions().Tree().Children)}}
	data, err := json.Marshal(tree)
	if err != nil {
		return ""
	}
	// a block of HTML of its own in the markdown: the page's Vue makes it
	return "\n<v-treeview :items='" + html.EscapeString(string(data)) + "' item-props density=\"compact\" open-on-click class=\"user-docs-perm-tree\"></v-treeview>\n"
}

// verbLabel is what a permission or an action asked of a resource is, in the
// language of ctx: the presets' permissions by their messages, an action by
// its title.
func verbLabel(msgs *Messages, a *presets.PermNodeAction, ctx context.Context) string {
	switch a.Name {
	case presets.PermList:
		return msgs.PermList
	case presets.PermGet:
		return msgs.PermGet
	case presets.PermCreate:
		return msgs.PermCreate
	case presets.PermUpdate:
		return msgs.PermUpdate
	case presets.PermDelete:
		return msgs.PermDelete
	case presets.PermDeleteWithRelated:
		return msgs.PermDeleteWithRelated
	}
	if a.Title != nil {
		if t := a.Title(ctx); t != "" {
			return t
		}
	}
	return presets.HumanizeString(strings.TrimLeft(a.Name, "@!"))
}
