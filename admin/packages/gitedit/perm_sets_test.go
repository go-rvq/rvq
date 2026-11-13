package gitedit

import (
	"sort"
	"testing"
	"testing/fstest"

	"github.com/go-rvq/rvq/x/perm"
)

// The sets of the page: what each gives, of the actions it has (the preview;
// the default files: maintainers).
func TestPermSets(t *testing.T) {
	var invalid error
	repo, _ := site(t)
	appOn(t, repo, &invalid, perm.New().AllowAll(), fstest.MapFS{"templates/a.gadx": {Data: []byte("a")}})
	t.Cleanup(func() { perm.SetSetExpander(nil) })
	exp := func(name string) []string {
		l, ok := perm.ExpandSet("admin:/site-files:", name)
		if !ok {
			return nil
		}
		sort.Strings(l)
		return l
	}
	has := func(name string, want ...string) {
		t.Helper()
		got := map[string]bool{}
		for _, l := range exp(name) {
			got[l] = true
		}
		for _, w := range want {
			if !got["admin:/site-files:"+w] {
				t.Errorf("%s: no %s in %q", name, w, exp(name))
			}
		}
	}
	has("readers", "@get", "!preview")
	has("authors", "@get", "!preview", "!create", "!edit", "!rename", "!move", "!delete", "!import", "!commit", "!update", "!discard", "!reset")
	has("publishers", "!create", "!publish")
	has("developers", "@get", "!git")
	has("sharers", "!share")
	has("reviewers", "!drafts")
	has("maintainers", "!init_defaults", "!reset_defaults")
	for _, l := range exp("authors") {
		if l == "admin:/site-files:!publish" || l == "admin:/site-files:!git" {
			t.Errorf("authors gives %s", l)
		}
	}
	// the general sets of a page too: managers, all
	if got := exp("managers"); len(got) != 1 || got[0] != "admin:/site-files:*" {
		t.Errorf("managers: %q", got)
	}
}
