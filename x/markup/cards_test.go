package markup

import (
	"strings"
	"testing"
)

// TestLinkCards: a line [[target]] — a path of the site, a record path of
// the admin — is a link card (marked, its target kept); in code, in a
// paragraph, or of another site or a script, it is not; the rich text's
// card kept; made Markdown again, its line.
func TestLinkCards(t *testing.T) {
	got := Render(Markdown, "Before\n[[/en-us/about]]\nafter\n\n[[admin:posts:<01a0_en-US>:]]\n\n```\n[[/in-code]]\n```\n\nsee [[/inline]] here\n\n[[//evil.com/x]]\n\n[[javascript:alert(1)]]\n")
	for _, want := range []string{
		`<div class="mdkit-link-card" data-href="/en-us/about"><a href="/en-us/about" rel="nofollow">/en-us/about</a></div>`,
		`<div class="mdkit-link-card" data-href="admin:posts:&lt;01a0_en-US&gt;:">admin:posts:&lt;01a0_en-US&gt;:</div>`,
		"[[/in-code]]", "see [[/inline]] here", "[[//evil.com/x]]", "[[javascript:alert(1)]]",
		"<p>Before</p>", "<p>after</p>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("no %s in\n%s", want, got)
		}
	}
	if strings.Count(got, "mdkit-link-card") != 2 {
		t.Errorf("cards: %d\n%s", strings.Count(got, "mdkit-link-card"), got)
	}

	rich := Render(Rich, `<div class="mdkit-link-card" data-href="/en-us/about"><a href="/en-us/about">About</a></div><div class="mdkit-link-card" data-href="javascript:x">y</div>`)
	if !strings.Contains(rich, `data-href="/en-us/about"`) || strings.Contains(rich, "javascript") {
		t.Errorf("rich: %s", rich)
	}
	if md := HTMLToMarkdown(`<p>a</p><div class="mdkit-link-card" data-href="/en-us/about"><a href="/en-us/about">About</a></div>`); !strings.Contains(md, "\n\n[[/en-us/about]]\n") {
		t.Errorf("to Markdown: %q", md)
	}
}
