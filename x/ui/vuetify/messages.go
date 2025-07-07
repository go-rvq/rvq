package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VMessagesBuilder struct {
	VTagBuilder[*VMessagesBuilder]
}

func VMessages(children ...h.HTMLComponent) *VMessagesBuilder {
	return VTag(&VMessagesBuilder{}, "v-messages", children...)
}

func (b *VMessagesBuilder) Active(v bool) (r *VMessagesBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VMessagesBuilder) Color(v string) (r *VMessagesBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VMessagesBuilder) Messages(v interface{}) (r *VMessagesBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VMessagesBuilder) Transition(v interface{}) (r *VMessagesBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VMessagesBuilder) On(name string, value string) (r *VMessagesBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VMessagesBuilder) Bind(name string, value string) (r *VMessagesBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
