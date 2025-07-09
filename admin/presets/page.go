package presets

import (
	"context"
	"errors"
	"net/http"
	"path"
	"strings"

	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
)

type PageBuilder struct {
	path         string
	fullPath     string
	methods      []string
	handler      http.Handler
	menuGroup    string
	verififer    *perm.PermVerifierBuilder
	baseVerifier *perm.Verifier
	titleFunc    func(ctx context.Context) string
	autoPerm     bool
	notInMenu    bool
	menuIcon     string
	postBuild    []func(ph *PageHandler)
	pageHandler  *PageHandler
}

func Page(path string) *PageBuilder {
	return &PageBuilder{path: path}
}

func (b *PageBuilder) Methods(methods ...string) *PageBuilder {
	for i, method := range methods {
		methods[i] = strings.ToUpper(method)
	}
	b.methods = append(b.methods, methods...)
	return b
}

func (b *PageBuilder) GetMethods() []string {
	return b.methods
}

func (b *PageBuilder) Handler(h http.Handler) *PageBuilder {
	b.handler = h
	return b
}

func (b *PageBuilder) GetHandler() http.Handler {
	return b.handler
}

func (b *PageBuilder) MenuGroup(menuGroup string) *PageBuilder {
	b.menuGroup = menuGroup
	return b
}

func (b *PageBuilder) GetMenuGroup() string {
	return b.menuGroup
}

func (b *PageBuilder) Perm(v *perm.PermVerifierBuilder) *PageBuilder {
	b.verififer = v
	return b
}

func (b *PageBuilder) GetVerifier() *perm.PermVerifierBuilder {
	return b.verififer
}

func (b *PageBuilder) BaseVerifier(verifier *perm.Verifier) *PageBuilder {
	b.baseVerifier = verifier
	return b
}

func (b *PageBuilder) GetBaseVerifier() *perm.Verifier {
	return b.baseVerifier
}

func (b *PageBuilder) TitleFunc(f func(ctx context.Context) string) *PageBuilder {
	b.titleFunc = f
	return b
}

func (b *PageBuilder) StringTitle(t string) *PageBuilder {
	b.titleFunc = func(context.Context) string {
		return t
	}
	return b
}

func (b *PageBuilder) GetTitleFunc() func(ctx context.Context) string {
	return b.titleFunc
}

func (b *PageBuilder) TTitle(ctx context.Context) string {
	if b.titleFunc == nil {
		return ""
	}
	return b.titleFunc(ctx)
}

func (b *PageBuilder) InMenu(v bool) (r *PageBuilder) {
	b.notInMenu = !v
	return b
}

func (b *PageBuilder) IsInMenu() bool {
	return !b.notInMenu
}

func (b *PageBuilder) MenuIcon(v string) (r *PageBuilder) {
	b.menuIcon = v
	return b
}

func (b *PageBuilder) GetMenuIcon() string {
	return b.menuIcon
}

func (b *PageBuilder) VerifierWithBase(base *perm.Verifier, r *http.Request) *perm.Verifier {
	return b.verififer.Build(base.Spawn().WithReq(r).Do(PermFromRequest(r)))
}

func (b *PageBuilder) Verifier(r *http.Request) *perm.Verifier {
	return b.VerifierWithBase(b.baseVerifier, r)
}

func (b *PageBuilder) AutoPerm() *PageBuilder {
	b.autoPerm = true
	if b.verififer == nil {
		b.verififer = perm.PermVerifier()
	}
	return b
}

func (b *PageBuilder) IsAutoPerm() *PageBuilder {
	b.autoPerm = true
	return b
}

func (b *PageBuilder) HandlerFromPageFunc(wrap func(f web.PageFunc) http.Handler, f web.PageFunc) *PageBuilder {
	b.handler = wrap(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		if b.verififer != nil && b.Verifier(ctx.R).Denied() {
			err = perm.PermissionDenied
			return
		}
		return f(ctx)
	})
	return b
}

func (b *PageBuilder) PageHandler() *PageHandler {
	return b.pageHandler
}

func (b *PageBuilder) PostBuild(ph ...func(ph *PageHandler)) *PageBuilder {
	b.postBuild = append(b.postBuild, ph...)
	return b
}

func (b *PageBuilder) FullPath() string {
	return b.fullPath
}

func (b *PageBuilder) Build(prefix string) *PageHandler {
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
	return ph
}

type PresetsPageBuilder struct {
	p       *Builder
	page    *PageBuilder
	events  web.EventsHub
	actions []*ActionBuilder
	layout  *LayoutConfig
	pf      web.PageFunc
	handler func(f web.PageFunc) http.Handler
}

func PresetsPage(b *Builder, page *PageBuilder) *PresetsPageBuilder {
	return &PresetsPageBuilder{p: b, page: page}
}

