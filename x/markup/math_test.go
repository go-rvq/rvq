package markup

import (
	"strings"
	"testing"
)

// TestMath: the formulas of a text, as GitHub's — inline $…$ and $`…`$, a
// block of $$ lines, $$…$$ alone on a line, a ```math fence — made the rich
// text editor's HTML (data-type, data-latex); money, a $ before a digit, an
// escaped one, code: not formulas; kept in the rich text; made Markdown
// again.
func TestMath(t *testing.T) {
	got := Render(Markdown, "Inline $a^2 + b_1$ and $`x*y`$.\n\n"+
		"It costs $5 and $10, \\$x$ not.\n\n"+
		"$$\n\\frac{1}{2}\n$$\n\n"+
		"$$E = mc^2$$\n\n"+
		"```math\n\\sqrt{x}\n```\n\n"+
		"`$code$` and\n\n```\n$$\nnot\n$$\n```\n")
	for _, want := range []string{
		`<span data-type="inline-math" data-latex="a^2 + b_1">a^2 + b_1</span>`,
		`<span data-type="inline-math" data-latex="x*y">x*y</span>`,
		`<div data-type="block-math" data-latex="\frac{1}{2}">\frac{1}{2}</div>`,
		`<div data-type="block-math" data-latex="E = mc^2">E = mc^2</div>`,
		`<div data-type="block-math" data-latex="\sqrt{x}">\sqrt{x}</div>`,
		"It costs $5 and $10",
		"<code>$code$</code>",
		"$$\nnot\n$$",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("no %s in\n%s", want, got)
		}
	}
	if n := strings.Count(got, "data-type="); n != 5 {
		t.Errorf("formulas: %d\n%s", n, got)
	}

	rich := Render(Rich, `<p>a <span data-type="inline-math" data-latex="x^2">x^2</span></p><div data-type="block-math" data-latex="y">y</div><span data-type="evil">z</span>`)
	if !strings.Contains(rich, `data-latex="x^2"`) || !strings.Contains(rich, `data-type="block-math"`) || strings.Contains(rich, "evil") {
		t.Errorf("rich: %s", rich)
	}
	md := HTMLToMarkdown(`<p>a <span data-type="inline-math" data-latex="x^2">x^2</span></p><div data-type="block-math" data-latex="y">y</div>`)
	if !strings.Contains(md, "$x^2$") || !strings.Contains(md, "$$\ny\n$$") {
		t.Errorf("to Markdown: %q", md)
	}
}

// TestDiagramsAndDetails: a ```mermaid fence stays code of its language (the
// page draws it); a <details> folds, its summary, its Markdown; both kept
// in the rich text and made Markdown again.
func TestDiagramsAndDetails(t *testing.T) {
	got := Render(Markdown, "```mermaid\ngraph TD\n  A-->B\n```\n\n<details>\n<summary>More</summary>\n\nSome **bold**.\n\n</details>\n")
	for _, want := range []string{`<pre><code class="language-mermaid">graph TD`, "<details>", "<summary>More</summary>", "<strong>bold</strong>"} {
		if !strings.Contains(got, want) {
			t.Errorf("no %s in\n%s", want, got)
		}
	}
	// the rich text's: its parts kept (the editor reads them back)
	rich := Render(Rich, `<details data-type="details" open><summary>S</summary><div data-type="detailsContent"><p>x</p></div></details>`)
	if !strings.Contains(rich, `data-type="detailsContent"`) || !strings.Contains(rich, "<summary>S</summary>") {
		t.Errorf("rich: %s", rich)
	}
	md := HTMLToMarkdown(`<pre><code class="language-mermaid">graph TD
A--&gt;B</code></pre><details open><summary>More</summary><div data-type="detailsContent"><p>Some <strong>bold</strong>.</p></div></details>`)
	for _, want := range []string{"```mermaid\ngraph TD\nA-->B\n```", "<details open>\n<summary>More</summary>\n\nSome **bold**.\n\n</details>"} {
		if !strings.Contains(md, want) {
			t.Errorf("to Markdown: no %q in %q", want, md)
		}
	}
}
