// Package originui draws an origin (origin.Origin) in the admin: what is known
// of it and, its coordinates known, the place on a map (geomap).
package originui

import (
	"strconv"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/geomap"
	"github.com/go-rvq/rvq/admin/origin"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// Labels are the words of Body, in the admin's language.
type Labels struct {
	IP, Browser, Place, Coordinates string
	// NoMap is said in place of the map, its coordinates not known.
	NoMap string
}

// Body is what is known of o — the address, the browser, the place, its
// coordinates — and the place on a map, its coordinates known.
func Body(l Labels, o origin.Origin) h.HTMLComponent {
	row := func(label, value string) h.HTMLComponent {
		if value == "" {
			value = "—"
		}
		return h.Tr(h.Th(label).Style("text-align:left;white-space:nowrap"), h.Td(h.Text(value)))
	}
	var coords string
	var points []geomap.Point
	if o.Latitude != nil && o.Longitude != nil {
		coords = FormatCoords(*o.Latitude, *o.Longitude)
		label := o.Place()
		if label == "" {
			label = o.IP
		}
		points = append(points, geomap.Point{Lat: *o.Latitude, Lng: *o.Longitude, Kind: "place", Count: 1, Label: label})
	}
	body := h.HTMLComponents{
		v.VTable(h.Tbody(
			row(l.IP, o.IP),
			row(l.Browser, o.UserAgent),
			row(l.Place, o.Place()),
			row(l.Coordinates, coords),
		)).Density(v.DensityCompact),
	}
	if len(points) > 0 {
		body = append(body, h.Div(geomap.SizedMap(points, "320px", 10)).Class("mt-4"))
	} else {
		body = append(body, h.Div(h.Text(l.NoMap)).Class("mt-4 text-medium-emphasis text-body-2"))
	}
	return body
}

// Link is o on one line, a link that runs open — the script that opens the
// action showing it —, with attr (a data- attribute naming it) set to id; a
// dash when nothing is known.
func Link(o origin.Origin, open, attr, id string) h.HTMLComponent {
	text := o.String()
	if text == "" {
		return h.Span("—").Class("text-medium-emphasis")
	}
	return h.A(
		v.VIcon("mdi-map-marker-outline").Size(v.SizeSmall).Class("me-1"),
		h.Text(text),
	).Href("javascript:void(0)").
		Attr("title", o.UserAgent).
		Attr(attr, id).
		Attr("@click.stop", open)
}

// FormatCoords is a latitude and a longitude, to four places.
func FormatCoords(lat, lng float64) string {
	return strconv.FormatFloat(lat, 'f', 4, 64) + ", " + strconv.FormatFloat(lng, 'f', 4, 64)
}
