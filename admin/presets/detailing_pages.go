package presets

import (
	"errors"
	"net/http"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
)

func (d *DetailingBuilder) setupPagesRegistrator() {
	d.pagesRegistrator = NewPagesRegistrator(
		d.mb.p,
		func() string {
			return d.mb.itemRoutePath
		},
		func(pf web.PageFunc, do DoPageBuilder) http.Handler {
			return d.mb.p.Wrap(pf, do)
		},
		RequestPermVerifierFunc(func(r *http.Request) *perm.Verifier {
			return d.mb.permissioner.Reader(r, r.Context().Value(ModelIDKey).(model.ID), r.Context().Value(ParentsModelIDKey).(model.IDSlice)...)
		}),
	).
		Model(d.mb).
		LayoutFunc(func(config *LayoutConfig, f func(ctx *web.EventContext) (r web.PageResponse, err error)) web.PageFunc {
			return d.mb.p.detailLayoutFunc(d.mb.BindPageFunc(f), config)
		}).
		RecordIdLayoutFunc(func(config *LayoutConfig, f func(ctx *web.EventContext) (r web.PageResponse, err error)) web.PageFunc {
			return d.mb.p.detailLayoutFunc(d.mb.BindPageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
				var mid ID

				if !d.mb.singleton {
					id := ctx.Param(ParamID)

					if id == "" {
						err = MustGetMessages(ctx.Context()).ErrEmptyParamID
						return
					}

					if mid, err = d.mb.ParseRecordID(id); err != nil {
						return
					}

					parentsID := ParentsModelID(ctx.R)

					if !IsSkipAutoBreadcrumb(ctx) {
						var (
							bc     = GetOrInitBreadcrumbs(ctx.R)
							record = d.mb.NewModel()
							label  string
						)

						mid.SetTo(record)

						if label, err = d.mb.RecordTitleFetch(record, ctx); err != nil {
							return
						}

						bc.Append(&Breadcrumb{
							Label: label,
							URI:   d.mb.Info().DetailingHref(mid.String(), parentsID...),
						})
					}
				}

				ctx.WithContextValue(ModelIDKey, mid)

				return f(ctx)
			}), config)
		})
}

func (b *DetailingBuilder) PagesRegistrator() *PagesRegistrator {
	return b.pagesRegistrator
}

func (b *DetailingBuilder) BuildPage(vf *perm.PermVerifierBuilder, builder func(ctx *web.EventContext, obj any, mid model.ID) (r web.PageResponse, err error)) func(ctx *web.EventContext) (r web.PageResponse, err error) {
	if vf == nil {
		vf = perm.PermVerifier()
	}
	return b.mb.BindPageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		var mid ID

		obj := b.mb.NewModel()
		msgr := MustGetMessages(ctx.Context())

		if !b.mb.singleton {
			id := ctx.Param(ParamID)
			if id == "" {
				err = msgr.ErrEmptyParamID
				return
			}

			if mid, err = b.mb.ParseRecordID(id); err != nil {
				return
			}
		}

		v := vf.Build(b.mb.permissioner.Reader(ctx.R, mid, ParentsModelID(ctx.R)...))

		if v.Denied() {
			r.Body = h.Div(h.Text(MustGetMessages(ctx.Context()).ErrPermissionDenied.Error()))
			return
		}

		err = b.GetFetchFunc()(obj, mid, ctx)
		if err != nil {
			if errors.Is(err, ErrRecordNotFound) {
				return b.mb.p.DefaultNotFoundPageFunc(ctx)
			}
			return
		}

		return builder(ctx, obj, mid)
	})
}
