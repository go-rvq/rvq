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
