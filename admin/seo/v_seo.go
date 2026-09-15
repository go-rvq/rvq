package seo

import (
	h "github.com/go-rvq/htmlgo"
)

type VSeoBuilder struct {
	tag *h.HTMLTagBuilder
}

func VSeo(children ...h.HTMLComponent) (r *VSeoBuilder) {
	r = &VSeoBuilder{
		tag: h.Tag("vx-send-variables").Children(children...),
	}
	return
}

func (b *VSeoBuilder) Value(v string) (r *VSeoBuilder) {
	b.tag.Attr(":value", h.JSONString(v))
	return b
}

// Template sets the JS function (as a bound `:template` prop) that wraps a
// clicked tag into the text inserted at the cursor. jsFunc is a JS expression
// evaluating to a `(tag) => string` function, e.g. `(tag) => "{" + tag + "}"`.
// Without it the component keeps its default `{{tag}}` form.
func (b *VSeoBuilder) Template(jsFunc string) (r *VSeoBuilder) {
	b.tag.Attr(":template", jsFunc)
	return b
}

func (b *VSeoBuilder) Placeholder(v string) (r *VSeoBuilder) {
	b.tag.Attr(":placeholder", h.JSONString(v))
	return b
}

func (b *VSeoBuilder) SetAttr(k string, v interface{}) {
	b.tag.SetAttr(k, v)
}

func (b *VSeoBuilder) Attr(vs ...interface{}) (r *VSeoBuilder) {
	b.tag.Attr(vs...)
	return b
}

func (b *VSeoBuilder) Children(children ...h.HTMLComponent) (r *VSeoBuilder) {
	b.tag.Children(children...)
	return b
}

func (b *VSeoBuilder) AppendChildren(children ...h.HTMLComponent) (r *VSeoBuilder) {
	b.tag.AppendChildren(children...)
	return b
}

func (b *VSeoBuilder) PrependChildren(children ...h.HTMLComponent) (r *VSeoBuilder) {
	b.tag.PrependChildren(children...)
	return b
}

func (b *VSeoBuilder) Class(names ...string) (r *VSeoBuilder) {
	b.tag.Class(names...)
	return b
}

func (b *VSeoBuilder) ClassIf(name string, add bool) (r *VSeoBuilder) {
	b.tag.ClassIf(name, add)
	return b
}

func (b *VSeoBuilder) Write(ctx *h.Context) (err error) {
	return b.tag.Write(ctx)
}
