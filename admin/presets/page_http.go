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
	path     string
	fullPath string
	methods  []string
	handler  http.Handler
	// b is the builder the page was registered with, and the owner of the menu
	// tree the page's entry lives in. It is nil until registration.
	b *Builder
	// menuGroupName is a group named before the page was registered, when there
	// is no builder yet to resolve it against. Registration resolves it.
	menuGroupName string
	verififer     *perm.PermVerifierBuilder
	baseVerifier  RequestPermVerifier
	titleFunc     func(ctx context.Context) string
	subTitleFunc  func(ctx context.Context) string
	// descriptionFunc is what the page is, in the language of the request
	descriptionFunc func(ctx context.Context) string
	menuItemFunc    func(ctx *web.EventContext, uri string) h.HTMLComponent
	autoPerm        bool
	notInMenu       bool
	menuIcon        string
	postBuild       []func(ph *PageHandler)
	pageHandler     *PageHandler
	preWraper       web.PageFuncWrapper
	wraper          web.PageFuncWrapper
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

// MenuGroup puts the page in the group of that name, creating the group when it
// does not exist yet. A page built before it is registered has no builder to
// resolve the name against, so the name waits until registration.
func (b *HttpPageBuilder) MenuGroup(menuGroup string) *HttpPageBuilder {
	if menuGroup == "" {
		return b
	}
	b.menuGroupName = menuGroup
	if b.b != nil {
		return b.SetMenuGroup(b.b.MenuGroup(menuGroup))
	}
	return b
}

// registerMenu gives the page its entry in the menu tree, resolving a group
// named before the page had a builder.
func (b *HttpPageBuilder) registerMenu(bd *Builder) {
	b.b = bd
	if _, err := bd.RegisterMenuItem(MenuItemPage, b.path, b); err != nil {
		panic(err)
	}
	if b.menuGroupName != "" {
		b.SetMenuGroup(bd.MenuGroup(b.menuGroupName))
	}
}

// SetMenuGroup puts the page in the group, moving its menu entry with it. A
// page with no entry in the tree — a model's page, which is a child of the
// listing or of the detailing — has nothing to move: it keeps the name, which
// is what shapes its URL.
func (b *HttpPageBuilder) SetMenuGroup(g *MenuGroupBuilder) *HttpPageBuilder {
	if g == nil {
		return b
	}
	b.menuGroupName = g.Path()
	if b.b == nil {
		return b
	}
	if err := g.b.MoveMenuItem(PageItem(b.path), g); err != nil {
		panic(err)
	}
	return b
}

// GetMenuGroupBuilder is where the page's entry sits in the menu tree, or nil
// at the root. Read from the tree, never copied onto the page: a copy would go
// stale the moment something moved the entry.
func (b *HttpPageBuilder) GetMenuGroupBuilder() *MenuGroupBuilder {
	if b.b == nil {
		return nil
	}
	return b.b.menuGroupOf(PageItem(b.path))
}

// GetMenuGroup returns the page's menu-group path — "a/b/c" for a page in a
// group c nested in b nested in a. It is the page's URL prefix.
func (b *HttpPageBuilder) GetMenuGroup() string {
	return path.Join(b.menuGroupPathNames()...)
}

// menuGroupPathNames is where the page sits, group by group. A page of the
// builder reads it from the tree, which is the only place that knows about
// moves. A page of a model has no entry there — it is a child of the listing
// or of the detailing, not a side-menu item — so it keeps the name it was
// given, which still shapes its URL and its permission.
func (b *HttpPageBuilder) menuGroupPathNames() []string {
	if b.b != nil {
		// In the tree, and the tree is the only authority: a page moved back to
		// the root has no group, whatever name it was given before.
		return b.GetMenuGroupBuilder().PathNames()
	}
	if b.menuGroupName == "" {
		return nil
	}
	return strings.Split(strings.Trim(b.menuGroupName, "/"), "/")
}

// menuGroupURIPathNames are the URI segments of the groups the page sits in
// (MenuGroupBuilder.URIName): what its URL is made of, while its permission
// keeps the groups' names (menuGroupPathNames).
func (b *HttpPageBuilder) menuGroupURIPathNames() []string {
	if b.b != nil {
		return b.GetMenuGroupBuilder().URIPathNames()
	}
	return b.menuGroupPathNames()
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

func (b *HttpPageBuilder) DescriptionFunc(f func(ctx context.Context) string) *HttpPageBuilder {
	b.descriptionFunc = f
	return b
}

func (b *HttpPageBuilder) GetDescriptionFunc() func(ctx context.Context) string {
	return b.descriptionFunc
}

// TDescription is what the page is, in the language of ctx; "" when it has
// none.
func (b *HttpPageBuilder) TDescription(ctx context.Context) string {
	if b.descriptionFunc == nil {
		return ""
	}
	return b.descriptionFunc(ctx)
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
	groups := b.menuGroupPathNames()

	if b.autoPerm {
		// The whole chain of groups, not only the innermost: the permission
		// follows the same path the URL does.
		var parts []string
		for _, g := range groups {
			parts = append(parts, GroupPermPart(g))
		}
		parts = append(parts, b.path)
		unique := b.UniquePermName()
		b.verififer.Func(func(v *perm.Verifier) *perm.Verifier {
			// and its unique name, which decides before its groups
			base := v.ResourceParts()
			return v.On(parts...).Prefer(append(base, unique)...)
		})
	}

	if b.verififer != nil && b.titleFunc != nil {
		b.verififer.Title(b.titleFunc)
	}

	b.fullPath = path.Join("/", prefix, path.Join(b.menuGroupURIPathNames()...), b.path)
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

// Path returns the page's path segment.
func (b *HttpPageBuilder) Path() string { return b.path }

// AutoPermEnabled reports whether the page uses the automatic path-based
// permission (its path is used as the permission segment).
func (b *HttpPageBuilder) AutoPermEnabled() bool { return b.autoPerm }

// HasCustomVerifier reports whether the page has its own (developer-defined)
// verifier — i.e. its own permission check mechanism, not the automatic
// path-based one.
func (b *HttpPageBuilder) HasCustomVerifier() bool { return b.verififer != nil && !b.autoPerm }
