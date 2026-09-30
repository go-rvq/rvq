package admin

import (
	"net"
	"net/http"
	"strings"

	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
)

// OriginFunc is where the request r — a save, a Record — came from: the
// revision's Origin.
type OriginFunc func(r *http.Request) histmodels.Origin

var originFunc OriginFunc = DefaultOrigin

// SetOriginFunc sets how the revisions learn where their author was: an
// application that locates addresses (a GeoIP database) fills the place in —
// DefaultOrigin gives the address and the browser. nil: DefaultOrigin.
func SetOriginFunc(f OriginFunc) {
	if f == nil {
		f = DefaultOrigin
	}
	originFunc = f
}

// DefaultOrigin is the address r came from — behind a reverse proxy, the first
// of X-Forwarded-For (or X-Real-Ip) —, and its browser; no place.
func DefaultOrigin(r *http.Request) histmodels.Origin {
	return histmodels.Origin{IP: ClientIP(r), UserAgent: r.UserAgent()}
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

// originOf is where r came from; zero with no request.
func originOf(r *http.Request) histmodels.Origin {
	if r == nil {
		return histmodels.Origin{}
	}
	return originFunc(r)
}
