package vuetifyx

import (
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type (
	VXAdvancedCardTagGetter[T any] interface {
		v.VTagBuilderGetter[T]
		GetVXAdvancedCardTagBuilder() *VXAdvancedCardTagBuilder[T]
	}

	VXAdvancedCardTagger interface {
		h.HTMLComponent
		SetDensity(s string)
		SetToolbarProps(s string)
		SetContainerProps(s string)
		SetMainMenuUp(v bool)
		SetSecondaryMenuUp(v bool)
		SetTitle(s string)
		SetMainMenuTitle(s string)
		SetSecondaryMenuTitle(s string)
		SetSlotMainMenu(child ...h.HTMLComponent)
		SetSlotMainMenuContainer(child ...h.HTMLComponent)
		SetSlotSecondaryMenu(child ...h.HTMLComponent)
		SetSlotSecondaryMenuContainer(child ...h.HTMLComponent)
		SetSlotPrependToolbar(child ...h.HTMLComponent)
		SetSlotAppendToolbar(child ...h.HTMLComponent)
		SetSlotHeader(child ...h.HTMLComponent)
		SetSlotPrependHeader(child ...h.HTMLComponent)
		SetSlotAppendHeader(child ...h.HTMLComponent)
		SetSlotTop(child ...h.HTMLComponent)
		SetSlotBottom(child ...h.HTMLComponent)
		SetScopedSlotBottom(scope string, child ...h.HTMLComponent)
		SetSlotBody(child ...h.HTMLComponent)
		SetScopedSlotBody(scope string, child ...h.HTMLComponent)
		SetSlotPortals(child ...h.HTMLComponent)
	}

	VXAnyAdvancedCardTagGetter interface {
		GetVXAnyAdvancedCardTagBuilder() *VXAnyAdvancedCardTagBuilder
	}
)

func (b *VXAdvancedCardTagBuilder[T]) SetDensity(s string) {
	b.Attr("density", s)
}

func (b *VXAdvancedCardTagBuilder[T]) SetToolbarProps(s string) {
	b.Attr(":toolbar-props", s)
}

func (b *VXAdvancedCardTagBuilder[T]) SetContainerProps(s string) {
	b.Attr(":container-props", s)
}

func (b *VXAdvancedCardTagBuilder[T]) SetMainMenuUp(v bool) {
	b.Attr("main-menu-up", v)
}

func (b *VXAdvancedCardTagBuilder[T]) SetSecondaryMenuUp(v bool) {
	b.Attr("secondary-menu-up", v)
}

func (b *VXAdvancedCardTagBuilder[T]) SetTitle(s string) {
	b.Attr("title", s)
}

func (b *VXAdvancedCardTagBuilder[T]) SetMainMenuTitle(s string) {
	b.Attr("main-menu-title", s)
}

func (b *VXAdvancedCardTagBuilder[T]) SetSecondaryMenuTitle(s string) {
	b.Attr("secondary-menu-title", s)
}

func (b *VXAdvancedCardTagBuilder[T]) SetSlotMainMenu(child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:mainMenu", true).Attr())
}

func (b *VXAdvancedCardTagBuilder[T]) SetSlotMainMenuContainer(child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:mainMenuContainer", true).Attr())
}

func (b *VXAdvancedCardTagBuilder[T]) SetSlotSecondaryMenu(child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:secondaryMenu", true).Attr())
}

func (b *VXAdvancedCardTagBuilder[T]) SetSlotSecondaryMenuContainer(child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:secondaryMenuContainer", true).Attr())
}

func (b *VXAdvancedCardTagBuilder[T]) SetSlotPrependToolbar(child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:prependToolbar", true).Attr())
}

func (b *VXAdvancedCardTagBuilder[T]) SetSlotAppendToolbar(child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:appendToolbar", true).Attr())
}

func (b *VXAdvancedCardTagBuilder[T]) SetSlotHeader(child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:header", true).Attr())
}

func (b *VXAdvancedCardTagBuilder[T]) SetSlotPrependHeader(child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:prependHeader", true).Attr())
}

func (b *VXAdvancedCardTagBuilder[T]) SetSlotAppendHeader(child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:appendHeader", true).Attr())
}

func (b *VXAdvancedCardTagBuilder[T]) SetSlotTop(child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:top", true).Attr())
}

func (b *VXAdvancedCardTagBuilder[T]) SetSlotBottom(child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:bottom", "{isActive}").Attr())
}

func (b *VXAdvancedCardTagBuilder[T]) SetScopedSlotBottom(scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:bottom", scope).Attr())
}

func (b *VXAdvancedCardTagBuilder[T]) SetSlotBody(child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:body", true).Attr())
}

func (b *VXAdvancedCardTagBuilder[T]) SetScopedSlotBody(scope string, child ...h.HTMLComponent) {
	b.AppendChild(web.Slot(child...).Name("body").Scope(scope))
}

func (b *VXAdvancedCardTagBuilder[T]) SetSlotPortals(child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:portals", true).Attr())
}
