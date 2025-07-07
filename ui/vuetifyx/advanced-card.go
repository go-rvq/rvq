package vuetifyx

import (
	"github.com/qor5/web/v3/tag"
	v "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type (
	VXAdvancedCardTagBuilder[T any] struct {
		v.VTagBuilder[T]
	}

	VXAnyAdvancedCardTagBuilder struct {
		B h.HTMLComponent
	}
)

func VXAdvancedCardTag[T VXAdvancedCardTagGetter[T]](dot T, name string, children ...h.HTMLComponent) T {
	vtb := dot.GetVXAdvancedCardTagBuilder()
	vtb.TagBuilder = *tag.NewTag(dot, name, children...).GetTagBuilder()
	return dot
}

func (b *VXAdvancedCardTagBuilder[T]) GetVXAdvancedCardTagBuilder() *VXAdvancedCardTagBuilder[T] {
	return b
}

func (b *VXAdvancedCardTagBuilder[T]) GetVXAnyAdvancedCardTagBuilder() *VXAnyAdvancedCardTagBuilder {
	return &VXAnyAdvancedCardTagBuilder{B: any(b.Dot()).(h.HTMLComponent)}
}

func (b *VXAdvancedCardTagBuilder[T]) Density(s string) T {
	b.SetDensity(s)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) ToolbarProps(s string) T {
	b.SetToolbarProps(s)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) ContainerProps(s string) T {
	b.SetContainerProps(s)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) MainMenuUp(v bool) T {
	b.SetMainMenuUp(v)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SecondaryMenuUp(v bool) T {
	b.SetSecondaryMenuUp(v)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) Title(s string) T {
	b.SetTitle(s)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) MainMenuTitle(s string) T {
	b.SetMainMenuTitle(s)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SecondaryMenuTitle(s string) T {
	b.SetSecondaryMenuTitle(s)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SlotMainMenu(child ...h.HTMLComponent) T {
	b.SetSlotMainMenu(child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SlotMainMenuContainer(child ...h.HTMLComponent) T {
	b.SetSlotMainMenuContainer(child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SlotSecondaryMenu(child ...h.HTMLComponent) T {
	b.SetSlotSecondaryMenu(child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SlotSecondaryMenuContainer(child ...h.HTMLComponent) T {
	b.SetSlotSecondaryMenuContainer(child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SlotPrependToolbar(child ...h.HTMLComponent) T {
	b.SetSlotPrependToolbar(child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SlotAppendToolbar(child ...h.HTMLComponent) T {
	b.SetSlotAppendToolbar(child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SlotHeader(child ...h.HTMLComponent) T {
	b.SetSlotHeader(child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SlotPrependHeader(child ...h.HTMLComponent) T {
	b.SetSlotPrependHeader(child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SlotAppendHeader(child ...h.HTMLComponent) T {
	b.SetSlotAppendHeader(child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SlotTop(child ...h.HTMLComponent) T {
	b.SetSlotTop(child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SlotBottom(child ...h.HTMLComponent) T {
	b.SetSlotBottom(child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) ScopedSlotBottom(scope string, child ...h.HTMLComponent) T {
	b.SetScopedSlotBottom(scope, child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SlotBody(child ...h.HTMLComponent) T {
	b.SetSlotBody(child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) ScopedSlotBody(scope string, child ...h.HTMLComponent) T {
	b.SetScopedSlotBody(scope, child...)
	return b.Dot()
}

func (b *VXAdvancedCardTagBuilder[T]) SlotPortals(child ...h.HTMLComponent) T {
	b.SetSlotPortals(child...)
	return b.Dot()
}
