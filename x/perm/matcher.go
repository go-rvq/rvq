package perm

import (
	"regexp"
	"strings"
	"sync"

	"github.com/gobwas/glob"
	"github.com/ory/ladon"
)

// PathMatcher matches the resources of the policies as globs: "*" any text
// (":" too), "{a,b}" either. "<*>" is one record — its id, between "<" and
// ">", no ":" in it —: "admin:posts:<*>:@get" is the reading of any post,
// not of what is under one ("admin:posts:<7>:comments:<1>:@get").
type PathMatcher struct{}

func (m *PathMatcher) Matches(p ladon.Policy, haystack []string, needle string) (ok bool, err error) {
	for _, h := range haystack {
		if strings.Contains(h, RecordAny) {
			re, err := recordPattern(h)
			if err != nil {
				return false, err
			}
			if re.MatchString(needle) {
				return true, nil
			}
			continue
		}
		g, err := glob.Compile(h)
		if err != nil {
			return false, err
		}
		if g.Match(needle) {
			return true, nil
		}
	}
	return
}

// RecordAny is the record of any id in a resource (PathMatcher).
const RecordAny = "<*>"

var recordPatterns sync.Map // pattern → *regexp.Regexp

// recordPattern is the glob p — with "<*>" — as a regular expression.
func recordPattern(p string) (*regexp.Regexp, error) {
	if re, ok := recordPatterns.Load(p); ok {
		return re.(*regexp.Regexp), nil
	}
	var b strings.Builder
	b.WriteString("^")
	for i := 0; i < len(p); {
		switch {
		case strings.HasPrefix(p[i:], RecordAny):
			b.WriteString("<[^:>]*>")
			i += len(RecordAny)
		case p[i] == '*':
			b.WriteString(".*")
			i++
		case p[i] == '?':
			b.WriteString(".")
			i++
		case p[i] == '{':
			j := strings.IndexByte(p[i:], '}')
			if j < 0 {
				b.WriteString(regexp.QuoteMeta(p[i:]))
				i = len(p)
				continue
			}
			alts := strings.Split(p[i+1:i+j], ",")
			for k, a := range alts {
				alts[k] = regexp.QuoteMeta(a)
			}
			b.WriteString("(?:" + strings.Join(alts, "|") + ")")
			i += j + 1
		default:
			b.WriteString(regexp.QuoteMeta(p[i : i+1]))
			i++
		}
	}
	b.WriteString("$")
	re, err := regexp.Compile(b.String())
	if err != nil {
		return nil, err
	}
	recordPatterns.Store(p, re)
	return re, nil
}
