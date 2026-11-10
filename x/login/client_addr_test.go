package login

import (
	"net/http/httptest"
	"testing"
)

// The address of the client: behind the loopback or a trusted proxy, the
// last of X-Forwarded-For not of one of them (what a client wrote before it
// is not taken); from anyone else, the connection's.
func TestClientAddr(t *testing.T) {
	old := TrustedProxies
	TrustedProxies = "172.16.0.0/12"
	t.Cleanup(func() { TrustedProxies = old })
	for _, c := range []struct{ remote, fwd, want string }{
		{"127.0.0.1:4000", "203.0.113.9", "203.0.113.9"},
		{"127.0.0.1:4000", "1.2.3.4, 203.0.113.9, 172.18.0.5", "203.0.113.9"},
		{"172.18.0.5:80", "198.51.100.7", "198.51.100.7"},
		{"203.0.113.9:5000", "1.2.3.4", "203.0.113.9:5000"},
		{"127.0.0.1:4000", "", "127.0.0.1:4000"},
	} {
		r := httptest.NewRequest("GET", "/", nil)
		r.RemoteAddr = c.remote
		if c.fwd != "" {
			r.Header.Set("X-Forwarded-For", c.fwd)
		}
		if got := ClientAddr(r); got != c.want {
			t.Errorf("%s %q: %s, want %s", c.remote, c.fwd, got, c.want)
		}
	}
}
