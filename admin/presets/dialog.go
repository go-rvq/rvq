package presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/web/tag"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

type DialogBuilder struct {
	width             string
	height            string
	maxHeight         string
	targetPortal      string
	contentPortalName string
	scrollabled       bool
	closerProvided    bool
	wrap              func(comp *v.VDialogBuilder)
	rootWrap          func(comp h.HTMLComponent) h.HTMLComponent
}

func Dialog(portalName string) *DialogBuilder {
	return &DialogBuilder{targetPortal: portalName}
}

func (p *DialogBuilder) Width() string {
	return p.width
}

func (p *DialogBuilder) SetWidth(width string) *DialogBuilder {
	p.width = width
	return p
}

func (p *DialogBuilder) SetValidWidth(width string) *DialogBuilder {
	if width != "" {
		p.width = width
	}
	return p
}

func (p *DialogBuilder) Height() string {
	return p.height
}

func (p *DialogBuilder) SetHeight(height string) *DialogBuilder {
	p.height = height
	return p
}

func (p *DialogBuilder) SetValidHeight(height string) *DialogBuilder {
	if height != "" {
		p.height = height
	}
	return p
}

func (p *DialogBuilder) MaxHeight() string {
	return p.maxHeight
}

func (p *DialogBuilder) SetMaxHeight(maxHeight string) *DialogBuilder {
	p.maxHeight = maxHeight
	return p
}

func (p *DialogBuilder) SetValidMaxHeight(height string) *DialogBuilder {
	if height != "" {
		p.maxHeight = height
	}
	return p
}

func (p *DialogBuilder) TargetPortal() string {
	return p.targetPortal
}

func (p *DialogBuilder) SetTargetPortal(portalName string) *DialogBuilder {
	p.targetPortal = portalName
	return p
}

func (p *DialogBuilder) SetValidTargetPortalName(portalName string) *DialogBuilder {
	if portalName != "" {
		p.targetPortal = portalName
	}
	return p
}

func (p *DialogBuilder) ContentPortalName() string {
	return p.contentPortalName
}

func (p *DialogBuilder) SetContentPortalName(contentPortalName string) *DialogBuilder {
	p.contentPortalName = contentPortalName
	return p
}

func (p *DialogBuilder) ValidContentPortalName(portalName string) *DialogBuilder {
	if portalName != "" {
		p.contentPortalName = portalName
	}
	return p
}

func (p *DialogBuilder) Wrap(wrap func(comp *v.VDialogBuilder)) *DialogBuilder {
	p.wrap = wrap
	return p
}

func (p *DialogBuilder) RootWrap(wrap func(comp h.HTMLComponent) h.HTMLComponent) *DialogBuilder {
	p.rootWrap = wrap
	return p
}

func (p *DialogBuilder) SetScrollable(s bool) *DialogBuilder {
	p.scrollabled = s
	return p
}

// SetCloserProvided marks that the caller already provides `closer` in the
// content's scope (a FormHost owns it). The dialog then binds to that closer
// instead of creating its own child one — so turning it off destroys the host's
// form, and turning it back on re-opens a fresh one.
func (p *DialogBuilder) SetCloserProvided(v bool) *DialogBuilder {
	p.closerProvided = v
	return p
}

// closerScope wraps comp in a scope. When the caller already provides the closer
// (a FormHost owns it), the scope is created WITHOUT one — the content still gets
// its regular form/locals scope, but binds to the host's closer, so turning that
// off destroys the host's form instead of just hiding a child overlay.
func (p *DialogBuilder) closerScope(comp h.HTMLComponent) h.HTMLComponent {
	if p.closerProvided {
		sb, ok := comp.(*web.ScopeBuilder)
		if !ok {
			sb = web.Scope(comp)
		}
		// Re-provide the closer already in scope (the host's one, a reactive proxy
		// → go-plaid-scope adopts it as is) instead of creating a child closer, and
		// keep the usual slot vars for the content.
		return sb.Closer().Attr(":closer", "closer")
	}
	return web.CloserScope(comp, true)
}

func (p *DialogBuilder) Component(comp h.HTMLComponent) h.HTMLComponent {
	if fvc := FirstValidComponent(comp); fvc != nil {
		switch t := fvc.(type) {
		case *v.VCardBuilder:
			t.SetAttr("style", "max-height:inherit")
		}
	}

	d := v.VDialog(comp).
		Attr("v-model", "closer.show").
		Fullscreen("closer.fullscreen")

	if p.width != "" {
		d.Width(web.Var("closer.fullscreen ? '100%' : " + p.width))
	}

	if p.height != "" {
		d.Height(web.Var("closer.fullscreen ? '100%' : " + p.height))
	}

	if p.height != "" {
		d.Height(web.Var("closer.fullscreen ? '100%' : " + p.height))
	}

	if p.maxHeight != "" {
		d.MaxHeight(web.Var("closer.fullscreen ? null : " + p.maxHeight))
	}

	if p.scrollabled {
		d.Scrollable(true)
	}

	if p.wrap != nil {
		p.wrap(d)
	}

	comp = p.closerScope(d)

	if p.rootWrap != nil {
		comp = p.rootWrap(comp)
	}

	return comp
}

func (p *DialogBuilder) Respond(ctx *web.EventContext, r *web.EventResponse, comp h.HTMLComponent) {
	// a form host owning the closer asks (via the request) not to create another
	// one, so its `<scope>.show` really controls this overlay.
	if CloserProvided(ctx) {
		p.closerProvided = true
	}

	if ac, _ := web.Unscoped(comp).(vx.VXAdvancedCloseCardTagger); ac != nil {
		ac.SetVModel("closer.show")
		if acd, ok := ac.(vx.VXAdvancedExpandCloseCardTagger); ok {
			acd.SetDialogOrDrawer(true)
			if p.width != "" {
				acd.SetWidth(p.width)
			}
		}
	} else {
		for _, f := range GetRespondDialogHandlers(ctx) {
			f(p)
		}
		comp = p.Component(comp)
	}
	r.UpdatePortal(p.targetPortal, p.closerScope(comp))
}

func (b *Builder) dialog(ctx *web.EventContext, r *web.EventResponse, comp h.HTMLComponent, width string) {
	p := b.Dialog()
	if width != "" {
		p.SetWidth(width)
	}
	p.Respond(ctx, r, comp)
}

func (b *Builder) Dialog(width ...string) *DialogBuilder {
	w := b.rightDrawerWidth
	if len(width) > 0 && width[0] != "" {
		w = width[0]
	}
	return Dialog(actions.Dialog.PortalName()).
		SetContentPortalName(actions.Dialog.ContentPortalName()).
		SetValidWidth(w)
}

func (b *Builder) DialogPortal(portal string, width ...string) *DialogBuilder {
	w := b.rightDrawerWidth
	if len(width) > 0 && width[0] != "" {
		w = width[0]
	}
	if portal != "" {
		return Dialog(portal).
			SetValidWidth(w)
	}
	return b.Dialog(b.rightDrawerWidth).SetValidWidth(w)
}
