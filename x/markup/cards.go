package markup

import (
	stdhtml "html"
	"regexp"
	"strings"
)

// A link card: a page of the site, a record of the admin, or a page of
// another site, linked on a line of its own — in Markdown `[[/en-us/about]]`
// (a record: `[[admin:posts:<id>:]]`, a path of presets' ParseRecordPath;
// another site's: `[[https://example.com/a]]`) —, made
// HTML a link marked (its class LinkCardClass, its target data-href) that
// the script of the page makes a card: its title, its cover, its summary,
// in the language of the page, by who may see it (else the link stays).

// LinkCardClass is the class of a link card.
const LinkCardClass = "mdkit-link-card"

var (
	// cardLineRe is a line of a link card in Markdown
	cardLineRe = regexp.MustCompile(`^[ ]{0,3}\[\[([^\s\[\]]+)\]\][ \t]*$`)
	// cardTargetRe is what a card links: a path of the site (not of another
	// site: //host), a record path of the admin (admin:…:, its action
	// after), an address of another site (http or https), never a script
	cardTargetRe = regexp.MustCompile(`^(/|/[^/\\\s"'<>][^\s"'<>]*|admin:[^\s"'\\]*:(![A-Za-z0-9_]+)?|https?://[^/\s"'<>\\][^\s"'<>\\]*)$`)
)

// LinkCardTarget says target may be linked by a card.
func LinkCardTarget(target string) bool {
	return cardTargetRe.MatchString(target)
}

// LinkCardHTML is the HTML of the card of target: a link marked — to the
// page, a path of the site or of another site (in a new tab); a record's,
// its path as text (the script makes it the admin's link) —.
func LinkCardHTML(target string) string {
	t := stdhtml.EscapeString(target)
	inner := t
	switch {
	case strings.HasPrefix(target, "/"):
		inner = `<a href="` + t + `">` + t + `</a>`
	case strings.HasPrefix(target, "http://"), strings.HasPrefix(target, "https://"):
		// another site's: in a new tab, telling it nothing of this one
		inner = `<a href="` + t + `" target="_blank" rel="noopener noreferrer nofollow">` + t + `</a>`
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
