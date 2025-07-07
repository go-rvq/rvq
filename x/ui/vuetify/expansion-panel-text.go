package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VExpansionPanelTextBuilder struct {
	VTagBuilder[*VExpansionPanelTextBuilder]
}

func VExpansionPanelText(children ...h.HTMLComponent) *VExpansionPanelTextBuilder {
	return VTag(&VExpansionPanelTextBuilder{}, "v-expansion-panel-text", children...)
}

func (b *VExpansionPanelTextBuilder) Eager(v bool) (r *VExpansionPanelTextBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VExpansionPanelTextBuilder) On(name string, value string) (r *VExpansionPanelTextBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VExpansionPanelTextBuilder) Bind(name string, value string) (r *VExpansionPanelTextBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
