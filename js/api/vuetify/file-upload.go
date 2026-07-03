package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VFileUploadBuilder struct {
	VTagBuilder[*VFileUploadBuilder]
}

func VFileUpload(children ...h.HTMLComponent) *VFileUploadBuilder {
	return VTag(&VFileUploadBuilder{}, "v-file-upload", children...)
}

func (b *VFileUploadBuilder) Length(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":length", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) Tag(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) Name(v string) (r *VFileUploadBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VFileUploadBuilder) Title(v string) (r *VFileUploadBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VFileUploadBuilder) Theme(v string) (r *VFileUploadBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VFileUploadBuilder) Subtitle(v string) (r *VFileUploadBuilder) {
	b.Attr("subtitle", v)
	return b
}

func (b *VFileUploadBuilder) Disabled(v bool) (r *VFileUploadBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VFileUploadBuilder) Multiple(v bool) (r *VFileUploadBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VFileUploadBuilder) Border(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) Density(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) Height(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) MaxHeight(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) MaxWidth(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) MinHeight(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) MinWidth(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) Width(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) Elevation(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) Rounded(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) Tile(v bool) (r *VFileUploadBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VFileUploadBuilder) Color(v string) (r *VFileUploadBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VFileUploadBuilder) Opacity(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) ModelValue(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) Scrim(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":scrim", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) CloseDelay(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":close-delay", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) OpenDelay(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":open-delay", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) Location(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) Clearable(v bool) (r *VFileUploadBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VFileUploadBuilder) Position(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) Icon(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) Thickness(v interface{}) (r *VFileUploadBuilder) {
	b.Attr(":thickness", h.JSONString(v))
	return b
}

func (b *VFileUploadBuilder) ShowSize(v bool) (r *VFileUploadBuilder) {
	b.Attr(":show-size", fmt.Sprint(v))
	return b
}

func (b *VFileUploadBuilder) BrowseText(v string) (r *VFileUploadBuilder) {
	b.Attr("browse-text", v)
	return b
}

func (b *VFileUploadBuilder) DividerText(v string) (r *VFileUploadBuilder) {
	b.Attr("divider-text", v)
	return b
}

func (b *VFileUploadBuilder) HideBrowse(v bool) (r *VFileUploadBuilder) {
	b.Attr(":hide-browse", fmt.Sprint(v))
	return b
}

func (b *VFileUploadBuilder) On(name string, value string) (r *VFileUploadBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFileUploadBuilder) Bind(name string, value string) (r *VFileUploadBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VFileUploadBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VFileUploadBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VFileUploadBuilder) Slot(name string, child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VFileUploadBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VFileUploadBuilder) SetSlotBrowse(child ...h.HTMLComponent) {
	b.SetSlot("browse", child...)
}

func (b *VFileUploadBuilder) SetScopedSlotBrowse(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("browse", scope, child...)
}

func (b *VFileUploadBuilder) SlotBrowse(child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetSlotBrowse(child...)
	return b
}

func (b *VFileUploadBuilder) ScopedSlotBrowse(scope string, child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetScopedSlotBrowse(scope, child...)
	return b
}

func (b *VFileUploadBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VFileUploadBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VFileUploadBuilder) SlotDefault(child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VFileUploadBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VFileUploadBuilder) SetSlotIcon(child ...h.HTMLComponent) {
	b.SetSlot("icon", child...)
}

func (b *VFileUploadBuilder) SetScopedSlotIcon(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("icon", scope, child...)
}

func (b *VFileUploadBuilder) SlotIcon(child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetSlotIcon(child...)
	return b
}

func (b *VFileUploadBuilder) ScopedSlotIcon(scope string, child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetScopedSlotIcon(scope, child...)
	return b
}

func (b *VFileUploadBuilder) SetSlotInput(child ...h.HTMLComponent) {
	b.SetSlot("input", child...)
}

func (b *VFileUploadBuilder) SetScopedSlotInput(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("input", scope, child...)
}

func (b *VFileUploadBuilder) SlotInput(child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetSlotInput(child...)
	return b
}

func (b *VFileUploadBuilder) ScopedSlotInput(scope string, child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetScopedSlotInput(scope, child...)
	return b
}

func (b *VFileUploadBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VFileUploadBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VFileUploadBuilder) SlotItem(child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VFileUploadBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VFileUploadBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VFileUploadBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VFileUploadBuilder) SlotTitle(child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VFileUploadBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VFileUploadBuilder) SetSlotDivider(child ...h.HTMLComponent) {
	b.SetSlot("divider", child...)
}

func (b *VFileUploadBuilder) SetScopedSlotDivider(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("divider", scope, child...)
}

func (b *VFileUploadBuilder) SlotDivider(child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetSlotDivider(child...)
	return b
}

func (b *VFileUploadBuilder) ScopedSlotDivider(scope string, child ...h.HTMLComponent) (r *VFileUploadBuilder) {
	b.SetScopedSlotDivider(scope, child...)
	return b
}
