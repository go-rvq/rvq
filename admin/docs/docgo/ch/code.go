package ch

import (
	h "github.com/go-rvq/htmlgo"
)

type CodeBuilder struct {
	tag *h.HTMLTagBuilder
}

func Code(code string) (r *CodeBuilder) {
	r = &CodeBuilder{
		tag: h.Tag("highlightjs").Attr(":language", h.JSONString("go")),
	}
	r.Code(code)
	return
}

func (b *CodeBuilder) Code(v string) (r *CodeBuilder) {
	b.tag.Attr(":code", h.JSONString(v))
	return b
}

func (b *CodeBuilder) Language(v string) (r *CodeBuilder) {
	b.tag.Attr(":language", h.JSONString(v))
	return b
}

func (b *CodeBuilder) Write(ctx *h.Context) error {
	return b.tag.Write(ctx)
}
