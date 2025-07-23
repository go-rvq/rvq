package presets

import (
	"context"
	"net/http"
	"path"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
)

type HttpPageBuilder struct {
	path         string
	fullPath     string
	methods      []string
	handler      http.Handler
	menuGroup    string
	verififer    *perm.PermVerifierBuilder
	baseVerifier RequestPermVerifier
	titleFunc    func(ctx context.Context) string
	subTitleFunc func(ctx context.Context) string
	menuItemFunc func(ctx *web.EventContext, uri string) h.HTMLComponent
	autoPerm     bool
	notInMenu    bool
	menuIcon     string
	postBuild    []func(ph *PageHandler)
	pageHandler  *PageHandler
	preWraper    web.PageFuncWrapper
	wraper       web.PageFuncWrapper
}

func HttpPage(pth string) *HttpPageBuilder {
	return &HttpPageBuilder{
		path:      path.Join("/", pth),
		preWraper: web.PageFuncDefaultWrap,
		wraper:    web.PageFuncDefaultWrap,
	}
}

func (b *HttpPageBuilder) Methods(methods ...string) *HttpPageBuilder {
	for i, method := range methods {
		methods[i] = strings.ToUpper(method)
	}
	b.methods = append(b.methods, methods...)
	return b
}

func (b *HttpPageBuilder) GetMethods() []string {
	return b.methods
}

func (b *HttpPageBuilder) Handler(h http.Handler) *HttpPageBuilder {
	b.handler = h
	return b
}

func (b *HttpPageBuilder) GetHandler() http.Handler {
	return b.handler
}

func (b *HttpPageBuilder) MenuGroup(menuGroup string) *HttpPageBuilder {
	b.menuGroup = menuGroup
	return b
}

func (b *HttpPageBuilder) GetMenuGroup() string {
	return b.menuGroup
}

func (b *HttpPageBuilder) Perm(v *perm.PermVerifierBuilder) *HttpPageBuilder {
	b.verififer = v
	return b
}

func (b *HttpPageBuilder) GetVerifier() *perm.PermVerifierBuilder {
	return b.verififer
}

func (b *HttpPageBuilder) BaseVerifier(verifier RequestPermVerifier) *HttpPageBuilder {
	b.baseVerifier = verifier
	return b
}

func (b *HttpPageBuilder) GetBaseVerifier() RequestPermVerifier {
	return b.baseVerifier
}

func (b *HttpPageBuilder) TitleFunc(f func(ctx context.Context) string) *HttpPageBuilder {
	b.titleFunc = f
	return b
}

func (b *HttpPageBuilder) StringTitle(t string) *HttpPageBuilder {
	b.titleFunc = func(context.Context) string {
		return t
	}
	return b
}

func (b *HttpPageBuilder) SubTitleFunc(f func(ctx context.Context) string) *HttpPageBuilder {
	b.subTitleFunc = f
	return b
}

func (b *HttpPageBuilder) MenuItem(f func(ctx *web.EventContext, uri string) h.HTMLComponent) *HttpPageBuilder {
	b.menuItemFunc = f
	return b
}

func (b *HttpPageBuilder) GetTitleFunc() func(ctx context.Context) string {
	return b.titleFunc
}

func (b *HttpPageBuilder) TTitle(ctx context.Context) string {
	if b.titleFunc == nil {
		return ""
	}
	return b.titleFunc(ctx)
}

func (b *HttpPageBuilder) InMenu(v bool) (r *HttpPageBuilder) {
	b.notInMenu = !v
	return b
}

func (b *HttpPageBuilder) IsInMenu() bool {
	return !b.notInMenu
}

func (b *HttpPageBuilder) MenuIcon(v string) (r *HttpPageBuilder) {
	b.menuIcon = v
	return b
}

func (b *HttpPageBuilder) GetMenuIcon() string {
	return b.menuIcon
}

func (b *HttpPageBuilder) VerifierWithBase(base *perm.Verifier, r *http.Request) *perm.Verifier {
	return b.verififer.BuildDo(base.Spawn().WithReq(r), PermFromRequest(r))
}

func (b *HttpPageBuilder) Verifier(r *http.Request) *perm.Verifier {
	return b.VerifierWithBase(b.baseVerifier.Verifier(r), r)
}

func (b *HttpPageBuilder) AutoPerm() *HttpPageBuilder {
	b.autoPerm = true
	if b.verififer == nil {
		b.verififer = perm.PermVerifier()
	}
	return b
}

func (b *HttpPageBuilder) IsAutoPerm() *HttpPageBuilder {
	b.autoPerm = true
	return b
}

func (b *HttpPageBuilder) PreWrap(f func(old web.PageFuncWrapper) web.PageFuncWrapper) *HttpPageBuilder {
	b.preWraper = f(b.preWraper)
	return b
}

func (b *HttpPageBuilder) Wrap(f func(old web.PageFuncWrapper) web.PageFuncWrapper) *HttpPageBuilder {
	b.wraper = f(b.wraper)
	return b
}

func (b *HttpPageBuilder) HandlerFromPageFunc(wrap func(f web.PageFunc) http.Handler, f web.PageFunc) *HttpPageBuilder {
	f = b.wraper(b.preWraper(f))
	b.handler = wrap(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		if b.verififer != nil && b.Verifier(ctx.R).Denied() {
			err = perm.PermissionDenied
			return
		}
		return f(ctx)
	})
	return b
}

func (b *HttpPageBuilder) PageHandler() *PageHandler {
	return b.pageHandler
}

func (b *HttpPageBuilder) PostBuild(ph ...func(ph *PageHandler)) *HttpPageBuilder {
	b.postBuild = append(b.postBuild, ph...)
	return b
}

func (b *HttpPageBuilder) FullPath() string {
	return b.fullPath
}

func (b *HttpPageBuilder) Build(prefix string) *PageHandler {
	if b.autoPerm {
		var parts []string
		if len(b.menuGroup) > 0 {
			parts = append(parts, b.menuGroup)
		}
		parts = append(parts, b.path)
		b.verififer.Func(func(v *perm.Verifier) *perm.Verifier {
			return v.On(parts...)
		})
	}

	if b.verififer != nil && b.titleFunc != nil {
		b.verififer.Title(b.titleFunc)
	}

	b.fullPath = path.Join("/", prefix, b.menuGroup, b.path)
	ph := NewPageHandler(b.fullPath, b.handler, b.methods...)
	for _, f := range b.postBuild {
		f(ph)
	}
	b.pageHandler = ph

	if b.verififer != nil {
		for _, method := range b.methods {
			b.verififer.Action(PermFromHttpMethod(method), func(context context.Context) string {
				// TODO: translate method action
				return method
			})
		}
	}

	return ph
}
