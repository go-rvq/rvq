package presets

import (
	"net/http"

	"github.com/go-rvq/rvq/web"
)

type wrapedResponseWriter struct {
	web.ResponseWriter
}

func (rw *wrapedResponseWriter) WriteHeader(code int) {
	if code == http.StatusNotFound {
		// default 404 will use http.Error to set Content-Type to text/plain,
		// So we have to set it to html before WriteHeader
		rw.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.SetStatusCode(code)
		// prevent header sent
		return
	}

	rw.ResponseWriter.WriteHeader(code)
}

func (rw *wrapedResponseWriter) Write(b []byte) (int, error) {
	// don't write content, because we use customized page body
	if !rw.Writed() && rw.StatusCode() == http.StatusNotFound {
		return 0, nil
	}

	return rw.ResponseWriter.Write(b)
}

// forcedStatusWriter answers with a fixed status code, whatever the handler
// writing through it asks for. The not-found page renders through the normal
// page pipeline, which writes 200; the response has to stay a 404.
type forcedStatusWriter struct {
	http.ResponseWriter
	code    int
	written bool
}

func (w *forcedStatusWriter) WriteHeader(int) {
	if w.written {
		return
	}
	w.written = true
	w.ResponseWriter.WriteHeader(w.code)
}

func (w *forcedStatusWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.WriteHeader(w.code)
	}
	return w.ResponseWriter.Write(b)
}
