package web

import (
	"bufio"
	"net"
	"net/http"
)

type ResposeWriten interface {
	Writed() bool
}

type ResponseWriter interface {
	http.ResponseWriter
	ResposeWriten
	StatusCode() int
	SetStatusCode(code int)
	Unwrap() http.ResponseWriter
}

type ResponseWriteHijacker interface {
	ResponseWriter
	http.Hijacker
}

func WrapResponseWriter(w http.ResponseWriter) ResponseWriter {
	if w, ok := w.(ResponseWriter); ok {
		return w
	}

	rw := &wrapedResponseWriter{ResponseWriter: w}

	if _, ok := w.(http.Hijacker); ok {
		return &hijackedResponseWriterWrapper{ResponseWriter: rw}
	}

	return rw
}

type wrapedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	writed     bool
}

func (w *wrapedResponseWriter) SetStatusCode(code int) {
	w.statusCode = code
}

func (w *wrapedResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w *wrapedResponseWriter) Writed() bool {
	return w.writed
}

func (w *wrapedResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.writed = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *wrapedResponseWriter) Write(b []byte) (int, error) {
	if !w.writed {
		if w.statusCode > 0 {
			w.WriteHeader(w.statusCode)
		} else {
			w.writed = true
		}
	}
	return w.ResponseWriter.Write(b)
}

func (w *wrapedResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

type hijackedResponseWriterWrapper struct {
	ResponseWriter
}

func (w *wrapedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
