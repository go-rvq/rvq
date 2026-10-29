package markup

import (
	"strings"
	"testing"
)

// TestTemplateHints: the hints of a template (<!-- … -->) are the
// placeholders of empty paragraphs for the rich text editor, where they
// were, and nothing of what is sent.
func TestTemplateHints(t *testing.T) {
	md := "## What happened\n<!-- Describe the steps -->\n\n## Expected\n<!-- What should\nhave happened -->\n"
	got := TemplateHTML(md)
	want := `<h2>What happened</h2>` + "\n" + `<p data-hint="Describe the steps"></p>` + "\n" +
		`<h2>Expected</h2>` + "\n" + `<p data-hint="What should` + "\n" + `have happened"></p>` + "\n"
	if got != want {
		t.Fatalf("markdown template:\n%q\nwant\n%q", got, want)
	}

	// the editor's HTML; a hint with what must be escaped
	got = TemplateHTML(`<h2>Steps</h2><!-- "a" & <b> --><p>text</p>`)
	if want := `<h2>Steps</h2><p data-hint="&#34;a&#34; &amp; &lt;b&gt;"></p><p>text</p>`; got != want {
		t.Fatalf("rich template: %q, want %q", got, want)
	}

	// within a line of text: no paragraph, no token left
	if got := TemplateHTML("Version: <!-- x.y --> here"); strings.Contains(got, "HERMONHINT") || strings.Contains(got, "data-hint") {
		t.Fatalf("inline hint: %q", got)
	}
	// no hints: the same as Render
	if got, want := TemplateHTML("**a**"), Render("markdown", "**a**"); got != want {
		t.Fatalf("no hints: %q, want %q", got, want)
	}

	sent := "## What happened\n<!-- Describe the steps -->\n\nIt broke.\n\n## Expected\n<!-- What -->\n"
	if got, want := StripHints(sent), "## What happened\n\nIt broke.\n\n## Expected"; got != want {
		t.Fatalf("StripHints: %q, want %q", got, want)
	}
	if got := StripHints("no hints\n"); got != "no hints\n" {
		t.Fatalf("StripHints without hints: %q", got)
	}
}
