package vuetifyx

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

func HTMLList(tag *h.HTMLTagBuilder, itens ...any) h.HTMLComponent {
	tag.Style("list-style: inside")
	for _, iten := range itens {
		var c h.HTMLComponent
		switch t := iten.(type) {
		case h.HTMLComponent:
			c = t
		default:
			c = h.RawHTML(fmt.Sprint(t))
		}
		tg, _ := c.(h.TagGetter)
		if tg == nil || tg.GetHTMLTagBuilder().TagName != "li" {
			c = h.Li(c)
		}
		tag.Childs = append(tag.Childs, c)
	}
	return tag
}
