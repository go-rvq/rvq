package presets

import (
	"fmt"
	"sort"
	"strconv"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/js"
	"github.com/go-rvq/rvq/web/vue"
)

// FormHostBuilder renders a *form host*: the client-side owner of an edit/create
// overlay. The host declares one reactive scope variable — the form's `closer` —
// and guards the form with it:
//
//	<user-component :scope='{editing: {show:false}}' v-slot='{editing}'>
//	  … page content (a button just does `editing.show = true`) …
//	  <user-component v-if='editing.show' :scope='{form: {$parent: form}}' v-slot='{form}'>
//	    <go-plaid-portal portal-name='…' :scope='{"closer": editing}'/>
//	    <go-plaid-run-script :script='plaid()…eventFunc("presets_Edit")…'/>
//	  </user-component>
//	</user-component>
//
// The consequences, which are the whole point:
//
//   - `<scope>.show = true` mounts the guarded block, whose run-script loads the
//     form straight into the inner portal — ONE request, no wrapper round-trip;
//   - `<scope>.show = false` unmounts it, destroying the form and its `form`
//     scope completely (no stale state can leak into the next open);
//   - the loaded overlay binds to THIS closer (the portal seeds `closer` in the
//     content scope and ParamCloserProvided stops the responder from creating a
//     child closer), so closing/saving turns `<scope>.show` off — and any button
//     on the page can re-open a fresh form by turning it back on.
//
// The `form` scope is created once by the host, so the form's own re-renders
// (list-editor add/remove/sort, validation) never recreate it.
// The BASE scope variable names of the built-in form hosts. They are prefixed
// with `$presets` so they never collide with an application's own state on the
// shared `vars` object.
//
// These are the names of a TOP-LEVEL host only. A host rendered inside an
// overlay gets a unique name — see FormHostScope — so a button must ask the host
// for its reference (ShowExpr/OpenExpr, via the ctx helpers) instead of
// hardcoding one of these constants.
const (
	// DetailingEditScope hosts the edit form of a detailing page.
	DetailingEditScope = "$presetsEditing"
	// ListingNewScope hosts the create form of a listing.
	ListingNewScope = "$presetsCreating"
	// ListingItemEditScope / ListingItemDetailScope host the per-item edit and
	// detail overlays of a listing. There is ONE host for the whole listing (not
	// one per row): a row opens it with `<scope>.id = "<id>"; <scope>.show = true`
	// and the host's run-script loads that record.
	ListingItemEditScope   = "$presetsItemEditing"
	ListingItemDetailScope = "$presetsItemDetailing"
)

type FormHostBuilder struct {
	scope       string
	wrapsOpener bool
	show        bool
	vars        map[string]string
	portal      string
	load        *web.VueEventTagBuilder
	children    h.HTMLComponents
}

// FormHost starts a form host bound to the scope variable named scope (e.g.
// DetailingEditScope), loading the form into portal via the load event.
func FormHost(scope, portal string, load *web.VueEventTagBuilder) *FormHostBuilder {
	return &FormHostBuilder{scope: scope, portal: portal, load: load}
}

// WrapsOpener tells whether this host renders around its OPENER (the button)
// instead of around the content. A page needs it: the primary action is
// rendered by the LAYOUT into the app bar, outside the content, where the host's
// slot variable does not exist. In an event the button is inside the content the
// host wraps, so it already receives the variable from the scope.
func (b *FormHostBuilder) WrapsOpener(v bool) *FormHostBuilder {
	b.wrapsOpener = v
	return b
}

// OpenerWrapped reports the WrapsOpener setting (nil-safe).
func (b *FormHostBuilder) OpenerWrapped() bool {
	return b != nil && b.wrapsOpener
}

// Show sets the initial state: false (default) waits for a button to turn the
// scope on; true opens the form as soon as the host mounts.
func (b *FormHostBuilder) Show(v bool) *FormHostBuilder {
	b.show = v
	return b
}

// Children are rendered inside the host scope (so they can toggle it), before
// the guarded form block.
func (b *FormHostBuilder) Children(comps ...h.HTMLComponent) *FormHostBuilder {
	b.children = append(b.children, comps...)
	return b
}

// Var declares an extra reactive variable on the host scope (JS initializer),
// which the load event can read — e.g. the record id of the row being opened.
func (b *FormHostBuilder) Var(name, init string) *FormHostBuilder {
	if b.vars == nil {
		b.vars = map[string]string{}
	}
	b.vars[name] = init
	return b
}

// ShowExpr is the expression a button uses to open the form.
func (b *FormHostBuilder) ShowExpr() string {
	return fmt.Sprintf("%s.show = true", b.Ref())
}

// OpenExpr is the expression a row uses to open the host for a given record:
// it sets the vars (values are JS expressions) and turns the host on.
func (b *FormHostBuilder) OpenExpr(vars map[string]string) string {
	var s string
	for _, name := range sortedKeys(vars) {
		s += fmt.Sprintf("%s.%s = %s; ", b.Ref(), name, vars[name])
	}
	return s + b.ShowExpr()
}

// ScopeVarExpr references a variable of this host's state (e.g.
// "$presetsItemEditing.id"), for use inside the load event.
func (b *FormHostBuilder) ScopeVarExpr(name string) string {
	return b.Ref() + "." + name
}

func sortedKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Ref is the JS expression that addresses this host's state: the slot variable
// the host declares. It is per-render, so nested levels — a listing opened in a
// dialog, a record's detail opened from THAT listing — never share state.
func (b *FormHostBuilder) Ref() string {
	return b.scope
}

