package presets

import (
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	. "github.com/qor5/web/v3/tag"
	v "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type DialogBuilder struct {
	width             string
	height            string
	maxHeight         string
	targetPortal      string
	contentPortalName string
	scrollabled       bool
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

func (p *DialogBuilder) SetValidPortalName(portalName string) *DialogBuilder {
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

	comp = web.CloserScope(d, true)

	if p.rootWrap != nil {
		comp = p.rootWrap(comp)
	}

	return comp
}

func (p *DialogBuilder) Respond(ctx *web.EventContext, r *web.EventResponse, comp h.HTMLComponent) {
	for _, f := range GetRespondDialogHandlers(ctx) {
		f(p)
	}
	r.UpdatePortal(p.targetPortal, p.Component(comp))
}

func (b *Builder) dialog(ctx *web.EventContext, r *web.EventResponse, comp h.HTMLComponent, width string) {
	p := b.Dialog()
	if width != "" {
		p.SetWidth(width)
	}
	p.Respond(ctx, r, comp)
}

func (b *Builder) Dialog() *DialogBuilder {
	return Dialog(actions.Dialog.PortalName()).
		SetContentPortalName(actions.Dialog.ContentPortalName()).
		SetValidWidth(b.rightDrawerWidth)
}

func (b *Builder) DialogPortal(portal string) *DialogBuilder {
	return Dialog(portal).
		SetValidWidth(b.rightDrawerWidth)
}
