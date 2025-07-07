package tag

import (
	h "github.com/theplant/htmlgo"
)

func (t *TagBuilder[T]) GetChildren() []h.HTMLComponent {
	return t.tag.Childs
}

func (t *TagBuilder[T]) Children(c ...h.HTMLComponent) T {
	t.tag.Children(c...)
	return t.dot
}

func (t *TagBuilder[T]) AppendChild(c ...h.HTMLComponent) T {
	t.tag.AppendChildren(c...)
	return t.dot
}

func (t *TagBuilder[T]) PrependChild(c ...h.HTMLComponent) T {
	t.tag.PrependChildren(c...)
	return t.dot
}

func (t *TagBuilder[T]) GetAttr(key string) *h.Attr {
	return t.tag.Attrs.Get(key)
}

func (t *TagBuilder[T]) RemoveAttr(key ...string) T {
	t.tag.Attrs.RemoveMany(key...)
	return t.dot
}

func (t *TagBuilder[T]) SetAttr(k string, v interface{}) T {
	t.tag.SetAttr(k, v)
	return t.dot
}
func (t *TagBuilder[T]) Attr(vs ...interface{}) T {
	t.tag.Attr(vs...)
	return t.dot
}
func (t *TagBuilder[T]) AttrIf(key, value interface{}, add bool) T {
	t.tag.AttrIf(key, value, add)
	return t.dot
}

func (t *TagBuilder[T]) SetTag(v string) T {
	t.tag.Tag(v)
	return t.dot
}

func (t *TagBuilder[T]) OmitEndTag() T {
	t.tag.OmitEndTag()
	return t.dot
}

func (t *TagBuilder[T]) Text(v string) T {
	t.tag.Text(v)
	return t.dot
}

func (t *TagBuilder[T]) Class(names ...string) T {
	t.tag.Class(names...)
	return t.dot
}

func (t *TagBuilder[T]) ClassIf(name string, add bool) T {
	t.tag.ClassIf(name, add)
	return t.dot
}

func (t *TagBuilder[T]) Data(vs ...string) T {
	t.tag.Data(vs...)
	return t.dot
}

func (t *TagBuilder[T]) Id(v string) T {
	t.tag.Id(v)
	return t.dot
}

func (t *TagBuilder[T]) Href(v string) T {
	t.tag.Href(v)
	return t.dot
}

func (t *TagBuilder[T]) Rel(v string) T {
	t.tag.Rel(v)
	return t.dot
}

func (t *TagBuilder[T]) Title(v string) T {
	t.tag.Title(v)
	return t.dot
}

func (t *TagBuilder[T]) TabIndex(v int) T {
	t.tag.TabIndex(v)
	return t.dot
}

func (t *TagBuilder[T]) Required(v bool) T {
	t.tag.Required(v)
	return t.dot
}

func (t *TagBuilder[T]) Readonly(v bool) T {
	t.tag.Readonly(v)
	return t.dot
}

func (t *TagBuilder[T]) Role(v string) T {
	t.tag.Role(v)
	return t.dot
}

func (t *TagBuilder[T]) Alt(v string) T {
	t.tag.Alt(v)
	return t.dot
}

func (t *TagBuilder[T]) Target(v string) T {
	t.tag.Target(v)
	return t.dot
}

func (t *TagBuilder[T]) Name(v string) T {
	t.tag.Name(v)
	return t.dot
}

func (t *TagBuilder[T]) Value(v string) T {
	t.tag.Value(v)
	return t.dot
}

func (t *TagBuilder[T]) For(v string) T {
	t.tag.For(v)
	return t.dot
}

func (t *TagBuilder[T]) Style(v string) T {
	t.tag.Style(v)
	return t.dot
}

func (t *TagBuilder[T]) StyleIf(v string, add bool) T {
	t.tag.StyleIf(v, add)
	return t.dot
}

func (t *TagBuilder[T]) Type(v string) T {
	t.tag.Type(v)
	return t.dot
}

func (t *TagBuilder[T]) Placeholder(v string) T {
	t.tag.Placeholder(v)
	return t.dot
}

func (t *TagBuilder[T]) Src(v string) T {
	t.tag.Src(v)
	return t.dot
}

func (t *TagBuilder[T]) Property(v string) T {
	t.tag.Property(v)
	return t.dot
}

func (t *TagBuilder[T]) Action(v string) T {
	t.tag.Action(v)
	return t.dot
}

func (t *TagBuilder[T]) Method(v string) T {
	t.tag.Method(v)
	return t.dot
}

func (t *TagBuilder[T]) Content(v string) T {
	t.tag.Content(v)
	return t.dot
}

func (t *TagBuilder[T]) Charset(v string) T {
	t.tag.Charset(v)
	return t.dot
}

func (t *TagBuilder[T]) Disabled(v bool) T {
	t.tag.Disabled(v)
	return t.dot
}

func (t *TagBuilder[T]) Checked(v bool) T {
	t.tag.Checked(v)
	return t.dot
}
