package presets

import (
	"slices"
	"testing"
)

type hiddenFieldsRecord struct {
	ID        uint
	Title     string
	Body      string
	OnlineUrl string
}

func newHiddenFieldsBuilder() *FieldsBuilder {
	b := NewFieldsBuilder(nil).Model(&hiddenFieldsRecord{})
	b.Field("ID")
	b.Field("Title")
	b.Field("Body")
	return b.HiddenField("OnlineUrl")
}

// hiddenComponent reports whether the builder still carries the hidden-input
// component HiddenField installed for name. Losing it is what turns a hidden
// field into a visible one: the render falls back to a default field.
func hiddenComponent(b *FieldsBuilder, name string) bool {
	f := b.GetField(name)
	return f != nil && f.compFunc != nil
}

// Only used to drop the builders of hidden fields while keeping their names in
// the hidden list. The render walks that list outside the layout, found no
// builder, and made a default — visible — field, which is how OnlineUrl turned
// up on a "New Page" form that asked for three other fields.
func TestOnlyKeepsHiddenFields(t *testing.T) {
	b := newHiddenFieldsBuilder().Only("Title", "Body")

	if !slices.Contains(b.hiddenFields, "OnlineUrl") {
		t.Fatalf("hiddenFields = %v, want it to still hold OnlineUrl", b.hiddenFields)
	}
	if !hiddenComponent(b, "OnlineUrl") {
		t.Error("OnlineUrl lost its hidden component, so it would render as a normal field")
	}
	// The point of Only still holds.
	if b.GetField("ID") != nil {
		t.Error("ID survived Only(\"Title\", \"Body\")")
	}
}

// A hidden field named in Only is kept once, not twice.
func TestOnlyKeepsHiddenFieldOnce(t *testing.T) {
	b := newHiddenFieldsBuilder().Only("Title", "OnlineUrl")

	var n int
	for _, f := range b.fields {
		if f.name == "OnlineUrl" {
			n++
		}
	}
	if n != 1 {
		t.Errorf("OnlineUrl appears %d times in fields, want 1", n)
	}
}

// Excluding a hidden field has to drop it from the hidden list too: left there,
// the render rebuilds it from the defaults and it shows up visible.
func TestExceptDropsHiddenField(t *testing.T) {
	b := newHiddenFieldsBuilder().Except("OnlineUrl")

	if slices.Contains(b.hiddenFields, "OnlineUrl") {
		t.Errorf("hiddenFields = %v, want OnlineUrl gone", b.hiddenFields)
	}
	if b.GetField("OnlineUrl") != nil {
		t.Error("OnlineUrl survived Except")
	}
}

// Excluding something else leaves the hidden field alone.
func TestExceptKeepsOtherHiddenFields(t *testing.T) {
	b := newHiddenFieldsBuilder().Except("Body")

	if !slices.Contains(b.hiddenFields, "OnlineUrl") {
		t.Errorf("hiddenFields = %v, want it to still hold OnlineUrl", b.hiddenFields)
	}
	if !hiddenComponent(b, "OnlineUrl") {
		t.Error("OnlineUrl lost its hidden component")
	}
}
