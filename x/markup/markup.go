// Package markup is the text a user writes in a format — Markdown, or the
// HTML of a rich text editor (tiptap) —: shown safe (Render, RenderIn), as
// text (Text), made the other format (Convert); its footnotes mapped and
// relabeled (Footnotes, RelabelFootnotes: texts put together keep their
// notes apart); the hints of a template (TemplateHTML, StripHints). The
// fields of the admin's MarkdownBuilder, and of an application (hermon: its
// issues, the fields of its forms), are written so.
package markup

import (
	"bytes"
	"fmt"
	stdhtml "html"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// The formats of a text.
const (
	// Markdown is Markdown, as GitHub's (GFM, footnotes)
	Markdown = "markdown"
	// Rich is the HTML of a rich text editor
	Rich = "rich"
)

// policy is what a body keeps of HTML: the user generated content's — no
// script, style, on*, iframe nor javascript: —, the same for Markdown made
// HTML and for the rich text's; of the classes, only the footnotes' (their
// look, their balloon).
var policy = func() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(footnoteClassRe).OnElements("a", "div", "sup", "section")
	p.AllowAttrs("class").Matching(tableClassRe).OnElements("table")
	// a link card: its class, its target (a path of the site, a record's)
	p.AllowAttrs("class").Matching(regexp.MustCompile(`^` + LinkCardClass + `$`)).OnElements("div")
	p.AllowAttrs("data-href").Matching(cardTargetRe).OnElements("div")
	// the language of a block of code: a diagram (```mermaid) is drawn by
	// the page
	p.AllowAttrs("class").Matching(regexp.MustCompile(`^language-[A-Za-z0-9_+-]+$`)).OnElements("code")
	// a section that folds, as the rich text editor's (tiptap's details):
	// its parts
	p.AllowAttrs("data-type").Matching(regexp.MustCompile(`^(details|detailsSummary|detailsContent)$`)).OnElements("details", "summary", "div")
	// a formula (math.go): its kind, its LaTeX
	p.AllowAttrs("data-type").Matching(mathTypeRe).OnElements("span", "div")
	p.AllowAttrs("data-latex").Matching(regexp.MustCompile(`^[^\x00]*$`)).OnElements("span", "div")
	return p
}()

// footnoteClassRe are the classes of the footnotes (goldmark's, and the rich
// text editor's, the same HTML): a reference, the way back, the notes.
var footnoteClassRe = regexp.MustCompile(`^(footnote-ref|footnote-backref|footnotes)$`)

// tableClassRe is the class of a table not fluid (the table builder's): as
// wide as its columns, placed left, center or right (the admin's TipTap puts
// table-wrapper before it).
var tableClassRe = regexp.MustCompile(`^(table-wrapper )?mdkit-fit( mdkit-(left|center|right))?$`)

// md is Markdown as GitHub's (GFM) — its footnotes too: `text[^1]`, `[^1]:
// the note` —, its HTML kept — a table with merged cells is HTML (the
// editor's table builder) —, all of it then through the policy; the
// alignment of a table's column its cells' align (a style the policy would
// drop).
var md = goldmark.New(
	goldmark.WithExtensions(extension.Strikethrough, extension.Linkify, extension.TaskList,
		extension.NewTable(extension.WithTableCellAlignMethod(extension.TableCellAlignAttribute)),
		extension.Footnote),
	goldmark.WithRendererOptions(gmhtml.WithUnsafe()),
)

// Render is the body in HTML, safe to show (the site, the admin, the mails):
// Markdown made HTML (goldmark), or the rich text's HTML, both through the
// policy of user content (bluemonday).
func Render(format, body string) string {
	if format == Rich {
		return policy.Sanitize(body)
	}
	var buf bytes.Buffer
	if err := md.Convert([]byte(linkCards(inlineMath(mathBlocks(body)))), &buf); err != nil {
		return policy.Sanitize("<p>" + stdhtml.EscapeString(body) + "</p>")
	}
	return policy.Sanitize(buf.String())
}

// footnoteIDRe is an id of a footnote, or a link to one, in HTML rendered —
// the note's (fn:1), a reference's (fnref:1, fnref1:1: the second to it) —,
// with its prefix, when it has one (RenderIn).
var footnoteIDRe = regexp.MustCompile(`(id="|href="#)(?:[A-Za-z0-9_-]+-)?(fn(?:ref\d*)?:)`)

// prefixRe is what a prefix of ids may hold.
var prefixRe = regexp.MustCompile(`[^A-Za-z0-9_-]+`)

