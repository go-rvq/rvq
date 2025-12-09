package web

import (
	"context"
	"net/http"
)

const RequestMethodParam = ":method"

func ParseRequest(r *http.Request) *http.Request {
	q := UrlQueryFromRequest(r)
	if q == nil {
		q = &UrlQuery{Values: r.URL.Query()}
	}
	r = r.WithContext(context.WithValue(r.Context(), UrlQueryKey, q))

	if method := r.FormValue(RequestMethodParam); method != "" {
		r.Method = method
		delete(r.Form, RequestMethodParam)
	}
	return r
}
