package web

import "net/http"

const RequestMethodParam = ":method"

func ParseRequest(r *http.Request) {
	if method := r.FormValue(RequestMethodParam); method != "" {
		r.Method = method
		delete(r.Form, RequestMethodParam)
	}
}