// RenderIn is Render of a body shown with others in a page — an issue and
// its comments —: the ids of its footnotes, and the links to them, begin
// with prefix (the record's), so a note of one is not another's.
func RenderIn(format, body, prefix string) string {
	return Reprefix(Render(format, body), prefix)
}

// Reprefix is html rendered (Render, RenderIn) with the ids of its footnotes,
// and the links to them, beginning with prefix — another place of the same
// text in a page (the views of an answer) —; "", none.
func Reprefix(html, prefix string) string {
	if !strings.Contains(html, "fn") {
		return html
	}
	if prefix = prefixRe.ReplaceAllString(prefix, ""); prefix != "" {
		prefix += "-"
	}
	return footnoteIDRe.ReplaceAllString(html, "${1}"+prefix+"${2}")
}

var (
	tags   = regexp.MustCompile(`<[^>]*>`)
	blanks = regexp.MustCompile(`[ \t]*\n[ \t\n]*\n+`)
)

// The footnotes in a text: of Markdown, a definition ([^1]: …) and a
// reference ([^1]); of HTML (the rich text editor's, as goldmark's), a
// reference (sup > a.footnote-ref), a way back (a.footnote-backref) and a
// note (li id="fn:…").
var (
	mdFootnoteDefRe    = regexp.MustCompile(`(?m)^[ ]{0,3}\[\^([^\]\s]+)\]:[ \t]*`)
	mdFootnoteRefRe    = regexp.MustCompile(`\[\^([^\]\s]+)\]`)
	htmlFootnoteRefRe  = regexp.MustCompile(`(?is)<sup[^>]*>\s*<a[^>]*class="footnote-ref"[^>]*>(.*?)</a>\s*</sup>`)
	htmlFootnoteBackRe = regexp.MustCompile(`(?is)(&#160;|&nbsp;|\x{a0})?<a[^>]*class="footnote-backref"[^>]*>.*?</a>`)
	htmlFootnoteNoteRe = regexp.MustCompile(`(?i)<li[^>]*id="(?:[A-Za-z0-9_-]+-)?fn:([^"]+)"[^>]*>(?:\s*<p>)?`)
)

// Text is the body without markup: what the search finds and a mail's plain
// part says. Markdown is already text (its footnotes [n]); the rich text's
// HTML, its text (its footnotes [n]).
func Text(format, body string) string {
	if format != Rich {
		// its footnotes: [label] where they are, and before their note
		s := mdFootnoteDefRe.ReplaceAllString(body, "[$1] ")
		return strings.TrimSpace(mdFootnoteRefRe.ReplaceAllString(s, "[$1]"))
	}
	// its footnotes: [n] where they are and before their note, the ways back
	// none
	s := htmlFootnoteRefRe.ReplaceAllString(body, "[$1]")
	s = htmlFootnoteBackRe.ReplaceAllString(s, "")
	s = htmlFootnoteNoteRe.ReplaceAllString(s, "$0[$1] ")
	s = regexp.MustCompile(`(?i)<br\s*/?>|</(p|div|h[1-6]|li|blockquote|pre)>`).ReplaceAllString(s, "\n")
	s = stdhtml.UnescapeString(tags.ReplaceAllString(s, ""))
	return strings.TrimSpace(blanks.ReplaceAllString(s, "\n\n"))
}

// Convert is body, in the format from, in the format to: Markdown made HTML
// (safe), or HTML made Markdown — what Markdown has not (colors, alignment,
// underline) is lost —. The same format: as it is.
func Convert(from, to, body string) string {
	switch {
	case from == to:
		return body
	case to == Rich:
		return Render(Markdown, body)
	default:
		return HTMLToMarkdown(body)
	}
}

// HTMLToMarkdown is the HTML of the rich text editor as Markdown: paragraphs,
// headings, bold, italic, strikethrough, code, links, images, lists, quotes,
// rules and line breaks; any other element, its text.
func HTMLToMarkdown(s string) string {
	nodes, err := html.ParseFragment(strings.NewReader(s), &html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body})
	if err != nil {
		return Text(Rich, s)
	}
	var b strings.Builder
	w := &mdWriter{b: &b}
	for _, n := range nodes {
		w.block(n, "")
	}
	return strings.TrimSpace(blanks.ReplaceAllString(b.String(), "\n\n")) + "\n"
}

type mdWriter struct {
	b *strings.Builder
}

func attr(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

// attrOK is the attribute name of n, and whether it has it (one with no
// value: details open).
func attrOK(n *html.Node, name string) (string, bool) {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val, true
		}
	}
	return "", false
}

