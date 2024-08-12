package presets

import (
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type DialogBuilder struct {
	width             string
	height            string
	targetPortal      string
	contentPortalName string
	wrap              func(comp *v.VDialogBuilder)
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

func (p *DialogBuilder) Respond(r *web.EventResponse, comp h.HTMLComponent) {
	if fvc := web.FirstValidComponent(comp); fvc != nil {
		switch t := fvc.(type) {
		case *v.VCardBuilder:
			t.SetAttr("style", "height:inherit")
		}
	}
	if p.contentPortalName != "" {
		comp = web.Portal(comp).Name(p.contentPortalName).Style("height: inherit")
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

	if p.wrap != nil {
		p.wrap(d)
	}

	r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{
		Name: p.targetPortal,
		Body: web.CloserScope(d, true),
	})
}

func (b *Builder) dialog(r *web.EventResponse, comp h.HTMLComponent, width string) {
	p := b.Dialog()
	if width != "" {
		p.SetWidth(width)
	}
	p.Respond(r, comp)
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
