package presets

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// ListEditorItemAction contributes an extra control to an item's "Actions" cell
// (rendered before the built-in delete control).
type ListEditorItemAction func(ctx *ListEditorItemContext) h.HTMLComponent

// ListEditorTableBuilder is a ListEditorComponentBuilder that renders the list
// editor as a table: one column per (renderable) field — with grouped/nested
// column headers when the fields layout has sections — and, in edit mode, a
// trailing "Actions" column holding the delete/revert control plus any
// user-added actions.
//
// Item and DeletedItem each return a <tr>; Container wraps the rows in a
// <table> with the computed header. Field cells render with FieldContext.
// MustInput set, so only the bare input shows (the label lives in the header).
//
// Use it in place of the default list builder:
//
//	NewListEditor(field).
//	    ComponentBuilder(NewListEditorTableBuilder().
//	        Action(func(ctx *ListEditorItemContext) h.HTMLComponent { ... }))
//
// ListEditorTableLayoutFunc returns the fields layout (which drives the columns
// and their nested/grouped headers) for a given mode. Returning nil falls back
// to the nested fields' current layout.
type ListEditorTableLayoutFunc func(fb *FieldsBuilder) FieldsLayout

type ListEditorTableBuilder struct {
	actions      []ListEditorItemAction
	actionsLabel string
	density      string
	// viewLayoutFunc / editLayoutFunc / newLayoutFunc override the columns per
	// mode: view is used when the editor renders read-only (detail), new when the
	// parent record is being created, and edit otherwise. Each falls back to the
	// nested fields' current layout when nil.
	viewLayoutFunc ListEditorTableLayoutFunc
	editLayoutFunc ListEditorTableLayoutFunc
	newLayoutFunc  ListEditorTableLayoutFunc
}

func NewListEditorTableBuilder() *ListEditorTableBuilder {
	return &ListEditorTableBuilder{density: DensityCompact}
}

// ViewLayoutFunc sets the columns layout used in read-only (view) mode.
func (b *ListEditorTableBuilder) ViewLayoutFunc(f ListEditorTableLayoutFunc) *ListEditorTableBuilder {
	b.viewLayoutFunc = f
	return b
}

// EditLayoutFunc sets the columns layout used when editing an existing record.
func (b *ListEditorTableBuilder) EditLayoutFunc(f ListEditorTableLayoutFunc) *ListEditorTableBuilder {
	b.editLayoutFunc = f
	return b
}

// NewLayoutFunc sets the columns layout used when the parent record is being
// created (new).
func (b *ListEditorTableBuilder) NewLayoutFunc(f ListEditorTableLayoutFunc) *ListEditorTableBuilder {
	b.newLayoutFunc = f
	return b
}

func (b *ListEditorTableBuilder) LayoutFunc(f ListEditorTableLayoutFunc) *ListEditorTableBuilder {
	b.newLayoutFunc = f
	b.editLayoutFunc = f
	b.viewLayoutFunc = f
	return b
}

// Action appends extra per-item actions rendered in the actions column, before
// the built-in delete control.
func (b *ListEditorTableBuilder) Action(actions ...ListEditorItemAction) *ListEditorTableBuilder {
	b.actions = append(b.actions, actions...)
	return b
}

// ActionsLabel overrides the header label of the actions column.
func (b *ListEditorTableBuilder) ActionsLabel(v string) *ListEditorTableBuilder {
	b.actionsLabel = v
	return b
}

func (b *ListEditorTableBuilder) Density(v string) *ListEditorTableBuilder {
	b.density = v
	return b
}

// modeLayoutFunc picks the layout override for the field's current mode.
func (b *ListEditorTableBuilder) modeLayoutFunc(fc *FieldContext) ListEditorTableLayoutFunc {
	switch {
	case fc.ReadOnly || fc.Mode.Dot().Has(DETAIL):
		return b.viewLayoutFunc
	case fc.Mode.Dot().Has(NEW):
		return b.newLayoutFunc
	default:
		return b.editLayoutFunc
	}
}

// tree builds the renderable field tree of the nested fields, preserving groups
// so the header can render nested/spanned columns. The layout is chosen per mode
// (view / new / edit) when an override is set, otherwise the current layout.
func (b *ListEditorTableBuilder) tree(fc *FieldContext) FieldBuilderTreeNodes {
	fb := fc.Nested.FieldsBuilder()
	layout := fb.CurrentLayout()
	if f := b.modeLayoutFunc(fc); f != nil {
		if l := f(fb); l != nil {
			layout = l
		}
	}
	return fb.fields.FieldTreeLayout(fb.fields.FilterLayout(layout, FieldRenderable()))
}

