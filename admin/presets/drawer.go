package presets

import (
	"strings"

	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
)

type Drawer struct {
	location   string
	width      string
	portalName string
	safeClose  bool
	scrollable bool
	rootWrap   func(comp h.HTMLComponent) h.HTMLComponent
}

func NewDrawer(width string, portalName string) *Drawer {
	return &Drawer{location: v.LocationRight, width: width, portalName: portalName}
}

func (p *Drawer) Location() string {
	return p.location
}

func (p *Drawer) SetLocation(location string) *Drawer {
	p.location = location
	return p
}

func (p *Drawer) Left() *Drawer {
	return p.SetLocation(v.LocationLeft)
}

func (p *Drawer) Top() *Drawer {
	return p.SetLocation(v.LocationTop)
}

func (p *Drawer) Bottom() *Drawer {
	return p.SetLocation(v.LocationBottom)
}

func (p *Drawer) Right() *Drawer {
	return p.SetLocation(v.LocationRight)
}

func (p *Drawer) Start() *Drawer {
	return p.SetLocation(v.LocationStart)
}

func (p *Drawer) End() *Drawer {
	return p.SetLocation(v.LocationEnd)
}

func (p *Drawer) Width() string {
	return p.width
}

func (p *Drawer) SetWidth(width string) *Drawer {
	p.width = width
	return p
}

func (p *Drawer) SetValidWidth(width string) *Drawer {
	if width != "" {
		p.width = width
	}
	return p
}

func (p *Drawer) PortalName() string {
	return p.portalName
}

func (p *Drawer) SetPortalName(portalName string) *Drawer {
	p.portalName = portalName
	return p
}

func (p *Drawer) SetValidPortalName(portalName string) *Drawer {
	if portalName != "" {
		p.portalName = portalName
	}
	return p
}

func (p *Drawer) SafeClose() bool {
	return p.safeClose
}

func (p *Drawer) SetSafeClose(safeClose bool) *Drawer {
	p.safeClose = safeClose
	return p
}

func (p *Drawer) Scrollable() bool {
	return p.scrollable
}

func (p *Drawer) SetScrollable(scrollable bool) *Drawer {
	p.scrollable = scrollable
	return p
}

func (p *Drawer) RootWrap(wrap func(comp h.HTMLComponent) h.HTMLComponent) *Drawer {
	p.rootWrap = wrap
	return p
}

func (p *Drawer) Respond(r *web.EventResponse, comp h.HTMLComponent) {
	if ac, _ := web.Unscoped(comp).(vx.VXAdvancedCloseCardTagger); ac != nil {
		ac.SetVModel("closer.show")
	} else {
		drawer := v.VNavigationDrawer(
			// web.GlobalEvents().Attr("@keyup.esc", varName+" = false"),
			comp,
		).
			// Attr("@input", "plaidForm.dirty && vars.presetsRightDrawer == false && !confirm('You have unsaved changes on this form. If you close it, you will lose all unsaved changes. Are you sure you want to close it?') ? vars.presetsRightDrawer = true: vars.presetsRightDrawer = $event"). // remove because drawer plaidForm has to be reset when UpdateOverlayContent
			Attr("v-model", "closer.show").
			Location(p.location).
			Temporary(true).
			// Fixed(true).
			RawWidth(`closer.fullscreen ? null : `+h.JSONString(p.width)).
			RawClass(`closer.fullscreen ? 'v-navigation-drawer--fullscreen' : null`).
			Attr(":height", `"100%"`)

		if p.scrollable {
			drawer.Class("v-navigation-drawer--scrollable")
			drawer.Style("position:fixed; top:0; overflow-y:scroll;")
			switch p.location {
			case v.LocationLeft:
				drawer.Style(`left:0`)
			case v.LocationRight:
				drawer.Style(`right:0`)
			}
		}

		comp = drawer
	}

	comp = web.CloserScope(
		comp,
		true,
	)

	if p.rootWrap != nil {
		comp = p.rootWrap(comp)
	}

	r.UpdatePortal(p.portalName, comp)
}

func (p *Builder) Drawer(drawerMode actions.OverlayMode) *Drawer {
	return NewDrawer(p.rightDrawerWidth, drawerMode.PortalName()).
		SetLocation(strings.ToLower(strings.TrimRight(drawerMode.String(), "Drawer")))
}
