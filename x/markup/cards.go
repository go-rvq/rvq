package markup

import (
	stdhtml "html"
	"regexp"
	"strings"
)

// A link card: a page of the site, or a record of the admin, linked on a
// line of its own — in Markdown `[[/en-us/about]]` (a record:
// `[[admin:posts:<id>:]]`, a path of presets' ParseRecordPath) —, made
// HTML a link marked (its class LinkCardClass, its target data-href) that
// the script of the page makes a card: its title, its cover, its summary,
// in the language of the page, by who may see it (else the link stays).

// LinkCardClass is the class of a link card.
const LinkCardClass = "mdkit-link-card"

var (
	// cardLineRe is a line of a link card in Markdown
	cardLineRe = regexp.MustCompile(`^[ ]{0,3}\[\[([^\s\[\]]+)\]\][ \t]*$`)
	// cardTargetRe is what a card links: a path of the site (not of another
	// site: //host), or a record path of the admin (admin:…:, its action
	// after), never a script
	cardTargetRe = regexp.MustCompile(`^(/|/[^/\\\s"'<>][^\s"'<>]*|admin:[^\s"'\\]*:(![A-Za-z0-9_]+)?)$`)
)

// LinkCardTarget says target may be linked by a card.
func LinkCardTarget(target string) bool {
	return cardTargetRe.MatchString(target)
}

// LinkCardHTML is the HTML of the card of target: a link marked — to the
// page, a path of the site; a record's, its path as text (the script makes
// it the admin's link) —.
func LinkCardHTML(target string) string {
	t := stdhtml.EscapeString(target)
	inner := t
	if strings.HasPrefix(target, "/") {
		inner = `<a href="` + t + `">` + t + `</a>`
	}
	return `<div class="` + LinkCardClass + `" data-href="` + t + `">` + inner + `</div>`
}

// linkCards is the Markdown body with each line of a card — out of code —
// its HTML, a block of its own.
func linkCards(body string) string {
	if !strings.Contains(body, "[[") {
		return body
	}
	code := codeRanges(body)
	var b strings.Builder
	pos := 0
	for _, line := range strings.SplitAfter(body, "\n") {
		at := pos
		pos += len(line)
		m := cardLineRe.FindStringSubmatch(strings.TrimRight(line, "\r\n"))
		if m == nil || inRanges(code, at) || !LinkCardTarget(m[1]) {
			b.WriteString(line)
			continue
		}
		b.WriteString("\n" + LinkCardHTML(m[1]) + "\n\n")
	}
	return b.String()
}
