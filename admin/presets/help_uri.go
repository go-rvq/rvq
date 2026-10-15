package presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

// HelpURIAttr is the attribute of an overlay — a drawer, a dialog — that says
// the address that filled it (the event's request: its path and query): the
// help of what is in it is the help of that address (the user docs' help
// bar).
const HelpURIAttr = "data-help-uri"

// HelpURI is the address of the request of ctx, as an overlay it fills says
// it (HelpURIAttr).
func HelpURI(ctx *web.EventContext) string {
	if ctx == nil || ctx.R == nil {
		return ""
	}
	return ctx.R.URL.RequestURI()
}

// markHelpURI puts uri on comp (HelpURIAttr) — on its tag, out of a scope —:
// nothing when comp is no tag or uri is "".
func markHelpURI(comp h.HTMLComponent, uri string) {
	if uri == "" {
		return
	}
	if t, _ := web.Unscoped(comp).(h.TagGetter); t != nil {
		t.GetHTMLTagBuilder().Attr(HelpURIAttr, uri)
	}
}
