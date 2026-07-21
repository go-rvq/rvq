package admin

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

const (
	// nestedMenuPortal holds the selected organization's resource menu.
	nestedMenuPortal = "orgsNestedMenu"
	// nestedMenuEvent (re)loads the portal for a given organization.
	nestedMenuEvent = "orgs_load_nested_menu"
	// paramSelectedOrg carries the selected organization id.
	paramSelectedOrg = "__org__"
)

// NestedMenuPortal renders the portal that holds the selected organization's
// nested-resources menu. Its content is loaded server-side (via nestedMenuEvent)
// whenever the organization changes, so the menu can be filtered by the user's
// permissions on that specific organization — forbidden resources are never
// exposed. Loads on mount with the current vars.currentOrganization (empty until
// one is chosen).
func NestedMenuPortal(b *presets.Builder) presets.ComponentFunc {
	return func(ctx *web.EventContext) h.HTMLComponent {
		return web.Portal().Name(nestedMenuPortal).
			Loader(web.GET().EventFunc(nestedMenuEvent).Query(paramSelectedOrg, web.Var("vars."+OrganizationVar)))
	}
}

// NestedMenuContent builds the menu of the organization's child resources the
// current user is allowed to list (checked per that organization), with concrete
// links under /orgs/{orgID}/<group>/<resource>. Returns an empty node when no
// organization is selected or none is permitted.
func NestedMenuContent(b *presets.Builder, db *gorm.DB, ctx *web.EventContext, orgIDStr string) h.HTMLComponent {
	orgMb := OrganizacaoModel(b)
	if orgMb == nil || orgIDStr == "" {
		return h.Div()
	}
	orgModelID, err := orgMb.ParseRecordID(orgIDStr)
	if err != nil || orgModelID.IsZero() {
		return h.Div()
	}
	prefix := b.GetURIPrefix()

	// group the permitted children by menu-group segment, preserving order.
	var order []string
	groups := map[string][]*presets.ModelBuilder{}
	for _, c := range orgMb.Children() {
		// only resources the user may list in THIS organization (same rule the
		// main menu applies), so the menu never exposes forbidden resources.
		if c.Permissioner().Lister(ctx.R, orgModelID).Denied() {
			continue
		}
		g := c.MenuGroupName()
		if _, ok := groups[g]; !ok {
			order = append(order, g)
		}
		groups[g] = append(groups[g], c)
	}
	if len(order) == 0 {
		return h.Div()
	}

	var items h.HTMLComponents
	for _, g := range order {
		if g != "" {
			items = append(items, v.VListSubheader(h.Text(capitalize(g))))
		}
		for _, c := range groups[g] {
			seg := c.UriName()
			if g != "" {
				seg = g + "/" + seg
			}
			href := prefix + "/" + OrganizacaoURIName + "/" + orgIDStr + "/" + seg
			item := v.VListItem(v.VListItemTitle(h.Text(c.TTitlePlural(ctx.Context())))).Attr("href", href)
			if icon := c.GetMenuIcon(); icon != "" {
				item = item.PrependIcon(icon)
			}
			items = append(items, item)
		}
	}
	return v.VList(items...).Density(v.DensityCompact)
}

// configureNestedMenu registers the portal event and installs the portal below
// the selector.
func configureNestedMenu(b *presets.Builder, db *gorm.DB) {
	b.GetWebBuilder().RegisterEventFunc(nestedMenuEvent, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{
			Name: nestedMenuPortal,
			Body: NestedMenuContent(b, db, ctx, ctx.R.FormValue(paramSelectedOrg)),
		})
		return
	})
	b.AddMenuTopItemFunc("orgs_nested_menu", NestedMenuPortal(b))
}

// capitalize upper-cases the first rune of the menu-group segment for its header.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	if r[0] >= 'a' && r[0] <= 'z' {
		r[0] -= 'a' - 'A'
	}
	return string(r)
}