func (b *PresetsPageBuilder) Page() *PageBuilder {
	return b.page
}

func (b *PresetsPageBuilder) EventHandler(eventFuncId string, ef web.EventHandler) *PresetsPageBuilder {
	b.events.RegisterEventHandler(eventFuncId, ef)
	return b
}

func (b *PresetsPageBuilder) EventFunc(eventFuncId string, ef web.EventFunc) *PresetsPageBuilder {
	return b.EventHandler(eventFuncId, ef)
}

func (b *PresetsPageBuilder) PrivateEventFunc(eventFuncId string, ef web.EventFunc, vf ...func(v *perm.Verifier) *perm.Verifier) *PresetsPageBuilder {
	return b.EventFunc(eventFuncId, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		v := b.page.Verifier(ctx.R)
		for _, f := range vf {
			v = f(v)
		}
		if v.Denied() {
			err = perm.PermissionDenied
			return
		}
		return ef(ctx)
	})
}

func (b *PresetsPageBuilder) Raw(pf web.PageFunc) *PresetsPageBuilder {
	b.pf = pf
	b.handler = func(f web.PageFunc) http.Handler {
		return b.p.Wrap(f, b.buildPage)
	}
	return b
}

func (b *PresetsPageBuilder) Layout(pf web.PageFunc, layout ...*LayoutConfig) *PresetsPageBuilder {
	var lc *LayoutConfig
	if len(layout) > 0 {
		lc = layout[0]
	}
	if lc == nil {
		lc = b.p.homePageLayoutConfig
	}
	b.pf = pf
	b.handler = func(f web.PageFunc) http.Handler {
		return b.p.Wrap(b.p.layoutFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
			if r, err = f(ctx); err == nil {
				if r.PageTitle == "" {
					r.PageTitle = b.page.TTitle(ctx.Context())
				}
			}
			return
		}, lc), b.buildPage)
	}
	return b
}

func (b *PresetsPageBuilder) parseRequestAction(ctx *web.EventContext) (id string, action *ActionBuilder, err error) {
	action = getAction(b.actions, ctx.R.FormValue(ParamAction))
	if action == nil {
		err = errors.New("action required")
		return
	}

	var enabled bool
	if enabled, err = action.IsEnabled(id, ctx); err != nil {
		return
	}

	if !enabled {
		err = errors.New("action disabled")
		return
	}
	return
}

func (b *PresetsPageBuilder) actionForm(ctx *web.EventContext) (r web.EventResponse, err error) {
	var action *ActionBuilder
	if _, action, err = b.parseRequestAction(ctx); err != nil {
		return
	}
	err = action.View(nil, "", ctx, &r)
	return
}

func (b *PresetsPageBuilder) doAction(ctx *web.EventContext) (r web.EventResponse, err error) {
	var action *ActionBuilder
	if _, action, err = b.parseRequestAction(ctx); err != nil {
		return
	}

	_, err = action.Do(nil, "", ctx, &r)
	return
}

func (b *PresetsPageBuilder) buildPage(p *web.PageBuilder) {
	b.EventFunc(actions.Action, b.actionForm)
	b.EventFunc(actions.DoAction, b.doAction)
	p.MergeHub((&b.events).Wrap(WrapEventHandler))
}

func (b *PresetsPageBuilder) Action(name string) *ActionBuilder {
	action := Action(b.p, name)
	action.typ = ActionTypePage
	b.actions = append(b.actions, action)
	return action
}

func (b *PresetsPageBuilder) PrivateAction(name string, vf ...func(v *perm.Verifier) *perm.Verifier) *ActionBuilder {
	return b.Action(name).SetVerifier(func(ctx *web.EventContext) *perm.Verifier {
		v := b.page.Verifier(ctx.R)
		for _, f := range vf {
			v = f(v)
		}
		return v
	})
}

func (b *PresetsPageBuilder) Private() *PresetsPageBuilder {
	b.page.AutoPerm().BaseVerifier(b.p.verifier)
	b.page.PostBuild(func(*PageHandler) {
		b.page.BaseVerifier(b.p.verifier)
	})
	return b
}

func (b *Builder) AddPage(page *PageBuilder) *Builder {
	b.pages = append(b.pages, page)
	return b
}

func (b *Builder) CreatePage(page *PageBuilder) *PresetsPageBuilder {
	b.pages = append(b.pages, page)
	return PresetsPage(b, page)
}

func (b *Builder) GetPage(pth string) *PageBuilder {
	for _, page := range b.pages {
		if page.path == pth {
			return page
		}
	}
	return nil
}

func (b *PresetsPageBuilder) Build() {
	if b.pf != nil {
		b.page.HandlerFromPageFunc(b.handler, b.pf)
	}
}
