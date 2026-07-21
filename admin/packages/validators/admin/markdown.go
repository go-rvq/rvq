package admin

import (
	"html"
	"regexp"
	"strings"
)

// mdToHTML renders a small, safe subset of Markdown to HTML for the validator
// documentation shown in the detail view: ATX headings, fenced code blocks,
// unordered/ordered lists, blockquotes, paragraphs, and the inline spans
// **bold**, `code` and [text](url). All text is HTML-escaped; no raw HTML from
// the source is passed through. It is intentionally dependency-free.
func mdToHTML(md string) string {
	lines := strings.Split(strings.ReplaceAll(md, "\r\n", "\n"), "\n")
	var (
		out      strings.Builder
		inCode   bool
		codeBuf  []string
		listKind string // "ul" | "ol" | ""
		para     []string
	)

	flushPara := func() {
		if len(para) > 0 {
			out.WriteString("<p>" + inlineMD(strings.Join(para, " ")) + "</p>\n")
			para = nil
		}
	}
	closeList := func() {
		if listKind != "" {
			out.WriteString("</" + listKind + ">\n")
			listKind = ""
		}
	}

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			if inCode {
				out.WriteString("<pre><code>" + html.EscapeString(strings.Join(codeBuf, "\n")) + "</code></pre>\n")
				codeBuf = nil
				inCode = false
			} else {
				flushPara()
				closeList()
				inCode = true
			}
			continue
		}
		if inCode {
			codeBuf = append(codeBuf, line)
			continue
		}

		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "":
			flushPara()
			closeList()
		case headingRE.MatchString(trimmed):
			flushPara()
			closeList()
			m := headingRE.FindStringSubmatch(trimmed)
			level := len(m[1])
			out.WriteString("<h" + itoa(level) + ">" + inlineMD(m[2]) + "</h" + itoa(level) + ">\n")
		case ulRE.MatchString(trimmed):
			flushPara()
			if listKind != "ul" {
				closeList()
				out.WriteString("<ul>\n")
				listKind = "ul"
			}
			out.WriteString("<li>" + inlineMD(ulRE.FindStringSubmatch(trimmed)[1]) + "</li>\n")
		case olRE.MatchString(trimmed):
			flushPara()
			if listKind != "ol" {
				closeList()
				out.WriteString("<ol>\n")
				listKind = "ol"
			}
			out.WriteString("<li>" + inlineMD(olRE.FindStringSubmatch(trimmed)[1]) + "</li>\n")
		case strings.HasPrefix(trimmed, ">"):
			flushPara()
			closeList()
			out.WriteString("<blockquote>" + inlineMD(strings.TrimSpace(trimmed[1:])) + "</blockquote>\n")
		default:
			para = append(para, trimmed)
		}
	}
	if inCode {
		out.WriteString("<pre><code>" + html.EscapeString(strings.Join(codeBuf, "\n")) + "</code></pre>\n")
	}
	flushPara()
	closeList()
	return out.String()
}

var (
	headingRE = regexp.MustCompile(`^(#{1,6})\s+(.*)$`)
	ulRE      = regexp.MustCompile(`^[-*]\s+(.*)$`)
	olRE      = regexp.MustCompile(`^\d+\.\s+(.*)$`)
	boldRE    = regexp.MustCompile(`\*\*(.+?)\*\*`)
	codeRE    = regexp.MustCompile("`([^`]+)`")
	linkRE    = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
)

// inlineMD escapes text and applies the inline spans. Code spans are protected
// from the other transforms via placeholders.
func inlineMD(s string) string {
	var codes []string
	s = codeRE.ReplaceAllStringFunc(s, func(m string) string {
		inner := codeRE.FindStringSubmatch(m)[1]
		codes = append(codes, "<code>"+html.EscapeString(inner)+"</code>")
		return "\x00" + itoa(len(codes)-1) + "\x00"
	})
	s = html.EscapeString(s)
	s = boldRE.ReplaceAllString(s, "<strong>$1</strong>")
	s = linkRE.ReplaceAllStringFunc(s, func(m string) string {
		g := linkRE.FindStringSubmatch(m)
		return `<a href="` + html.EscapeString(g[2]) + `" target="_blank" rel="noreferrer">` + g[1] + `</a>`
	})
	for i, c := range codes {
		s = strings.ReplaceAll(s, "\x00"+itoa(i)+"\x00", c)
	}
	return s
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [8]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
