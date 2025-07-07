package web

import h "github.com/go-rvq/htmlgo"

func NoopLayoutFunc(in PageFunc) PageFunc {
	return in
}

func defaultLayoutFunc(in PageFunc) PageFunc {
	return func(ctx *EventContext) (r PageResponse, err error) {
		r, err = in(ctx)
		if ctx.W.Writed() || r.RedirectURL != "" {
			return
		}

		if err != nil {
			r.Body = h.Div(h.Text(err.Error()))
			err = nil
		} else if r.PageTitle != "" {
			ctx.Injector.Title(r.PageTitle)
		}

		var body []byte
		if body, err = r.Body.MarshalHTML(WrapEventContext(ctx.Context(), ctx)); err != nil {
			return
		}

		r.Body = h.HTMLComponents{
			h.RawHTML("<!DOCTYPE html>\n"),
			h.Tag("html").Children(
				h.Head(
					ctx.Injector.GetHeadHTMLComponent(),
				),
				h.Body(
					h.Div(
						// NOTES:
						// 1. put body on portal, because vue uses #app.innerHTML for build app template.
						// innerHTML replaces attributes names to kebab-case, bugging non kebab-case slots names.
						// 2. The main portal is anonymous to prevent cache.
						Portal().Raw(true).Content(string(body)),
					).Id("app").Attr("v-cloak", true),
					ctx.Injector.GetTailHTMLComponent(),
				).Class("front"),
			).Attr(ctx.Injector.HTMLLangAttrs()...),
		}
		return
	}
}
