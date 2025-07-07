package tag

import (
	"context"

	h "github.com/theplant/htmlgo"
)

type (
	TagGetter interface {
		GetHTMLTagBuilder() *h.HTMLTagBuilder
	}

	TagBuilderGetter[T any] interface {
		TagGetter
		GetTagBuilder() *TagBuilder[T]
	}

	TagBuilder[T any] struct {
		dot T
		tag *h.HTMLTagBuilder
	}
)

func NewTag[T TagBuilderGetter[T]](dot T, name string, children ...h.HTMLComponent) T {
	b := dot.GetTagBuilder()
	b.dot = dot
	b.tag = h.Tag(name)
	if len(children) > 0 {
		b.Children(children...)
	}
	return dot
}

func (t *TagBuilder[T]) GetHTMLTagBuilder() *h.HTMLTagBuilder {
	return t.tag
}

func (t *TagBuilder[T]) GetTagBuilder() *TagBuilder[T] {
	return t
}

func (t *TagBuilder[T]) Dot() T {
	return t.dot
}

func (t *TagBuilder[T]) MarshalHTML(ctx context.Context) ([]byte, error) {
	return t.tag.MarshalHTML(ctx)
}

func UpdatedDot[T TagBuilderGetter[T]](dot T) {
	dot.GetTagBuilder().dot = dot
}

func SetTag[T any](t TagBuilderGetter[T], tag *h.HTMLTagBuilder) {
	t.GetTagBuilder().tag = tag
}
