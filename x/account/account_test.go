package account

import (
	"context"
	"testing"
)

func TestSplit(t *testing.T) {
	for _, c := range []struct {
		base, p, path string
		n             int
		ok            bool
	}{
		{"/admin", "/admin/u/1/content/pages", "/admin/content/pages", 1, true},
		{"/admin", "/admin/u/2", "/admin", 2, true},
		{"/admin", "/admin/u/0/", "/admin/", 0, true},
		{"/admin/", "/admin/u/3/x", "/admin/x", 3, true},
		{"/admin", "/admin/content", "/admin/content", 0, false},
		{"/admin", "/admin/u/x/y", "/admin/u/x/y", 0, false},
		{"/admin", "/admin/u/01/y", "/admin/u/01/y", 0, false},
		{"/admin", "/admin/u/-1/y", "/admin/u/-1/y", 0, false},
		{"/admin", "/administer/u/1", "/administer/u/1", 0, false},
		{"/admin", "/auth/login", "/auth/login", 0, false},
		{"", "/u/4/a", "/a", 4, true},
	} {
		path, n, ok := Split(c.base, c.p)
		if path != c.path || n != c.n || ok != c.ok {
			t.Errorf("Split(%q, %q) = %q, %d, %v; want %q, %d, %v", c.base, c.p, path, n, ok, c.path, c.n, c.ok)
		}
	}
}

func TestURL(t *testing.T) {
	for _, c := range []struct {
		p    string
		n    int
		want string
	}{
		{"/admin/content/pages", 1, "/admin/u/1/content/pages"},
		{"/admin", 2, "/admin/u/2"},
		{"/admin/", 2, "/admin/u/2/"},
		{"/admin/u/3/x", 1, "/admin/u/3/x"},
		{"/admin/x", 0, "/admin/x"},
		{"/auth/login", 1, "/auth/login"},
		{"/administer", 1, "/administer"},
	} {
		if got := URL("/admin", c.p, c.n); got != c.want {
			t.Errorf("URL(%q, %d) = %q, want %q", c.p, c.n, got, c.want)
		}
	}
}

func TestIndexPrefixCookie(t *testing.T) {
	ctx := context.Background()
	if Index(ctx) != 0 {
		t.Fatal("no index: 0")
	}
	ctx = WithIndex(ctx, 2)
	if Index(ctx) != 2 || PrefixCtx(ctx, "/admin") != "/admin/u/2" {
		t.Fatalf("index 2: %d %q", Index(ctx), PrefixCtx(ctx, "/admin"))
	}
	if Prefix("/admin", 0) != "/admin" {
		t.Fatal("account 0: the base")
	}
	if CookieName("auth", 0) != "auth" || CookieName("auth", 3) != "auth_u3" {
		t.Fatal("cookie names")
	}
}
