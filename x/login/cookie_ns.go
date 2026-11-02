package login

import (
	"fmt"
	"strings"
)

// The namespace of the login's cookies: their names prefixed with it, so two
// applications on one host (one under the other's path: senir's admin at
// /!/, an app at /) keep their sessions, accounts and flashes apart — the
// browser sends the cookies of / to /!/ too. "" (the default): the names as
// they are.
var cookieNS string

// cookieName is the name of the login's cookie n in the namespace.
func cookieName(n string) string {
	if cookieNS == "" {
		return n
	}
	return cookieNS + "_" + n
}

func setCookieNamespace(ns string) { cookieNS = ns }

// nsOfPrefix is a namespace of a path prefix: its letters and digits, the
// other characters by their code ("/!" → "x21", "/admin" → "admin").
func nsOfPrefix(prefix string) string {
	var sb strings.Builder
	for _, r := range strings.Trim(prefix, "/") {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			sb.WriteRune(r)
		default:
			fmt.Fprintf(&sb, "x%x", r)
		}
	}
	return sb.String()
}
