package markup

import (
	"strings"
	"testing"
)

const footnotesMD = "Um[^a] dois[^fonte] tres[^a].\n\n[^a]: Nota **a**.\n[^fonte]: Nota b\n    continua.\n"

// TestFootnotes: the footnotes of Markdown rendered — the references, the
// notes at the end, their classes kept (only theirs) —; a body among others
// (RenderIn), its ids and links its own; as text, [n].
func TestFootnotes(t *testing.T) {
	got := Render(Markdown, footnotesMD)
	for _, want := range []string{
		`<sup id="fnref:1"><a href="#fn:1" class="footnote-ref"`, `<sup id="fnref:2">`, `<sup id="fnref1:1">`,
		`<div class="footnotes">`, `<li id="fn:1">`, `<strong>a</strong>`, `href="#fnref:1" class="footnote-backref"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("render: no %q in\n%s", want, got)
		}
	}
	if strings.Contains(got, "role=") {
		t.Errorf("render: a role kept\n%s", got)
	}
	// other classes: none
	if got := Render(Rich, `<a href="#x" class="evil">x</a><div class="footnotes x">y</div>`); strings.Contains(got, "class") {
		t.Errorf("other classes kept: %s", got)
	}

	// two bodies of a page: their own ids and links
	a := RenderIn(Markdown, footnotesMD, "c1")
	b := RenderIn(Markdown, footnotesMD, "c2")
	for _, want := range []string{`id="c1-fnref:1"`, `href="#c1-fn:1"`, `id="c1-fn:1"`, `href="#c1-fnref1:1"`} {
		if !strings.Contains(a, want) {
			t.Errorf("RenderIn: no %q in\n%s", want, a)
		}
	}
	if strings.Contains(a, `"fn:`) || strings.Contains(a, `"#fn:`) || !strings.Contains(b, `id="c2-fn:1"`) {
		t.Errorf("RenderIn: an id not prefixed\n%s", a)
	}
	// a prefix of what an id may hold
	if got := RenderIn(Markdown, footnotesMD, `x"><script>`); strings.Contains(got, "<script") || !strings.Contains(got, `id="xscript-fn:1"`) {
		t.Errorf("RenderIn: the prefix not cleaned\n%s", got)
	}
	// none: as Render
	if RenderIn(Markdown, footnotesMD, "") != got {
		t.Error("RenderIn without a prefix: not Render")
	}

	// as text: [label] in Markdown; [n] in the rich text's HTML, its ways back none
	if got := Text(Markdown, footnotesMD); got != "Um[a] dois[fonte] tres[a].\n\n[a] Nota **a**.\n[fonte] Nota b\n    continua." {
		t.Errorf("Text of Markdown: %q", got)
	}
	rich := Render(Markdown, footnotesMD)
	text := Text(Rich, rich)
	for _, want := range []string{"Um[1] dois[2] tres[1].", "[1] Nota a.", "[2] Nota b"} {
		if !strings.Contains(text, want) {
			t.Errorf("Text of HTML: no %q in %q", want, text)
		}
	}
	if strings.Contains(text, "↩") {
		t.Errorf("Text of HTML: a way back in %q", text)
	}
}

