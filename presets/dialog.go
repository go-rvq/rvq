package presets

import (
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type PortalDialog struct {
	width             string
	portalName        string
	contentPortalName string
}

func NewPortalDialog(width string, portalName string) *PortalDialog {
	return &PortalDialog{width: width, portalName: portalName, contentPortalName: portalName + "Content"}
}

func (p *PortalDialog) Width() string {
	return p.width
}

func (p *PortalDialog) SetWidth(width string) *PortalDialog {
	p.width = width
	return p
}

func (p *PortalDialog) SetValidWidth(width string) *PortalDialog {
	if width != "" {
		p.width = width
	}
	return p
}

func (p *PortalDialog) PortalName() string {
	return p.portalName
}

func (p *PortalDialog) SetPortalName(portalName string) *PortalDialog {
	p.portalName = portalName
	return p
}

func (p *PortalDialog) ContentPortalName() string {
	return p.contentPortalName
}

func (p *PortalDialog) SetContentPortalName(contentPortalName string) *PortalDialog {
	p.contentPortalName = contentPortalName
	return p
}

func (p *PortalDialog) Respond(r *web.EventResponse, comp h.HTMLComponent) {
	varName := "vars.presetsDialog"
	if p.portalName != actions.Dialog.PortalName() {
		varName = varName + p.portalName
	}

	r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{
		Name: p.portalName,
		Body: web.Scope(
			v.VDialog(
				web.Portal(comp).Name(p.contentPortalName),
			).
				Attr("v-model", varName).
				Width(p.width),
		).VSlot("{ form }"),
	})
	r.RunScript = "setTimeout(function(){ " + varName + " = true }, 100)"
}

func (b *Builder) dialog(r *web.EventResponse, comp h.HTMLComponent, width string) {
	p := b.Dialog()
	if width != "" {
		p.SetWidth(width)
	}
	p.Respond(r, comp)
}

func (b *Builder) Dialog() *PortalDialog {
	return NewPortalDialog(b.rightDrawerWidth, actions.Dialog.PortalName()).
		SetContentPortalName(actions.Dialog.ContentPortalName())
}
