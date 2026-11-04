package presets

import (
	"net/url"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

// DetailingLink is a link to the detail of the record id of the model — its
// address in the href (copied, opened in a tab), its click opening it in a
// dialog (actions.Detailing) in the portal targetPortal (empty: the dialog
// portal of the layout). parents are the ids of the parents of the model;
// query goes in both — the tab of the detail to open (FieldsTabParam), say.
// Its content: the icon "mdi-open-in-new", unless children are given.
func (mb *ModelBuilder) DetailingLink(ctx *web.EventContext, id string, parents IDSlice, targetPortal string,
	query url.Values, children ...h.HTMLComponent) *h.HTMLTagBuilder {
	href := mb.Info().DetailingHref(id, parents...)
	if len(query) > 0 {
		href += "?" + query.Encode()
	}
	if targetPortal == "" {
		targetPortal = actions.Dialog.PortalName()
	}
	open := web.Plaid().
		URL(mb.Info().ListingHref(parents...)).
		EventFunc(actions.Detailing).
		Query(ParamID, id).
		Query(ParamOverlay, actions.Dialog).
		Query(ParamTargetPortal, targetPortal)
	for k, vs := range query {
		for _, v := range vs {
			open.Query(k, v)
		}
	}
	if len(children) == 0 {
		title := MustGetMessages(ctx.Context()).DetailingOpenLink
		children = []h.HTMLComponent{VIcon("mdi-open-in-new").Size(SizeSmall)}
		return h.A(children...).Href(href).Attr("title", title).Attr("aria-label", title).
			Attr("@click.prevent.stop", open.Go()).Class("rvq-detailing-link text-decoration-none")
	}
	return h.A(children...).Href(href).Attr("@click.prevent.stop", open.Go()).Class("rvq-detailing-link text-decoration-none")
}
