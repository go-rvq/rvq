package presets

import (
	"errors"
	"iter"
	"net/http"
	"path"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

type (
	RequestPermVerifier interface {
		Verifier(r *http.Request) *perm.Verifier
	}

	RequestPermVerifierFunc func(r *http.Request) *perm.Verifier

	RequestPermVerifierMiddleware func(next RequestPermVerifier) RequestPermVerifier
)

func (f RequestPermVerifierFunc) Verifier(r *http.Request) *perm.Verifier {
	return f(r)
}

type PagesRegistrator struct {
	b                  *Builder
	prefixFunc         func() string
	wrapFunc           func(pf web.PageFunc, do DoPageBuilder) http.Handler
	baseVerifier       RequestPermVerifier
	layoutFunc         func(config *LayoutConfig, f func(ctx *web.EventContext) (r web.PageResponse, err error)) web.PageFunc
	recordIdLayoutFunc func(config *LayoutConfig, f func(ctx *web.EventContext) (r web.PageResponse, err error)) web.PageFunc
	httpPages          []*HttpPageBuilder
	pages              []*PageBuilder
	handlers           PageHandlers
	builded            bool
	mb                 *ModelBuilder
	fetchFunc          FetchFunc
}

func NewPagesRegistrator(
	b *Builder,
	prefixFunc func() string,
	wrapFunc func(pf web.PageFunc, do DoPageBuilder) http.Handler,
	baseVerifierFunc RequestPermVerifier,
) *PagesRegistrator {
	return &PagesRegistrator{
		b:            b,
		prefixFunc:   prefixFunc,
		wrapFunc:     wrapFunc,
		baseVerifier: baseVerifierFunc,
	}
}

func (pr *PagesRegistrator) Model(mb *ModelBuilder) *PagesRegistrator {
	pr.mb = mb
	return pr
}

func (pr *PagesRegistrator) LayoutFunc(f func(config *LayoutConfig, f func(ctx *web.EventContext) (r web.PageResponse, err error)) web.PageFunc) *PagesRegistrator {
	pr.layoutFunc = f
	return pr
}

func (pr *PagesRegistrator) RecordIdLayoutFunc(f func(config *LayoutConfig, f func(ctx *web.EventContext) (r web.PageResponse, err error)) web.PageFunc) *PagesRegistrator {
	pr.recordIdLayoutFunc = f
	return pr
}

func (pr *PagesRegistrator) FetchFunc(f FetchFunc) *PagesRegistrator {
	pr.fetchFunc = f
	return pr
}

func (pr *PagesRegistrator) New(page *HttpPageBuilder) *PageBuilder {
	pr.httpPages = append(pr.httpPages, page)
	pg := Page(pr, page)
	pr.pages = append(pr.pages, pg)
	return pg
}

func (pr *PagesRegistrator) Get(pth string) *PageBuilder {
	for _, page := range pr.pages {
		if page.page.path == pth {
			return page
		}
	}
	return nil
}

func (pr *PagesRegistrator) AddHttpPage(page ...*HttpPageBuilder) {
	pr.httpPages = append(pr.httpPages, page...)
}

func (pr *PagesRegistrator) GetHttpPage(pth string) *HttpPageBuilder {
	for _, page := range pr.httpPages {
		if page.path == pth {
			return page
		}
	}
	return nil
}

func (pr *PagesRegistrator) Wrap(pf web.PageFunc, do DoPageBuilder) http.Handler {
	return pr.wrapFunc(pf, do)
}

func (pr *PagesRegistrator) BaseVerifier(r *http.Request) *perm.Verifier {
	return pr.baseVerifier.Verifier(r)
}

func (pr *PagesRegistrator) Layout(config *LayoutConfig, f func(ctx *web.EventContext) (r web.PageResponse, err error)) web.PageFunc {
	return pr.layoutFunc(config, f)
}

func (pr *PagesRegistrator) RecordID(f func(ctx *web.EventContext, mid model.ID, parent []model.ID) (r web.PageResponse, err error)) func(ctx *web.EventContext) (r web.PageResponse, err error) {
	return func(ctx *web.EventContext) (r web.PageResponse, err error) {
		var mid ID

		msgr := MustGetMessages(ctx.Context())

		if !pr.mb.singleton {
			id := ctx.Param(ParamID)
			if id == "" {
				err = msgr.ErrEmptyParamID
				return
			}

			if mid, err = pr.mb.ParseRecordID(id); err != nil {
				return
			}
		}

		return f(ctx, mid, ParentsModelID(ctx.R))
	}
}

func (pr *PagesRegistrator) RecordIdLayout(config *LayoutConfig, f func(ctx *web.EventContext, mid model.ID, parent []model.ID) (r web.PageResponse, err error)) func(ctx *web.EventContext) (r web.PageResponse, err error) {
	return pr.recordIdLayoutFunc(config, pr.RecordID(f))
}

func (pr *PagesRegistrator) Record(f func(ctx *web.EventContext, mid model.ID, parent []model.ID, fetch func() (any, error)) (r web.PageResponse, err error)) func(ctx *web.EventContext) (r web.PageResponse, err error) {
	return pr.RecordID(func(ctx *web.EventContext, mid model.ID, parent []model.ID) (r web.PageResponse, err error) {
		return f(ctx, mid, ParentsModelID(ctx.R), func() (obj any, err error) {
			obj = pr.mb.NewModel()
			if err = pr.fetchFunc(obj, mid, ctx); errors.Is(err, ErrRecordNotFound) {
				err = nil
				obj = nil
			}
			return
		})
	})
}

func (pr *PagesRegistrator) RecordLayout(config *LayoutConfig, f func(ctx *web.EventContext, mid model.ID, parent []model.ID, fetch func() (any, error)) (r web.PageResponse, err error)) func(ctx *web.EventContext) (r web.PageResponse, err error) {
	return pr.mb.BindPageFunc(pr.recordIdLayoutFunc(config, pr.Record(f)))
}

func (pr *PagesRegistrator) Build() *PagesRegistrator {
	prefix := pr.prefixFunc()
	for _, page := range pr.httpPages {
		pr.handlers.Add(page.Build(prefix))
	}
	pr.builded = true
	return pr
}

func (pr *PagesRegistrator) SetupRoutes(mux *http.ServeMux, cb ...func(pattern string, ph *PageHandler)) *PagesRegistrator {
	pr.BuildedHandlers().SetupRoutes(mux, cb...)
	return pr
}

func (pr *PagesRegistrator) CheckBuild() *PagesRegistrator {
	if !pr.builded {
		panic("not builded")
	}
	return pr
}

func (pr *PagesRegistrator) BuildedHandlers() PageHandlers {
	return pr.CheckBuild().handlers
}

func (pr *PagesRegistrator) BuildedVerifiers() iter.Seq[*perm.PermVerifierBuilder] {
	return func(yield func(*perm.PermVerifierBuilder) bool) {
		for _, page := range pr.CheckBuild().httpPages {
			if !yield(page.GetVerifier()) {
				return
			}
		}
	}
}
func (pr *PagesRegistrator) MenuItems(ctx *web.EventContext, uri string) (inmenu h.HTMLComponents) {
	for _, page := range pr.httpPages {
		if !page.notInMenu && page.titleFunc != nil {
			if page.menuItemFunc == nil {
				item := v.VListItem(
					v.VListItemTitle(
						h.Text(page.titleFunc(ctx.Context())),
					),
					h.If(page.menuIcon != "", web.Slot(v.VIcon(page.menuIcon)).Name("prepend")),
					h.Iff(page.subTitleFunc != nil, func() h.HTMLComponent {
						return v.VListItemSubtitle(h.Text(page.subTitleFunc(ctx.Context())))
					}),
				).
					Slim(true).
					Tag("a").
					Attr("href", path.Join(uri, page.path))
				inmenu = append(inmenu, item)
			} else {
				inmenu = append(inmenu, page.menuItemFunc(ctx, path.Join(uri, page.path)))
			}
		}
	}
	return
}
