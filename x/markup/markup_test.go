package markup

import (
	"strings"
	"testing"
)

// The body, Markdown or rich text, shown safe: no script, on*, javascript:;
// its text for the search; one format made the other.
func TestMarkup(t *testing.T) {
	md := "## Steps\n\n1. Open **it**\n2. <script>alert(1)</script>Click [here](javascript:alert(1))\n3. [there](javascript:alert(2))\n\n> quoted"
	h := Render(Markdown, md)
	for _, want := range []string{"<h2", "<strong>it</strong>", "<ol>", "<blockquote>"} {
		if !strings.Contains(h, want) {
			t.Errorf("markdown: no %s in %s", want, h)
		}
	}
	// (its HTML kept, then sanitized: a line of HTML is a block of it, its
	// script dropped, the rest text — no link)
	for _, bad := range []string{"<script", `href="javascript`, "alert(1)</script"} {
		if strings.Contains(h, bad) {
			t.Errorf("markdown: %s in %s", bad, h)
		}
	}
	rich := `<p onclick="x()">Hi <strong>you</strong><img src=x onerror=alert(1)></p><script>bad()</script><a href="javascript:x()">l</a>`
	r := Render(Rich, rich)
	for _, bad := range []string{"onclick", "onerror", "<script", "javascript:"} {
		if strings.Contains(r, bad) {
			t.Errorf("rich: %s in %s", bad, r)
		}
	}
	if !strings.Contains(r, "<strong>you</strong>") {
		t.Errorf("rich lost its markup: %s", r)
	}
	if got := Text(Rich, "<h2>Title</h2><p>One &amp; two</p><p>Three</p>"); got != "Title\nOne & two\nThree" {
		t.Errorf("text: %q", got)
	}

	back := HTMLToMarkdown(`<h2>Steps</h2><p>Open <strong>it</strong> and <em>see</em> <a href="https://x.example">this</a></p><ul><li><p>one</p></li><li><p>two</p></li></ul><ol><li>a</li><li>b</li></ol><blockquote><p>quoted</p></blockquote><pre><code>x := 1</code></pre>`)
	for _, want := range []string{"## Steps", "Open **it** and *see* [this](https://x.example)", "- one\n- two", "1. a\n2. b", "> quoted", "```\nx := 1\n```"} {
		if !strings.Contains(back, want) {
			t.Errorf("to markdown: no %q in\n%s", want, back)
		}
	}
	if got := Convert(Markdown, Rich, "**x**"); !strings.Contains(got, "<strong>x</strong>") {
		t.Errorf("markdown to rich: %s", got)
	}
}

// TestMarkdownTables: a table of Markdown, its columns aligned (align, not
// a style the policy drops); one of HTML — merged cells, the table builder's
// — kept, sanitized.
func TestMarkdownTables(t *testing.T) {
	got := Render("markdown", "| A | B |\n| :-: | --: |\n| 1 | 2 |\n")
	for _, want := range []string{`<th align="center">A</th>`, `<td align="right">2</td>`} {
		if !strings.Contains(got, want) {
			t.Errorf("GFM table: no %q in %s", want, got)
		}
	}
	got = Render("markdown", "<table>\n<tr><th colspan=\"2\" onclick=\"x()\">AB</th></tr>\n<tr><td rowspan=\"2\">1</td><td>2</td></tr>\n</table>\n\n<script>x()</script>")
	for _, want := range []string{`<th colspan="2">AB</th>`, `<td rowspan="2">1</td>`} {
		if !strings.Contains(got, want) {
			t.Errorf("HTML table: no %q in %s", want, got)
		}
	}
	if strings.Contains(got, "onclick") || strings.Contains(got, "script") {
		t.Errorf("HTML table: not sanitized: %s", got)
	}
}

// TestTableLayout: the class of a table not fluid kept (only it).
func TestTableLayout(t *testing.T) {
	got := Render(Markdown, "<table class=\"mdkit-fit mdkit-center\"><tr><td>a</td></tr></table>")
	if !strings.Contains(got, `<table class="mdkit-fit mdkit-center">`) {
		t.Errorf("the class of a table not fluid: %s", got)
	}
	if got := Render(Rich, `<table class="table-wrapper mdkit-fit mdkit-right"><tr><td>a</td></tr></table>`); !strings.Contains(got, `class="table-wrapper mdkit-fit mdkit-right"`) {
		t.Errorf("the class of the admin's TipTap: %s", got)
	}
	if got := Render(Rich, `<table class="evil"><tr><td>a</td></tr></table>`); strings.Contains(got, "class") {
		t.Errorf("another class of a table kept: %s", got)
	}
}

// TestCellPlace: where the content of a cell goes (the table builder's align
// and valign) is kept, in Markdown and in the rich text; another value is not.
func TestCellPlace(t *testing.T) {
	for _, format := range []string{Markdown, Rich} {
		got := Render(format, `<table><tr><td align="center" valign="bottom">a</td><th valign="middle">b</th><td valign="evil">c</td></tr></table>`)
		for _, want := range []string{`<td align="center" valign="bottom">a</td>`, `<th valign="middle">b</th>`, `<td>c</td>`} {
			if !strings.Contains(got, want) {
				t.Errorf("%s: no %s in %s", format, want, got)
			}
		}
	}
}
