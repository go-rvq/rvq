package profile

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/dashboard"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

// DashboardField is the virtual field on the top of the profile that holds
// its widgets (Dashboard).
const DashboardField = "Dashboard"

// Dashboard makes the profile a dashboard: on its top, a virtual field holds
// the widgets ids of d (Panel) — those of the user's own, each one loaded for
// the user of the request, the ones they may see. Called again, the ids are
// added. After Install.
func (c *Builder) Dashboard(d *dashboard.Builder, ids ...string) *Builder {
	first := len(c.widgets) == 0
	c.dash = d
	c.widgets = append(c.widgets, ids...)
	if !first {
		return c
	}
	det := c.mb.Detailing()
	det.Prepend(DashboardField)
	det.Field(DashboardField).ComponentFunc(func(_ *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Div(
			h.Div(h.Text(GetMessages(ctx.Context()).ProfileDashboard)).Class("text-h6 mb-3"),
			c.dash.Panel(ctx, c.widgets...),
		).Class("mb-6").Attr("data-profile-dashboard", true)
	})
	return c
}
