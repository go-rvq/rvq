package perm

import "testing"

// "<*>" is one record: not what is under it; "*" and "{a,b}" as ever.
func TestPathMatcherRecordAny(t *testing.T) {
	m := &PathMatcher{}
	for _, c := range []struct {
		pattern, needle string
		ok              bool
	}{
		{"admin:posts:<*>:@get", "admin:posts:<7>:@get", true},
		{"admin:posts:<*>:@get", "admin:posts:<7>:comments:<1>:@get", false},
		{"admin:posts:<*>:#{Title,Body}:*", "admin:posts:<7>:#Title:@get", true},
		{"admin:posts:<*>:#{Title,Body}:*", "admin:posts:<7>:#Secret:@get", false},
		{"admin:posts:<*>:comments:*", "admin:posts:<7>:comments:<1>:@get", true},
		{"admin:posts:*", "admin:posts:<7>:comments:<1>:@get", true},
		{"admin:posts:{a,b}:*", "admin:posts:b:x", true},
		{"admin:p.s:<*>:@get", "admin:pxs:<1>:@get", false},
	} {
		got, err := m.Matches(nil, []string{c.pattern}, c.needle)
		if err != nil || got != c.ok {
			t.Errorf("%s ~ %s: %v %v, want %v", c.pattern, c.needle, got, err, c.ok)
		}
	}
}
