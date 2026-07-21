package admin

import (
	"strings"
	"testing"
)

func TestMdToHTML(t *testing.T) {
	md := "# Title\n\nSome **bold** and `code` and a [link](https://x).\n\n- one\n- two\n\n```\nplain\n```"
	got := mdToHTML(md)
	for _, want := range []string{
		"<h1>Title</h1>",
		"<strong>bold</strong>",
		"<code>code</code>",
		`<a href="https://x"`,
		"<ul>", "<li>one</li>",
		"<pre><code>plain</code></pre>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("mdToHTML missing %q in:\n%s", want, got)
		}
	}
	// HTML in the source is escaped, not passed through
	if strings.Contains(mdToHTML("<script>alert(1)</script>"), "<script>") {
		t.Errorf("raw HTML must be escaped")
	}
}
