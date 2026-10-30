package markup

import (
	stdhtml "html"
	"regexp"
	"strings"
)

// The formulas of a text, as GitHub's (LaTeX, drawn by KaTeX on the page:
// rvq mdkit's blocks): inline, `$…$` — the $ against the formula at both
// ends, the closing one not before a digit ($5 and $10 is money) —, or
// $`…`$; a block, `$$` on lines of their own around it, `$$…$$` alone on a
// line, or a ```math fence. Made the HTML the rich text editor's
// (tiptap's mathematics) has — span[data-type=inline-math],
// div[data-type=block-math], the formula in data-latex, and as its text
// until it is drawn —; out of code.

const (
	// InlineMath and BlockMath are the data-type of a formula.
	InlineMath = "inline-math"
	BlockMath  = "block-math"
)

var (
	mathTypeRe = regexp.MustCompile(`^(` + InlineMath + `|` + BlockMath + `)$`)
	// inlineMathRe: $`…`$, or $…$ — not $ $ (a space after the first or
	// before the last), the last not before a digit nor a $
	inlineMathRe = regexp.MustCompile("\\$`([^`\\n]+)`\\$|\\$([^\\s$](?:[^$\\n]*[^\\s$\\\\])?)\\$")
	blockMathLineRe = regexp.MustCompile(`^[ ]{0,3}\$\$(.*\S.*)\$\$[ \t]*$`)
	mathFenceRe     = regexp.MustCompile("^[ ]{0,3}(`{3,}|~{3,})[ \\t]*math[ \\t]*$")
)

// InlineMathHTML is the HTML of an inline formula.
func InlineMathHTML(latex string) string {
	l := stdhtml.EscapeString(latex)
	return `<span data-type="` + InlineMath + `" data-latex="` + l + `">` + l + `</span>`
}

// BlockMathHTML is the HTML of a formula of its own.
func BlockMathHTML(latex string) string {
	l := stdhtml.EscapeString(strings.TrimSpace(latex))
	return `<div data-type="` + BlockMath + `" data-latex="` + l + `">` + l + `</div>`
}

// mathBlocks is the Markdown body with its blocks of formulas — $$ lines,
// ```math fences — their HTML (out of the other code).
func mathBlocks(body string) string {
	if !strings.Contains(body, "$$") && !strings.Contains(body, "math") {
		return body
	}
	lines := strings.SplitAfter(body, "\n")
	var b strings.Builder
	fence := ""
	for i := 0; i < len(lines); i++ {
		line := strings.TrimRight(lines[i], "\r\n")
		if fence != "" {
			// a fence that is not math: as it is, to its end
			b.WriteString(lines[i])
			if t := strings.TrimLeft(line, " "); strings.HasPrefix(t, fence) && strings.TrimSpace(strings.TrimLeft(t, fence[:1])) == "" {
				fence = ""
			}
			continue
		}
		if m := mathFenceRe.FindStringSubmatch(line); m != nil {
			// a math fence: to the closing one
			var tex []string
			j := i + 1
			for ; j < len(lines); j++ {
				t := strings.TrimLeft(strings.TrimRight(lines[j], "\r\n"), " ")
				if strings.HasPrefix(t, m[1]) && strings.TrimSpace(strings.TrimLeft(t, m[1][:1])) == "" {
					break
				}
				tex = append(tex, strings.TrimRight(lines[j], "\r\n"))
			}
			b.WriteString("\n" + BlockMathHTML(strings.Join(tex, "\n")) + "\n\n")
			i = j
			continue
		}
		if m := fenceOpenRe.FindStringSubmatch(line); m != nil {
			fence = m[1]
			b.WriteString(lines[i])
			continue
		}
		if m := blockMathLineRe.FindStringSubmatch(line); m != nil {
			b.WriteString("\n" + BlockMathHTML(m[1]) + "\n\n")
			continue
		}
		if strings.TrimSpace(line) == "$$" {
			// $$ to the next $$ line
			var tex []string
			j := i + 1
			for ; j < len(lines) && strings.TrimSpace(lines[j]) != "$$"; j++ {
				tex = append(tex, strings.TrimRight(lines[j], "\r\n"))
			}
			if j < len(lines) {
				b.WriteString("\n" + BlockMathHTML(strings.Join(tex, "\n")) + "\n\n")
				i = j
				continue
			}
		}
		b.WriteString(lines[i])
	}
	return b.String()
}

// inlineMath is the Markdown body with its inline formulas their HTML (out
// of code, not \$).
func inlineMath(body string) string {
	if !strings.Contains(body, "$") {
		return body
	}
	code := codeRanges(body)
	var b strings.Builder
	last := 0
	for _, m := range inlineMathRe.FindAllStringSubmatchIndex(body, -1) {
		start, end := m[0], m[1]
		if start < last || inRanges(code, start) || (start > 0 && body[start-1] == '\\') {
			continue
		}
		// the closing $ before a digit: money, not a formula
		if end < len(body) && body[end] >= '0' && body[end] <= '9' {
			continue
		}
		tex := ""
		if m[2] >= 0 {
			tex = body[m[2]:m[3]]
		} else {
			tex = body[m[4]:m[5]]
		}
		b.WriteString(body[last:start])
		b.WriteString(InlineMathHTML(tex))
		last = end
	}
	b.WriteString(body[last:])
	return b.String()
}
