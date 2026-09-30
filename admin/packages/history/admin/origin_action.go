package admin

import (
	"context"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/origin/originui"
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// ActionOrigin is the revision's action — of its detail (DetailingBuilder), in
// the admin's dialog — that shows where its author made it from: the address,
// the browser, the place and — its coordinates known — the place on a map. The
// listing's Origin column and the detail's Origin field open it too.
const ActionOrigin = "origin"

// installOriginAction adds the Origin action to the revisions (child), and
// the Origin field of their detail.
func (mh *ModelHistory) installOriginAction(child *presets.ModelBuilder) {
	d := child.Detailing()
	d.Action(ActionOrigin).
		SetI18nLabel(func(c context.Context) string { return getMessages(c).Origin }).
		Icon("mdi-map-marker-outline").
		// what it shows: the detail's actions open in the admin's dialog
		ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
			hash, err := decodeHash(id)
			if err != nil {
				return nil, err
			}
			rev, err := mh.Revision(mh.recordKey(ctx), hash)
			if err != nil {
				return nil, err
			}
			return originBody(getMessages(ctx.Context()), rev.Origin), nil
		})

	d.Field("Origin").
		SetI18nLabel(func(c context.Context) string { return getMessages(c).Origin }).
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			rev := field.Obj.(*histmodels.Revision)
			return vx.VXReadonlyField(originLink(child, rev, ctx)).Label(field.Label)
		})
}

// originLink is a revision's origin on one line, opening the Origin action; a
// dash when nothing is known.
func originLink(child *presets.ModelBuilder, rev *histmodels.Revision, ctx *web.EventContext) h.HTMLComponent {
	open := web.Plaid().
		EventFunc(actions.Action).
		Query(presets.ParamID, rev.Hash.String()).
		Query(presets.ParamAction, ActionOrigin).
		Query(presets.ParamOverlay, actions.Dialog).
		// as the detail's menu opens its actions: the revisions', the id in
		// the query
		URL(child.Info().ListingHrefCtx(ctx)).
		Go()
	return originui.Link(rev.Origin, open, "data-revision-origin", rev.Hash.String())
}

// originBody is what is known of an origin and the place on a map
// (originui.Body).
func originBody(msgs *Messages, o histmodels.Origin) h.HTMLComponent {
	return originui.Body(originui.Labels{
		IP: msgs.OriginIP, Browser: msgs.OriginBrowser, Place: msgs.OriginPlace,
		Coordinates: msgs.OriginCoordinates, NoMap: msgs.OriginNoMap,
	}, o)
}