// block writes n, a block (or text at the top), each line with prefix (a
// quote's "> ").
func (w *mdWriter) block(n *html.Node, prefix string) {
	if n.Type == html.TextNode {
		if t := strings.TrimSpace(n.Data); t != "" {
			w.b.WriteString(prefix + t + "\n\n")
		}
		return
	}
	if n.Type != html.ElementNode {
		return
	}
	switch n.Data {
	case "div", "section":
		if t := attr(n, "data-type"); t == BlockMath || t == InlineMath {
			// a formula of its own
			w.b.WriteString(prefix + "$$\n" + prefix + attr(n, "data-latex") + "\n" + prefix + "$$\n\n")
			return
		}
		if attr(n, "class") == LinkCardClass && LinkCardTarget(attr(n, "data-href")) {
			// a link card: its line
			w.b.WriteString(prefix + "[[" + attr(n, "data-href") + "]]\n\n")
			return
		}
		if attr(n, "class") == "footnotes" {
			// the notes: their definitions ([^label]: …)
			w.footnotes(n, prefix)
			return
		}
		w.b.WriteString(prefix + w.inline(n) + "\n\n")
	case "p":
		w.b.WriteString(prefix + w.inline(n) + "\n\n")
	case "h1", "h2", "h3", "h4", "h5", "h6":
		w.b.WriteString(prefix + strings.Repeat("#", int(n.Data[1]-'0')) + " " + w.inline(n) + "\n\n")
	case "ul", "ol":
		i := 1
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode || c.Data != "li" {
				continue
			}
			mark := "- "
			if n.Data == "ol" {
				mark = fmt.Sprintf("%d. ", i)
				i++
			}
			w.b.WriteString(prefix + mark + strings.TrimSpace(w.inline(c)) + "\n")
		}
		w.b.WriteString("\n")
	case "blockquote":
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			w.block(c, prefix+"> ")
		}
	case "pre":
		code := textOf(n)
		// its language (a diagram: mermaid), the code's class
		lang := ""
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "code" {
				lang = strings.TrimPrefix(attr(c, "class"), "language-")
			}
		}
		if strings.ContainsAny(lang, " \t") {
			lang = ""
		}
		w.b.WriteString(prefix + "```" + lang + "\n" + strings.TrimRight(code, "\n") + "\n" + prefix + "```\n\n")
	case "details":
		// a section that folds: as HTML, its summary, its content Markdown
		// (between blank lines)
		open := ""
		if _, ok := attrOK(n, "open"); ok {
			open = " open"
		}
		w.b.WriteString(prefix + "<details" + open + ">\n")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "summary" {
				w.b.WriteString(prefix + "<summary>" + stdhtml.EscapeString(strings.TrimSpace(textOf(c))) + "</summary>\n\n")
				continue
			}
			if c.Type == html.ElementNode && c.Data == "div" && attr(c, "data-type") == "detailsContent" {
				for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
					w.block(cc, prefix)
				}
				continue
			}
			w.block(c, prefix)
		}
		w.b.WriteString(prefix + "</details>\n\n")
	case "hr":
		w.b.WriteString(prefix + "---\n\n")
	default:
		w.b.WriteString(prefix + w.inline(n) + "\n\n")
	}
}

// inline is the content of n as inline Markdown.
func (w *mdWriter) inline(n *html.Node) string {
	var b strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch {
		case c.Type == html.TextNode:
			b.WriteString(c.Data)
		case c.Type != html.ElementNode:
		case c.Data == "strong" || c.Data == "b":
			b.WriteString("**" + w.inline(c) + "**")
		case c.Data == "em" || c.Data == "i":
			b.WriteString("*" + w.inline(c) + "*")
		case c.Data == "s" || c.Data == "del":
			b.WriteString("~~" + w.inline(c) + "~~")
		case attr(c, "data-type") == InlineMath:
			// a formula in the line
			b.WriteString("$" + attr(c, "data-latex") + "$")
		case c.Data == "code":
			b.WriteString("`" + textOf(c) + "`")
		case c.Data == "a" && attr(c, "class") == "footnote-backref":
			// a way back of a note: none
		case c.Data == "sup" && footnoteRefLabel(c) != "":
			b.WriteString("[^" + footnoteRefLabel(c) + "]")
		case c.Data == "a":
			b.WriteString("[" + w.inline(c) + "](" + attr(c, "href") + ")")
		case c.Data == "img":
			b.WriteString("![" + attr(c, "alt") + "](" + attr(c, "src") + ")")
		case c.Data == "br":
			b.WriteString("  \n")
		case c.Data == "p" || c.Data == "ul" || c.Data == "ol":
			// a list item's paragraph or list: inline here
			b.WriteString(w.inline(c))
		default:
			b.WriteString(w.inline(c))
		}
	}
	return b.String()
}

