package printer

import (
	"net/http"
	"strings"

	h "github.com/go-rvq/htmlgo"
)

const PrintPath = "!print"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasSuffix(r.RequestURI, PrintPath) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if p, err := ParseRequest(r); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				h.Write(w, p.Component())
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
