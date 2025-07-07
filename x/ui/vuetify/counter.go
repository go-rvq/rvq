package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCounterBuilder struct {
	VTagBuilder[*VCounterBuilder]
}

func VCounter(children ...h.HTMLComponent) *VCounterBuilder {
	return VTag(&VCounterBuilder{}, "v-counter", children...)
}

func (b *VCounterBuilder) Active(v bool) (r *VCounterBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VCounterBuilder) Disabled(v bool) (r *VCounterBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCounterBuilder) Max(v interface{}) (r *VCounterBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VCounterBuilder) Value(v interface{}) (r *VCounterBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VCounterBuilder) Transition(v interface{}) (r *VCounterBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VCounterBuilder) On(name string, value string) (r *VCounterBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCounterBuilder) Bind(name string, value string) (r *VCounterBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
