package presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

// ListEditorCardsBuilder is a ListEditorComponentBuilder that renders the list
// editor as cards side by side (wrapping): each item a card with its fields —
// those of the layout, labelled — and, in edit mode, its delete control; a
// removed one, a card that offers to revert it. Files and images, say: the
// media box at the top of the card, its description under it.
//
//	NestedSlice(mb, fb).SetComponentBuilder(NewListEditorCardsBuilder().
//	    LayoutFunc(func(*FieldsBuilder) FieldsLayout { return FieldsLayout{"Media", "Description"} }))
type ListEditorCardsBuilder struct {
	// fields: the layout of the table builder (its modes) and its hidden
	// cells (the primary key)
	fields *ListEditorTableBuilder
	width  string
	after  ListEditorAfterFunc
	// floatingDelete: the delete control floats over the card's top right
	floatingDelete bool
}

func NewListEditorCardsBuilder() *ListEditorCardsBuilder {
	return &ListEditorCardsBuilder{fields: NewListEditorTableBuilder(), width: "240px"}
}

// LayoutFunc sets the fields of a card (in every mode).
func (b *ListEditorCardsBuilder) LayoutFunc(f ListEditorTableLayoutFunc) *ListEditorCardsBuilder {
	b.fields.LayoutFunc(f)
	return b
}

// Width is the width of a card ("240px").
func (b *ListEditorCardsBuilder) Width(v string) *ListEditorCardsBuilder {
	b.width = v
	return b
}

// FloatingDelete floats the delete control over the top right of the card —
// over its first field, an image, say — instead of a row of actions under it.
func (b *ListEditorCardsBuilder) FloatingDelete(v bool) *ListEditorCardsBuilder {
	b.floatingDelete = v
	return b
}

// AfterFunc sets what follows the cards, even with none (ListEditorAfterFunc).
func (b *ListEditorCardsBuilder) AfterFunc(f ListEditorAfterFunc) *ListEditorCardsBuilder {
	b.after = f
	return b
}

func (b *ListEditorCardsBuilder) After(ctx *ListEditorContainerContext) h.HTMLComponent {
	if b.after == nil {
		return nil
	}
	return b.after(ctx)
}

func (b *ListEditorCardsBuilder) Container(_ *ListEditorContainerContext, items h.HTMLComponents) h.HTMLComponent {
	return h.Div(items...).Class("d-flex flex-wrap ga-3 py-2").Attr("data-list-editor-cards", true)
}

// renderField renders a field of the item for a card: labelled, compact.
func (b *ListEditorCardsBuilder) renderField(ctx *ListEditorItemContext, f *FieldBuilder) h.HTMLComponent {
	fc := ctx.FieldContext

	parent := *fc
	parent.FormKey = ctx.ItemFormKey
	parent.Path = ctx.Path
	parent.Obj = ctx.Item

	fctx := f.NewContext(listEditorItemInfo(fc, ctx.ItemKey), fc.EventContext, &parent, ctx.Item)
	fctx.Mode = fc.Mode
	fctx.Density = DensityCompact
	if vErr, ok := fc.EventContext.Flash.(*web.ValidationErrors); ok && vErr != nil {
		fctx.Errors = vErr.GetFieldErrors(fctx.FormKey)
	}

	fb := fc.Nested.FieldsBuilder()
	fb.FieldToComponentSetup.Setup(fctx)
	f.ConfigureContext(fctx)

	if fctx.Disabled || !f.IsEnabled(fctx) {
		return nil
	}
	return f.ToComponent(fctx)
}

func (b *ListEditorCardsBuilder) Item(ctx *ListEditorItemContext) h.HTMLComponent {
	fc := ctx.FieldContext
	var body []h.HTMLComponent
	for _, f := range b.fields.tree(fc).Fields() {
		if c := b.renderField(ctx, f); c != nil {
			body = append(body, h.Div(c).Class("mb-1"))
		}
	}
	children := []h.HTMLComponent{VCardText(body...).Class("pa-3")}
	if !fc.ReadOnly {
		del := VBtn("").
			Color("error").
			Density(DensityCompact).
			Icon("mdi-delete").
			Attr("@click", ctx.DeleteExpr()).
			Attr("data-list-editor-card-delete", true)
		var actions h.HTMLComponent
		if b.floatingDelete {
			// over the top right of the first field (the card's padding in)
			// light, the icon red: seen over any image
			actions = del.Color("surface").Variant(VariantElevated).Size(SizeSmall).
				Style("position: absolute; top: 18px; right: 18px; z-index: 1; color: rgb(var(--v-theme-error)) !important")
		} else {
			actions = VCardActions(VSpacer(), del.Variant(VariantText)).Class("pt-0")
		}
		children = append(children,
			actions,
			// the hidden fields (the primary key): submitted with the item
			h.Div(b.fields.hiddenCells(ctx)...).Style("display:none"),
		)
	}
	card := VCard(children...).Variant(VariantOutlined).Width(b.width).Style("position: relative")
	return h.Div(card).Attr("v-show", "!"+ctx.DeletedCond()).Attr("data-list-editor-card", true)
}

func (b *ListEditorCardsBuilder) DeletedItem(ctx *ListEditorItemContext) h.HTMLComponent {
	fc := ctx.FieldContext
	msgr := MustGetMessages(fc.EventContext.Context())
	label := ctx.Label
	if label == "" {
		label = msgr.ListEditorDeletedItem
	}
	var actions h.HTMLComponent
	if !fc.ReadOnly {
		actions = VCardActions(
			VBtn(msgr.ListEditorRevertDeletion).Variant(VariantText).Color("primary").Density(DensityCompact).
				PrependIcon("mdi-undo-variant").Attr("@click", ctx.RevertExpr()),
			VBtn("").Variant(VariantText).Color("error").Density(DensityCompact).
				Icon("mdi-delete-forever").Attr("title", msgr.ListEditorRemoveItem).Attr("@click", ctx.PurgeExpr()),
		)
	}
	return h.Div(
		VCard(
			VCardText(
				VIcon("mdi-delete-outline").Class("mr-2").Color("error"),
				h.Span(label).Class("text-decoration-line-through text-medium-emphasis"),
			).Class("d-flex align-center pa-3"),
			actions,
		).Variant(VariantOutlined).Width(b.width),
	).Attr("v-show", ctx.DeletedVisibleCond())
}
