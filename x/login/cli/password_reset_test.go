package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// The new password read by lines (a script): the two equal, of at least
// MinPasswordLen; a short one, two different, or no line refused; asked on
// the prompt only when there is one.
func TestReadNewPassword(t *testing.T) {
	for _, c := range []struct {
		in, want, err string
	}{
		{"123456\n123456\n", "123456", ""},
		{"123456\r\n123456\r\n", "123456", ""},
		{"secret-pass\nsecret-pass", "secret-pass", ""},
		{"123\n123\n", "", "at least 6"},
		{"123456\n654321\n", "", "not equal"},
		{"123456\n", "", "no password"},
		{"", "", "no password"},
	} {
		got, err := readNewPassword(linesReader(strings.NewReader(c.in)), nil, "admin")
		switch {
		case c.err == "" && (err != nil || got != c.want):
			t.Errorf("%q: %q, %v", c.in, got, err)
		case c.err != "" && (err == nil || !strings.Contains(err.Error(), c.err)):
			t.Errorf("%q: want an error of %q, got %q, %v", c.in, c.err, got, err)
		}
	}
	if _, err := readNewPassword(linesReader(strings.NewReader("")), nil, "admin"); !errors.Is(err, errNoPassword) {
		t.Errorf("no input: %v", err)
	}

	var prompt bytes.Buffer
	if _, err := readNewPassword(linesReader(strings.NewReader("123456\n123456\n")), &prompt, "admin"); err != nil {
		t.Fatal(err)
	}
	if p := prompt.String(); !strings.Contains(p, `"admin" password`) || !strings.Contains(p, "Confirm password") {
		t.Errorf("the prompt: %q", p)
	}
}
