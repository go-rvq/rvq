package vuetifyx

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/x/ui/vuetify"
)

type VXReadonlyFieldBuilder struct {
	label    string
	value    interface{}
	children h.HTMLComponents
	checkbox bool
	icon     string
	color    string
}

func VXReadonlyField(children ...h.HTMLComponent) *VXReadonlyFieldBuilder {
	b := &VXReadonlyFieldBuilder{}
	if len(children) > 0 {
		b.children = children
	}
	return b
}

func (b *VXReadonlyFieldBuilder) Label(v string) *VXReadonlyFieldBuilder {
	b.label = v
	return b
}

func (b *VXReadonlyFieldBuilder) Value(v interface{}) *VXReadonlyFieldBuilder {
	b.value = v
	return b
}

func (b *VXReadonlyFieldBuilder) Children(children ...h.HTMLComponent) *VXReadonlyFieldBuilder {
	b.children = children
	return b
}

func (b *VXReadonlyFieldBuilder) Checkbox(v bool) *VXReadonlyFieldBuilder {
	b.checkbox = v
	return b
}

func (b *VXReadonlyFieldBuilder) Icon(v string) *VXReadonlyFieldBuilder {
	b.icon = v
	return b
}

func (b *VXReadonlyFieldBuilder) Color(v string) *VXReadonlyFieldBuilder {
	b.color = v
	return b
}

func (b *VXReadonlyFieldBuilder) Write(ctx *h.Context) (err error) {
	var vComp h.HTMLComponent
	if b.children != nil {
		vComp = b.children
	} else {
		if b.checkbox {
			ck := vuetify.VCheckbox().Value(b.value).
				Readonly(true).
				Ripple(false).
				HideDetails(true).
				Class("my-0 py-0")
			vComp = ck

			if b.color != "" {
				ck.Color(b.color)
			}
		} else if b.icon != "" {
			i := vuetify.VIcon(b.icon)
			vComp = i
			if b.color != "" {
				i.Color(b.color)
			}
		} else if b.value != nil {
			if b.color != "" {
				vComp = h.Span(fmt.Sprint(b.value)).Class("text-" + b.color)
			} else {
				vComp = h.Text(fmt.Sprint(b.value))
			}
		}
	}

	return h.Div(
		h.Label(b.label).Class("v-label theme--light text-caption"),
		h.Div(vComp).Class("pt-1"),
	).Class("mb-4").Write(ctx)
}