// hintDialog renders, when the field has a hint, a blue info ("i") icon that
// opens a dialog with the hint text — placed next to the column header label.
// Returns nil when there is no hint.
func (b *ListEditorTableBuilder) hintDialog(fctx *FieldContext) h.HTMLComponent {
	fctx.CheckHint()
	hint := fctx.Hint
	if hint == "" {
		return nil
	}
	return web.Scope(
		VIcon("mdi-information").
			Attr("size", "small").
			Color("info").
			Class("ml-1").
			Style("cursor:pointer").
			Attr("@click", "locals.hintOpen = true"),
		vx.VXDialog().
			Title(fctx.Label).
			SlotBody(h.Div(h.Text(hint)).Class("pa-4 text-body-2")).
			Closable(true).
			Width("420").
			VModel("locals.hintOpen"),
	).LocalsInit("{hintOpen: false}").Slot("{ locals }")
}

func listEditorTreeLeafCount(nodes FieldBuilderTreeNodes) (n int) {
	for _, node := range nodes {
		if node.IsTree {
			n += listEditorTreeLeafCount(node.Tree.Nodes)
		} else {
			n++
		}
	}
	return
}

func listEditorTreeDepth(nodes FieldBuilderTreeNodes) (d int) {
	for _, node := range nodes {
		depth := 1
		if node.IsTree {
			depth = 1 + listEditorTreeDepth(node.Tree.Nodes)
		}
		if depth > d {
			d = depth
		}
	}
	if d == 0 {
		d = 1
	}
	return
}

func (b *ListEditorTableBuilder) Container(ctx *ListEditorContainerContext, items h.HTMLComponents) h.HTMLComponent {
	var (
		fc       = ctx.FieldContext
		nodes    = b.tree(fc)
		depth    = listEditorTreeDepth(nodes)
		readOnly = fc.ReadOnly
		rows     = make([][]h.HTMLComponent, depth)
	)

	// walk the tree level by level: group nodes span their leaf count on their
	// level and recurse; leaf nodes span the remaining rows below them.
	var walk func(nodes FieldBuilderTreeNodes, level int)
	walk = func(nodes FieldBuilderTreeNodes, level int) {
		for _, node := range nodes {
			if node.IsTree {
				title := node.Tree.Name
				if node.Tree.Title != nil {
					title = node.Tree.Title(fc.EventContext.Context())
				}
				rows[level] = append(rows[level],
					h.Th(title).
						Attr("colspan", fmt.Sprint(listEditorTreeLeafCount(node.Tree.Nodes))).
						Class("text-center"))
				walk(node.Tree.Nodes, level+1)
			} else {
				fctx := node.Field.NewContext(fc.ModelInfo, fc.EventContext, nil, nil)
				rows[level] = append(rows[level],
					h.Tag("th").Attr("rowspan", fmt.Sprint(depth-level)).Children(
						h.Div(
							h.Span(fctx.Label),
							b.hintDialog(fctx),
						).Class("d-flex align-center"),
					))
			}
		}
	}
	walk(nodes, 0)

	if !readOnly {
		label := b.actionsLabel
		if label == "" {
			label = MustGetMessages(fc.EventContext.Context()).ListEditorActions
		}
		rows[0] = append(rows[0],
			h.Th(label).Attr("rowspan", fmt.Sprint(depth)).Class("text-right"))
	}

	headRows := make([]h.HTMLComponent, 0, depth)
	for _, r := range rows {
		headRows = append(headRows, h.Tr(r...))
	}

	return VTable(
		h.Thead(headRows...),
		h.Tbody(items...),
	).Attr("density", b.density)
}

