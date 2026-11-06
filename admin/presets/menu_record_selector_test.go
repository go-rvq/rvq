package presets

import (
	"testing"

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
