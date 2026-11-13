package presets

import (
	"strings"
	"testing"

	"github.com/go-rvq/rvq/x/i18n"
)

type renameItemModel struct{ ID uint }

// Renamed, an item keeps its place in the menu, under its new name.
func TestRowMenuRename(t *testing.T) {
	rmb := New(i18n.New()).Model(&renameItemModel{}).Listing().RowMenu()
	rmb.RowMenuItem("Edit")
	rmb.RowMenuItem("Localize")
	rmb.RowMenuItem("Delete")
	names := func() string {
		var ns []string
		for _, it := range rmb.Items() {
			ns = append(ns, it.Name())
		}
		return strings.Join(ns, ",")
	}
	before := names()
	if rmb.Rename("Localize", "Translate") == nil {
		t.Fatal("Localize not found")
	}
	if got, want := names(), strings.Replace(before, "Localize", "Translate", 1); got != want {
		t.Errorf("the menu: %s, want %s", got, want)
	}
	if rmb.Rename("Localize", "Other") != nil {
		t.Error("renamed an item that is gone")
	}
}