func (b *ListEditorTableBuilder) Item(ctx *ListEditorItemContext) h.HTMLComponent {
	var (
		fc     = ctx.FieldContext
		leaves = b.tree(fc).Fields()
		cells  = make([]h.HTMLComponent, 0, len(leaves)+1)
	)

	for _, f := range leaves {
		cells = append(cells, h.Td(b.renderCell(ctx, f, false)).Class("py-1"))
	}

	if !fc.ReadOnly {
		var actionComps []h.HTMLComponent
		for _, af := range b.actions {
			if c := af(ctx); c != nil {
				actionComps = append(actionComps, c)
			}
		}
		actionComps = append(actionComps,
			VBtn("").
				Color("error").
				Variant(VariantText).
				Density(DensityCompact).
				Icon("mdi-delete").
				// mark the item deleted on the reactive form (client-side), so it
				// can be reverted and its data is not lost.
				Attr("@click", ctx.DeleteExpr()))

		// hidden fields (e.g. the primary key) are kept in the DOM — invisibly —
		// so the row stays submittable (and, when deleted, the {ID,__deleted}
		// payload has the id) even though the row is v-show hidden.
		cell := h.Td(
			h.Div(actionComps...).Class("d-flex justify-end align-center"),
			h.Div(b.hiddenCells(ctx)...).Style("display:none"),
		).Class("py-1")
		cells = append(cells, cell)
	}

	return h.Tr(cells...).Attr("v-show", "!"+ctx.DeletedCond())
}

func (b *ListEditorTableBuilder) DeletedItem(ctx *ListEditorItemContext) h.HTMLComponent {
	var (
		fc    = ctx.FieldContext
		msgr  = MustGetMessages(fc.EventContext.Context())
		label = ctx.Label
	)
	if label == "" {
		label = msgr.ListEditorDeletedItem
	}

	cells := []h.HTMLComponent{
		h.Td(
			h.Div(
				VIcon("mdi-delete-outline").Class("mr-2").Color("error"),
				h.Span(label).Class("text-decoration-line-through text-medium-emphasis"),
			).Class("d-flex align-center"),
		).Attr("colspan", fmt.Sprint(listEditorTreeLeafCount(b.tree(fc)))),
	}

	if !fc.ReadOnly {
		cells = append(cells, h.Td(
			h.Div(
				VBtn(msgr.ListEditorRevertDeletion).
					Variant(VariantText).
					Color("primary").
					Density(DensityCompact).
					PrependIcon("mdi-undo-variant").
					Attr("@click", ctx.RevertExpr()),
			).Class("d-flex justify-end"),
		))
	}

	return h.Tr(cells...).Attr("v-show", ctx.DeletedCond())
}

// renderCell renders a single field of the item as the bare input for a table
// cell. It mirrors the field body rendering but sets MustInput so no label,
// hint or container is produced.
func (b *ListEditorTableBuilder) renderCell(ctx *ListEditorItemContext, f *FieldBuilder, hidden bool) h.HTMLComponent {
	fc := ctx.FieldContext

	parent := *fc
	parent.FormKey = ctx.ItemFormKey
	parent.Path = ctx.Path
	parent.Obj = ctx.Item

	fctx := f.NewContext(fc.ModelInfo, fc.EventContext, &parent, ctx.Item)
	fctx.Mode = fc.Mode
	fctx.MustInput = true
	fctx.Label = ""
	fctx.Density = b.density
	if vErr, ok := fc.EventContext.Flash.(*web.ValidationErrors); ok && vErr != nil {
		fctx.Errors = vErr.GetFieldErrors(fctx.FormKey)
	}

	fb := fc.Nested.FieldsBuilder()
	fb.FieldToComponentSetup.Setup(fctx)
	f.ConfigureContext(fctx)

	if fctx.Disabled || !f.IsEnabled(fctx) {
		return nil
	}

	comp := f.ToComponent(fctx)
	// belt-and-suspenders for components that don't route through FieldWithHint
	// (selects, custom inputs): shrink the row by applying the density and — only
	// when there is no error to show — hiding the details area.
	if tg, ok := comp.(h.TagGetter); ok && tg.GetHTMLTagBuilder() != nil {
		tag := tg.GetHTMLTagBuilder()
		if len(fctx.Errors) == 0 {
			tag.Attr("hide-details", true)
		}
		if b.density != "" {
			tag.Attr("density", b.density)
		}
	}
	return comp
}

// hiddenCells renders the nested fields' hidden inputs (e.g. the primary key)
// so they are submitted with the row.
func (b *ListEditorTableBuilder) hiddenCells(ctx *ListEditorItemContext) (comps []h.HTMLComponent) {
	fb := ctx.FieldContext.Nested.FieldsBuilder()
	for _, name := range fb.hiddenFields {
		if c := b.renderCell(ctx, fb.GetFieldOrDefault(name), true); c != nil {
			comps = append(comps, c)
		}
	}
	return
}
