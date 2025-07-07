package web

import (
	"github.com/go-rvq/rvq/web/tag"
	"github.com/go-rvq/rvq/web/vue"
	h "github.com/theplant/htmlgo"
)

func Unscoped(comp h.HTMLComponent) (r h.HTMLComponent) {
	tag.Walk(comp, func(c h.HTMLComponent) (state tag.WalkState) {
		switch c.(type) {
		case *ScopeBuilder:
		case *vue.UserComponentBuilder:
		default:
			r = c
			state = tag.SkipAll
		}
		return
	})
	return
}