// TestFootnotesMap: the map of the footnotes of a text — in Markdown (its
// code none of them), in the rich text's HTML —, and its labels changed; two
// texts put together (a book) with their labels made their own: their notes
// apart.
func TestFootnotesMap(t *testing.T) {
	md := "Um[^a] dois[^fonte] tres[^a] e `[^code]`.\n\n```\n[^x]: no\n```\n\n[^a]: Nota **a**.\n[^fonte]: Nota b\n    continua.\n[^sobra]: Sem referência.\n"
	got := Footnotes(Markdown, md)
	want := []Footnote{
		{Label: "a", Number: 1, Refs: 2, Text: "Nota **a**.", Defined: true},
		{Label: "fonte", Number: 2, Refs: 1, Text: "Nota b\ncontinua.", Defined: true},
		{Label: "sobra", Text: "Sem referência.", Defined: true},
	}
	if len(got) != len(want) {
		t.Fatalf("Footnotes of Markdown: %+v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Footnotes of Markdown %d: %+v, want %+v", i, got[i], want[i])
		}
	}
	// a reference with no note
	if fs := Footnotes(Markdown, "x[^nada]"); len(fs) != 1 || fs[0].Defined || fs[0].Number != 1 {
		t.Errorf("a reference with no note: %+v", fs)
	}

	re := RelabelFootnotes(Markdown, md, map[string]string{"a": "c1-a", "fonte": "c1 fonte!"})
	if !strings.Contains(re, "Um[^c1-a] dois[^c1fonte] tres[^c1-a]") || !strings.Contains(re, "[^c1-a]: Nota") ||
		!strings.Contains(re, "[^c1fonte]: Nota b") || !strings.Contains(re, "`[^code]`") || !strings.Contains(re, "[^x]: no") ||
		!strings.Contains(re, "[^sobra]: Sem") {
		t.Errorf("RelabelFootnotes of Markdown:\n%s", re)
	}

	// the rich text's HTML (goldmark's, its ids the labels)
	html := `<p>Um<sup id="fnref:a"><a href="#fn:a" class="footnote-ref">1</a></sup> e<sup id="fnref1:a"><a href="#fn:a" class="footnote-ref">1</a></sup></p>` +
		`<div class="footnotes"><hr><ol><li id="fn:a"><p>Nota <strong>a</strong>.&#160;<a href="#fnref:a" class="footnote-backref">↩</a></p></li></ol></div>`
	hf := Footnotes(Rich, html)
	if len(hf) != 1 || hf[0].Label != "a" || hf[0].Refs != 2 || hf[0].Number != 1 || hf[0].Text != "<p>Nota <strong>a</strong>.</p>" {
		t.Errorf("Footnotes of HTML: %+v", hf)
	}
	rh := RelabelFootnotes(Rich, html, map[string]string{"a": "c2-a"})
	for _, want := range []string{`id="fnref:c2-a"`, `id="fnref1:c2-a"`, `href="#fn:c2-a"`, `id="fn:c2-a"`, `href="#fnref:c2-a"`} {
		if !strings.Contains(rh, want) {
			t.Errorf("RelabelFootnotes of HTML: no %q in %s", want, rh)
		}
	}
	if strings.Contains(rh, `:a"`) {
		t.Errorf("RelabelFootnotes of HTML: a label kept: %s", rh)
	}

	// two chapters, the same labels: each its own, then together
	ch1 := "Um[^1].\n\n[^1]: Do primeiro.\n"
	ch2 := "Dois[^1].\n\n[^1]: Do segundo.\n"
	book := RelabelFootnotes(Markdown, ch1, FootnoteLabels(Markdown, ch1, "cap1")) + "\n" +
		RelabelFootnotes(Markdown, ch2, FootnoteLabels(Markdown, ch2, "cap2"))
	fs := Footnotes(Markdown, book)
	if len(fs) != 2 || fs[0].Label != "cap1-1" || fs[0].Text != "Do primeiro." || fs[1].Label != "cap2-1" || fs[1].Number != 2 {
		t.Errorf("a book of two chapters: %+v", fs)
	}
	if out := Render(Markdown, book); !strings.Contains(out, "Do primeiro.") || !strings.Contains(out, "Do segundo.") || !strings.Contains(out, `id="fn:2"`) {
		t.Errorf("a book of two chapters, rendered:\n%s", out)
	}
	if len(FootnoteLabels(Markdown, ch1, "")) != 0 {
		t.Error("FootnoteLabels without a prefix: none")
	}
}

// TestConvertFootnotes: a text with notes made the other format and back:
// the references and the notes kept (the labels of Markdown made numbers by
// the HTML, its ids by order).
func TestConvertFootnotes(t *testing.T) {
	html := `<p>Um<sup id="fnref:fonte"><a href="#fn:fonte" class="footnote-ref">1</a></sup> e<sup id="fnref1:fonte"><a href="#fn:fonte" class="footnote-ref">1</a></sup>.</p>` +
		`<div class="footnotes"><hr><ol><li id="fn:fonte"><p>Nota <strong>a</strong>.&#160;<a href="#fnref:fonte" class="footnote-backref">↩</a></p></li></ol></div>`
	md := Convert(Rich, Markdown, html)
	if md != "Um[^fonte] e[^fonte].\n\n---\n\n[^fonte]: Nota **a**.\n" && md != "Um[^fonte] e[^fonte].\n\n[^fonte]: Nota **a**.\n" {
		t.Errorf("rich → Markdown: %q", md)
	}
	fs := Footnotes(Markdown, md)
	if len(fs) != 1 || fs[0].Label != "fonte" || fs[0].Refs != 2 || fs[0].Text != "Nota **a**." {
		t.Errorf("rich → Markdown, its notes: %+v", fs)
	}
	back := Convert(Markdown, Rich, md)
	if fs := Footnotes(Rich, back); len(fs) != 1 || fs[0].Refs != 2 || !strings.Contains(fs[0].Text, "<strong>a</strong>") {
		t.Errorf("Markdown → rich, its notes: %+v\n%s", fs, back)
	}
}
