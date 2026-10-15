// Package portal is the Go API of <rvq-portal> (js/portal; its README).
package portal

import (
	h "github.com/go-rvq/htmlgo"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// Builder is <rvq-portal> (js/portal: @rvq/portal; its README) — not
// corejs's go-plaid-portal (web.Portal): a place of the page whose content
// is a Vue template given at run time, compiled and shown there, again each
// time it changes.
//
// The content (Model, bound by v-model) is:
//
//   - a template (a string), or {template, options: {data, setup, mounted,
//     unmounted, resetScroll}} — options.data variables of that template,
//     beside Data's;
//   - "" (or {template: ""}): nothing shown;
//   - null/undefined: the initial content again — the default slot, or
//     nothing.
//
// Mounted, it shows the initial content: the model's, when a template, else
// the default slot.
//
// A template of the content sees only its variables — Data's, its
// options.data, Vue and window — and its portal (vx$portal); not the scope of
// the page where the portal is (vars, locals, form, plaid…). It gets that
// scope two ways:
//
//   - Data: variables passed on by name (`{vars, locals}`);
//   - the default slot, scoped (Slot): written in the page, it sees all the
//     page's scope, its props the portal's variables (Data). It is the
//     initial content, and it stays the default slot of any content given: a
//     template of the content shows it where it says <slot></slot>.
//
// Its content set by an event (plaid): the response's body — asked with
// loadPortalBody(true), so it is not taken for the page's — into Model's
// expression:
//
//	plaid().vars(vars).locals(locals).eventFunc("x").
//	    loadPortalBody(true).go().then(r => vars.x = r.body)
type Builder struct {
	v.VTagBuilder[*Builder]
}

// Portal is a <rvq-portal>; children, its default slot as written (a
// template of the page: see Slot and DefaultSlot).
func Portal(children ...h.HTMLComponent) *Builder {
	return v.VTag(&Builder{}, "rvq-portal", children...)
}

// DefaultSlot is the initial content as a template of the portal (v-pre):
// it sees only the portal's variables (Data), as a content given would — not
// the scope of the page.
func (b *Builder) DefaultSlot(children ...h.HTMLComponent) *Builder {
	return b.Children(h.Template(children...).Attr("v-pre", true))
}

// Slot is the default slot, scoped: written in the page — all its scope —,
// its props (scope, a destructuring as "{x, y}"; "" none) the portal's
// variables (Data). It is the initial content, and the default slot of any
// content given (<slot></slot> in its template).
func (b *Builder) Slot(scope string, children ...h.HTMLComponent) *Builder {
	t := h.Template(children...)
	if scope == "" {
		t.Attr("#default", true)
	} else {
		t.Attr("#default", scope)
	}
	return b.Children(t)
}

// Model is the content: a JS expression (v-model) — a template, {template,
// options}, "" (empty), null (the initial content). The page sets it (an
// event's body: plaid's loadPortalBody), and it is shown.
func (b *Builder) Model(expr string) *Builder {
	return b.Attr("v-model", expr)
}

// Loading is a JS expression (v-model:loading): true, the slot "loading"
// shows in the place of the content; a content set makes it false.
func (b *Builder) Loading(expr string) *Builder {
	return b.Attr("v-model:loading", expr)
}

// Data are the variables of the content's templates — and the props of the
// default slot —: a JS expression of an object (":data"), `{vars, locals}`.
func (b *Builder) Data(expr string) *Builder {
	return b.Attr(":data", expr)
}

// Name registers the portal under name in the app
// (app.config.globalProperties.vx$portals[name]): its content (.content, the
// model's ref) and its elements.
func (b *Builder) Name(name string) *Builder {
	return b.Attr("name", name)
}

// Tag is the element the content is put in (a div): the one that takes the
// portal's id, class, style and attrs. None: no element of its own.
func (b *Builder) Tag(name string) *Builder {
	return b.Attr("tag", name)
}

// Source is the content of the portal from the server: token is the opaque
// key of it (a Source's Token) — the page names no address; the runtime of
// the site (rvq-portal's bundle) asks for it at its fixed path when the
// portal comes into view.
func (b *Builder) Source(token string) *Builder {
	return b.Attr("source", token)
}
