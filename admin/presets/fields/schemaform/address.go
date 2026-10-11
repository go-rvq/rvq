//go:build !formscore

package schemaform

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

// MapsKeyFunc is the Google Maps key (Maps JavaScript + Places) an address is
// found and shown with, for the request; "" leaves the address only typed.
type MapsKeyFunc func(ctx *web.EventContext) string

// MapsKey is the key of every Builder that has not its own (Builder.MapsKey):
// an application sets it once (hermon: the one of its SEO settings).
var MapsKey MapsKeyFunc

// MapsKey sets the Google Maps key of the addresses of this builder's forms
// (see MapsKeyFunc); without it, the package's MapsKey.
func (b *Builder) MapsKey(f MapsKeyFunc) *Builder {
	b.mapsKey = f
	return b
}

// mapsKeyOf is the key of ctx: the builder's, else the package's.
func (b *Builder) mapsKeyOf(ctx *web.EventContext) string {
	f := MapsKey
	if b != nil && b.mapsKey != nil {
		f = b.mapsKey
	}
	if f == nil || ctx == nil {
		return ""
	}
	return f(ctx)
}

// AddressComponentFunc is "address": the address typed and found, the place on
// a map below it (<vx-address-field>).
func AddressComponentFunc(c *Context) h.HTMLComponent {
	msgs := c.Messages()
	readOnly := c.Form != nil && (c.Form.ReadOnly || !c.Form.Mode.IsWrite())
	return h.Tag("vx-address-field").
		Attr("label", c.Label()).
		Attr("hint", c.Hint()).
		Attr("placeholder", c.Placeholder()).
		Attr("api-key", c.Builder.mapsKeyOf(c.Event)).
		Attr("no-key-text", msgs.AddressNoKey).
		Attr("load-error-text", msgs.AddressLoadError).
		Attr(":required", c.Field.Required()).
		Attr(":readonly", readOnly).
		Attr("data-schemaform-address", c.Path).
		Attr("v-model", c.Value)
}

// AddressDisplayFunc shows an address: its words, a link to the place on
// Google Maps when it has its coordinates.
func AddressDisplayFunc(c *Context) h.HTMLComponent {
	if isEmpty(c.Data) {
		return emptyDisplay(c)
	}
	m, _ := c.Data.(map[string]any)
	if m == nil {
		return h.Text(displayText(c.Data))
	}
	text := displayText(m["formatted"])
	lat, okLat := coordinate(m["lat"])
	lng, okLng := coordinate(m["lng"])
	if !okLat || !okLng {
		return h.Span(text)
	}
	href := fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%s,%s",
		strconv.FormatFloat(lat, 'f', -1, 64), strconv.FormatFloat(lng, 'f', -1, 64))
	if id := displayText(m["placeId"]); id != "" {
		href += "&query_place_id=" + url.QueryEscape(id)
	}
	return h.Span("").Children(
		h.Text(text+" "),
		h.A(h.I("").Class("mdi mdi-map-marker-outline")).Href(href).Target("_blank").Attr("rel", "noopener").
			Attr("title", c.Messages().AddressOnMap).Attr("@click.stop", ""),
	)
}

// coordinate is v as a float: a number decoded, or its text.
func coordinate(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
		return f, err == nil
	}
	return 0, false
}
