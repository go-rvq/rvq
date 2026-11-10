package presets

import (
	"testing"

	"github.com/go-rvq/rvq/x/i18n"
)

type leParent struct {
	ID    uint
	Items []*leItem
}

type leItem struct {
	ID      uint
	Agencia string
}

// The fields of the items of a list editor of a nested slice are of the
// items' model — their labels <Item><Field> —, a child of the parent's (an
// item of the slice for a row); without a nested model, the parent's.
func TestListEditorItemInfo(t *testing.T) {
	b := New(i18n.New())
	parent := b.Model(&leParent{}).URIName("le-parents")
	item := NewModelBuilder(b, &leItem{}).Label("Forma")
	fb := &item.Editing("Agencia").FieldsBuilder
	obj := &leParent{Items: []*leItem{{Agencia: "1"}}}
	fc := &FieldContext{Name: "Items", Obj: obj, ModelInfo: parent.Info(), Nested: NestedSlice(item, fb)}

	if info := listEditorItemInfo(fc, -1); info.mb != item {
		t.Fatalf("the header: of %q, not of the items' model", info.mb.label)
	} else if p, pobj := info.Parent(); p == nil || p.mb != parent || pobj != obj {
		t.Errorf("not a child of the parent: %v %v", p, pobj)
	}
	if info := listEditorItemInfo(fc, 0); info.mb != item || info.index != 0 || info.slice == nil {
		t.Errorf("a row: %+v", info)
	}
	fc.Nested = nil
	if info := listEditorItemInfo(fc, -1); info != parent.Info() {
		t.Errorf("no nested model: not the parent's")
	}
}
