package presets

import (
	"net/url"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	. "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type RowMenuItemDeleteBuilder struct {
	modelInfo *ModelInfo
	url       string
	urlValues url.Values
	wrapEvent func(rctx *RecordMenuItemContext, e *web.VueEventTagBuilder)
}

func NewRowMenuItemDeleteBuilder(modelInfo *ModelInfo) *RowMenuItemDeleteBuilder {
	return &RowMenuItemDeleteBuilder{modelInfo: modelInfo}
}

func (r *RowMenuItemDeleteBuilder) Url() string {
	return r.url
}

func (r *RowMenuItemDeleteBuilder) SetUrl(url string) *RowMenuItemDeleteBuilder {
	r.url = url
	return r
}

func (r *RowMenuItemDeleteBuilder) UrlValues() url.Values {
	return r.urlValues
}

func (r *RowMenuItemDeleteBuilder) SetUrlValues(urlValues url.Values) *RowMenuItemDeleteBuilder {
	r.urlValues = urlValues
	return r
}

func (r *RowMenuItemDeleteBuilder) WrapEvent() func(rctx *RecordMenuItemContext, e *web.VueEventTagBuilder) {
	return r.wrapEvent
}

func (r *RowMenuItemDeleteBuilder) SetWrapEvent(wrapEvent func(rctx *RecordMenuItemContext, e *web.VueEventTagBuilder)) *RowMenuItemDeleteBuilder {
	r.wrapEvent = wrapEvent
	return r
}

func (r *RowMenuItemDeleteBuilder) Build() RecordMenuItemFunc {
	return func(rctx *RecordMenuItemContext) h.HTMLComponent {
		var (
			ctx = rctx.Ctx
			obj = rctx.Obj
			id  = rctx.ID
			mi  = r.modelInfo
		)
		msgr := MustGetMessages(ctx.R)
		if mi.mb.Info().Verifier().Do(PermDelete).ObjectOn(obj).WithReq(ctx.R).IsAllowed() != nil {
			return nil
		}

		onclick := web.Plaid().
			EventFunc(actions.DeleteConfirmation).
			Queries(r.urlValues).
			Query(ParamID, id).
			Query(ParamTargetPortal, rctx.SharedRecordPortal).
			URL(r.url)

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
