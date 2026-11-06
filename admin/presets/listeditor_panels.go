package presets

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

// ListEditorPanelsBuilder is a ListEditorComponentBuilder that renders the
// list editor as expansion panels: each item a panel — its label (and its
// delete) in the title, its form (the item's fields, as the list's normal
// mode) inside. The table's on a phone (ListEditorTableBuilder.Mobile): a row
// of inputs side by side does not fit there.
type ListEditorPanelsBuilder struct {
	// foot follows the panels (the table's tfoot, its totals): nil, none
	foot func(ctx *ListEditorContainerContext) h.HTMLComponent
}

func NewListEditorPanelsBuilder() *ListEditorPanelsBuilder {
	return &ListEditorPanelsBuilder{}
}

func (b *ListEditorPanelsBuilder) Container(ctx *ListEditorContainerContext, items h.HTMLComponents) h.HTMLComponent {
	panels := VExpansionPanels(items...).Multiple(true).Variant("accordion").Attr("data-list-editor-panels", true)
	if b.foot == nil {
		return panels
	}
	if foot := b.foot(ctx); foot != nil {
		return h.Components(panels, foot)
	}
	return panels
}

// title is what a panel's title says: the item's label, else its place.
func (b *ListEditorPanelsBuilder) title(ctx *ListEditorItemContext) string {
	if ctx.Label != "" {
		return ctx.Label
	}
	return fmt.Sprintf("#%d", ctx.ItemKey+1)
}

func (b *ListEditorPanelsBuilder) Item(ctx *ListEditorItemContext) h.HTMLComponent {
	fc := ctx.FieldContext
	title := VExpansionPanelTitle(
		h.Div(
			// its place in the list (the form's __pos: the order, sorted)
			VChip(
				h.Span("").Attr("v-text", fmt.Sprintf("(form[%q] ?? %d) + 1", ctx.ItemPositionFormKey, ctx.ItemKey)),
			).Size(SizeSmall).Variant(VariantTonal).Color("primary").Label(true).
				Class("mr-2 flex-shrink-0").
				Attr("data-list-editor-panel-index", true),
			h.Span(b.title(ctx)).Class("text-truncate"),
			VSpacer(),
			h.If(!fc.ReadOnly,
				VBtn("").
					Color("error").
					Variant(VariantText).
					Density(DensityCompact).
					Icon("mdi-delete").
					Class("mr-2").
					// the panel stays as it is (closed or open)
					Attr("@click.stop", ctx.DeleteExpr()).
					Attr("data-list-editor-panel-delete", true),
			),
		).Class("d-flex align-center flex-grow-1").Style("min-width: 0"),
	)
	return VExpansionPanel(
		title,
		// eager: its inputs mounted even closed — each puts its value in the
		// form when mounted (v-assign), and a panel never opened would post
		// none (a save would empty its fields)
		VExpansionPanelText(ctx.Body()).Eager(true),
	).Attr("v-show", "!"+ctx.DeletedCond()).Attr("data-list-editor-panel", true)
}

func (b *ListEditorPanelsBuilder) DeletedItem(ctx *ListEditorItemContext) h.HTMLComponent {
	fc := ctx.FieldContext
	msgr := MustGetMessages(fc.EventContext.Context())
	label := ctx.Label
	if label == "" {
		label = msgr.ListEditorDeletedItem
	}
	return VExpansionPanel(
		VExpansionPanelTitle(
			h.Div(
				VIcon("mdi-delete-outline").Class("mr-2").Color("error"),
				h.Span(label).Class("text-decoration-line-through text-medium-emphasis text-truncate"),
				VSpacer(),
				h.If(!fc.ReadOnly,
					h.Div(
						VBtn("").Variant(VariantText).Color("primary").Density(DensityCompact).
							Icon("mdi-undo-variant").Attr("title", msgr.ListEditorRevertDeletion).
							Attr("@click.stop", ctx.RevertExpr()),
						VBtn("").Variant(VariantText).Color("error").Density(DensityCompact).
							Icon("mdi-delete-forever").Attr("title", msgr.ListEditorRemoveItem).
							Attr("@click.stop", ctx.PurgeExpr()),
					).Class("d-flex ga-1 mr-2"),
				),
			).Class("d-flex align-center flex-grow-1").Style("min-width: 0"),
		).HideActions(true),
	).Readonly(true).Attr("v-show", ctx.DeletedVisibleCond())
}

// ListEditorMobileBuilder is a ListEditorComponentBuilder with another one
// for a phone's screen (narrower than md): the list editor renders both, the
// browser shows the one of its width.
type ListEditorMobileBuilder interface {
	Mobile() ListEditorComponentBuilder
}
