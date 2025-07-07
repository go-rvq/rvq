package presets

import (
	"net/http"
	"path"
	"slices"
	"strings"
)

type PageHandler struct {
	path    string
	methods []string
	handler http.Handler
}

func (h *PageHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.handler.ServeHTTP(writer, request)
}

func (h *PageHandler) Path() string {
	return h.path
}

func (h *PageHandler) Methods() []string {
	return h.methods
}

func (h *PageHandler) Handler() http.Handler {
	return h.handler
}

func (h *PageHandler) WrapHandler(f func(old http.Handler) http.Handler) {
	h.handler = f(h.handler)
}

func NewPageHandler(path string, handler http.Handler, methods ...string) *PageHandler {
	for i, method := range methods {
		methods[i] = strings.ToUpper(method)
	}
	return &PageHandler{path: path, handler: handler, methods: methods}
}

func (h *PageHandler) SetupRoute(prefix string, mux *http.ServeMux, cb ...func(pattern string, ph *PageHandler)) {
	pth := path.Join(prefix, h.path)
	if len(h.methods) == 0 {
		pattern := pth

		mux.Handle(pattern, h)

		for _, f := range cb {
			f(pattern, h)
		}
	} else {
		for _, method := range h.methods {
			pattern := method + " " + pth

			mux.Handle(pattern, h)

			for _, f := range cb {
				f(pattern, h)
			}
		}
	}
}

func (h *PageHandler) ContainsMethods(methods ...string) (contained []string) {
	for _, m := range methods {
		if slices.Contains(h.methods, m) {
			contained = append(contained, m)
		}
	}
	return
}

type PageHandlers []*PageHandler

func (p PageHandlers) SetupRoutes(prefix string, mux *http.ServeMux, cb ...func(pattern string, ph *PageHandler)) {
	for _, h := range p {
		h.SetupRoute(prefix, mux, cb...)
	}
}

func (p *PageHandlers) Add(ph *PageHandler) {
	for _, h := range *p {
		if h.path == ph.path {
			duplicateMethods := ph.ContainsMethods(h.methods...)
			if (len(h.methods) == 0 && len(ph.methods) == 0) || len(duplicateMethods) > 0 {
				var s string
				if len(duplicateMethods) > 0 {
					s = strings.Join(duplicateMethods, ",") + " "
				}
				panic("duplicated path: " + s + h.path)
			}
		}
	}
	*p = append(*p, ph)
}
