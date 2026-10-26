package presets

import (
	"reflect"
	"testing"
)

// A select with no labels of its own shows its keys — no labels function
// is no crash (the mail sender's Sender: "GMAIL", "SMTP").
func TestSelectConfigKeyLabels(t *testing.T) {
	keys := []string{"GMAIL", "SMTP"}
	if got := (&SelectConfig{}).KeyLabels(nil, keys); !reflect.DeepEqual(got, keys) {
		t.Errorf("no labels function: %v", got)
	}
	withHints := &SelectConfig{KeyLabelsAndHintsFunc: func(_ *FieldContext, key []string) (r [][2]string) {
		for _, k := range key {
			r = append(r, [2]string{"L" + k, "H"})
		}
		return
	}}
	if got := withHints.KeyLabels(nil, keys); !reflect.DeepEqual(got, []string{"LGMAIL", "LSMTP"}) {
		t.Errorf("labels and hints: %v", got)
	}
}
