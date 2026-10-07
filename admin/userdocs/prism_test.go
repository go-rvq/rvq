package userdocs

import (
	"bytes"
	"strings"
	"testing"
)

// A fence renders as the admin's code viewer — <vx-code>, highlighted by
// Prism with its plugins —: its code an attribute (nothing in it read as
// markup, nor by Vue), its language, the words of its copy button. One with
// no language has none.
func TestFencesAsCodeViewer(t *testing.T) {
	src := "```gad\nx := {a: 1} // <b>{{ v }}</b> & \"q\"\n```\n\n```yaml\nk: v\n```\n\n```\nplain\n```\n"
	var out bytes.Buffer
	md := newMarkdown(func(d string) string { return d }, codeLabels{"Copiar", "Copiado!", "Ctrl+C"})
	if err := md.Convert([]byte(src), &out); err != nil {
		t.Fatal(err)
	}
	html := out.String()
	for _, want := range []string{
		`<vx-code code="x := {a: 1} // &lt;b&gt;{{ v }}&lt;/b&gt; &amp; &#34;q&#34;" language="gad" copy-text="Copiar" copied-text="Copiado!" copy-error-text="Ctrl+C"></vx-code>`,
		`<vx-code code="k: v" language="yaml"`,
		`<vx-code code="plain" copy-text=`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("no %s in:\n%s", want, html)
		}
	}
	if strings.Contains(html, "<pre") {
		t.Errorf("a fence rendered as a <pre>:\n%s", html)
	}
}
