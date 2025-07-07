package perm

import (
	"github.com/gobwas/glob"
	"github.com/ory/ladon"
)

type PathMatcher struct{}

func (m *PathMatcher) Matches(p ladon.Policy, haystack []string, needle string) (ok bool, err error) {
	var g glob.Glob
	for _, h := range haystack {
		if g, err = glob.Compile(h); err != nil {
			return
		}
		if g.Match(needle) {
			return true, nil
		}
	}
	return
}
