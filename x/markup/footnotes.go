package markup

import (
	"regexp"
	"sort"
	"strings"
)

// The footnotes of a text, as a map — each note, its label, its number, its
// references, its text — (Footnotes), and the text with its labels changed
// by a map (RelabelFootnotes): texts put together (the chapters of a book)
// keep their notes apart, each text's labels made its own first
// (FootnoteLabels).
//
// A label is what a note is named by: in Markdown, `text[^label]` and
// `[^label]: the note`; in the rich text editor's HTML (the HTML goldmark
// makes, its ids the labels), `<sup id="fnref:label"><a href="#fn:label"
// class="footnote-ref">1</a></sup>` and `<li id="fn:label">the note …</li>`.
// It holds letters, digits, "_", "." and "-" (what an id the policy keeps
// may hold); its number is the order of its first reference.

// Footnote is a footnote of a text.
type Footnote struct {
	// Label is what it is named by
	Label string
	// Number is how it is shown: the order of its first reference; 0, none
	Number int
	// Refs is how many references it has
	Refs int
	// Text is its note, in the text's format (Markdown; HTML); "" none
	Text string
	// Defined says it has a note
	Defined bool
}

// labelRe is what a label may not hold.
var labelRe = regexp.MustCompile(`[^A-Za-z0-9_.-]+`)

// CleanFootnoteLabel is s made a label: what a label may not hold, out.
func CleanFootnoteLabel(s string) string { return labelRe.ReplaceAllString(s, "") }

// mdSpan is a label in a Markdown text: where it is ([start, end) of the
// label itself), whether of a note's definition.
type mdSpan struct {
	start, end int
	label      string
	def        bool
}

var (
	fenceOpenRe  = regexp.MustCompile("^[ ]{0,3}(`{3,}|~{3,})")
	mdDefLineRe  = regexp.MustCompile(`^[ ]{0,3}\[\^([^\]\s]+)\]:[ \t]?`)
	mdRefLabelRe = regexp.MustCompile(`\[\^([^\]\s]+)\]`)
)

// codeRanges are the parts of a Markdown text that are code — its fenced
// blocks, its spans (a run of backticks to the next run as long) —: their
// labels are none.
func codeRanges(s string) (out [][2]int) {
	// the blocks: from a fence to one as long of its mark, or the end
	pos := 0
	fence := ""
	start := 0
	for _, line := range strings.SplitAfter(s, "\n") {
		trimmed := strings.TrimLeft(line, " ")
		if fence == "" {
			if m := fenceOpenRe.FindStringSubmatch(line); m != nil {
				fence, start = m[1], pos
			}
		} else if strings.HasPrefix(trimmed, fence) && strings.TrimSpace(strings.TrimLeft(trimmed, fence[:1])) == "" {
			out = append(out, [2]int{start, pos + len(line)})
			fence = ""
		}
		pos += len(line)
	}
	if fence != "" {
		out = append(out, [2]int{start, len(s)})
	}
	// the spans, out of the blocks
	for i := 0; i < len(s); {
		if s[i] != '`' || inRanges(out, i) {
			i++
			continue
		}
		n := 0
		for i+n < len(s) && s[i+n] == '`' {
			n++
		}
		end := -1
		for j := i + n; j < len(s); {
			if s[j] != '`' {
				j++
				continue
			}
			k := 0
			for j+k < len(s) && s[j+k] == '`' {
				k++
			}
			if k == n {
				end = j + k
				break
			}
			j += k
		}
		if end < 0 {
			i += n
			continue
		}
		out = append(out, [2]int{i, end})
		i = end
	}
	return
}

func inRanges(rs [][2]int, i int) bool {
	for _, r := range rs {
		if i >= r[0] && i < r[1] {
			return true
		}
	}
	return false
}

// mdSpans are the labels of a Markdown text, in order; and the text of each
// note (its definition's line and those that go on with it, indented).
func mdSpans(s string) (spans []mdSpan, notes map[string]string) {
	notes = map[string]string{}
	code := codeRanges(s)
	defStarts := map[int]bool{}
	lines := strings.SplitAfter(s, "\n")
	pos := 0
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if m := mdDefLineRe.FindStringSubmatchIndex(line); m != nil && !inRanges(code, pos) {
			label := line[m[2]:m[3]]
			spans = append(spans, mdSpan{pos + m[2], pos + m[3], label, true})
			defStarts[pos+m[2]-2] = true
			text := []string{strings.TrimRight(line[m[1]:], "\n")}
			// the lines that go on with it: indented, or blank before indented
			for i+1 < len(lines) {
				next := lines[i+1]
				blank := strings.TrimSpace(next) == ""
				if blank && i+2 < len(lines) && strings.HasPrefix(lines[i+2], "    ") {
					text = append(text, "")
				} else if !blank && (strings.HasPrefix(next, "    ") || strings.HasPrefix(next, "\t")) {
					text = append(text, strings.TrimRight(strings.TrimPrefix(strings.TrimPrefix(next, "    "), "\t"), "\n"))
				} else {
					break
				}
				pos += len(lines[i])
				i++
			}
			if _, ok := notes[label]; !ok {
				notes[label] = strings.TrimSpace(strings.Join(text, "\n"))
			}
		}
		pos += len(lines[i])
	}
	for _, m := range mdRefLabelRe.FindAllStringSubmatchIndex(s, -1) {
		if inRanges(code, m[0]) || defStarts[m[0]] {
			continue
		}
		// [^x]: at the start of a line is a definition (taken above)
		if len(s) > m[1] && s[m[1]] == ':' {
			continue
		}
		spans = append(spans, mdSpan{m[2], m[3], s[m[2]:m[3]], false})
	}
	sort.Slice(spans, func(i, j int) bool { return spans[i].start < spans[j].start })
	return
}

