package web

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web/js"
)

type PortalBuilder struct {
	tag    *h.HTMLTagBuilder
	form   string
	locals string
	scope  js.Object
}

func Portal(children ...h.HTMLComponent) (r *PortalBuilder) {
	r = &PortalBuilder{
		tag: h.Tag("go-plaid-portal").Children(children...),
		scope: js.Object{
			"presetsListing": js.Raw("presetsListing"),
			// The refresh hooks travel with the scope the way presetsListing
			// does (see Plaid and presets.PostSaveScript): a portal that dropped
			// them would hand what it renders the root's EMPTY list, and a save
			// made in there — an EDIT opened from a detail, say — would refresh
			// nothing its openers show. A form host replaces it with its own,
			// longer list (presets.FormHostBuilder).
			"onSaveCallbacks": js.Raw("onSaveCallbacks"),
			//	"presetsDetailing": js.Raw("presetsDetailing"),
			//	"presetsCreating":  js.Raw("presetsCreating"),
			//	"presetsEditing":   js.Raw("presetsEditing"),
		},
	}
	r.Visible("true")
	return
}

func (b *PortalBuilder) Inline() *PortalBuilder {
	b.tag.Style("display: inline-block")
	return b
}

func (b *PortalBuilder) Class(names ...string) *PortalBuilder {
	b.tag.Class(names...)
	return b
}

func (b *PortalBuilder) Scope(name string, value any) *PortalBuilder {
	b.scope[name] = value
	return b
}

func (b *PortalBuilder) Raw(v bool) (r *PortalBuilder) {
	b.tag.SetAttr("raw", v)
	return b
}

func (b *PortalBuilder) Loader(v *VueEventTagBuilder) (r *PortalBuilder) {
	b.tag.SetAttr(":loader", v.String())
	return b
}

func (b *PortalBuilder) LoaderString(v string) (r *PortalBuilder) {
	b.tag.SetAttr(":loader", v)
	return b
}

func (b *PortalBuilder) Content(v string) (r *PortalBuilder) {
	b.tag.SetAttr(":content", h.JSONString(v))
	return b
}

func (b *PortalBuilder) Visible(v string) (r *PortalBuilder) {
	b.tag.Attr(":visible", v)
	return b
}

func (b *PortalBuilder) Name(v string) (r *PortalBuilder) {
	b.tag.Attr("portal-name", v)
	return b
}

func (b *PortalBuilder) Form(v string) (r *PortalBuilder) {
	b.form = v
	return b
}

func (b *PortalBuilder) AutoReloadInterval(v interface{}) (r *PortalBuilder) {
	b.tag.Attr(":auto-reload-interval", v)
	return b
}

func (b *PortalBuilder) Style(v string) (r *PortalBuilder) {
	b.tag.Style(v)
	return b
}

func (b *PortalBuilder) Children(comps ...h.HTMLComponent) (r *PortalBuilder) {
	b.tag.Children(comps...)
	return b
}

func (b *PortalBuilder) LoadWhenParentVisible() (r *PortalBuilder) {
	b.Visible("parent.isVisible")
	return b
}

func (b *PortalBuilder) ParentForceUpdateAfterLoaded() (r *PortalBuilder) {
	b.tag.Attr(":after-loaded", "parent.forceUpdate")
	return b
}

func (b *PortalBuilder) Write(ctx *h.Context) (err error) {
	if b.form == "" {
		b.form = "form"
	}
	if b.locals == "" {
		b.locals = "locals"
	}
	b.tag.Attr(":form", b.form)
	b.tag.Attr(":locals", b.locals)
	if len(b.scope) > 0 {
		b.tag.Attr(":scope", b.scope.String())
	}
	return b.tag.Write(ctx)
}
