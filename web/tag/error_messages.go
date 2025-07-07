package tag

import (
	"strings"

	h "github.com/theplant/htmlgo"
)

func stringsTrim(vs ...string) (r []string) {
	for _, v := range vs {
		if cv := strings.TrimSpace(v); len(cv) > 0 {
			r = append(r, cv)
		}
	}
	return
}

func (t *TagBuilder[T]) ErrorMessages(vs ...string) T {
	cvs := stringsTrim(vs...)
	if len(cvs) > 0 {
		t.SetAttr(":error-messages", h.JSONString(cvs))
	}
	return t.dot
}

func (t *TagBuilder[T]) Errors(err ...error) T {
	var s []string
	for _, e := range err {
		if e != nil {
			s = append(s, e.Error())
		}
	}
	return t.ErrorMessages(s...)
}
