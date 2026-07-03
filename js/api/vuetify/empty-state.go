package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VEmptyStateBuilder struct {
	VTagBuilder[*VEmptyStateBuilder]
}

func VEmptyState(children ...h.HTMLComponent) *VEmptyStateBuilder {
	return VTag(&VEmptyStateBuilder{}, "v-empty-state", children...)
}

func (b *VEmptyStateBuilder) Title(v string) (r *VEmptyStateBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VEmptyStateBuilder) Theme(v string) (r *VEmptyStateBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VEmptyStateBuilder) Size(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) Text(v string) (r *VEmptyStateBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VEmptyStateBuilder) BgColor(v string) (r *VEmptyStateBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VEmptyStateBuilder) Height(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) MaxHeight(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) MaxWidth(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) MinHeight(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) MinWidth(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) Width(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) Color(v string) (r *VEmptyStateBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VEmptyStateBuilder) Image(v string) (r *VEmptyStateBuilder) {
	b.Attr("image", v)
	return b
}

func (b *VEmptyStateBuilder) Href(v string) (r *VEmptyStateBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VEmptyStateBuilder) Icon(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) Headline(v string) (r *VEmptyStateBuilder) {
	b.Attr("headline", v)
	return b
}

func (b *VEmptyStateBuilder) ActionText(v string) (r *VEmptyStateBuilder) {
	b.Attr("action-text", v)
	return b
}

func (b *VEmptyStateBuilder) Justify(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":justify", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) TextWidth(v interface{}) (r *VEmptyStateBuilder) {
	b.Attr(":text-width", h.JSONString(v))
	return b
}

func (b *VEmptyStateBuilder) To(v string) (r *VEmptyStateBuilder) {
	b.Attr("to", v)
	return b
}

func (b *VEmptyStateBuilder) On(name string, value string) (r *VEmptyStateBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VEmptyStateBuilder) Bind(name string, value string) (r *VEmptyStateBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VEmptyStateBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VEmptyStateBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VEmptyStateBuilder) Slot(name string, child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VEmptyStateBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VEmptyStateBuilder) SetSlotActions(child ...h.HTMLComponent) {
	b.SetSlot("actions", child...)
}

func (b *VEmptyStateBuilder) SetScopedSlotActions(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("actions", scope, child...)
}

func (b *VEmptyStateBuilder) SlotActions(child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetSlotActions(child...)
	return b
}

func (b *VEmptyStateBuilder) ScopedSlotActions(scope string, child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetScopedSlotActions(scope, child...)
	return b
}

func (b *VEmptyStateBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VEmptyStateBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VEmptyStateBuilder) SlotDefault(child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VEmptyStateBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VEmptyStateBuilder) SetSlotHeadline(child ...h.HTMLComponent) {
	b.SetSlot("headline", child...)
}

func (b *VEmptyStateBuilder) SetScopedSlotHeadline(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("headline", scope, child...)
}

func (b *VEmptyStateBuilder) SlotHeadline(child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetSlotHeadline(child...)
	return b
}

func (b *VEmptyStateBuilder) ScopedSlotHeadline(scope string, child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetScopedSlotHeadline(scope, child...)
	return b
}

func (b *VEmptyStateBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VEmptyStateBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VEmptyStateBuilder) SlotTitle(child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VEmptyStateBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VEmptyStateBuilder) SetSlotMedia(child ...h.HTMLComponent) {
	b.SetSlot("media", child...)
}

func (b *VEmptyStateBuilder) SetScopedSlotMedia(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("media", scope, child...)
}

func (b *VEmptyStateBuilder) SlotMedia(child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetSlotMedia(child...)
	return b
}

func (b *VEmptyStateBuilder) ScopedSlotMedia(scope string, child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetScopedSlotMedia(scope, child...)
	return b
}

func (b *VEmptyStateBuilder) SetSlotText(child ...h.HTMLComponent) {
	b.SetSlot("text", child...)
}

func (b *VEmptyStateBuilder) SetScopedSlotText(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("text", scope, child...)
}

func (b *VEmptyStateBuilder) SlotText(child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetSlotText(child...)
	return b
}

func (b *VEmptyStateBuilder) ScopedSlotText(scope string, child ...h.HTMLComponent) (r *VEmptyStateBuilder) {
	b.SetScopedSlotText(scope, child...)
	return b
}
