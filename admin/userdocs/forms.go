package userdocs

import (
	"strings"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

// The forms of a model documented, each a node of it (MODEL/forms/NAME): the
// form of a new record, the one of an edit, the detail.
const (
	FormNew    = "new"
	FormEdit   = "edit"
	FormDetail = "detail"
)

// formNodes are the nodes of the forms of mb there are, its node id.
func formNodes(mb *presets.ModelBuilder, id string, msgs *Messages) (nodes []*Node) {
	if !mb.GetSingleton() && !mb.CreatingDisabled() {
		nodes = append(nodes, &Node{ID: id + "/forms/" + FormNew, Title: msgs.FormNew, Icon: "mdi-plus-box-outline"})
	}
	if !mb.EditingDisabled() {
		nodes = append(nodes, &Node{ID: id + "/forms/" + FormEdit, Title: msgs.FormEdit, Icon: "mdi-pencil-box-outline"})
	}
	if mb.HasDetailing() {
		nodes = append(nodes, &Node{ID: id + "/forms/" + FormDetail, Title: msgs.FormDetail, Icon: "mdi-card-text-outline"})
	}
	return
}

// formFields are the fields of the form of mb: of a new record, of an edit,
// of the detail.
func formFields(mb *presets.ModelBuilder, form string) *presets.FieldsBuilder {
	switch form {
	case FormNew:
		return &mb.Editing().CreatingBuilder().FieldsBuilder
	case FormEdit:
		return &mb.Editing().FieldsBuilder
	case FormDetail:
		if mb.HasDetailing() {
			return &mb.Detailing().FieldsBuilder
		}
	}
	return nil
}

// fieldsTable is the table of the fields of form of mb, in the order shown:
// each its label and its hint, in the language of ctx; a field shown
// depending on the record marked; a field made of others, followed by them.
func fieldsTable(mb *presets.ModelBuilder, form string, ctx *web.EventContext) string {
	msgs := GetMessages(ctx.Context())
	fb := formFields(mb, form)
	if fb == nil {
		return msgs.NoFields
	}
	mode := presets.FieldModeStack{presets.EDIT}
	switch form {
	case FormNew:
		mode = presets.FieldModeStack{presets.NEW}
	case FormDetail:
		mode = presets.FieldModeStack{presets.DETAIL}
	}
	var rows []string
	var walk func(fb *presets.FieldsBuilder, info *presets.ModelInfo, obj any, prefix string)
	walk = func(fb *presets.FieldsBuilder, info *presets.ModelInfo, obj any, prefix string) {
		// those the form shows, of a new record
		for _, name := range fb.ShownFields(info, obj, mode, ctx) {
			f := fb.GetField(name)
			if f == nil {
				continue
			}
			label := f.ContextLabel(info, ctx.Context())
			if label == "" {
				label = presets.HumanizeString(name)
			}
			desc := f.ContextHint(info, ctx.Context())
			if f.Enabled() != nil {
				desc = strings.TrimSpace(desc + " *" + msgs.Conditional + "*")
			}
			rows = append(rows, "| "+cell(prefix+label)+" | "+cell(desc)+" |")
			if n := f.GetNested(); n != nil && n.FieldsBuilder() != nil {
				sub, subObj := info, any(nil)
				if m := n.Model(); m != nil {
					sub, subObj = m.Info(), m.NewModel()
				}
				walk(n.FieldsBuilder(), sub, subObj, prefix+label+" › ")
			}
		}
	}
	walk(fb, mb.Info(), mb.NewModel(), "")
	if len(rows) == 0 {
		return msgs.NoFields
	}
	return "| " + msgs.Field + " | " + msgs.Description + " |\n| --- | --- |\n" + strings.Join(rows, "\n") + "\n"
}

// cell is s in a cell of a table of markdown: one line, its bars escaped.
func cell(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	return strings.ReplaceAll(s, "|", `\|`)
}

// detailMenu is the menu of a record in the detail of mb, in markdown: the
// models nested in it, its actions — each available by a rule of its
// record with the rule, as its document says it (RuleHeading) —, its pages;
// each a link to its document.
func (r *renderer) detailMenu(mb *presets.ModelBuilder) string {
	b, ctx := r.b, r.ctx
	msgs := GetMessages(ctx.Context())
	locale := b.locale(ctx)
	id := b.modelNode(mb)
	link := func(title, node string) string {
		return "- [" + title + "](" + b.DocHref(ctx, node) + ")"
	}
	var out []string
	section := func(title string, items []string) {
		if len(items) > 0 {
			out = append(out, "**"+title+"**\n\n"+strings.Join(items, "\n"))
		}
	}

	var children []string
	for _, c := range mb.Children() {
		children = append(children, link(c.TTitleAuto(ctx.Context()), id+"/children/"+childKey(c)))
	}
	section(msgs.Children, children)

	var acts []string
	seen := map[string]bool{}
	for _, a := range detailingActions(mb) {
		seen[a.Name()] = true
		item := link(a.RequestTitle(mb, ctx.Context()), id+"/actions/"+a.Name())
		if a.GetEnabled() != nil || a.GetEnabledObj() != nil {
			rule, ok := b.rule(locale, id+"/actions/"+a.Name(), msgs)
			if !ok {
				rule = msgs.RuleUndocumented
			}
			item += " — " + rule
		}
		acts = append(acts, item)
	}
	if mb.HasDetailing() {
		for _, it := range mb.Detailing().RowMenu().Items() {
			if it.Child() == nil && !seen[it.Name()] && (it.Name() != "Delete" || !mb.DeletingDisabled()) {
				seen[it.Name()] = true
				acts = append(acts, link(it.TTitle(ctx.Context()), id+"/actions/"+it.Name()))
			}
		}
	}
	section(msgs.Actions, acts)

	var pages []string
	for _, p := range modelPages(mb) {
		pages = append(pages, link(p.TTitle(ctx.Context()), id+"/pages/"+pageKey(p.Path())))
	}
	section(msgs.Pages, pages)

	if len(out) == 0 {
		return msgs.NothingInMenu
	}
	return strings.Join(out, "\n\n") + "\n"
}

// rule is what the document of node says of when it is available: the text
// under its heading RuleHeading, in one paragraph — false when it says
// nothing.
func (b *Builder) rule(locale, node string, msgs *Messages) (string, bool) {
	_, _, content, err := b.find(locale, node)
	if err != nil || content == "" {
		return "", false
	}
	var lines []string
	in := false
	for _, l := range strings.Split(content, "\n") {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "#") {
			if in {
				break
			}
			in = strings.TrimSpace(strings.TrimLeft(t, "#")) == msgs.RuleHeading
			continue
		}
		if in && t != "" {
			lines = append(lines, t)
		}
	}
	if len(lines) == 0 {
		return "", false
	}
	return strings.Join(lines, " "), true
}

// modelOf is the model of node — a model's, or of a form, an action, a page
// of it —, nil when none is.
func (b *Builder) modelOf(node string) *presets.ModelBuilder {
	node = strings.TrimSuffix(node, "/permissions")
	for _, sep := range []string{"/forms/", "/actions/", "/pages/"} {
		if i := strings.LastIndex(node, sep); i >= 0 {
			node = node[:i]
			break
		}
	}
	var find func(models []*presets.ModelBuilder) *presets.ModelBuilder
	find = func(models []*presets.ModelBuilder) *presets.ModelBuilder {
		for _, mb := range models {
			if b.modelNode(mb) == node {
				return mb
			}
			// the nested ones, not always among the builder's
			if c := find(mb.Children()); c != nil {
				return c
			}
		}
		return nil
	}
	return find(b.p.Models())
}
