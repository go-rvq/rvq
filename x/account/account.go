// Package account is the index of the account a request is for, when a
// browser holds several logged-in accounts (as Gmail's /mail/u/N/): the
// account 0 is the one of the plain paths (/admin/…), the account N>0 the one
// of the paths with the prefix /u/N (/admin/u/1/…). The login (x/login) reads
// and writes the cookies of that account; the admin (presets) makes its links
// with that prefix.
//
// Only the index is here, no session: who the account is, its cookies say.
package account

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

// Header is the header a page of the account N sends with its requests (its
// events, its portals): the account of a request whose path has no /u/N.
const Header = "X-Rvq-Account"

// Segment is the segment of the paths before the index: /u/N.
const Segment = "u"

type ctxKey struct{}

// WithIndex is ctx of the account n.
func WithIndex(ctx context.Context, n int) context.Context {
	return context.WithValue(ctx, ctxKey{}, n)
}

// Index is the account of ctx: 0 when none was said.
func Index(ctx context.Context) int {
	n, _ := ctx.Value(ctxKey{}).(int)
	return n
}

// IndexOf is the account of r.
func IndexOf(r *http.Request) int {
	return Index(r.Context())
}

// Prefix is base for the account n: base itself for the account 0,
// base/u/N for another ("/admin" → "/admin/u/1").
func Prefix(base string, n int) string {
	if n <= 0 {
		return base
	}
	return strings.TrimSuffix(base, "/") + "/" + Segment + "/" + strconv.Itoa(n)
}

// PrefixCtx is base for the account of ctx.
func PrefixCtx(ctx context.Context, base string) string {
	return Prefix(base, Index(ctx))
}

// URL is the path p — one under base (base, base/…) — for the account n;
// any other path as it is. A path that already has its /u/N is the account
// it says: it is not changed.
func URL(base, p string, n int) string {
	base = strings.TrimSuffix(base, "/")
	if n <= 0 || base == "" {
		return p
	}
	rest, ok := under(base, p)
	if !ok {
		return p
	}
	if _, _, has := splitIndex(rest); has {
		return p
	}
	return Prefix(base, n) + rest
}

// Split takes the /u/N of p, a path under base: the path without it, and N
// (0, false when p has none).
func Split(base, p string) (path string, n int, ok bool) {
	base = strings.TrimSuffix(base, "/")
	rest, under := under(base, p)
	if !under {
		return p, 0, false
	}
	after, n, ok := splitIndex(rest)
	if !ok {
		return p, 0, false
	}
	return base + after, n, true
}

// under is what p has after base ("" or "/…"), when p is base or under it.
func under(base, p string) (rest string, ok bool) {
	if base == "" {
		return p, strings.HasPrefix(p, "/")
	}
	if p == base {
		return "", true
	}
	if strings.HasPrefix(p, base+"/") {
		return p[len(base):], true
	}
	return "", false
}

// splitIndex reads "/u/N" at the start of rest ("/u/1", "/u/1/…"): what
// follows it ("" or "/…") and N.
func splitIndex(rest string) (after string, n int, ok bool) {
	pre := "/" + Segment + "/"
	if !strings.HasPrefix(rest, pre) {
		return rest, 0, false
	}
	s := rest[len(pre):]
	end := strings.IndexByte(s, '/')
	if end < 0 {
		end = len(s)
	}
	n, err := strconv.Atoi(s[:end])
	if err != nil || n < 0 || s[:end] != strconv.Itoa(n) {
		return rest, 0, false
	}
	return s[end:], n, true
}

// CookieName is the name of the cookie base for the account n: base for the
// account 0 (the names of a single account, as ever), base_uN for another.
func CookieName(base string, n int) string {
	if n <= 0 {
		return base
	}
	return base + "_" + Segment + strconv.Itoa(n)
}
