package profile

import (
	"net/url"
	"testing"
)

func TestPageURI(t *testing.T) {
	for in, want := range map[string]string{
		"/admin/shared/confirm?__execute_event__=shared_confirm&t=x": "/admin/shared/confirm?t=x",
		"/admin/orgs?__execute_event__=presets_Update&__reload__=1":  "/admin/orgs",
		"/admin/orgs?page=2": "/admin/orgs?page=2",
	} {
		u, _ := url.Parse(in)
		if got := pageURI(u); got != want {
			t.Errorf("pageURI(%q) = %q, want %q", in, got, want)
		}
	}
}
