package markup

import (
	"strings"
	"testing"
)

// The buttons and the tabs: kept by the policy — their kinds, place, look,
// targets —, what is not theirs dropped (a color, an icon, a target not
// known; a script); in Markdown a block of HTML, a tab's content Markdown;
// made Markdown and back the same.
func TestWidgets(t *testing.T) {
	buttons := `<div data-type="buttons" data-align="center">` +
		`<a data-type="button" href="/en-us/about" data-href="/en-us/about" title="Who we are" data-color="primary" data-icon-left="mdi-information"></a>` +
		`<a data-type="button" href="https://x.test/" data-icon-right="mdi-open-in-new" target="_blank" rel="noopener">Site</a></div>`
	got := Render(Rich, buttons)
	for _, want := range []string{`data-type="buttons"`, `data-align="center"`, `data-type="button"`, `data-href="/en-us/about"`,
		`title="Who we are"`, `data-color="primary"`, `data-icon-left="mdi-information"`, `data-icon-right="mdi-open-in-new"`, `target="_blank"`, `>Site</a>`} {
		if !strings.Contains(got, want) {
			t.Errorf("buttons: no %s\n%s", want, got)
		}
	}
	bad := Render(Rich, `<div data-type="buttons" data-align="left"><a data-type="button" href="javascript:alert(1)" data-color="pink" data-icon-left="x y" data-href="//evil" target="_top" onclick="x()">X</a></div>`)
	for _, not := range []string{"javascript", "pink", "x y", "evil", "_top", "onclick", `data-align="left"`} {
		if strings.Contains(bad, not) {
			t.Errorf("buttons kept %s: %s", not, bad)
		}
	}

	tabs := `<div data-type="tabs" data-align="end"><div data-type="tab" data-title="One &#34;1&#34;" data-color="success" data-icon-right="mdi-star">` +
		`<p><strong>bold</strong></p><ul><li><p>a</p></li></ul></div><div data-type="tab" data-title="Two"><p>two</p></div></div>`
	got = Render(Rich, tabs)
	for _, want := range []string{`data-type="tabs"`, `data-align="end"`, `data-type="tab"`, `data-title="One &#34;1&#34;"`, `data-color="success"`, `<strong>bold</strong>`} {
		if !strings.Contains(got, want) {
			t.Errorf("tabs: no %s\n%s", want, got)
		}
	}

	// made Markdown: the tabs' HTML lines, their contents Markdown
	md := Convert(Rich, Markdown, buttons+tabs)
	for _, want := range []string{"<div data-type=\"buttons\" data-align=\"center\">",
		"<div data-type=\"tabs\" data-align=\"end\">\n<div data-type=\"tab\" data-title=\"One &#34;1&#34;\" data-color=\"success\" data-icon-right=\"mdi-star\">\n\n**bold**\n\n- a\n\n</div>\n<div data-type=\"tab\" data-title=\"Two\">\n\ntwo\n\n</div>\n</div>"} {
		if !strings.Contains(md, want) {
			t.Errorf("markdown: no %q\n%s", want, md)
		}
	}
	// and its Markdown rendered: the tabs, their contents HTML again
	back := Render(Markdown, md)
	for _, want := range []string{`data-type="tabs"`, `data-title="One &#34;1&#34;"`, `<strong>bold</strong>`, `<p>two</p>`, `data-href="/en-us/about"`} {
		if !strings.Contains(back, want) {
			t.Errorf("rendered back: no %s\n%s", want, back)
		}
	}
}
