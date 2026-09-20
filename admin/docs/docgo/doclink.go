package docgo

import (
	. "github.com/go-rvq/htmlgo"
)

type DocLinkBuilder struct {
	doc *DocBuilder
	tag *HTMLTagBuilder
}

func DocLink(doc *DocBuilder) (r *DocLinkBuilder) {
	return &DocLinkBuilder{
		doc: doc,
		tag: A(),
	}
}

func (b *DocLinkBuilder) Write(ctx *Context) error {
	if b.doc == nil {
		return nil
	}
	b.tag.Text(b.doc.title).Href(b.doc.GetPageURL())
	return b.tag.Write(ctx)
}
