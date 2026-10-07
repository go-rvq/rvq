package userdocs

import (
	"bytes"
	"strings"
	"testing"
)

// A fence renders as the block Prism highlights — `<code class="language-…">`
// —, whatever its language: Gad's (gad, gadx, gadt) and the common ones; one
// with none stays plain text.
func TestFencesForPrism(t *testing.T) {
	var src strings.Builder
	langs := []string{"gad", "gadx", "gadt", "go", "json", "bash", "yaml", "html", "css", "js", "ts", "sql", "diff", "python", "toml", "dockerfile"}
	for _, l := range langs {
		src.WriteString("```" + l + "\nx\n```\n\n")
	}
	src.WriteString("```\nplain\n```\n")
	var out bytes.Buffer
	if err := newMarkdown(func(d string) string { return d }).Convert([]byte(src.String()), &out); err != nil {
		t.Fatal(err)
	}
	html := out.String()
	for _, l := range langs {
		if !strings.Contains(html, `<pre><code class="language-`+l+`">`) {
			t.Errorf("no block of %s:\n%s", l, html)
		}
	}
	if !strings.Contains(html, "<pre><code>plain") {
		t.Errorf("a fence of no language:\n%s", html)
	}
}
