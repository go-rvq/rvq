package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	h "github.com/theplant/htmlgo"
)

var Default = New()

func Page(pf PageFunc, efs ...interface{}) (p *PageBuilder) {
	p = &PageBuilder{
		b: Default,
	}
	p.pageRenderFunc = pf
	p.RegisterEventHandler("__reload__", EventFunc(reload))
	p.EventFuncs(efs...)
	return
}

type PageBuilder struct {
	EventsHub
	b              *Builder
	pageRenderFunc PageFunc
	maxFormSize    int64
}

func (b *Builder) Page(pf PageFunc) (p *PageBuilder) {
	p = Page(pf).Builder(b)
	return
}

func (p *PageBuilder) Builder(v *Builder) (r *PageBuilder) {
	p.b = v
	r = p
	return
}

func (p *PageBuilder) Wrap(middlewares ...func(in PageFunc) PageFunc) (r *PageBuilder) {
	pf := p.pageRenderFunc
	for _, m := range middlewares {
		pf = m(pf)
	}
	p.pageRenderFunc = pf
	r = p
	return
}

func (p *PageBuilder) MaxFormSize(v int64) (r *PageBuilder) {
	p.maxFormSize = v
	r = p
	return
}

func (p *PageBuilder) EventFuncs(vs ...interface{}) (r *PageBuilder) {
	p.addMultipleEventFuncs(vs...)
	return p
}

func (p *PageBuilder) EventFunc(name string, ef EventFunc) (r *PageBuilder) {
	p.RegisterEventHandler(name, ef)
	return p
}

func (p *PageBuilder) MergeHub(hub *EventsHub) (r *PageBuilder) {
	p.EventsHub.eventFuncs = append(hub.eventFuncs, p.EventsHub.eventFuncs...)
	return p
}

func (p *PageBuilder) render(
	w ResponseWriter,
	r *http.Request,
	c context.Context,
	head *PageInjector,
	event bool,
) (pager *PageResponse, body string) {
	if p.pageRenderFunc == nil {
		return
	}
	rf := p.pageRenderFunc
	if !event {
		rf = p.b.layoutFunc(p.pageRenderFunc)
	}

	ctx := MustGetEventContext(c)

	ctx.R = r
	ctx.W = w
	ctx.Injector = head

	pr, err := rf(ctx)
	if err != nil {
		panic(err)
	}

	if w.Writed() {
		return
	}

	pager = &pr

	if pager.Body == nil {
		return
	}

	// fmt.Println("eventFuncs count: ", len(p.eventFuncs))
	b, err := pager.Body.MarshalHTML(c)
	if err != nil {
		panic(err)
	}
	body = string(b)

	return
}

func (p *PageBuilder) index(w ResponseWriter, r *http.Request) {
	var (
		err     error
		inj     = &PageInjector{}
		ctx     = new(EventContext)
		c       = WrapEventContext(r.Context(), ctx)
		_, body = p.render(w, r, c, inj, false)
	)

	if body == "" || w.Writed() {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if _, err = fmt.Fprintln(w, body); err != nil {
		panic(err)
	}
}

func (p *PageBuilder) parseForm(r *http.Request) *multipart.Form {
	maxSize := p.maxFormSize
	if maxSize == 0 {
		maxSize = 128 << 20 // 128MB
	}

	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		panic(err)
	}

	return r.MultipartForm
}

const EventFuncIDName = "__execute_event__"

func (p *PageBuilder) executeEvent(w ResponseWriter, r *http.Request) {
	ctx := new(EventContext)
	ctx.R = r
	ctx.W = w
	ctx.Injector = &PageInjector{}

	c := WrapEventContext(r.Context(), ctx)

	eventFuncID := r.FormValue(EventFuncIDName)

	// for server side restart and lost all the eventFuncs,
	// but user keep clicking page without refresh page to call p.render to fill up eventFuncs
	// because default added reload
	if len(p.eventFuncs) <= 1 &&
		p.eventHandleById(eventFuncID) == nil &&
		p.b.eventHandleById(eventFuncID) == nil {
		log.Println("Re-render because event funcs gone, might server restarted")
		p.render(w, r, c, &PageInjector{}, true)
	}

	ef := p.eventHandleById(eventFuncID)
	if ef == nil {
		ef = p.b.eventHandleById(eventFuncID)
	}

	if ef == nil {
		log.Printf("event %s not found in %s\n", eventFuncID, p.EventsHub.String())
		http.NotFound(w, r)
		return
	}

	er, err := ef.Handle(ctx)
	if err != nil {
		panic(err)
	}

	if er.Reload {
		pr, body := p.render(w, r, c, &PageInjector{}, true)
		if w.Writed() {
			return
		}
		er.Body = h.RawHTML(body)
		if len(er.PageTitle) == 0 && pr != nil {
			er.PageTitle = pr.PageTitle
		}
	}

	er.Body = h.RawHTML(h.MustString(er.Body, c))

	for _, up := range er.UpdatePortals {
		up.Body = h.RawHTML(h.MustString(up.Body, c))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(er)
	if err != nil {
		panic(err)
	}
}

func reload(*EventContext) (r EventResponse, err error) {
	r.Reload = true
	return
}

func (p *PageBuilder) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	w := WrapResponseWriter(rw)
	if strings.Index(r.URL.String(), EventFuncIDName) >= 0 {
		p.executeEvent(w, r)
		return
	}
	p.index(w, r)
}
