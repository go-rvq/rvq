package vuetifyx

import (
	"fmt"

	v "github.com/go-rvq/rvq/x/ui/vuetify"

	h "github.com/go-rvq/htmlgo"
)

type DetailInfoBuilder struct {
	columns []h.HTMLComponent
	classes []string
}

func DetailInfo(columns ...h.HTMLComponent) (r *DetailInfoBuilder) {
	r = &DetailInfoBuilder{}
	r.columns = columns
	return
}

func (b *DetailInfoBuilder) Write(ctx *h.Context) (err error) {
	row := v.VRow()
	for _, col := range b.columns {
		row.AppendChild(v.VCol(col).Md(true))
	}

	return v.VContainer(row).Class(b.classes...).Write(ctx)
}

func (b *DetailInfoBuilder) Class(v ...string) (r *DetailInfoBuilder) {
	b.classes = v
	return b
}

type DetailFieldBuilder struct {
	label         string
	labelMinWidth string
	icon          h.HTMLComponent
	children      []h.HTMLComponent
}

func DetailField(children ...h.HTMLComponent) (r *DetailFieldBuilder) {
	r = &DetailFieldBuilder{
		labelMinWidth: "180px",
	}
	r.Children(children...)
	return
}

func (b *DetailFieldBuilder) Children(comps ...h.HTMLComponent) (r *DetailFieldBuilder) {
	b.children = comps
	return b
}

func (b *DetailFieldBuilder) Label(v string) (r *DetailFieldBuilder) {
	b.label = v
	return b
}

func (b *DetailFieldBuilder) LabelMinWidth(v string) (r *DetailFieldBuilder) {
	b.labelMinWidth = v
	return b
}

func (b *DetailFieldBuilder) Icon(v h.HTMLComponent) (r *DetailFieldBuilder) {
	b.icon = v
	return b
}

func (b *DetailFieldBuilder) Write(ctx *h.Context) (err error) {
	ki := h.Tag("div").
		Children(
			h.Tag("label").
				Text(b.label).
				Class("blue-grey--text lighten-3").Style(fmt.Sprintf("min-width: %s", b.labelMinWidth)),
		).Class("d-flex pb-2")

	if b.icon != nil {
		ki.AppendChildren(b.icon)
	}

	ki.AppendChildren(b.children...)
	return ki.Write(ctx)
}

type DetailColumnBuilder struct {
	key      string
	children []h.HTMLComponent
	header   string
}

func DetailColumn(children ...h.HTMLComponent) (r *DetailColumnBuilder) {
	r = &DetailColumnBuilder{}
	r.Children(children...)
	return
}

func (b *DetailColumnBuilder) Children(comps ...h.HTMLComponent) (r *DetailColumnBuilder) {
	b.children = comps
	return b
}

func (b *DetailColumnBuilder) Header(text string) (r *DetailColumnBuilder) {
	b.header = text
	return b
}

func (b *DetailColumnBuilder) Append(label string, comp h.HTMLComponent) (r *DetailColumnBuilder) {
	b.AppendIcon(label, nil, comp)
	return b
}

func (b *DetailColumnBuilder) AppendIcon(label string, icon h.HTMLComponent, comp h.HTMLComponent) (r *DetailColumnBuilder) {
	b.children = append(b.children, DetailField(comp).Label(label).Icon(icon))
	return b
}

func (b *DetailColumnBuilder) Write(ctx *h.Context) (err error) {
	detailInfoBody := h.Tag("div")
	if len(b.header) > 0 {
		detailInfoBody.AppendChildren(
			h.Tag("h5").
				Text(b.header).
				Class("subtitle-2 pb-2"),
		)
	}
	detailInfoBody.AppendChildren(b.children...)

	return detailInfoBody.Write(ctx)
}

type OptionalTextBuilder struct {
	text      string
	zeroLabel string
}

func OptionalText(text string) (r *OptionalTextBuilder) {
	r = &OptionalTextBuilder{text: text}
	return
}

func (b *OptionalTextBuilder) ZeroLabel(label string) (r *OptionalTextBuilder) {
	b.zeroLabel = label
	return b
}

func (b *OptionalTextBuilder) Write(ctx *h.Context) (err error) {
	var body h.HTMLComponent

	if len(b.text) > 0 {
		body = h.Tag("span").
			Text(b.text)
	} else {
		body = h.Tag("span").
			Class("grey--text lighten-5").
			Text(b.zeroLabel)
	}

	return body.Write(ctx)
}
