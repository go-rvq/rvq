package web

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web/tag"
)

type SlotBuilder struct {
	tag.TagBuilder[*SlotBuilder]
	scope string
	name  string
}

func Slot(children ...h.HTMLComponent) (r *SlotBuilder) {
	return tag.NewTag(&SlotBuilder{}, "template", children...)
}

func (b *SlotBuilder) Scope(v string) (r *SlotBuilder) {
	b.scope = v
	return b
}

func (b *SlotBuilder) Name(v string) (r *SlotBuilder) {
	b.name = v
	return b
}

func (b *SlotBuilder) Write(ctx *h.Context) (err error) {
	if len(b.name) == 0 {
		panic("Slot(...).Name(name) required")
	}

	attrName := fmt.Sprintf("v-slot:%s", b.name)
	if len(b.scope) == 0 {
		b.Attr(attrName, true)
	} else {
		b.Attr(attrName, b.scope)
	}
	return b.TagBuilder.Write(ctx)
}
