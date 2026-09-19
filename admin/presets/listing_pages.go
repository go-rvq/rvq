package presets

import (
	"net/http"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
)

func (b *ListingBuilder) AddPageFunc(vf *perm.PermVerifierBuilder, path string, f web.PageFunc, methods ...string) (ph *PageHandler) {
	if vf != nil && !vf.Valid() {
		vf.Path(path)
	}
	ph = NewPageHandler(path, b.mb.p.Wrap(b.mb.p.layoutFunc(b.mb.BindVerifiedPageFunc(vf, f), b.mb.layoutConfig)), methods...)
	b.pages.Add(ph)
	return
}

func (b *ListingBuilder) AddRawPageFunc(path string, f web.PageFunc, methods ...string) (ph *PageHandler) {
	ph = NewPageHandler(path, b.mb.p.Wrap(f), methods...)
	b.pages.Add(ph)
	return
}

func (b *ListingBuilder) setupPagesRegistrator() {
	b.pagesRegistrator = NewPagesRegistrator(
		b.mb.p,
		func() string {
			return b.mb.routePath
		},
		func(pf web.PageFunc, do DoPageBuilder) http.Handler {
			return b.mb.p.Wrap(pf, do)
		},
		RequestPermVerifierFunc(func(r *http.Request) *perm.Verifier {
			return b.mb.permissioner.ReqLister(r)
		}),
	).
		Model(b.mb).
		LayoutFunc(func(config *LayoutConfig, f func(ctx *web.EventContext) (r web.PageResponse, err error)) web.PageFunc {
			return b.mb.p.detailLayoutFunc(b.mb.BindPageFunc(f), config)
		})
}

func (b *ListingBuilder) PagesRegistrator() *PagesRegistrator {
	return b.pagesRegistrator
}

// registeredPagePaths are the paths of the pages registered on the listing, by
// either door: AddPageFunc/AddRawPageFunc and the registrator.
func (b *ListingBuilder) registeredPagePaths() (paths []string) {
	for _, ph := range b.pages {
		paths = append(paths, ph.Path())
	}
	for _, p := range b.pagesRegistrator.httpPages {
		paths = append(paths, p.Path())
	}
	return
}
