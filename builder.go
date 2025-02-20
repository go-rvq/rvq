package web

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	h "github.com/theplant/htmlgo"
)

type Builder struct {
	EventsHub
	layoutFunc LayoutFunc
}

func New() (b *Builder) {
	b = new(Builder)
	b.layoutFunc = defaultLayoutFunc
	return
}

func (b *Builder) LayoutFunc(mf LayoutFunc) (r *Builder) {
	if mf == nil {
		panic("layout func is nil")
	}
	b.layoutFunc = mf
	return b
}

func (p *Builder) EventFuncs(vs ...interface{}) (r *Builder) {
	p.addMultipleEventFuncs(vs...)
	return p
}

type ComponentsPack string

func ComponentsPackFromBytes(b []byte, e ...error) (p ComponentsPack) {
	for _, err := range e {
		if err != nil {
			panic(err)
		}
	}
	return ComponentsPack(b)
}

var startTime = time.Now()

func PacksHandler(contentType string, packs ...ComponentsPack) http.Handler {
	return Default.PacksHandler(contentType, packs...)
}

func (b *Builder) PacksHandler(contentType string, packs ...ComponentsPack) http.Handler {
	buf := bytes.NewBuffer(nil)
	for _, pk := range packs {
		// buf = append(buf, []byte(fmt.Sprintf("\n// pack %d\n", i+1))...)
		// buf = append(buf, []byte(fmt.Sprintf("\nconsole.log('pack %d, length %d');\n", i+1, len(pk)))...)
		buf.WriteString(string(pk))
		// fmt.Println(contentType)
		if strings.Contains(strings.ToLower(contentType), "javascript") {
			buf.WriteString(";")
		}
		buf.WriteString("\n\n")
	}

	body := bytes.NewReader(buf.Bytes())

	return gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		http.ServeContent(w, r, "", startTime, body)
	}))
}

func NoopLayoutFunc(in PageFunc) PageFunc {
	return in
}

func defaultLayoutFunc(in PageFunc) PageFunc {
	return func(ctx *EventContext) (r PageResponse, err error) {
		if r, err = in(ctx); err != nil {
			r.Body = h.Div(h.Text(err.Error()))
			err = nil
		} else if w, _ := ctx.ResponseWriter().(interface{ Writed() bool }); w != nil && w.Writed() {
			return
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
