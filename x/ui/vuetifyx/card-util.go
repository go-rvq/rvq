package vuetifyx

import (
	v "github.com/go-rvq/rvq/x/ui/vuetify"

	h "github.com/go-rvq/htmlgo"
)

type CardBuilder struct {
	children   []h.HTMLComponent
	systemBar  []h.HTMLComponent
	header     []h.HTMLComponent
	actions    []h.HTMLComponent
	classNames []string
	variant    string
}

func Card(children ...h.HTMLComponent) (r *CardBuilder) {
	r = &CardBuilder{}
	r.Children(children...)
	return
}

func (b *CardBuilder) Children(comps ...h.HTMLComponent) (r *CardBuilder) {
	b.children = comps
	return b
}

func (b *CardBuilder) Actions(actions ...h.HTMLComponent) (r *CardBuilder) {
	b.actions = actions
	return b
}

func (b *CardBuilder) Header(header ...h.HTMLComponent) (r *CardBuilder) {
	b.header = header
	return b
}

func (b *CardBuilder) HeaderTitle(title string) (r *CardBuilder) {
	b.header = []h.HTMLComponent{h.Text(title)}
	return b
}

func (b *CardBuilder) SystemBar(systemBar ...h.HTMLComponent) (r *CardBuilder) {
	b.systemBar = systemBar
	return b
}

func (b *CardBuilder) Class(names ...string) (r *CardBuilder) {
	b.classNames = names
	return b
}

func (b *CardBuilder) Variant(v string) (r *CardBuilder) {
	b.variant = v
	return b
}

func (b *CardBuilder) Write(ctx *h.Context) (err error) {
	var sb h.HTMLComponent
	var hr h.HTMLComponent
	if len(b.systemBar) > 0 {
		sb = v.VSystemBar(b.systemBar...).Class("mx-2 pt-4").Color("white").Height(32)
	}
	if len(b.children) > 0 {
		empty := true
		for _, c := range b.children {
			if c != nil {
				empty = false
			}
		}
		if !empty {
			hr = v.VDivider()
		}
	}

	return v.VCard(
		sb,
		v.VToolbar(
			v.VToolbarTitle("", b.header...),
			v.VSpacer(),
		).Flat(true).AppendChild(b.actions...).Color("white"),
		hr,
	).
		Variant(b.variant).
		Class(b.classNames...).
		AppendChild(b.children...).
		Write(ctx)
}
