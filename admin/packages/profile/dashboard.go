package profile

import (
	"context"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/dashboard"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

const (
	// PanelProfile is the panel of the profile (Dashboard): its widgets on
	// its top.
	PanelProfile = "profile"
	// GroupProfile is the group of the widgets of the profile: of the user's
	// own (another package adds its by d.GetGroup(GroupProfile)).
	GroupProfile = "profile"
	// DashboardField is the virtual field on the top of the profile that
	// holds its panel (Dashboard).
	DashboardField = "Dashboard"
)

// Dashboard makes the profile a dashboard of d: the panel PanelProfile on its
// top, in a virtual field, with the widgets of the group GroupProfile, which
// it makes — other packages put theirs in the panel (d.GetPanel) or in the
// group (d.GetGroup). Each user customizes it, as any panel. After Install.
func (c *Builder) Dashboard(d *dashboard.Builder) *dashboard.Panel {
	title := func(ctx context.Context) string { return GetMessages(ctx).ProfileDashboard }
	g := d.Group(GroupProfile).Title(title).Icon("mdi-account-circle-outline")
	pn := d.Panel(PanelProfile).Title(title).Group(g)
	det := c.mb.Detailing()
	det.Prepend(DashboardField)
	det.Field(DashboardField).ComponentFunc(func(_ *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Div(pn.Body(ctx)).Class("mb-6").Attr("data-profile-dashboard", true)
	})
	return pn
}
