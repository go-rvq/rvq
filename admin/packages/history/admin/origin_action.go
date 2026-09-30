package admin

import (
	"context"
	"strconv"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/geomap"
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
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
	text := rev.Origin.String()
	if text == "" {
		return h.Span("—").Class("text-medium-emphasis")
	}
	open := web.Plaid().
		EventFunc(actions.Action).
		Query(presets.ParamID, rev.Hash.String()).
		Query(presets.ParamAction, ActionOrigin).
		Query(presets.ParamOverlay, actions.Dialog).
		// as the detail's menu opens its actions: the revisions', the id in
		// the query
		URL(child.Info().ListingHrefCtx(ctx)).
		Go()
	return h.A(
		v.VIcon("mdi-map-marker-outline").Size(v.SizeSmall).Class("me-1"),
		h.Text(text),
	).Href("javascript:void(0)").
		Attr("title", rev.Origin.UserAgent).
		Attr("data-revision-origin", rev.Hash.String()).
		Attr("@click.stop", open)
}

// originBody is what is known of an origin — the address, the browser, the
// place, its coordinates — and the place on a map, its coordinates known.
func originBody(msgs *Messages, o histmodels.Origin) h.HTMLComponent {
	row := func(label, value string) h.HTMLComponent {
		if value == "" {
			value = "—"
		}
		return h.Tr(h.Th(label).Style("text-align:left;white-space:nowrap"), h.Td(h.Text(value)))
	}
	var coords string
	var points []geomap.Point
	if o.Latitude != nil && o.Longitude != nil {
		coords = formatCoords(*o.Latitude, *o.Longitude)
		label := o.Place()
		if label == "" {
			label = o.IP
		}
		points = append(points, geomap.Point{Lat: *o.Latitude, Lng: *o.Longitude, Kind: "place", Count: 1, Label: label})
	}
	body := h.HTMLComponents{
		v.VTable(h.Tbody(
			row(msgs.OriginIP, o.IP),
			row(msgs.OriginBrowser, o.UserAgent),
			row(msgs.OriginPlace, o.Place()),
			row(msgs.OriginCoordinates, coords),
		)).Density(v.DensityCompact),
	}
	if len(points) > 0 {
		body = append(body, h.Div(geomap.SizedMap(points, "320px", 10)).Class("mt-4"))
	} else {
		body = append(body, h.Div(h.Text(msgs.OriginNoMap)).Class("mt-4 text-medium-emphasis text-body-2"))
	}
	return body
}

func formatCoords(lat, lng float64) string {
	return strconv.FormatFloat(lat, 'f', 4, 64) + ", " + strconv.FormatFloat(lng, 'f', 4, 64)
}