// The footnotes in the HTML of the rich text editor (goldmark's), by their
// ids: a reference (fnref:label, fnref2:label), a link to a note (#fn:label),
// a note (li id="fn:label"), a way back to a reference (#fnref:label).
var (
	htmlFnAttrRe = regexp.MustCompile(`(id="fnref\d*:|href="#fn:|id="fn:|href="#fnref\d*:)([^"]+)"`)
	htmlRefIDRe  = regexp.MustCompile(`<sup[^>]*\bid="fnref\d*:([^"]+)"`)
	htmlNoteRe   = regexp.MustCompile(`(?is)<li[^>]*\bid="fn:([^"]+)"[^>]*>(.*?)</li>`)
)

// Footnotes is the map of the footnotes of body in format: the referenced
// ones by their numbers, then those with no reference, in the order of their
// notes; each with its text.
func Footnotes(format, body string) []Footnote {
	byLabel := map[string]*Footnote{}
	var order []string
	note := func(label string) *Footnote {
		f, ok := byLabel[label]
		if !ok {
			f = &Footnote{Label: label}
			byLabel[label] = f
			order = append(order, label)
		}
		return f
	}
	number := 0
	ref := func(label string) {
		f := note(label)
		f.Refs++
		if f.Number == 0 {
			number++
			f.Number = number
		}
	}
	if format == Rich {
		for _, m := range htmlRefIDRe.FindAllStringSubmatch(body, -1) {
			ref(m[1])
		}
		for _, m := range htmlNoteRe.FindAllStringSubmatch(body, -1) {
			f := note(m[1])
			f.Defined = true
			f.Text = strings.TrimSpace(htmlFootnoteBackRe.ReplaceAllString(m[2], ""))
		}
	} else {
		spans, notes := mdSpans(body)
		for _, sp := range spans {
			if sp.def {
				f := note(sp.label)
				f.Defined = true
				f.Text = notes[sp.label]
			} else {
				ref(sp.label)
			}
		}
	}
	out := make([]Footnote, 0, len(order))
	for _, l := range order {
		out = append(out, *byLabel[l])
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i].Number, out[j].Number
		switch {
		case a == 0 && b == 0:
			return false
		case a == 0:
			return false
		case b == 0:
			return true
		}
		return a < b
	})
	return out
}

// RelabelFootnotes is body in format with the labels of its footnotes changed
// as labels says (old → new; the others as they are) — in its references and
// in its notes —, each new one made a label (CleanFootnoteLabel).
func RelabelFootnotes(format, body string, labels map[string]string) string {
	if len(labels) == 0 {
		return body
	}
	to := func(l string) (string, bool) {
		n, ok := labels[l]
		if !ok {
			return "", false
		}
		if n = CleanFootnoteLabel(n); n == "" {
			return "", false
		}
		return n, true
	}
	if format == Rich {
		return htmlFnAttrRe.ReplaceAllStringFunc(body, func(m string) string {
			sm := htmlFnAttrRe.FindStringSubmatch(m)
			if n, ok := to(sm[2]); ok {
				return sm[1] + n + `"`
			}
			return m
		})
	}
	spans, _ := mdSpans(body)
	var b strings.Builder
	last := 0
	for _, sp := range spans {
		n, ok := to(sp.label)
		if !ok {
			continue
		}
		b.WriteString(body[last:sp.start])
		b.WriteString(n)
		last = sp.end
	}
	b.WriteString(body[last:])
	return b.String()
}

// FootnoteLabels is the map that makes the labels of the footnotes of body
// its own — prefix and the label ("cap2-1") —: for a text put with others
// (RelabelFootnotes).
func FootnoteLabels(format, body, prefix string) map[string]string {
	prefix = CleanFootnoteLabel(prefix)
	out := map[string]string{}
	for _, f := range Footnotes(format, body) {
		if prefix != "" {
			out[f.Label] = prefix + "-" + f.Label
		}
	}
	return out
}
