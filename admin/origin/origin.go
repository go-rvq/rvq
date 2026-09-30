// Package origin is where a request came from — the address, the browser
// and, when the application locates addresses, the place —: where the author
// of a revision was (the history), where a record was deleted from (the
// trash).
package origin

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Origin is where a request came from: its address and browser and — when
// the application locates addresses (SetFunc) — where the address is,
// approximate (a city's). Zero when nothing is known: no request.
type Origin struct {
	IP          string `gorm:"type:varchar(64)"`
	UserAgent   string
	Country     string `gorm:"type:varchar(2)"` // ISO code: "BR"
	CountryName string
	Region      string
	City        string
	Latitude    *float64
	Longitude   *float64
}

// Place is where the origin is in words — "Viçosa, Minas Gerais, Brazil" —;
// "" when it is not known.
func (o Origin) Place() string {
	country := o.CountryName
	if country == "" {
		country = o.Country
	}
	var parts []string
	for _, s := range []string{o.City, o.Region, country} {
		if s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, ", ")
}

// String is the origin on one line: its place and its address — "Viçosa,
// Minas Gerais, Brazil (200.1.2.3)" —, the address alone when the place is
// not known; "" when nothing is.
func (o Origin) String() string {
	switch place := o.Place(); {
	case place != "" && o.IP != "":
		return fmt.Sprintf("%s (%s)", place, o.IP)
	case place != "":
		return place
	}
	return o.IP
}

// Func is where the request r came from.
type Func func(r *http.Request) Origin

var fn Func = Default

// SetFunc sets how the origin of a request is learned: an application that
// locates addresses (a GeoIP database) fills the place in — Default gives the
// address and the browser. nil: Default.
func SetFunc(f Func) {
	if f == nil {
		f = Default
	}
	fn = f
}

// Default is the address r came from — behind a reverse proxy, the first
// of X-Forwarded-For (or X-Real-Ip) —, and its browser; no place.
func Default(r *http.Request) Origin {
	return Origin{IP: ClientIP(r), UserAgent: r.UserAgent()}
}

// ClientIP is the address r came from, as far as this server can tell: behind
// a reverse proxy the connection is the proxy's, and the first address of
// X-Forwarded-For (or X-Real-Ip) is the client it saw. No port.
func ClientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		first, _, _ := strings.Cut(fwd, ",")
		return strings.TrimSpace(first)
	}
	if ip := strings.TrimSpace(r.Header.Get("X-Real-Ip")); ip != "" {
		return ip
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// Of is where r came from (the Func set); zero with no request.
func Of(r *http.Request) Origin {
	if r == nil {
		return Origin{}
	}
	return fn(r)
}
