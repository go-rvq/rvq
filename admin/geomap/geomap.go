// Package geomap is a map of places in the admin — where a record's visitors
// were, where a revision was made from —: the
// rvq-geo-map Vue component, on Leaflet (BSD 2-Clause, embedded) and
// OpenStreetMap's tiles.
package geomap

import (
	"embed"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

//go:embed assets
var assets embed.FS

func pack(name string) web.ComponentsPack {
	b, err := assets.ReadFile("assets/" + name)
	if err != nil {
		panic(err)
	}
	return web.ComponentsPack(b)
}

// Install adds Leaflet and the component to the admin's pages.
func Install(b *presets.Builder) {
	b.ExtraAsset("/rvq-leaflet.css", "text/css", pack("leaflet.css"))
	b.ExtraAsset("/rvq-leaflet.js", "text/javascript", pack("leaflet.js"))
	b.ExtraAsset("/rvq-geomap.js", "text/javascript", pack("geomap.js"))
}

// Point is a place on the map: how many of a kind — "like", "comment" — came
// from it, and what its label says.
type Point struct {
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	Kind  string  `json:"kind"`
	Count int64   `json:"count"`
	Label string  `json:"label"`
}

// Map is the component drawing points.
func Map(points []Point) h.HTMLComponent {
	return h.Tag("rvq-geo-map").Attr(":points", h.JSONString(points))
}

// SizedMap is Map of the given height ("240px"), zoomed to zoom when it has a
// single point (a city: 10).
func SizedMap(points []Point, height string, zoom int) h.HTMLComponent {
	return h.Tag("rvq-geo-map").Attr(":points", h.JSONString(points)).Attr("height", height).Attr(":zoom", zoom)
}
