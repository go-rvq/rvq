package presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

// Built-in side-menu item names.
const (
	// PreMenuItemLanguageSwitch is the admin-language selector, registered as a
	// pre-menu item by New.
	PreMenuItemLanguageSwitch = "LanguageSwitch"
	// PreMenuItemMenuFilter is the side-menu text filter, registered as a pre-menu
	// item by New.
	PreMenuItemMenuFilter = "MenuFilter"
)

// SideMenuItem is a component rendered before (pre) or after (post) the model
// menu in the navigation drawer. Handler may return nil to render nothing (e.g.
// the language selector with a single language), and Enabled toggles it off
// without removing it.
type SideMenuItem struct {
	Name    string
	Handler ComponentFunc
	Enabled bool
}

// AddPreMenuItem appends items shown before the model menu.
func (b *Builder) AddPreMenuItem(items ...*SideMenuItem) *Builder {
	b.preMenuItems = append(b.preMenuItems, items...)
	return b
}

// GetPreMenuItem returns the pre-menu item with the given name (nil if absent).
// The returned pointer can be mutated (e.g. toggle Enabled).
func (b *Builder) GetPreMenuItem(name string) *SideMenuItem {
	return findSideMenuItem(b.preMenuItems, name)
}

// SetPreMenuOrder sets the order (by name) the pre-menu items render in; names
// not listed keep their insertion order, after the listed ones.
func (b *Builder) SetPreMenuOrder(order []string) *Builder {
	b.preMenuOrder = order
	return b
}

func (b *Builder) preMenu(ctx *web.EventContext) h.HTMLComponent {
	return renderSideMenu(ctx, b.preMenuItems, b.preMenuOrder)
}

// AddPostMenuItem appends items shown after the model menu.
func (b *Builder) AddPostMenuItem(items ...*SideMenuItem) *Builder {
	b.postMenuItems = append(b.postMenuItems, items...)
	return b
}

// GetPostMenuItem returns the post-menu item with the given name (nil if absent).
func (b *Builder) GetPostMenuItem(name string) *SideMenuItem {
	return findSideMenuItem(b.postMenuItems, name)
}

// SetPostMenuOrder sets the order (by name) the post-menu items render in.
func (b *Builder) SetPostMenuOrder(order []string) *Builder {
	b.postMenuOrder = order
	return b
}

func (b *Builder) postMenu(ctx *web.EventContext) h.HTMLComponent {
	return renderSideMenu(ctx, b.postMenuItems, b.postMenuOrder)
}

func findSideMenuItem(items []*SideMenuItem, name string) *SideMenuItem {
	for _, it := range items {
		if it.Name == name {
			return it
		}
	}
	return nil
}

// renderSideMenu renders the enabled items: those named in order first (in that
// order), then the remaining ones in insertion order. A nil handler result is
// skipped, so an item can hide itself.
func renderSideMenu(ctx *web.EventContext, items []*SideMenuItem, order []string) h.HTMLComponent {
	var comps h.HTMLComponents
	seen := make(map[string]bool, len(items))
	byName := make(map[string]*SideMenuItem, len(items))
	for _, it := range items {
		byName[it.Name] = it
	}
	emit := func(it *SideMenuItem) {
		if it == nil || !it.Enabled || it.Handler == nil || seen[it.Name] {
			return
		}
		seen[it.Name] = true
		if c := it.Handler(ctx); c != nil {
			comps = append(comps, c)
		}
	}
	for _, name := range order {
		emit(byName[name])
	}
	for _, it := range items {
		emit(it)
	}
	return comps
}
