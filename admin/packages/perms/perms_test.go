package perms

import (
	"reflect"
	"testing"
)

func TestParseSubjects(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"", nil},
		{"admin", []string{"admin"}},
		{"admin, editor", []string{"admin", "editor"}},
		{"admin\neditor\nviewer", []string{"admin", "editor", "viewer"}},
		{" admin ; editor ,, admin ", []string{"admin", "editor"}}, // trims + de-dups
		{"a\r\nb", []string{"a", "b"}},
	}
	for _, c := range cases {
		if got := parseSubjects(c.in); !reflect.DeepEqual(got, c.want) {
			t.Errorf("parseSubjects(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}
