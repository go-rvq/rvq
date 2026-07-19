package activity

import (
	"context"
	"net"
	"net/http"
	"strings"
)

// RequestInfo carries the origin of the request that triggered an activity
// record: client IP address and browser (user agent) information.
type RequestInfo struct {
	IP        string
	UserAgent string
}

// RequestInfoSetter is optionally implemented by log models that persist the
// request origin of the action (see ActivityLog.IP / ActivityLog.UserAgent).
type RequestInfoSetter interface {
	SetIP(string)
	SetUserAgent(string)
}

type requestInfoContextKey struct{}

// RequestInfoFromRequest extracts the client IP (honouring X-Forwarded-For /
// X-Real-IP) and the user agent of r.
func RequestInfoFromRequest(r *http.Request) *RequestInfo {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// first address is the originating client
		if i := strings.IndexByte(ip, ','); i >= 0 {
			ip = ip[:i]
		}
		ip = strings.TrimSpace(ip)
	}
	if ip == "" {
		ip = r.Header.Get("X-Real-Ip")
	}
	if ip == "" {
		ip = r.RemoteAddr
		if host, _, err := net.SplitHostPort(ip); err == nil {
			ip = host
		}
	}
	return &RequestInfo{IP: ip, UserAgent: r.UserAgent()}
}

// ContextWithRequestInfo stores the origin of r into ctx so activity records
// created down the call chain persist IP and user agent.
func ContextWithRequestInfo(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestInfoContextKey{}, RequestInfoFromRequest(r))
}

// RequestInfoFromContext returns the request info stored by
// ContextWithRequestInfo, or nil.
func RequestInfoFromContext(ctx context.Context) *RequestInfo {
	info, _ := ctx.Value(requestInfoContextKey{}).(*RequestInfo)
	return info
}
