package presets

import (
	"net/url"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	v "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

const defaultBulkActionDialogWidth = "600"

type BulkActionLinkHandler func(baseModel *ModelBuilder, ctx *web.EventContext, r *web.EventResponse, q url.Values, selectedIds []string)

type BulkActionBuilder struct {
	NameLabel

	buttonCompFunc                 ComponentFunc
	updateFunc                     BulkActionUpdateFunc
	compFunc                       BulkActionComponentFunc
	selectedIdsProcessorFunc       BulkActionSelectedIdsProcessorFunc
	selectedIdsProcessorNoticeFunc BulkActionSelectedIdsProcessorNoticeFunc
	componentHandler               BulkActionComponentHandler
	doBtnLabel                     func(ctx *web.EventContext) string

	linkHandler BulkActionLinkHandler

	dialogWidth   string
	buttonColor   string
	fullComponent bool

	l *ListingBuilder
}

func getBulkAction(actions []*BulkActionBuilder, name string) *BulkActionBuilder {
	for _, f := range actions {
		if f.name == name {
			return f
		}
	}
	return nil
}

func (b *BulkActionBuilder) Label(v string) *BulkActionBuilder {
	b.label = v
	return b
}

func (b *BulkActionBuilder) SetI18nLabel(i18nLabel func(ctx web.ContextValuer) string) *BulkActionBuilder {
	b.NameLabel.SetI18nLabel(i18nLabel)
	return b
}

func (b *BulkActionBuilder) ButtonCompFunc(v ComponentFunc) (r *BulkActionBuilder) {
	b.buttonCompFunc = v
	return b
}

func (b *BulkActionBuilder) UpdateFunc(v BulkActionUpdateFunc) (r *BulkActionBuilder) {
	b.updateFunc = v
	return b
}

func (b *BulkActionBuilder) ComponentFunc(v BulkActionComponentFunc) (r *BulkActionBuilder) {
	b.compFunc = v
	return b
}

func (b *BulkActionBuilder) SelectedIdsProcessorFunc(v BulkActionSelectedIdsProcessorFunc) (r *BulkActionBuilder) {
	b.selectedIdsProcessorFunc = v
	return b
}

func (b *BulkActionBuilder) SelectedIdsProcessorNoticeFunc(v BulkActionSelectedIdsProcessorNoticeFunc) (r *BulkActionBuilder) {
	b.selectedIdsProcessorNoticeFunc = v
	return b
}

func (b *BulkActionBuilder) DialogWidth(v string) (r *BulkActionBuilder) {
	b.dialogWidth = v
	return b
}

func (b *BulkActionBuilder) ButtonColor(v string) (r *BulkActionBuilder) {
	b.buttonColor = v
	return b
}

func (b *BulkActionBuilder) LinkHandler() BulkActionLinkHandler {
	return b.linkHandler
}

func (b *BulkActionBuilder) SetLinkHandler(linkHandler BulkActionLinkHandler) *BulkActionBuilder {
	b.linkHandler = linkHandler
	return b
}

func (b *BulkActionBuilder) SetFullComponent(v bool) (r *BulkActionBuilder) {
	b.fullComponent = v
	return b
}

func (b *BulkActionBuilder) SetDoBtnLabel(f func(ctx *web.EventContext) string) *BulkActionBuilder {
	b.doBtnLabel = f
	return b
}

func (b *BulkActionBuilder) ComponentHandler() BulkActionComponentHandler {
	return b.componentHandler
}

func (b *BulkActionBuilder) SetComponentHandler(componentHandler BulkActionComponentHandler) *BulkActionBuilder {
	b.componentHandler = componentHandler
	return b
}

func (b *BulkActionBuilder) DoBtnLabel(ctx *web.EventContext) string {
	if b.doBtnLabel != nil {
		return b.doBtnLabel(ctx)
	}
	return MustGetMessages(ctx.Context()).Execute
}

func (b *BulkActionBuilder) RequestTitle(ctx web.ContextValuer) (label string) {
	label = b.labelKey

	if label == "" {
		if b.i18nLabel != nil {
			return b.i18nLabel(ctx)
		}
		if label = b.label; label == "" {
			label = b.name
		}
	}

	return i18n.Translate(b.l.mb.BulkActionTranslator(), ctx.Context(), label)
}

func (b *BulkActionBuilder) Component(selectedIds []string, overlay actions.OverlayMode, ctx *web.EventContext) h.HTMLComponent {
	var body h.HTMLComponent
	var err error
	if b.fullComponent {
		body, err = b.compFunc(selectedIds, ctx)
	} else if b.compFunc != nil {
		body, err = b.compFunc(selectedIds, ctx)
	}

	if err != nil {
		body = v.VAlert(h.RawHTML(err.Error())).
			Type(v.ColorError).
			Variant(v.VariantTonal).
			Density(v.DensityCompact)
	}

	cb := &ContentComponentBuilder{
		Overlay: &ContentComponentBuilderOverlay{
			Mode:       overlay,
			Scrollable: true,
		},
		Title: b.RequestTitle(ctx),
		Body:  body,
	}

	if err == nil && b.updateFunc != nil {
		cb.BottomActions = append(cb.BottomActions, v.VBtn("").
			Color("primary").
			Variant(v.VariantFlat).
			Attr(":disabled", "isFetching").
			Attr(":loading", "isFetching").
			Density("comfortable").
			Children(v.VIcon("mdi-run-fast"), h.Text(b.DoBtnLabel(ctx))).
			Attr("@click", web.Plaid().
				EventFunc(actions.DoBulkAction).
				Query(ParamBulkActionName, b.name).
				MergeQuery(true).
				URL(ctx.R.RequestURI).
				Go()))
	}

	if ctx.Flash != nil {
		cb.Notice(ctx.Flash)
	}

	WithRespondDialogHandlers(ctx, func(d *DialogBuilder) {
		d.SetScrollable(true)
	})

	if b.componentHandler != nil {
		b.componentHandler(cb, ctx)
	}

	return cb.BuildOverlay()
}

func (b *BulkActionBuilder) Button(mb *ModelBuilder, ctx *web.EventContext) h.HTMLComponent {
	if b.buttonCompFunc != nil {
		return b.buttonCompFunc(ctx)
	}
	return b.DefaultButton(mb, ctx, b.RequestTitle(ctx), b.DefaulButtonOnClick(ctx).Go())
}

func (b *BulkActionBuilder) DefaulButtonOnClick(ctx *web.EventContext) *web.VueEventTagBuilder {
	onclick := web.Plaid().EventFunc(actions.OpenBulkActionDialog).
		Queries(url.Values{bulkPanelOpenParamName: []string{b.name}}).
		MergeQuery(true)
	if IsInDialog(ctx) {
		onclick.URL(ctx.R.RequestURI).
			Query(ParamOverlay, actions.Dialog)
	}
	return onclick
}

func (b *BulkActionBuilder) DefaultButton(mb *ModelBuilder, ctx *web.EventContext, text, onClick string) h.HTMLComponent {
	buttonColor := b.buttonColor
	if buttonColor == "" {
		buttonColor = v.ColorSecondary
	}

	return v.VBtn(text).
		Color(buttonColor).
		Variant(v.VariantFlat).
		// Size(SizeSmall).
		Class("ml-2").
		Attr("@click", onClick)
}

func (b *BulkActionBuilder) component(selectedIds []string, ctx *web.EventContext) h.HTMLComponent {
	comp := b.Component(selectedIds, actions.OverlayMode(ctx.Param(ParamOverlay)), ctx)
	return web.Scope(comp).FormInit()
}

func (b *BulkActionBuilder) View(baseModel *ModelBuilder, selectedIds []string, ctx *web.EventContext, r *web.EventResponse) (err error) {
	if b.linkHandler != nil {
		q := ctx.R.URL.Query()
		q.Del("__execute_event__")
		q.Del(bulkPanelOpenParamName)
		q.Del(ParamBulkActionName)

		b.linkHandler(baseModel, ctx, r, q, selectedIds)
		return
	}

	baseModel.p.dialog(ctx, r, b.component(selectedIds, ctx), b.dialogWidth)
	return
}

func (b *BulkActionBuilder) Do(selectedIds []string, ctx *web.EventContext) (err error) {
	return b.updateFunc(selectedIds, ctx)
}

func (b *ListingBuilder) BulkAction(name string) (r *BulkActionBuilder) {
	builder := getBulkAction(b.bulkActions, name)
	if builder != nil {
		return builder
	}

	r = &BulkActionBuilder{
		l: b,
	}
	r.name = name
	r.buttonColor = "black"
	r.dialogWidth = defaultBulkActionDialogWidth
	b.bulkActions = append(b.bulkActions, r)
	return
}