// footnoteLabelRe is the label of a note in an id or a link to it
// ("fn:label", "#c1-fn:label").
var footnoteLabelRe = regexp.MustCompile(`fn(?:ref\d*)?:(.+)$`)

// footnoteRefLabel is the label of the note n (a sup) refers to: its link's
// (a.footnote-ref); "" not a reference.
func footnoteRefLabel(n *html.Node) string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "a" && attr(c, "class") == "footnote-ref" {
			if m := footnoteLabelRe.FindStringSubmatch(attr(c, "href")); m != nil {
				return m[1]
			}
		}
	}
	return ""
}

// footnotes writes the notes of n (div.footnotes) as their definitions.
func (w *mdWriter) footnotes(n *html.Node, prefix string) {
	var walk func(*html.Node)
	walk = func(x *html.Node) {
		for c := x.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}
			if c.Data == "li" {
				if m := footnoteLabelRe.FindStringSubmatch(attr(c, "id")); m != nil {
					text := strings.TrimSpace(strings.TrimRight(strings.TrimSpace(w.inline(c)), "\u00a0"))
					w.b.WriteString(prefix + "[^" + m[1] + "]: " + strings.ReplaceAll(text, "\n", "\n    ") + "\n")
				}
				continue
			}
			walk(c)
		}
	}
	walk(n)
	w.b.WriteString("\n")
}

func textOf(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var b strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b.WriteString(textOf(c))
	}
	return b.String()
}

// The hints of a template of a body, as GitHub's: comments — `<!-- Say what
// happened -->` — that tell the user what to write where, and need not be
// erased: Markdown keeps them in its box and does not show them; the rich
// text editor shows each as the placeholder of an empty paragraph (its
// data-hint, TemplateHTML); what is sent is kept without them (StripHints).
var (
	hintRe = regexp.MustCompile(`(?s)<!--\s*(.*?)\s*-->`)
	// a hint on lines of its own (a paragraph of the template)
	blockHintRe = regexp.MustCompile(`(?ms)^[ \t]*<!--\s*(.*?)\s*-->[ \t]*$`)
)

// StripHints is body without its comments (the hints of its template), the
// blank lines they leave squeezed.
func StripHints(body string) string {
	if !strings.Contains(body, "<!--") {
		return body
	}
	s := hintRe.ReplaceAllString(body, "")
	return strings.TrimSpace(blanks.ReplaceAllString(s, "\n\n"))
}

// hintToken is the stand-in of the n-th hint while the template is rendered:
// a word Markdown and the policy keep as it is.
func hintToken(n int) string { return fmt.Sprintf("HERMONHINT%dX", n) }

// TemplateHTML is the template of a body (Markdown, or the editor's HTML)
// for the rich text editor: its HTML (Render), each hint on lines of its
// own — or between tags — an empty paragraph that tells it (`<p
// data-hint="…"></p>`), placed where it was; one within a line of text, none.
func TemplateHTML(body string) string {
	var hints []string
	token := func(re *regexp.Regexp) func(string) string {
		return func(m string) string {
			hints = append(hints, strings.TrimSpace(re.FindStringSubmatch(m)[1]))
			return "\n\n" + hintToken(len(hints)-1) + "\n\n"
		}
	}
	src := blockHintRe.ReplaceAllStringFunc(body, token(blockHintRe))
	if strings.Contains(hintRe.ReplaceAllString(src, ""), "<") {
		// the editor's HTML: between tags, too
		src = hintRe.ReplaceAllStringFunc(src, token(hintRe))
	} else {
		src = hintRe.ReplaceAllString(src, "")
	}
	format := Markdown
	if strings.Contains(src, "<") {
		format = Rich
		for i := range hints {
			src = strings.Replace(src, "\n\n"+hintToken(i)+"\n\n", "<p>"+hintToken(i)+"</p>", 1)
		}
	}
	out := Render(format, src)
	for i, hint := range hints {
		tok := hintToken(i)
		out = strings.Replace(out, "<p>"+tok+"</p>", `<p data-hint="`+stdhtml.EscapeString(hint)+`"></p>`, 1)
		// within a paragraph of text (no line of its own): none
		out = strings.Replace(out, tok, "", 1)
	}
	return out
}
