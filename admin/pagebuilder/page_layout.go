package pagebuilder

import (
	"fmt"
	"strings"

	"github.com/go-rvq/rvq/web"
	h "github.com/theplant/htmlgo"
)

func DefaultPageLayoutFunc(body h.HTMLComponent, input *PageLayoutInput, ctx *web.EventContext) h.HTMLComponent {
	var freeStyleCss h.HTMLComponent
	if len(input.FreeStyleCss) > 0 {
		freeStyleCss = h.Style(strings.Join(input.FreeStyleCss, "\n"))
	}

	head := h.Components(
		h.Meta().Attr("charset", "utf-8"),
		input.SeoTags,
		input.CanonicalLink,
		h.Meta().Attr("http-equiv", "X-UA-Compatible").Content("IE=edge"),
		h.Meta().Content("true").Name("HandheldFriendly"),
		h.Meta().Content("yes").Name("apple-mobile-web-app-capable"),
		h.Meta().Content("black").Name("apple-mobile-web-app-status-bar-style"),
		h.Meta().Name("format-detection").Content("telephone=no"),
		h.Meta().Name("viewport").Content("width=device-width, initial-scale=1"),
		h.If(len(input.EditorCss) > 0, input.EditorCss...),
		freeStyleCss,
		// RawHTML(dataLayer),
		input.StructuredData,
		scriptWithCodes(input.FreeStyleTopJs),
	)
	ctx.Injector.HTMLLang(input.LocaleCode)
	ctx.Injector.HeadHTML(h.MustString(head, nil))

	return h.Body(
		// It's required as the body first element!
		h.If(input.Header != nil, input.Header),
		body,
		h.If(input.Footer != nil, input.Footer),
		scriptWithCodes(input.FreeStyleBottomJs),
	).Attr("data-site-domain", ctx.R.Header.Get("Qor5-Site-Domain"))
}

func scriptWithCodes(jscodes []string) h.HTMLComponent {
	var js h.HTMLComponent
	if len(jscodes) > 0 {
		js = h.Script(fmt.Sprintf(`
try {
	%s
} catch (error) {
	console.log(error);
}
`, strings.Join(jscodes, "\n")))
	}
	return js
}
