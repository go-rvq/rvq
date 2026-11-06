package presets

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

type recSelA struct{ ID uint }
type recSelB struct{ ID uint }
type recSelC struct{ ID uint }

// The record chosen in the menu: on for a model of no parent, a child of a
// singleton, or a child of a model that has it too.
func TestHasMenuRecordSelector(t *testing.T) {
	p := New(i18n.New())
	top := p.Model(&recSelA{})
	child := NewModelBuilder(p, &recSelB{})
	top.AddChild(child)
	grand := NewModelBuilder(p, &recSelC{})
	child.AddChild(grand)

	if top.HasMenuRecordSelector() {
		t.Error("not asked: off")
	}
	top.MenuRecordSelector(true)
	if !top.HasMenuRecordSelector() {
		t.Error("of no parent: on")
	}
	child.MenuRecordSelector(true)
	grand.MenuRecordSelector(true)
	if !child.HasMenuRecordSelector() || !grand.HasMenuRecordSelector() {
		t.Error("under models that have it: on")
	}
	// its parent without it: off, and so its children's
	top.MenuRecordSelector(false)
	if child.HasMenuRecordSelector() || grand.HasMenuRecordSelector() {
		t.Error("under a parent without it: off")
	}
	// a child of a singleton: on
	top.Singleton(true)
	if !child.HasMenuRecordSelector() || !grand.HasMenuRecordSelector() {
		t.Error("under a singleton: on")
	}
	// External: its choice is its owner's, no action of its own
	ext := p.Model(&recSelC{}).URIName("ext")
	ext.MenuRecordSelectorConfig().External = true
	ext.MenuRecordSelector(true)
	if !ext.HasMenuRecordSelector() || getAction(ext.Listing().actions, ActionChooseMenuRecord) != nil {
		t.Error("external: on, without the action that chooses")
	}
	if getAction(child.Listing().actions, ActionChooseMenuRecord) == nil {
		t.Error("the action that chooses its record")
	}
}

// The side menu: a chip of an item at its right; a chip inline, in the place
// of its title, its clear at the right of the item (out of the chip).
func TestMenuChipInlineAndClear(t *testing.T) {
	p := New(i18n.New())
	chip := &MenuChip{Title: "Obra 1", Inline: true, Clear: &MenuChipAction{Icon: "mdi-close", Event: "clear"}}
	p.AddMenuRootFunc(func(ctx *web.EventContext) ([]*MenuNode, []string, string) {
		n := &MenuNode{Title: "Obra 1", Value: "g:x"}
		n.SetChip(chip)
		return []*MenuNode{n}, nil, ""
	})
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/", nil)}
	got := h.MustString(p.CreateMenus(ctx), context.TODO())
	title := got[strings.Index(got, "v-slot:title"):]
	title = title[:strings.Index(title, "</template>")]
	if !strings.Contains(title, "item.props.chip.inline") || !strings.Contains(title, "rvq-menu-chip") ||
		strings.Contains(title, "rvq-menu-chip-clear") {
		t.Errorf("the title: the chip, inline, without its clear:\n%s", title)
	}
	appendSlot := got[strings.Index(got, "v-slot:append"):]
	if !strings.Contains(appendSlot, "rvq-menu-chip-clear") || !strings.Contains(appendSlot, "!item.props.chip.inline") {
		t.Errorf("the right of the item: its chip (not inline) or its clear:\n%.800s", appendSlot)
	}
	// the node: its chip inline, its clear
	if !strings.Contains(got, `"inline":true`) || !strings.Contains(got, `"clear":{"icon":"mdi-close"`) {
		t.Errorf("the node of the menu:\n%.600s", got)
	}
}