// init is the JS initializer of this host's state.
func (b *FormHostBuilder) init() string {
	s := fmt.Sprintf("{show:%v", b.show)
	for _, name := range sortedKeys(b.vars) {
		s += fmt.Sprintf(", %s:%s", name, b.vars[name])
	}
	return s + "}"
}

// block is the guarded form block: mounting it loads the form into the host's
// portal, unmounting it destroys the form and its `form` scope.
func (b *FormHostBuilder) block() h.HTMLComponent {
	ref := b.Ref()

	// the load event owns this host's closer: seed it in the portal content scope
	// and tell the responder not to create another one.
	load := b.load.Clone().
		Scope(vue.Var("{closer: "+ref+"}")).
		Query(ParamCloserProvided, "true")

	// The guarded block gets its own child `form` scope. It is a go-plaid-scope
	// (not a user-component): the latter renders its children inside a
	// <template v-slot>, and a <template> nested in another one is not compiled
	// by the runtime template parser — the content would render inert.
	return web.Scope(
		web.Portal().Name(b.portal).Scope("closer", js.Raw(ref)),
		web.RunScript(load.Go()),
	).
		FormInit().
		// `?.` because the state is assigned on mount: the guard must not throw
		// while the host itself is still being set up.
		Attr("v-if", ref+"?.show")
}

func (b *FormHostBuilder) Component() h.HTMLComponent {
	return FormHosts(b.children, b)
}

// FormHosts renders several hosts as ONE component wrapping children: their
// state variables are declared together and their guarded blocks are siblings.
//
// They must not be nested one inside the other. A host that keeps its state in a
// slot variable renders a `<template v-slot>`, and a `<template v-slot>` inside
// another one is NOT compiled by the runtime template parser — the inner content
// (the whole listing) would render inert.
func FormHosts(children h.HTMLComponents, hosts ...*FormHostBuilder) h.HTMLComponent {
	var (
		uc     = vue.UserComponent()
		blocks h.HTMLComponents
	)

	for _, host := range hosts {
		if host == nil {
			continue
		}
		uc.ScopeVar(host.scope, host.init())
		blocks = append(blocks, host.block())
	}

	return uc.AppendChild(append(children, blocks...)...)
}

func (b *FormHostBuilder) Write(ctx *h.Context) error {
	return b.Component().Write(ctx)
}

// CloserProvided reports whether the request comes from a form host that already
// owns the overlay's closer (see ParamCloserProvided).
func CloserProvided(ctx *web.EventContext) bool {
	return ctx != nil && ctx.R != nil && ctx.R.FormValue(ParamCloserProvided) == "true"
}

// ItemFormHosts are the hosts a listing renders once and its buttons reuse: the
// per-item edit/detail overlays every row opens (instead of each row carrying
// its own overlay plaid), and the create form of the New button.
//
// They are published on the context because their scope names are not fixed —
// a listing rendered inside an overlay gets unique ones (see FormHostScope).
type ItemFormHosts struct {
	Edit   *FormHostBuilder
	Detail *FormHostBuilder
	New    *FormHostBuilder
}

// OpenNewExpr is what the New button runs to open the listing's create form. It
// returns "" when the listing did not host it (the caller falls back).
func (h *ItemFormHosts) OpenNewExpr() string {
	if h == nil || h.New == nil {
		return ""
	}
	return h.New.ShowExpr()
}

// OpenEditExpr / OpenDetailExpr are what a row's click handler runs to open the
// listing's shared overlay for that record. They return "" when the listing did
// not host that overlay (the caller then falls back to its own event).
func (h *ItemFormHosts) OpenEditExpr(id string) string {
	if h == nil || h.Edit == nil {
		return ""
	}
	return h.Edit.OpenExpr(map[string]string{"id": strconv.Quote(id)})
}

func (h *ItemFormHosts) OpenDetailExpr(id string) string {
	if h == nil || h.Detail == nil {
		return ""
	}
	return h.Detail.OpenExpr(map[string]string{"id": strconv.Quote(id)})
}

// WithItemFormHosts publishes the listing's per-item hosts so the rows (data
// table cells, row menu items) rendered later can target them.
func WithItemFormHosts(ctx *web.EventContext, hosts *ItemFormHosts) {
	ctx.WithContextValue(ctxItemFormHosts, hosts)
}

// GetItemFormHosts returns the hosts published by the enclosing listing, or nil
// when the component is rendered outside one.
func GetItemFormHosts(ctx *web.EventContext) *ItemFormHosts {
	h, _ := ctx.ContextValue(ctxItemFormHosts).(*ItemFormHosts)
	return h
}

// WithDetailingEditHost publishes the detailing's edit host, so the Edit button
// (and any other button of that page) targets the right variable — the detailing
// keeps the well-known name only when it is not itself inside an overlay.
func WithDetailingEditHost(ctx *web.EventContext, host *FormHostBuilder) {
	ctx.WithContextValue(ctxDetailingEditHost, host)
}

// GetDetailingEditHost returns the edit host published by the enclosing
// detailing, or nil outside one.
func GetDetailingEditHost(ctx *web.EventContext) *FormHostBuilder {
	h, _ := ctx.ContextValue(ctxDetailingEditHost).(*FormHostBuilder)
	return h
}

// ShowExprOr is ShowExpr, falling back to the bare scope name when the host was
// not published (component rendered outside any host).
func (b *FormHostBuilder) ShowExprOr(scope string) string {
	if b == nil {
		return scope + ".show = true"
	}
	return b.ShowExpr()
}
