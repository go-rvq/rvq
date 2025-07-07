package presets

import (
	"net/url"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

type DeletingMenuItemBuilder struct {
	modelInfo *ModelInfo
	url       string
	urlValues url.Values
	wrapEvent func(rctx *RecordMenuItemContext, e *web.VueEventTagBuilder)
}

func NewDeletingMenuItemBuilder(modelInfo *ModelInfo) *DeletingMenuItemBuilder {
	return &DeletingMenuItemBuilder{modelInfo: modelInfo}
}

func (r *DeletingMenuItemBuilder) Url() string {
	return r.url
}

func (r *DeletingMenuItemBuilder) SetUrl(url string) *DeletingMenuItemBuilder {
	r.url = url
	return r
}

func (r *DeletingMenuItemBuilder) UrlValues() url.Values {
	return r.urlValues
}

func (r *DeletingMenuItemBuilder) SetUrlValues(urlValues url.Values) *DeletingMenuItemBuilder {
	r.urlValues = urlValues
	return r
}

func (r *DeletingMenuItemBuilder) WrapEvent() func(rctx *RecordMenuItemContext, e *web.VueEventTagBuilder) {
	return r.wrapEvent
}

func (r *DeletingMenuItemBuilder) SetWrapEvent(wrapEvent func(rctx *RecordMenuItemContext, e *web.VueEventTagBuilder)) *DeletingMenuItemBuilder {
	r.wrapEvent = wrapEvent
	return r
}

func (r *DeletingMenuItemBuilder) Build() RecordMenuItemFunc {
	return func(rctx *RecordMenuItemContext) h.HTMLComponent {
		var (
			ctx = rctx.Ctx
			obj = rctx.Obj
			id  = rctx.ID
			mi  = r.modelInfo
		)

		if mi.mb.deletingDisabled || !mi.mb.CanDeleteObj(obj, ctx) {
			return nil
		}

		msgr := MustGetMessages(ctx.Context())

		onclick := web.Plaid().
			EventFunc(actions.DeleteConfirmation).
			Queries(r.urlValues).
			Query(ParamID, id).
			Query(ParamTargetPortal, rctx.TempPortal).
			URL(mi.ListingHrefCtx(ctx))

		if r.wrapEvent != nil {
			r.wrapEvent(rctx, onclick)
		}

		return VListItem(
			web.Slot(
				VIcon("mdi-delete").
					Color(ColorError),
			).Name("prepend"),
			VListItemTitle(h.Text(msgr.Delete)).Class("text-error"),
		).Attr("@click", onclick.Go())
	}
}
