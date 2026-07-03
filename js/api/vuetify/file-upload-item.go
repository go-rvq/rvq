package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VFileUploadItemBuilder struct {
	VTagBuilder[*VFileUploadItemBuilder]
}

func VFileUploadItem(children ...h.HTMLComponent) *VFileUploadItemBuilder {
	return VTag(&VFileUploadItemBuilder{}, "v-file-upload-item", children...)
}

func (b *VFileUploadItemBuilder) Replace(v bool) (r *VFileUploadItemBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VFileUploadItemBuilder) Link(v bool) (r *VFileUploadItemBuilder) {
	b.Attr(":link", fmt.Sprint(v))
	return b
}

func (b *VFileUploadItemBuilder) Tag(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) Nav(v bool) (r *VFileUploadItemBuilder) {
	b.Attr(":nav", fmt.Sprint(v))
	return b
}

func (b *VFileUploadItemBuilder) Title(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":title", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) Theme(v string) (r *VFileUploadItemBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VFileUploadItemBuilder) Value(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) Exact(v bool) (r *VFileUploadItemBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VFileUploadItemBuilder) Subtitle(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":subtitle", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) BaseColor(v string) (r *VFileUploadItemBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VFileUploadItemBuilder) ActiveColor(v string) (r *VFileUploadItemBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VFileUploadItemBuilder) ActiveClass(v string) (r *VFileUploadItemBuilder) {
	b.Attr("active-class", v)
	return b
}

func (b *VFileUploadItemBuilder) Disabled(v bool) (r *VFileUploadItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VFileUploadItemBuilder) Lines(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":lines", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) Slim(v bool) (r *VFileUploadItemBuilder) {
	b.Attr(":slim", fmt.Sprint(v))
	return b
}

func (b *VFileUploadItemBuilder) Border(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) Density(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) Height(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) MaxHeight(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) MaxWidth(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) MinHeight(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) MinWidth(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) Width(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) Elevation(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) Rounded(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) Tile(v bool) (r *VFileUploadItemBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VFileUploadItemBuilder) Color(v string) (r *VFileUploadItemBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VFileUploadItemBuilder) Variant(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) AppendIcon(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) PrependIcon(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":prepend-icon", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) Clearable(v bool) (r *VFileUploadItemBuilder) {
	b.Attr(":clearable", fmt.Sprint(v))
	return b
}

func (b *VFileUploadItemBuilder) Active(v bool) (r *VFileUploadItemBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VFileUploadItemBuilder) Href(v string) (r *VFileUploadItemBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VFileUploadItemBuilder) To(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) Ripple(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":ripple", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) ShowSize(v bool) (r *VFileUploadItemBuilder) {
	b.Attr(":show-size", fmt.Sprint(v))
	return b
}

func (b *VFileUploadItemBuilder) File(v interface{}) (r *VFileUploadItemBuilder) {
	b.Attr(":file", h.JSONString(v))
	return b
}

func (b *VFileUploadItemBuilder) FileIcon(v string) (r *VFileUploadItemBuilder) {
	b.Attr("file-icon", v)
	return b
}

func (b *VFileUploadItemBuilder) AppendAvatar(v string) (r *VFileUploadItemBuilder) {
	b.Attr("append-avatar", v)
	return b
}

func (b *VFileUploadItemBuilder) PrependAvatar(v string) (r *VFileUploadItemBuilder) {
	b.Attr("prepend-avatar", v)
	return b
}

func (b *VFileUploadItemBuilder) On(name string, value string) (r *VFileUploadItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFileUploadItemBuilder) Bind(name string, value string) (r *VFileUploadItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VFileUploadItemBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VFileUploadItemBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VFileUploadItemBuilder) Slot(name string, child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VFileUploadItemBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VFileUploadItemBuilder) SetSlotClear(child ...h.HTMLComponent) {
	b.SetSlot("clear", child...)
}

func (b *VFileUploadItemBuilder) SetScopedSlotClear(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("clear", scope, child...)
}

func (b *VFileUploadItemBuilder) SlotClear(child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetSlotClear(child...)
	return b
}

func (b *VFileUploadItemBuilder) ScopedSlotClear(scope string, child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetScopedSlotClear(scope, child...)
	return b
}

func (b *VFileUploadItemBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VFileUploadItemBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VFileUploadItemBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VFileUploadItemBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VFileUploadItemBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VFileUploadItemBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VFileUploadItemBuilder) SlotAppend(child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VFileUploadItemBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}

func (b *VFileUploadItemBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VFileUploadItemBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VFileUploadItemBuilder) SlotDefault(child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VFileUploadItemBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VFileUploadItemBuilder) SetSlotTitle(child ...h.HTMLComponent) {
	b.SetSlot("title", child...)
}

func (b *VFileUploadItemBuilder) SetScopedSlotTitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("title", scope, child...)
}

func (b *VFileUploadItemBuilder) SlotTitle(child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetSlotTitle(child...)
	return b
}

func (b *VFileUploadItemBuilder) ScopedSlotTitle(scope string, child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetScopedSlotTitle(scope, child...)
	return b
}

func (b *VFileUploadItemBuilder) SetSlotSubtitle(child ...h.HTMLComponent) {
	b.SetSlot("subtitle", child...)
}

func (b *VFileUploadItemBuilder) SetScopedSlotSubtitle(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("subtitle", scope, child...)
}

func (b *VFileUploadItemBuilder) SlotSubtitle(child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetSlotSubtitle(child...)
	return b
}

func (b *VFileUploadItemBuilder) ScopedSlotSubtitle(scope string, child ...h.HTMLComponent) (r *VFileUploadItemBuilder) {
	b.SetScopedSlotSubtitle(scope, child...)
	return b
}
